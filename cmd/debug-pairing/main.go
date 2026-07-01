package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "modernc.org/sqlite"
)

func main() {
	session := flag.String("session", "passkey-debug", "debug session name")
	storeDir := flag.String("store", "./state/debug-pairing", "directory for the isolated SQLite session store")
	qrPNG := flag.String("qr-png", "./state/debug-pairing/qr.png", "path where the latest QR PNG is written")
	logLevel := flag.String("log-level", "DEBUG", "whatsmeow log level")
	timeout := flag.Duration("timeout", 0, "optional runtime limit, for example 5m; 0 waits until interrupt")
	dumpQR := flag.Bool("dump-qr", false, "print raw QR payloads to stdout")
	refreshVersion := flag.Bool("refresh-version", true, "try to fetch and apply the latest WhatsApp Web version before connecting")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if *timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}

	if err := run(ctx, strings.TrimSpace(*session), *storeDir, *qrPNG, strings.ToUpper(*logLevel), *dumpQR, *refreshVersion); err != nil {
		fmt.Fprintf(os.Stderr, "debug-pairing error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, session, storeDir, qrPNG, logLevel string, dumpQR, refreshVersion bool) error {
	if session == "" {
		return fmt.Errorf("--session cannot be empty")
	}
	sessionDir := filepath.Join(storeDir, session)
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		return fmt.Errorf("create session dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(qrPNG), 0o755); err != nil {
		return fmt.Errorf("create QR dir: %w", err)
	}

	dbLog := waLog.Stdout("DebugPairing/Database", logLevel, true)
	clientLog := waLog.Stdout("DebugPairing/Client", logLevel, true)
	qrLog := waLog.Stdout("DebugPairing/QR", logLevel, true)

	dbPath := filepath.Join(sessionDir, "session.db")
	dsn := fmt.Sprintf(
		"file:%s?cache=shared&_pragma=foreign_keys(1)&_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)",
		dbPath,
	)
	container, err := sqlstore.New(ctx, "sqlite", dsn, dbLog)
	if err != nil {
		return fmt.Errorf("open whatsmeow store: %w", err)
	}
	defer container.Close()

	device, err := container.GetFirstDevice(ctx)
	if err != nil {
		return fmt.Errorf("get first device: %w", err)
	}
	if device == nil {
		device = container.NewDevice()
	}

	if refreshVersion {
		refreshWhatsAppWebVersion(ctx)
	}

	client := whatsmeow.NewClient(device, clientLog)
	client.EnableAutoReconnect = true
	client.AutoTrustIdentity = true
	client.PrePairCallback = func(jid types.JID, platform, businessName string) bool {
		fmt.Printf("\n[%s] PRE_PAIR_CALLBACK jid=%s platform=%q business_name=%q\n", now(), jid, platform, businessName)
		printState(client)
		return true
	}
	client.AddEventHandler(func(evt any) {
		dumpEvent(client, evt)
	})

	fmt.Printf("[%s] debug session=%s store=%s qr_png=%s wa_version=%s\n", now(), session, dbPath, qrPNG, store.GetWAVersion().String())
	printState(client)

	if client.Store != nil && client.Store.ID == nil {
		qrCh, err := client.GetQRChannel(ctx)
		if err != nil {
			return fmt.Errorf("get QR channel: %w", err)
		}
		go watchQR(qrCh, qrPNG, dumpQR, qrLog)
	} else {
		fmt.Printf("[%s] store already has device ID; connecting without a fresh QR\n", now())
	}

	if err := client.Connect(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	fmt.Printf("[%s] connect started; scan the QR PNG if this is a fresh session, then leave this process running\n", now())

	<-ctx.Done()
	fmt.Printf("[%s] shutting down: %v\n", now(), ctx.Err())
	client.Disconnect()
	return nil
}

func refreshWhatsAppWebVersion(ctx context.Context) {
	previous := store.GetWAVersion()
	latest, err := whatsmeow.GetLatestVersion(ctx, nil)
	if err != nil {
		fmt.Printf("[%s] could not refresh WhatsApp Web version; keeping %s: %v\n", now(), previous.String(), err)
		return
	}
	store.SetWAVersion(*latest)
	fmt.Printf("[%s] WhatsApp Web version %s -> %s\n", now(), previous.String(), latest.String())
}

func watchQR(ch <-chan whatsmeow.QRChannelItem, qrPNG string, dumpQR bool, log waLog.Logger) {
	for item := range ch {
		fmt.Printf("\n[%s] QR_CHANNEL event=%q timeout=%s error=%v\n", now(), item.Event, item.Timeout, item.Error)
		if item.Event != whatsmeow.QRChannelEventCode {
			continue
		}
		if dumpQR {
			fmt.Printf("[%s] QR_RAW %s\n", now(), item.Code)
		}
		if err := qrcode.WriteFile(item.Code, qrcode.Medium, 512, qrPNG); err != nil {
			fmt.Printf("[%s] QR_PNG_WRITE_ERROR path=%s err=%v\n", now(), qrPNG, err)
			continue
		}
		log.Infof("Wrote QR PNG to %s", qrPNG)
		fmt.Printf("[%s] QR_PNG_READY %s\n", now(), qrPNG)
	}
	fmt.Printf("[%s] QR_CHANNEL closed\n", now())
}

func dumpEvent(client *whatsmeow.Client, evt any) {
	typ := "<nil>"
	if evt != nil {
		typ = reflect.TypeOf(evt).String()
	}

	fmt.Printf("\n[%s] EVENT %s\n", now(), typ)
	if data, err := json.MarshalIndent(evt, "", "  "); err == nil && string(data) != "{}" {
		fmt.Printf("json=%s\n", data)
	}
	fmt.Printf("go=%#v\n", evt)
	printImportantEventDetails(evt)
	printState(client)
}

func printImportantEventDetails(evt any) {
	switch v := evt.(type) {
	case *events.PairSuccess:
		fmt.Printf("pair_success id=%s lid=%s platform=%q business_name=%q\n", v.ID, v.LID, v.Platform, v.BusinessName)
	case *events.PairError:
		fmt.Printf("pair_error id=%s lid=%s platform=%q business_name=%q error=%v\n", v.ID, v.LID, v.Platform, v.BusinessName, v.Error)
	case *events.PairPasskeyRequest:
		data, _ := json.Marshal(v.PublicKey)
		fmt.Printf("pair_passkey_request public_key=%s\n", data)
	case *events.PairPasskeyConfirmation:
		fmt.Printf("pair_passkey_confirmation code=%s skip_handoff_ux=%t\n", v.Code, v.SkipHandoffUX)
	case *events.PairPasskeyError:
		fmt.Printf("pair_passkey_error continuation=%t error=%v\n", v.Continuation, v.Error)
	case *events.ConnectFailure:
		fmt.Printf("connect_failure reason=%s message=%q raw=%s\n", v.Reason.String(), v.Message, v.Raw)
	case *events.StreamError:
		fmt.Printf("stream_error code=%q raw=%s\n", v.Code, v.Raw)
	case *events.LoggedOut:
		fmt.Printf("logged_out on_connect=%t reason=%s\n", v.OnConnect, v.Reason.String())
	case *events.TemporaryBan:
		fmt.Printf("temporary_ban code=%s expire=%s\n", v.Code.String(), v.Expire)
	case *events.QR:
		fmt.Printf("qr_event codes=%d\n", len(v.Codes))
	case *events.Connected:
		fmt.Println("connected_event authenticated=true")
	case *events.Disconnected:
		fmt.Println("disconnected_event websocket_closed=true")
	case *events.ClientOutdated:
		fmt.Printf("client_outdated current_version=%s\n", store.GetWAVersion().String())
	}
}

func printState(client *whatsmeow.Client) {
	var id, lid string
	if client != nil && client.Store != nil {
		if client.Store.ID != nil {
			id = client.Store.ID.String()
		}
		if !client.Store.LID.IsEmpty() {
			lid = client.Store.LID.String()
		}
	}
	fmt.Printf(
		"state connected=%t logged_in=%t store_id=%q store_lid=%q\n",
		client != nil && client.IsConnected(),
		client != nil && client.IsLoggedIn(),
		id,
		lid,
	)
}

func now() string {
	return time.Now().Format(time.RFC3339Nano)
}
