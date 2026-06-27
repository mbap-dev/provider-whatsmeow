// Copyright (c) 2021 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package whatsmeow

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	waBinary "go.mau.fi/whatsmeow/binary"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func (cli *Client) handleCallEvent(ctx context.Context, node *waBinary.Node) {
	defer cli.maybeDeferredAck(ctx, node)()

	if len(node.GetChildren()) != 1 {
		cli.dispatchEvent(&events.UnknownCallEvent{Node: node})
		return
	}
	ag := node.AttrGetter()
	child := node.GetChildren()[0]
	cag := child.AttrGetter()
	basicMeta := types.BasicCallMeta{
		From:        ag.JID("from"),
		Timestamp:   ag.UnixTime("t"),
		CallCreator: cag.JID("call-creator"),
		CallID:      cag.String("call-id"),
		GroupJID:    cag.OptionalJIDOrEmpty("group-jid"),
	}
	if basicMeta.CallCreator.Server == types.HiddenUserServer {
		basicMeta.CallCreatorAlt = cag.OptionalJIDOrEmpty("caller_pn")
	} else {
		// This may not actually exist
		basicMeta.CallCreatorAlt = cag.OptionalJIDOrEmpty("caller_lid")
	}
	switch child.Tag {
	case "offer":
		cli.dispatchEvent(&events.CallOffer{
			BasicCallMeta: basicMeta,
			CallRemoteMeta: types.CallRemoteMeta{
				RemotePlatform: ag.String("platform"),
				RemoteVersion:  ag.String("version"),
			},
			Data: &child,
		})
	case "offer_notice":
		cli.dispatchEvent(&events.CallOfferNotice{
			BasicCallMeta: basicMeta,
			Media:         cag.String("media"),
			Type:          cag.String("type"),
			Data:          &child,
		})
	case "relaylatency":
		cli.dispatchEvent(&events.CallRelayLatency{
			BasicCallMeta: basicMeta,
			Data:          &child,
		})
	case "accept":
		cli.dispatchEvent(&events.CallAccept{
			BasicCallMeta: basicMeta,
			CallRemoteMeta: types.CallRemoteMeta{
				RemotePlatform: ag.String("platform"),
				RemoteVersion:  ag.String("version"),
			},
			Data: &child,
		})
	case "preaccept":
		cli.dispatchEvent(&events.CallPreAccept{
			BasicCallMeta: basicMeta,
			CallRemoteMeta: types.CallRemoteMeta{
				RemotePlatform: ag.String("platform"),
				RemoteVersion:  ag.String("version"),
			},
			Data: &child,
		})
	case "transport":
		cli.dispatchEvent(&events.CallTransport{
			BasicCallMeta: basicMeta,
			CallRemoteMeta: types.CallRemoteMeta{
				RemotePlatform: ag.String("platform"),
				RemoteVersion:  ag.String("version"),
			},
			Data: &child,
		})
	case "terminate":
		cli.dispatchEvent(&events.CallTerminate{
			BasicCallMeta: basicMeta,
			Reason:        cag.String("reason"),
			Data:          &child,
		})
	case "reject":
		cli.dispatchEvent(&events.CallReject{
			BasicCallMeta: basicMeta,
			Data:          &child,
		})
	default:
		cli.dispatchEvent(&events.UnknownCallEvent{Node: node})
	}
}

// OfferCall offers a call to a user.
func (cli *Client) OfferCall(ctx context.Context, callTo types.JID, video bool) (string, error) {
	return cli.OfferCallWithDeviceMode(ctx, callTo, video, "")
}

