module your.org/provider-whatsmeow

go 1.24.0

toolchain go1.24.6

require (
	github.com/gorilla/mux v1.8.1
	github.com/hajimehoshi/go-mp3 v0.3.4
	github.com/nyaruka/phonenumbers v1.2.1
	github.com/pion/rtp v1.8.7
	github.com/pion/webrtc/v3 v3.3.6
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	go.mau.fi/whatsmeow v0.0.0-20250905121447-8d6da61ecbfa
	google.golang.org/protobuf v1.36.8
	layeh.com/gopus v0.0.0-20210501142526-1ee02d434e32
	modernc.org/sqlite v1.38.2
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/beeper/argo-go v1.1.2 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/elliotchance/orderedmap/v3 v3.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/petermattis/goid v0.0.0-20250813065127-a731cc31b4fe // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/vektah/gqlparser/v2 v2.5.30 // indirect
	go.mau.fi/libsignal v0.2.0 // indirect
	go.mau.fi/util v0.9.0 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/exp v0.0.0-20250813145105-42675adae3e6 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	modernc.org/libc v1.66.3 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)

// Use the GitHub mirror for whatsmeow while keeping the canonical import path.
replace go.mau.fi/whatsmeow => github.com/tulir/whatsmeow v0.0.0-20250905121447-8d6da61ecbfa
