package main

/*
#cgo CFLAGS: -std=gnu99
#cgo LDFLAGS: -lopus

#include <string.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdint.h>

#include <opus/opus.h>

//encoder = Opus.opus_encoder_create(freq, ch, Opus.OPUS_APPLICATION_AUDIO, out error);
*/
import "C"
import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	rtc "github.com/chobie/datachannel-go"
)

const (
	sampleRate = 48000
	frequency  = 440.0 // A4
	duration   = 1.0   // 1 second
	twoPi      = 2 * math.Pi
	maxInt16   = 1<<15 - 1
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
var peerConnectionMap2 = make(map[string]*Peer)
var datachannelConnectionMap = make(map[int]*DataChannel)
var config = rtc.RtcConfiguration{}
var webSocket int

func generateSineWave() []C.float {
	numSamples := int(sampleRate * duration)
	buf := make([]C.float, numSamples*2)
	for i := 0; i < numSamples; i++ {
		// Calculate the sample value
		t := float64(i) / sampleRate
		sampleValue := float32(math.Sin(twoPi*frequency*t)) * 0.1 // 0.1 for volume

		// Write the sample to the buffer twice (for stereo interleaved format)
		buf[i] = C.float(sampleValue)
		buf[i*2+1] = C.float(sampleValue)
	}
	return buf
}

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

func generateRandomU32() uint32 {
	for {
		randomBytes := make([]byte, 4)
		_, err := rand.Read(randomBytes)
		if err != nil {
			panic(err)
		}
		result := binary.LittleEndian.Uint32(randomBytes)
		// To ensure it's in the range [1, uint32 max value - 1], if it's 0 or uint32 max value, generate again.
		if result != 0 && result != ^uint32(0) {
			return result
		}
	}
}

func createPeerConnection(config rtc.RtcConfiguration, webSocket int, id string, cb func(*Peer)) *Peer {
	pc := rtc.RtcCreatePeerConnection(&config)
	peer := &Peer{
		Id: id,
		Pc: pc,
	}
	peerConnectionMap[peer.Pc] = peer
	peerConnectionMap2[id] = peer

	rtc.RtcSetUserPointer(peer.Pc, unsafe.Pointer(&pc))

	rtc.RtcSetTrackCallback(peer.Pc, func(pc, tr int, ptr unsafe.Pointer) {
		fmt.Printf("Track Opened!\n")
	})

	if cb != nil {
		cb(peer)
	}

	rtc.RtcSetLocalDescriptionCallback(peer.Pc, func(pc int, sdp string, descriptionType string, ptr unsafe.Pointer) {
		fmt.Printf("onLocalDescription\n")
		peer, _ := peerConnectionMap[*(*int)(ptr)]

		dict := make(map[string]string)
		dict["id"] = peer.Id
		dict["type"] = descriptionType
		dict["description"] = sdp
		fmt.Printf("desc %s\n", sdp)

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
		fmt.Printf("StateChanged %d\n", state)
	})

	rtc.RtcSetGatheringStateChangeCallback(peer.Pc, func(pc int, state int, ptr unsafe.Pointer) {
		fmt.Printf("Gahtering %d\n", state)
	})

	return peer
}

var Track int

func main() {

	sigTrack := make(chan int, 1)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	rtc.RtcPreload()

	rtc.RtcInitLogger(rtc.RTC_LOG_DEBUG, func(level rtc.RtcLogLevel, message string) {
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

			fmt.Printf("received: %s type: %s\n", receievedId, receievedType)
			id, _ := strconv.Atoi(receievedId)
			var peer *Peer
			if target, ok := peerConnectionMap2[receievedId]; ok {
				peer = target
			} else if receievedType == "offer" {
				fmt.Printf("Answering to %d\n", id)
				peer = createPeerConnection(config, webSocket, receievedId, func(peer *Peer) {
				})
			} else {
				fmt.Printf("RETURN! %s not found\n", receievedId)
				return
			}

			if receievedType == "offer" || receievedType == "answer" {
				sdp := dict["description"].(string)
				fmt.Printf("SetRemoteDescription: %s\n", sdp)

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

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter target ID: ")
	targetId, _ := reader.ReadString('\n')
	targetId = strings.Trim(targetId, "\n")
	if targetId == localId {
		panic("targetId must be different from localId")
	}

	fmt.Println(targetId)

	targetPc := createPeerConnection(config, webSocket, targetId, func(pc *Peer) {
		mediaSdp := "audio 9 UDP/TLS/RTP/SAVPF 111\r\na=mid:0\r\na=rtpmap:111 opus/48000/2\r\na=sendrecv\r\na=fmtp:111; cbr=1;maxaveragebitrate=128000;minptime=10;streo=1;useinbandfec=0\r\n"

		ssrc := generateRandomU32()
		// pinit := rtc.RtcTrackInit{
		// 	PayloadType: 111,
		// 	Codec:       rtc.RTC_CODEC_OPUS,
		// 	Direction:   rtc.RTC_DIRECTION_SENDONLY,
		// 	MID:         "audio",
		// 	MSID:        "audio-0",
		// 	SSRC:        ssrc,
		// }

		// tr := rtc.RtcAddTrackEx(pc.Pc, pinit)
		tr := rtc.RtcAddTrack(pc.Pc, mediaSdp)

		if tr < 0 {
			panic("failed to add track")
		}
		Track = tr
		fmt.Printf("addTrack %d\n", tr)

		init := rtc.RtcPacketizationHandlerInit{
			Ssrc:        ssrc,
			ClockRate:   48000,
			PayloadType: 111,
		}
		rtc.RtcSetOpenCallback(Track, func(i int, p unsafe.Pointer) {
			fmt.Printf("Track Opend!\n")
			sigTrack <- i
		})
		rtc.RtcSetClosedCallback(Track, func(i int, p unsafe.Pointer) {
			fmt.Printf("TrackClosed\n")
		})

		rtc.RtcSetOpusPacketizationHandler(Track, init)
		rtc.RtcChainRtcpSrReporter(Track)
		rtc.RtcChainRtcpNackResponder(Track, 1000)
	})

	// dc := rtc.RtcCreateDataChannel(targetPc.Pc, "test")
	// fmt.Printf("Datachannel Created %d\n", dc)
	// rtc.RtcSetClosedCallback(dc, func(i int, p unsafe.Pointer) {
	// 	fmt.Printf("DataChannelClosed")
	// })
	// rtc.RtcSetMessageCallback(dc, func(id int, message []byte, size int, ptr unsafe.Pointer) {
	// 	fmt.Printf("DataChannelMessage")
	// })
	// rtc.RtcSetOpenCallback(dc, func(i int, p unsafe.Pointer) {
	// 	fmt.Printf("DataChannelOpen")
	// 	rtc.RtcSendMessage(dc, []byte("Hello"), -1)
	// })
	rtc.RtcSetLocalDescription(targetPc.Pc, "offer")

	fmt.Printf("Waiting track open\n")
	//<-sigTrack
Loop:
	for {
		fmt.Printf("track wait\n")
		if rtc.RtcIsOpen(Track) {
			break Loop
		}
		time.Sleep(1 * time.Second)

		select {
		case _ = <-sigCh:
			break Loop
		default:
		}

	}
	fmt.Printf("sending sound to %d\n", targetPc.Pc)

	wave := generateSineWave()
	var e C.int
	encoder := C.opus_encoder_create(sampleRate, 2, C.OPUS_APPLICATION_AUDIO, &e)

	frameSize := 240
	waveSlice := (*C.float)(C.malloc(C.size_t(frameSize) * C.size_t(unsafe.Sizeof(C.float(0)))))
	timestamp := 0
	should_break := false
	for !should_break {
		select {
		case <-sigCh:
			should_break = true
			break
		default:
		}

		remaing := len(wave) / 2
		offset := 0

		for remaing > 0 {
			for i := 0; i < frameSize; i++ {
				pointerOffset := unsafe.Pointer(uintptr(unsafe.Pointer(waveSlice)) + uintptr(i)*unsafe.Sizeof(C.float(0)))
				*((*C.float)(pointerOffset)) = wave[offset+i]
			}

			var output = [4000]byte{}
			bytes_encoded := int(C.opus_encode_float(encoder, waveSlice, C.int(240), (*C.uchar)(&output[0]), C.int(4000)))

			//timestamp := rtc.RtcGetCurrentTrackTimestamp(Track)
			timestamp += 240
			//fmt.Printf("%d %d\n", bytes_encoded, timestamp)
			rtc.RtcSetTrackRtpTimestamp(Track, uint32(timestamp))
			rtc.RtcSendMessage(Track, output[0:bytes_encoded], bytes_encoded)

			remaing -= 240
		}

		time.Sleep(1 * time.Second)
	}
	C.free(unsafe.Pointer(waveSlice))

	fmt.Printf("waiting signal\n")
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