// OfferCallWithDeviceMode offers a call to a user.
//
// deviceMode controls call offer encryption fanout:
//   - "direct" emits a modern call offer with a direct <enc> child and LID
//     addressing when possible, matching call offers received from mobile apps.
//   - "peer" encrypts only to the current sender device and target primary JID,
//     matching the public OfferCall snippet.
//   - "all" expands all sender/target devices like regular messages.
func (cli *Client) OfferCallWithDeviceMode(ctx context.Context, callTo types.JID, video bool, deviceMode string) (string, error) {
	ownID := cli.getOwnID()
	if ownID.IsEmpty() {
		return "", ErrNotLoggedIn
	}
	if callTo.IsEmpty() {
		return "", fmt.Errorf("empty call recipient")
	}
	callTo = callTo.ToNonAD()
	deviceMode = strings.ToLower(strings.TrimSpace(deviceMode))
	if deviceMode == "" {
		deviceMode = "direct"
	}
	callID := strings.ToUpper(hex.EncodeToString(randomBytes(16)))
	plaintext, dsmPlaintext, err := marshalMessage(callTo, &waE2E.Message{Call: &waE2E.Call{CallKey: randomBytes(32)}})
	if err != nil {
		return "", fmt.Errorf("failed to marshal call: %w", err)
	}
	var callDevices []types.JID
	var directEnc *waBinary.Node
	var directIdentity *waBinary.Node
	ownPN := ownID.ToNonAD()
	callFrom := ownPN
	callerPN := ownPN
	if ownLID := cli.getOwnLID().ToNonAD(); !ownLID.IsEmpty() {
		callFrom = ownLID
	}
	callToWire := callTo
	if deviceMode == "direct" {
		callToWire = cli.resolveCallWireJID(ctx, callTo)
		directEnc, directIdentity, err = cli.encryptDirectCallOffer(ctx, callToWire, callID, plaintext)
		if err != nil {
			return "", err
		}
	}
	switch deviceMode {
	case "direct":
		callDevices = []types.JID{callToWire}
	case "peer":
		callDevices = []types.JID{ownID, callTo}
	case "all":
		callDevices, err = cli.GetUserDevices(ctx, []types.JID{ownID.ToNonAD(), callTo})
		if err != nil {
			return "", fmt.Errorf("failed to get call devices: %w", err)
		}
	default:
		return "", fmt.Errorf("unsupported call device_mode %q", deviceMode)
	}
	offerContent := []waBinary.Node{
		{Tag: "audio", Attrs: waBinary.Attrs{"enc": "opus", "rate": "16000"}},
		{Tag: "audio", Attrs: waBinary.Attrs{"enc": "opus", "rate": "8000"}},
	}
	if video {
		offerContent = append(offerContent, waBinary.Node{
			Tag: "video",
			Attrs: waBinary.Attrs{
				"orientation":        "0",
				"screen_width":       "1080",
				"screen_height":      "2340",
				"device_orientation": "0",
				"enc":                "vp8",
				"dec":                "vp8",
			},
		})
	}
	offerContent = append(offerContent, waBinary.Node{Tag: "capability", Attrs: waBinary.Attrs{"ver": "1"}, Content: []byte{1, 5, 255, 9, 228, 250, 19}})
	if deviceMode == "direct" {
		offerContent = append(offerContent, *directEnc)
		if directIdentity != nil {
			offerContent = append(offerContent, *directIdentity)
		}
		offerContent = append(offerContent,
			waBinary.Node{Tag: "encopt", Attrs: waBinary.Attrs{"keygen": "2"}},
			waBinary.Node{Tag: "metadata", Attrs: waBinary.Attrs{"peer_abtest_bucket": "provider_whatsmeow_fakecall", "peer_abtest_bucket_id_list": "0"}},
			waBinary.Node{Tag: "net", Attrs: waBinary.Attrs{"medium": "3"}},
			waBinary.Node{Tag: "rte", Content: randomBytes(6)},
			waBinary.Node{Tag: "uploadfieldstat"},
			waBinary.Node{Tag: "voip_settings", Attrs: waBinary.Attrs{"uncompressed": "1"}, Content: []byte(`{"options":{"caller_timeout":"90","disable_p2p":"1"},"voip_settings_version":{"version_number":"152090"}}`)},
		)
	} else {
		destinationNode, includeIdentity, err := cli.encryptMessageForDevices(ctx, callDevices, callID, plaintext, dsmPlaintext, nil)
		if err != nil {
			return "", fmt.Errorf("failed to encrypt call offer: %w", err)
		}
		if includeIdentity {
			destinationNode = append(destinationNode, cli.makeDeviceIdentityNode())
		}
		offerContent = append(offerContent,
			waBinary.Node{Tag: "destination", Content: destinationNode},
			waBinary.Node{Tag: "encopt", Attrs: waBinary.Attrs{"keygen": "2"}},
			waBinary.Node{Tag: "net", Attrs: waBinary.Attrs{"medium": "3"}},
		)
	}
	offerAttrs := waBinary.Attrs{"call-id": callID, "call-creator": callFrom}
	if deviceMode == "direct" {
		offerAttrs["caller_pn"] = callerPN
		offerAttrs["joinable"] = "1"
		if strings.HasPrefix(callerPN.User, "55") {
			offerAttrs["caller_country_code"] = "BR"
		}
	}
	node := waBinary.Node{
		Tag:   "call",
		Attrs: waBinary.Attrs{"e": "0", "id": cli.GenerateMessageID(), "from": callFrom, "to": callToWire, "platform": "web", "t": fmt.Sprintf("%d", time.Now().Unix()), "version": "2"},
		Content: []waBinary.Node{{
			Tag:     "offer",
			Attrs:   offerAttrs,
			Content: offerContent,
		}},
	}
	err = cli.sendNode(ctx, node)
	if err != nil {
		return "", err
	}
	return callID, nil
}

