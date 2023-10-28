package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"unsafe"

	rtc "github.com/chobie/datachannel-go"
)

type Peer struct {
	Id string
	Pc int
	Dc int
}

type DataChannel struct {
	Dc    int
	Label string
}

var peerConnectionMap = make(map[int]*Peer)
var datachannelConnectionMap = make(map[int]*DataChannel)
var config = rtc.RtcConfiguration{}
var webSocket int

var characters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomId(length int) string {
	var id bytes.Buffer
	buffer := make([]byte, 1)

	for i := 0; i < length; i++ {
		_, err := rand.Read(buffer)
		if err != nil {
			panic(err)
		}
		index := int(buffer[0]) % len(characters)
		id.WriteByte(characters[index])
	}

	return id.String()
}

func createPeerConnection(config rtc.RtcConfiguration, webSocket int, id string) *Peer {
	pc := rtc.RtcCreatePeerConnection(&config)
	peer := &Peer{
		Id: id,
		Pc: pc,
	}
	peerConnectionMap[peer.Pc] = peer

	rtc.RtcSetUserPointer(peer.Pc, unsafe.Pointer(&pc))

	rtc.RtcSetLocalDescriptionCallback(peer.Pc, func(pc int, sdp string, descriptionType string, ptr unsafe.Pointer) {
		fmt.Printf("onLocalDescription\n")
		peer, _ := peerConnectionMap[*(*int)(ptr)]

		dict := make(map[string]string)
		dict["id"] = peer.Id
		dict["type"] = descriptionType
		dict["description"] = sdp

		json, err := json.Marshal(dict)
		if err != nil {
			return
		}
		rtc.RtcSendMessage(webSocket, json, -1)
		fmt.Printf("sending ws to %s\n", json)
	})

	rtc.RtcSetLocalCandidateCallback(peer.Pc, func(pc int, cand string, mid string, ptr unsafe.Pointer) {
		fmt.Printf("onLocalCandidate\n")
		peer, _ := peerConnectionMap[*(*int)(ptr)]

		dict := make(map[string]string)
		dict["id"] = peer.Id
		dict["type"] = "candidate"
		dict["candidate"] = cand
		dict["mid"] = mid

		json, err := json.Marshal(dict)
		if err != nil {
			return
		}

		rtc.RtcSendMessage(webSocket, json, -1)
		fmt.Printf("sending ws to %s\n", json)
	})

	rtc.RtcSetStateChangeCallback(peer.Pc, func(pc int, state int, ptr unsafe.Pointer) {
		fmt.Printf("StateChanged")
	})

	rtc.RtcSetGatheringStateChangeCallback(peer.Pc, func(pc int, state int, ptr unsafe.Pointer) {
		fmt.Printf("Gahtering")
	})

	rtc.RtcSetDataChannelCallback(peer.Pc, func(pc int, dc int, ptr unsafe.Pointer) {
		peer, _ := peerConnectionMap[*(*int)(ptr)]
		peer.Dc = dc

		var bufferSize = 256
		buffer := make([]byte, bufferSize)
		length := rtc.RtcGetDataChannelLabel(dc, buffer, bufferSize)
		if length < 0 {
			fmt.Printf("rtcGetDataChannelLabel failed\n")
			return
		}

		label := string(buffer)
		fmt.Printf("DataChannel from %s received with label %s\n", peer.Id, label)

		rtc.RtcSetClosedCallback(dc, func(id int, ptr unsafe.Pointer) {
			fmt.Printf("DataChannel from {%d} closed", id)
		})
		rtc.RtcSetMessageCallback(dc, func(id int, message []byte, size int, ptr unsafe.Pointer) {
			if size < 0 {
				gomessage := string(message)
				fmt.Printf("Message from %d receive: %s\n", id, gomessage)
			}
		})
		rtc.RtcSetOpenCallback(dc, func(id int, ptr unsafe.Pointer) {
			rtc.RtcSendMessage(id, []byte("hello from moemo"), -1)
		})

		channel := &DataChannel{
			Dc:    int(dc),
			Label: label,
		}

		datachannelConnectionMap[peer.Pc] = channel
		peerConnectionMap[peer.Pc] = peer

	})
	return peer
}

func main() {

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	rtc.RtcPreload()

	rtc.RtcInitLogger(rtc.RTC_LOG_VERBOSE, func(level rtc.RtcLogLevel, message string) {
		fmt.Printf("%d %s\n", level, message)
	})

	var localId = randomId(4)
	var url = "ws://localhost:8000/" + localId
	webSocket = rtc.RtcCreateWebSocket(url)
	rtc.RtcSetOpenCallback(webSocket, func(id int, ptr unsafe.Pointer) {
		fmt.Printf("WebSocket connected, signaling ready\n")
	})
	rtc.RtcSetErrorCallback(webSocket, func(id int, message string, ptr unsafe.Pointer) {
		fmt.Printf("errorCallbackWS: %d %s\n", id, message)
	})
	rtc.RtcSetMessageCallback(webSocket, func(id int, message []byte, size int, ptr unsafe.Pointer) {
		if size > 0 {
			// binary
			fmt.Printf("received binary message\n")
		} else {
			fmt.Printf("Message: %s\n", string(message))
			var dict map[string]interface{}
			if err := json.Unmarshal([]byte(message), &dict); err != nil {
				fmt.Printf("Json decode failed")
				return
			}

			var receievedId string
			var receievedType string

			if val, ok := dict["id"]; !ok {
				fmt.Print("id not found")
				return
			} else {
				receievedId = val.(string)
			}

			if val, ok := dict["type"]; !ok {
				fmt.Print("type not found")
				return
			} else {
				receievedType = val.(string)
			}

			switch receievedType {
			case "offer":
				break
			case "answer":
				break
			}

			fmt.Printf("received: %s\n", receievedId)
			id, _ := strconv.Atoi(receievedId)
			var peer *Peer
			if target, ok := peerConnectionMap[id]; ok {
				peer = target
			} else if receievedType == "offer" {
				fmt.Printf("Answering to %d\n", id)
				peer = createPeerConnection(config, webSocket, receievedId)
			} else {
				return
			}

			if receievedType == "offer" || receievedType == "answer" {
				sdp := dict["description"].(string)
				rtc.RtcSetRemoteDescription(int(peer.Pc), sdp, receievedType)
			} else if receievedType == "candidate" {
				sdp := dict["candidate"].(string)
				mid := dict["mid"].(string)
				rtc.RtcAddRemoteCandidate(int(peer.Pc), sdp, mid)
			}
		}
	})
	rtc.RtcSetClosedCallback(webSocket, func(id int, ptr unsafe.Pointer) {
		fmt.Printf("WebSocket closed\n")
	})

	fmt.Printf("WebSocket URL is %s\n", url)
	fmt.Printf("Waiting for signaling to be connected..\n")

	_ = <-sigCh

	// cleanup datachannel
	for _, v := range datachannelConnectionMap {
		fmt.Printf("closing datachannelConnection %d\n", int(v.Dc))
		rtc.RtcClose(int(v.Dc))
		rtc.RtcDelete(int(v.Dc))
	}

	// cleanup peerconnection
	for _, v := range peerConnectionMap {
		fmt.Printf("closing peerConnection %d\n", int(v.Pc))
		rtc.RtcClosePeerConnection(int(v.Pc))
		rtc.RtcDeletePeerConnection(int(v.Pc))
	}

	rtc.RtcClose(webSocket)
	rtc.RtcDeleteWebSocket(webSocket)

	rtc.RtcCleanup()
}