func (cli *Client) resolveCallWireJID(ctx context.Context, callTo types.JID) types.JID {
	callTo = callTo.ToNonAD()
	if callTo.Server == types.HiddenUserServer {
		return callTo
	}
	if callTo.Server != types.DefaultUserServer || cli.Store == nil || cli.Store.LIDs == nil {
		return callTo
	}
	lid, err := cli.Store.LIDs.GetLIDForPN(ctx, callTo)
	if err == nil && !lid.IsEmpty() {
		return lid.ToNonAD()
	}
	info, err := cli.GetUserInfo(ctx, []types.JID{callTo})
	if err != nil {
		cli.Log.Debugf("Failed to fetch LID for call recipient %s: %v", callTo, err)
		return callTo
	}
	if userInfo, ok := info[callTo]; ok && !userInfo.LID.IsEmpty() {
		return userInfo.LID.ToNonAD()
	}
	return callTo
}

// ResolveCallRecipient resolves the JID that should be used for call signaling.
func (cli *Client) ResolveCallRecipient(ctx context.Context, callTo types.JID) types.JID {
	return cli.resolveCallWireJID(ctx, callTo)
}

func (cli *Client) encryptDirectCallOffer(ctx context.Context, to types.JID, callID string, plaintext []byte) (*waBinary.Node, *waBinary.Node, error) {
	nodes, includeIdentity, err := cli.encryptMessageForDevices(ctx, []types.JID{to}, callID, plaintext, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt direct call offer: %w", err)
	}
	if len(nodes) == 0 {
		return nil, nil, fmt.Errorf("failed to encrypt direct call offer: no recipient nodes")
	}
	children := nodes[0].GetChildren()
	if len(children) == 0 || children[0].Tag != "enc" {
		return nil, nil, fmt.Errorf("failed to encrypt direct call offer: unexpected node %s", nodes[0].String())
	}
	enc := children[0]
	var identity *waBinary.Node
	if includeIdentity {
		node := cli.makeDeviceIdentityNode()
		identity = &node
	}
	return &enc, identity, nil
}

func randomBytes(length int) []byte {
	data := make([]byte, length)
	if _, err := rand.Read(data); err != nil {
		panic(fmt.Errorf("failed to generate random bytes: %w", err))
	}
	return data
}

// RejectCall reject an incoming call.
func (cli *Client) RejectCall(ctx context.Context, callFrom types.JID, callID string) error {
	ownID := cli.getOwnID()
	if ownID.IsEmpty() {
		return ErrNotLoggedIn
	}
	ownID, callFrom = ownID.ToNonAD(), callFrom.ToNonAD()
	rejectNode := waBinary.Node{
		Tag:     "reject",
		Attrs:   waBinary.Attrs{"call-id": callID, "call-creator": callFrom, "count": "0"},
		Content: nil,
	}
	if token, err := cli.ensureTCToken(ctx, callFrom); err != nil {
		cli.Log.Warnf("Failed to get privacy token for call reject to %s: %v", callFrom, err)
	} else if len(token) > 0 {
		rejectNode.Content = []waBinary.Node{{
			Tag:     "tctoken",
			Content: token,
		}}
	}
	return cli.sendNode(ctx, waBinary.Node{
		Tag:     "call",
		Attrs:   waBinary.Attrs{"id": cli.GenerateMessageID(), "from": ownID, "to": callFrom},
		Content: []waBinary.Node{rejectNode},
	})
}
