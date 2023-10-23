/**
 * Copyright (c) 2023 Shuhei Tanuma
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package api

/*
#cgo CFLAGS: -std=gnu99
#cgo LDFLAGS: -L. -ldatachannel
#include <rtc/rtc.h>
#include <unistd.h>
#include <stdlib.h>

typedef const char cchar_t;

// forward declarations

void go_rtcLogCallbackFunc(rtcLogLevel level, cchar_t* message);
void go_rtcOpenCallbackFunc(int id, void* ptr);
void go_rtcErrorCallbackFunc(int id, cchar_t* message, void* ptr);
void go_rtcClosedCallbackFunc(int id, void *ptr);
void go_rtcMessageCallbackFunc(int id, cchar_t* message, int size, void *ptr);
void go_rtcDescriptionCallbackFunc(int pc, cchar_t* sdp, cchar_t* descriptionType, void *ptr);
void go_rtcSetLocalCandidateCallback(int pc, cchar_t* cand, cchar_t* mid, void* ptr);
void go_rtcSetStateChangeCallback(int pc, rtcState state, void* ptr);
void go_rtcSetGatheringStateChangeCallback(int pc, rtcGatheringState state, void* ptr);
void go_rtcSetDataChannelCallback(int pc, int dc, void* ptr);
void go_rtcIceStateChangeCallbackFunc(int pc, int dc, void* ptr);
void go_rtcSignalingStateCallbackFunc(int pc, rtcSignalingState state, void* ptr);
void go_rtcTrackCallbackFunc(int pc, rtcSignalingState state, void* ptr);
void go_rtcInterceptorCallbackFunc(int pc, rtcSignalingState state, void* ptr);
void go_rtcSetBufferedAmountLowCallback(int pc, void* ptr);
void go_rtcSetAvailableCallback(int pc, void* ptr);
void go_rtcCreateWebSocketServer(int wsserver, int ws, void* ptr);
*/
import "C"
import (
	"sync"
	"unsafe"
)

// libdatachannel API

type RtcState int

const (
	RTC_NEW          RtcState = C.RTC_NEW
	RTC_CONNECTING            = C.RTC_CONNECTING
	RTC_CONNECTED             = C.RTC_CONNECTED
	RTC_DISCONNECTED          = C.RTC_DISCONNECTED
	RTC_FAILED                = C.RTC_FAILED
	RTC_CLOSED                = C.RTC_CLOSED
)

type RtcIceState int

const (
	RTC_ICE_NEW          RtcIceState = C.RTC_ICE_NEW
	RTC_ICE_CHECKING                 = C.RTC_ICE_CHECKING
	RTC_ICE_CONNECTED                = C.RTC_ICE_CONNECTED
	RTC_ICE_COMPLETED                = C.RTC_ICE_COMPLETED
	RTC_ICE_FAILED                   = C.RTC_ICE_FAILED
	RTC_ICE_DISCONNECTED             = C.RTC_ICE_DISCONNECTED
	RTC_ICE_CLOSED                   = C.RTC_ICE_CLOSED
)

type RtcGatheringState int

const (
	RTC_GATHERING_NEW        RtcGatheringState = C.RTC_GATHERING_NEW
	RTC_GATHERING_INPROGRESS                   = C.RTC_GATHERING_INPROGRESS
	RTC_GATHERING_COMPLETE                     = C.RTC_GATHERING_COMPLETE
)

type RtcSignalingState int

const (
	RTC_SIGNALING_STABLE               RtcSignalingState = C.RTC_SIGNALING_STABLE
	RTC_SIGNALING_HAVE_LOCAL_OFFER                       = C.RTC_SIGNALING_HAVE_LOCAL_OFFER
	RTC_SIGNALING_HAVE_REMOTE_OFFER                      = C.RTC_SIGNALING_HAVE_REMOTE_OFFER
	RTC_SIGNALING_HAVE_LOCAL_PRANSWER                    = C.RTC_SIGNALING_HAVE_LOCAL_PRANSWER
	RTC_SIGNALING_HAVE_REMOTE_PRANSWER                   = C.RTC_SIGNALING_HAVE_REMOTE_PRANSWER
)

type RtcLogLevel int

const (
	RTC_LOG_NONE    RtcLogLevel = C.RTC_LOG_NONE
	RTC_LOG_FATAL               = C.RTC_LOG_FATAL
	RTC_LOG_ERROR               = C.RTC_LOG_ERROR
	RTC_LOG_WARNING             = C.RTC_LOG_WARNING
	RTC_LOG_INFO                = C.RTC_LOG_INFO
	RTC_LOG_DEBUG               = C.RTC_LOG_DEBUG
	RTC_LOG_VERBOSE             = C.RTC_LOG_VERBOSE
)

type RtcCertificateType int

const (
	RTC_CERTIFICATE_DEFAULT RtcCertificateType = C.RTC_CERTIFICATE_DEFAULT // ECDSA
	RTC_CERTIFICATE_ECDSA                      = C.RTC_CERTIFICATE_ECDSA
	RTC_CERTIFICATE_RSA                        = C.RTC_CERTIFICATE_RSA
)

type RtcCodec int

const (
	// video
	RTC_CODEC_H264 RtcCodec = C.RTC_CODEC_H264
	RTC_CODEC_VP8           = C.RTC_CODEC_VP8
	RTC_CODEC_VP9           = C.RTC_CODEC_VP9
	RTC_CODEC_H265          = C.RTC_CODEC_H265
	RTC_CODEC_AV1           = C.RTC_CODEC_AV1

	// audio
	RTC_CODEC_OPUS = C.RTC_CODEC_OPUS
	RTC_CODEC_PCMU = C.RTC_CODEC_PCMU
	RTC_CODEC_PCMA = C.RTC_CODEC_PCMA
	RTC_CODEC_AAC  = C.RTC_CODEC_AAC
)

type RtcDirection int

const (
	RTC_DIRECTION_SENDONLY RtcDirection = C.RTC_DIRECTION_SENDONLY
	RTC_DIRECTION_UNKNOWN               = C.RTC_DIRECTION_UNKNOWN
	RTC_DIRECTION_RECVONLY              = C.RTC_DIRECTION_RECVONLY
	RTC_DIRECTION_SENDRECV              = C.RTC_DIRECTION_SENDRECV
	RTC_DIRECTION_INACTIVE              = C.RTC_DIRECTION_INACTIVE
)

type RtcTransportPolicy int

const (
	RTC_TRANSPORT_POLICY_ALL   RtcTransportPolicy = C.RTC_TRANSPORT_POLICY_ALL
	RTC_TRANSPORT_POLICY_RELAY                    = C.RTC_TRANSPORT_POLICY_RELAY
)

const (
	RTC_ERR_SUCCESS   = C.RTC_ERR_SUCCESS
	RTC_ERR_INVALID   = C.RTC_ERR_INVALID   // invalid argument
	RTC_ERR_FAILURE   = C.RTC_ERR_FAILURE   // runtime error
	RTC_ERR_NOT_AVAIL = C.RTC_ERR_NOT_AVAIL // element not available
	RTC_ERR_TOO_SMALL = C.RTC_ERR_TOO_SMALL // buffer too small
)

type RTCLogCallbackFunc func(RtcLogLevel, string)

var logCallback RTCLogCallbackFunc

//export go_rtcLogCallbackFunc
func go_rtcLogCallbackFunc(level C.rtcLogLevel, message *C.cchar_t) {
	logCallback(RtcLogLevel(level), C.GoString(message))
}

type RtcDescriptionCallbackFunc func(pc int, sdp string, descriptionType string, ptr unsafe.Pointer)

var descriptionCallbackMap = make(map[int]RtcDescriptionCallbackFunc)
var descriptionCallbackMapLock = &sync.Mutex{}

//export go_rtcDescriptionCallbackFunc
func go_rtcDescriptionCallbackFunc(pc C.int, sdp *C.cchar_t, descriptionType *C.cchar_t, ptr unsafe.Pointer) {
	descriptionCallbackMapLock.Lock()
	cb, ok := descriptionCallbackMap[int(pc)]
	descriptionCallbackMapLock.Unlock()
	if ok {
		cb(int(pc), C.GoString(sdp), C.GoString(descriptionType), ptr)
	}
}

type RtcCandidateCallbackFunc func(pc int, sdp string, descriptionType string, ptr unsafe.Pointer)

var localCandidateCallbackMap = make(map[int]RtcCandidateCallbackFunc)
var localCandidateCallbackMapLock = &sync.Mutex{}

//export go_rtcSetLocalCandidateCallback
func go_rtcSetLocalCandidateCallback(pc C.int, cand *C.cchar_t, mid *C.cchar_t, ptr unsafe.Pointer) {
	localCandidateCallbackMapLock.Lock()
	cb, ok := localCandidateCallbackMap[int(pc)]
	localCandidateCallbackMapLock.Unlock()
	if ok {
		cb(int(pc), C.GoString(cand), C.GoString(mid), ptr)
	}
}

type RtcStateChangeCallbackFunc func(pc int, state int, ptr unsafe.Pointer)

var rtcStateChangeCallbackMap = make(map[int]RtcStateChangeCallbackFunc)
var rtcStateChangeCallbackMapLock = &sync.Mutex{}

//export go_rtcSetStateChangeCallback
func go_rtcSetStateChangeCallback(pc C.int, state C.rtcState, ptr unsafe.Pointer) {
	rtcStateChangeCallbackMapLock.Lock()
	cb, ok := rtcStateChangeCallbackMap[int(pc)]
	rtcStateChangeCallbackMapLock.Unlock()
	if ok {
		cb(int(pc), int(state), ptr)
	}
}

type RtcIceStateChangeCallbackFunc func(pc int, state int, ptr unsafe.Pointer)

var rtcIceStateChangeCallbackFunccMap = make(map[int]RtcIceStateChangeCallbackFunc)
var rtcIceStateChangeCallbackFunccMapLock = &sync.Mutex{}

//export go_rtcIceStateChangeCallbackFunc
func go_rtcIceStateChangeCallbackFunc(pc C.int, dc C.int, ptr unsafe.Pointer) {
	rtcIceStateChangeCallbackFunccMapLock.Lock()
	cb, ok := rtcIceStateChangeCallbackFunccMap[int(pc)]
	rtcIceStateChangeCallbackFunccMapLock.Unlock()
	if ok {
		cb(int(pc), int(dc), ptr)
	}
}

type RtcGatheringStateCallbackFunc func(pc int, state int, ptr unsafe.Pointer)

var rtcGatheringStateCallbackFuncMap = make(map[int]RtcGatheringStateCallbackFunc)
var rtcGatheringStateCallbackFuncMapLock = &sync.Mutex{}

//export go_rtcSetGatheringStateChangeCallback
func go_rtcSetGatheringStateChangeCallback(pc C.int, state C.rtcGatheringState, ptr unsafe.Pointer) {
	rtcGatheringStateCallbackFuncMapLock.Lock()
	cb, ok := rtcGatheringStateCallbackFuncMap[int(pc)]
	rtcGatheringStateCallbackFuncMapLock.Unlock()
	if ok {
		cb(int(pc), int(state), ptr)
	}
}

type RtcSignalingStateCallbackFunc func(pc int, state int, ptr unsafe.Pointer)

var rtcSignalingStateCallbackFuncMap = make(map[int]RtcSignalingStateCallbackFunc)
var rtcSignalingStateCallbackFuncLock = &sync.Mutex{}

//export go_rtcSignalingStateCallbackFunc
func go_rtcSignalingStateCallbackFunc(pc C.int, state C.rtcSignalingState, ptr unsafe.Pointer) {
	rtcSignalingStateCallbackFuncLock.Lock()
	cb, ok := rtcSignalingStateCallbackFuncMap[int(pc)]
	rtcSignalingStateCallbackFuncLock.Unlock()
	if ok {
		cb(int(pc), int(state), ptr)
	}
}

type RtcDataChannelCallbackFunc func(pc int, dc int, ptr unsafe.Pointer)

var rtcDataChannelCallbackFuncMap = make(map[int]RtcDataChannelCallbackFunc)
var rtcDataChannelCallbackFuncMapLock = &sync.Mutex{}

//export go_rtcSetDataChannelCallback
func go_rtcSetDataChannelCallback(pc C.int, dc C.int, ptr unsafe.Pointer) {
	rtcDataChannelCallbackFuncMapLock.Lock()
	cb, ok := rtcDataChannelCallbackFuncMap[int(pc)]
	rtcDataChannelCallbackFuncMapLock.Unlock()
	if ok {
		cb(int(pc), int(dc), ptr)
	}
}

type RtcTrackCallbackFunc func(pc int, tr int, ptr unsafe.Pointer)

var rtcTrackCallbackFuncMap = make(map[int]RtcTrackCallbackFunc)
var rtcTrackCallbackFuncLock = &sync.Mutex{}

//export go_rtcTrackCallbackFunc
func go_rtcTrackCallbackFunc(pc C.int, state C.rtcSignalingState, ptr unsafe.Pointer) {
	rtcTrackCallbackFuncLock.Lock()
	cb, ok := rtcTrackCallbackFuncMap[int(pc)]
	rtcTrackCallbackFuncLock.Unlock()
	if ok {
		cb(int(pc), int(state), ptr)
	}
}

type RtcOpenCallbackFunc func(int, unsafe.Pointer)

var openCallbackMap = make(map[int]RtcOpenCallbackFunc)
var openCallbackMapLock = &sync.Mutex{}

//export go_rtcOpenCallbackFunc
func go_rtcOpenCallbackFunc(id C.int, ptr unsafe.Pointer) {
	openCallbackMapLock.Lock()
	cb, ok := openCallbackMap[int(id)]
	openCallbackMapLock.Unlock()
	if ok {
		cb(int(id), ptr)
	}
}

type RtcClosedCallbackFunc func(int, unsafe.Pointer)

var closedCallbackMap = make(map[int]RtcClosedCallbackFunc)
var closedCallbackMapLock = &sync.Mutex{}

//export go_rtcClosedCallbackFunc
func go_rtcClosedCallbackFunc(id C.int, ptr unsafe.Pointer) {
	closedCallbackMapLock.Lock()
	cb, ok := closedCallbackMap[int(id)]
	closedCallbackMapLock.Unlock()
	if ok {
		cb(int(id), ptr)
	}
}

type RtcErrorCallbackFunc func(int, string, unsafe.Pointer)

var errorCallbackMap = make(map[int]RtcErrorCallbackFunc)
var errorCallbackMapLock = &sync.Mutex{}

//export go_rtcErrorCallbackFunc
func go_rtcErrorCallbackFunc(id C.int, message *C.cchar_t, ptr unsafe.Pointer) {
	errorCallbackMapLock.Lock()
	cb, ok := errorCallbackMap[int(id)]
	errorCallbackMapLock.Unlock()
	if ok {
		cb(int(id), C.GoString(message), ptr)
	}
}

type RtcMessageCallbackFunc func(id int, message []byte, size int, ptr unsafe.Pointer)

var messageCallbackMap = make(map[int]RtcMessageCallbackFunc)
var messageCallbackMapLock = &sync.Mutex{}

func convertToGoBytes(cStr *C.cchar_t, length int) []byte {
	if length <= 0 {
		return []byte(C.GoString(cStr)) // NULL終端文字列の場合
	} else {
		return C.GoBytes(unsafe.Pointer(cStr), C.int(length)) // バイト列の場合
	}
}

//export go_rtcMessageCallbackFunc
func go_rtcMessageCallbackFunc(id C.int, message *C.cchar_t, size C.int, ptr unsafe.Pointer) {
	messageCallbackMapLock.Lock()
	cb, ok := messageCallbackMap[int(id)]
	messageCallbackMapLock.Unlock()
	if ok {
		cb(int(id), convertToGoBytes(message, int(size)), int(size), ptr)
	}
}

type RtcInterceptorCallbackFunc func(pc int, message []byte, size int, ptr unsafe.Pointer)

var rtcInterceptorCallbackFuncMap = make(map[int]RtcInterceptorCallbackFunc)
var rtcInterceptorCallbackFuncLock = &sync.Mutex{}

//export go_rtcInterceptorCallbackFunc
func go_rtcInterceptorCallbackFunc(pc C.int, state C.rtcSignalingState, ptr unsafe.Pointer) {
	rtcTrackCallbackFuncLock.Lock()
	cb, ok := rtcTrackCallbackFuncMap[int(pc)]
	rtcTrackCallbackFuncLock.Unlock()
	if ok {
		cb(int(pc), int(state), ptr)
	}
}

type RtcBufferedAmountLowCallbackFunc func(id int, ptr unsafe.Pointer)

var rtcBufferedAmountLowCallbackFuncMap = make(map[int]RtcBufferedAmountLowCallbackFunc)
var rtcBufferedAmountLowCallbackFuncLock = &sync.Mutex{}

//export go_rtcSetBufferedAmountLowCallback
func go_rtcSetBufferedAmountLowCallback(pc C.int, ptr unsafe.Pointer) {
	rtcBufferedAmountLowCallbackFuncLock.Lock()
	cb, ok := rtcBufferedAmountLowCallbackFuncMap[int(pc)]
	rtcBufferedAmountLowCallbackFuncLock.Unlock()
	if ok {
		cb(int(pc), ptr)
	}
}

type RtcAvailableCallbackFunc func(id int, ptr unsafe.Pointer)

var rtcAvailableCallbackFuncMap = make(map[int]RtcBufferedAmountLowCallbackFunc)
var rtcAvailableCallbackFuncLock = &sync.Mutex{}

//export go_rtcSetAvailableCallback
func go_rtcSetAvailableCallback(pc C.int, ptr unsafe.Pointer) {
	rtcAvailableCallbackFuncLock.Lock()
	cb, ok := rtcAvailableCallbackFuncMap[int(pc)]
	rtcAvailableCallbackFuncLock.Unlock()
	if ok {
		cb(int(pc), ptr)
	}
}

// Log

// NULL cb on the first call will log to stdout
func RtcInitLogger(level RtcLogLevel, cb RTCLogCallbackFunc) {
	logCallback = cb

	if cb == nil {
		C.rtcInitLogger(C.rtcLogLevel(int(level)), C.rtcLogCallbackFunc(unsafe.Pointer(nil)))
	} else {
		C.rtcInitLogger(C.rtcLogLevel(int(level)), C.rtcLogCallbackFunc(C.go_rtcLogCallbackFunc))
	}
}

// User pointer
// Currenty, ptr expects primitive pointer type such as int, int32, int64, etc.
// Due to cgo pointer to pointer problem. Would like to save *Peer. but don't know generic solution.
func RtcSetUserPointer(id int, ptr unsafe.Pointer) {
	C.rtcSetUserPointer(C.int(id), ptr)
}

func RtcGetUserPointer(id int) unsafe.Pointer {
	return unsafe.Pointer(C.rtcGetUserPointer(C.int(id)))
}

// PeerConnection
type RtcConfiguration struct {
	config C.rtcConfiguration
}

func RtcCreatePeerConnection(config *RtcConfiguration) int {
	return int(C.rtcCreatePeerConnection(&config.config))
}

func RtcClosePeerConnection(id int) {
	C.rtcClosePeerConnection(C.int(id))
}

func RtcDeletePeerConnection(id int) {
	C.rtcDeletePeerConnection(C.int(id))
}

func RtcSetLocalDescription(pc int, descriptionType string) int {
	return int(C.rtcSetLocalDescription(C.int(pc), C.CString(descriptionType)))
}

func RtcSetLocalDescriptionCallback(id int, cb RtcDescriptionCallbackFunc) int {
	descriptionCallbackMapLock.Lock()
	descriptionCallbackMap[id] = cb
	descriptionCallbackMapLock.Unlock()

	return int(C.rtcSetLocalDescriptionCallback(C.int(id), C.rtcDescriptionCallbackFunc(C.go_rtcDescriptionCallbackFunc)))
}

func RtcSetLocalCandidateCallback(id int, cb RtcCandidateCallbackFunc) int {
	localCandidateCallbackMapLock.Lock()
	localCandidateCallbackMap[id] = cb
	localCandidateCallbackMapLock.Unlock()

	return int(C.rtcSetLocalCandidateCallback(C.int(id), C.rtcCandidateCallbackFunc(C.go_rtcSetLocalCandidateCallback)))
}

func RtcSetStateChangeCallback(id int, cb RtcStateChangeCallbackFunc) int {
	rtcStateChangeCallbackMapLock.Lock()
	rtcStateChangeCallbackMap[id] = cb
	rtcStateChangeCallbackMapLock.Unlock()

	return int(C.rtcSetStateChangeCallback(C.int(id), C.rtcStateChangeCallbackFunc(C.go_rtcSetStateChangeCallback)))
}

func RtcSetIceStateChangeCallback(id int, cb RtcIceStateChangeCallbackFunc) int {
	rtcIceStateChangeCallbackFunccMapLock.Lock()
	rtcIceStateChangeCallbackFunccMap[id] = cb
	rtcIceStateChangeCallbackFunccMapLock.Unlock()

	return int(C.rtcSetIceStateChangeCallback(C.int(id), C.rtcIceStateChangeCallbackFunc(C.go_rtcIceStateChangeCallbackFunc)))
}

func RtcSetGatheringStateChangeCallback(id int, cb RtcGatheringStateCallbackFunc) int {
	rtcGatheringStateCallbackFuncMapLock.Lock()
	rtcGatheringStateCallbackFuncMap[id] = cb
	rtcGatheringStateCallbackFuncMapLock.Unlock()

	return int(C.rtcSetGatheringStateChangeCallback(C.int(id), C.rtcGatheringStateCallbackFunc(C.go_rtcSetGatheringStateChangeCallback)))
}

func RtcSetSignalingStateChangeCallback(id int, cb RtcSignalingStateCallbackFunc) int {
	rtcSignalingStateCallbackFuncLock.Lock()
	rtcSignalingStateCallbackFuncMap[id] = cb
	rtcSignalingStateCallbackFuncLock.Unlock()

	return int(C.rtcSetSignalingStateChangeCallback(C.int(id), C.rtcSignalingStateCallbackFunc(C.go_rtcSignalingStateCallbackFunc)))
}

func RtcSetRemoteDescription(id int, sdb string, descriptionType string) int {
	return int(C.rtcSetRemoteDescription(C.int(id), C.CString(sdb), C.CString(descriptionType)))
}

func RtcAddRemoteCandidate(id int, sdp string, mid string) int {
	return int(C.rtcAddRemoteCandidate(C.int(id), C.CString(sdp), C.CString(mid)))
}

func RtcGetLocalDescription(pc int, buffer []byte, size int) int {
	return int(C.rtcGetLocalDescription(C.int(pc), (*C.char)(unsafe.Pointer(&buffer[0])), C.int(size)))
}

func RtcGetRemoteDescription(pc int, buffer []byte, size int) int {
	return int(C.rtcGetRemoteDescription(C.int(pc), (*C.char)(unsafe.Pointer(&buffer[0])), C.int(size)))
}

func RtcGetLocalDescriptionType(pc int, buffer []byte, size int) int {
	return int(C.rtcGetLocalDescriptionType(C.int(pc), (*C.char)(unsafe.Pointer(&buffer[0])), C.int(size)))
}

func RtcGetRemoteDescriptionType(pc int, buffer []byte, size int) int {
	return int(C.rtcGetRemoteDescriptionType(C.int(pc), (*C.char)(unsafe.Pointer(&buffer[0])), C.int(size)))
}

func RtcGetLocalAddress(pc int, buffer []byte, size int) int {
	return int(C.rtcGetLocalAddress(C.int(pc), (*C.char)(unsafe.Pointer(&buffer[0])), C.int(size)))
}

func RtcGetRemoteAddress(pc int, buffer []byte, size int) int {
	return int(C.rtcGetRemoteAddress(C.int(pc), (*C.char)(unsafe.Pointer(&buffer[0])), C.int(size)))
}

func RtcGetSelectedCandidatePair(pc int, local []byte, localSize int, remote []byte, remoteSize int) int {
	return int(C.rtcGetSelectedCandidatePair(C.int(pc), (*C.char)(unsafe.Pointer(&local[0])), C.int(localSize), (*C.char)(unsafe.Pointer(&remote[0])), C.int(remoteSize)))
}

func RtcGetMaxDataChannelStream(pc int) int {
	return int(C.rtcGetMaxDataChannelStream(C.int(pc)))
}

func RtcGetRemoteMaxMessageSize(pc int) int {
	return int(C.rtcGetRemoteMaxMessageSize(C.int(pc)))
}

// DataChannel, Track, and WebSocket common API

func RtcSetOpenCallback(id int, cb RtcOpenCallbackFunc) int {
	openCallbackMapLock.Lock()
	openCallbackMap[id] = cb
	openCallbackMapLock.Unlock()

	return int(C.rtcSetOpenCallback(C.int(id), C.rtcOpenCallbackFunc(C.go_rtcOpenCallbackFunc)))
}

func RtcSetClosedCallback(id int, cb RtcClosedCallbackFunc) {
	closedCallbackMapLock.Lock()
	closedCallbackMap[id] = cb
	closedCallbackMapLock.Unlock()

	C.rtcSetClosedCallback(C.int(id), C.rtcClosedCallbackFunc(C.go_rtcClosedCallbackFunc))
}

func RtcSetErrorCallback(id int, cb RtcErrorCallbackFunc) {
	errorCallbackMapLock.Lock()
	errorCallbackMap[id] = cb
	errorCallbackMapLock.Unlock()

	C.rtcSetErrorCallback(C.int(id), C.rtcErrorCallbackFunc(C.go_rtcErrorCallbackFunc))
}

func RtcSetMessageCallback(id int, cb RtcMessageCallbackFunc) {
	messageCallbackMapLock.Lock()
	messageCallbackMap[id] = cb
	messageCallbackMapLock.Unlock()

	C.rtcSetMessageCallback(C.int(id), C.rtcMessageCallbackFunc(C.go_rtcMessageCallbackFunc))
}

func RtcSendMessage(id int, message []byte, length int) {
	C.rtcSendMessage(C.int(id), (*C.char)(unsafe.Pointer(&message[0])), C.int(length))
}

func RtcClose(id int) int {
	return int(C.rtcClose(C.int(id)))
}

func RtcDelete(id int) int {
	return int(C.rtcDelete(C.int(id)))
}

func RtcIsOpen(id int) bool {
	return bool(C.rtcIsOpen(C.int(id)))
}

func RtcIsClosed(id int) bool {
	return bool(C.rtcIsClosed(C.int(id)))
}

func RtcMaxMessageSize(id int) int {
	return int(C.rtcMaxMessageSize(C.int(id)))
}

func RtcGetBufferedAmount(id int) int {
	return int(C.rtcGetBufferedAmount(C.int(id)))
}

func RtcSetBufferedAmountLowThreshold(id int, amount int) int {
	return int(C.rtcSetBufferedAmountLowThreshold(C.int(id), C.int(amount)))
}

func RtcSetBufferedAmountLowCallback(id int, cb RtcBufferedAmountLowCallbackFunc) {
	rtcBufferedAmountLowCallbackFuncLock.Lock()
	rtcBufferedAmountLowCallbackFuncMap[id] = cb
	rtcBufferedAmountLowCallbackFuncLock.Unlock()

	C.rtcSetBufferedAmountLowCallback(C.int(id), C.rtcBufferedAmountLowCallbackFunc(C.go_rtcSetBufferedAmountLowCallback))
}

// DataChannel, Track, and WebSocket common extended API

func RtcGetAvailableAmount(id int) int {
	return int(C.rtcGetAvailableAmount(C.int(id)))
}

func RtcSetAvailableCallback(id int, cb RtcBufferedAmountLowCallbackFunc) {
	rtcAvailableCallbackFuncLock.Lock()
	rtcAvailableCallbackFuncMap[id] = cb
	rtcAvailableCallbackFuncLock.Unlock()

	C.rtcSetAvailableCallback(C.int(id), C.rtcAvailableCallbackFunc(C.go_rtcSetAvailableCallback))
}

func RtcReceiveMessage(id int, buffer []byte, size *int) int {
	return int(C.rtcReceiveMessage(C.int(id), (*C.char)(unsafe.Pointer(&buffer[0])), (*C.int)(unsafe.Pointer(size))))
}

// DataChannel
type RtcReliability struct {
	Unordered         bool
	Unreliable        bool
	MaxPacketLifeTime int // ignored if reliable
	MaxRetransmits    int // ignored if reliable
}

type RtcDataChannelInit struct {
	Reliability  RtcReliability
	Protocol     string
	Negotiated   bool
	ManualStream bool
	Stream       uint16
}

func convertRtcDataChannelInit(goInit *RtcDataChannelInit) *C.rtcDataChannelInit {
	reliability := convertRtcReliability(goInit.Reliability)
	cInit := &C.rtcDataChannelInit{
		reliability:  *reliability,
		negotiated:   C.bool(goInit.Negotiated),
		manualStream: C.bool(goInit.ManualStream),
		stream:       C.uint16_t(goInit.Stream),
	}

	if goInit.Protocol != "" {
		cInit.protocol = C.CString(goInit.Protocol)
	} else {
		cInit.protocol = C.CString("")
	}

	return cInit
}

func freeRtcDataChannelInit(cInit *C.rtcDataChannelInit) {
	if cInit != nil {
		C.free(unsafe.Pointer(cInit.protocol))
	}
}

func convertRtcReliability(goReliability RtcReliability) *C.rtcReliability {
	return &C.rtcReliability{
		unordered:         C.bool(goReliability.Unordered),
		unreliable:        C.bool(goReliability.Unreliable),
		maxPacketLifeTime: C.int(goReliability.MaxPacketLifeTime),
		maxRetransmits:    C.int(goReliability.MaxRetransmits),
	}
}

func RtcSetDataChannelCallback(id int, cb RtcDataChannelCallbackFunc) int {
	rtcDataChannelCallbackFuncMapLock.Lock()
	rtcDataChannelCallbackFuncMap[id] = cb
	rtcDataChannelCallbackFuncMapLock.Unlock()

	return int(C.rtcSetDataChannelCallback(C.int(id), C.rtcDataChannelCallbackFunc(C.go_rtcSetDataChannelCallback)))
}

func RtcCreateDataChannel(id int, label string) int {
	return int(C.rtcCreateDataChannel(C.int(id), C.CString(label)))
}

func RtcCreateDataChannelEx(pc int, label string, init *RtcDataChannelInit) int {
	cInit := convertRtcDataChannelInit(init)
	ret := C.rtcCreateDataChannelEx(C.int(pc), C.CString(label), cInit)
	freeRtcDataChannelInit(cInit)

	return int(ret)
}

func RtcDeleteDataChannel(dc int) int {
	return int(C.rtcDeleteDataChannel(C.int(dc)))
}

func RtcGetDataChannelStream(dc int) int {
	return int(C.rtcGetDataChannelStream(C.int(dc)))
}

func RtcGetDataChannelLabel(dc int, buffer []byte, size int) int {
	return int(C.rtcGetDataChannelLabel(C.int(dc), (*C.char)(unsafe.Pointer(&buffer[0])), C.int(size)))
}

func RtcGetDataChannelProtocol(dc int, buffer []byte, size int) int {
	return int(C.rtcGetDataChannelProtocol(C.int(dc), (*C.char)(unsafe.Pointer(&buffer[0])), C.int(size)))
}

func RtcGetDataChannelReliability(dc int, reliability *RtcReliability) int {
	creliability := convertRtcReliability(*reliability)
	return int(C.rtcGetDataChannelReliability(C.int(dc), creliability))
}

// Track

type RtcTrackInit struct {
	Direction   RtcDirection
	Codec       RtcCodec
	PayloadType int
	SSRC        uint32
	MID         string
	Name        string
	MSID        string
	TrackID     string
	Profile     string
}

func convertRtcTrackInit(goInit RtcTrackInit) *C.rtcTrackInit {
	cInit := &C.rtcTrackInit{
		direction:   C.rtcDirection(goInit.Direction),
		codec:       C.rtcCodec(goInit.Codec),
		payloadType: C.int(goInit.PayloadType),
		ssrc:        C.uint32_t(goInit.SSRC),
	}

	// Optional string fields
	if goInit.MID != "" {
		cInit.mid = C.CString(goInit.MID)
	}
	if goInit.Name != "" {
		cInit.name = C.CString(goInit.Name)
	}
	if goInit.MSID != "" {
		cInit.msid = C.CString(goInit.MSID)
	}
	if goInit.TrackID != "" {
		cInit.trackId = C.CString(goInit.TrackID)
	}
	if goInit.Profile != "" {
		cInit.profile = C.CString(goInit.Profile)
	}

	return cInit
}

func freeCTrackInit(cInit *C.rtcTrackInit) {
	if cInit.mid != nil {
		C.free(unsafe.Pointer(cInit.mid))
	}
	if cInit.name != nil {
		C.free(unsafe.Pointer(cInit.name))
	}
	if cInit.msid != nil {
		C.free(unsafe.Pointer(cInit.msid))
	}
	if cInit.trackId != nil {
		C.free(unsafe.Pointer(cInit.trackId))
	}
	if cInit.profile != nil {
		C.free(unsafe.Pointer(cInit.profile))
	}

}

func RtcSetTrackCallback(id int, cb RtcTrackCallbackFunc) {
	rtcTrackCallbackFuncLock.Lock()
	rtcTrackCallbackFuncMap[id] = cb
	rtcTrackCallbackFuncLock.Unlock()

	C.rtcSetTrackCallback(C.int(id), C.rtcTrackCallbackFunc(C.go_rtcTrackCallbackFunc))
}

func RtcAddTrack(pc int, mediaDescriptionSdp string) int {
	return int(C.rtcAddTrack(C.int(pc), C.CString(mediaDescriptionSdp)))
}

func RtcAddTrackEx(pc int, init RtcTrackInit) int {
	cconfig := convertRtcTrackInit(init)

	return int(C.rtcAddTrackEx(C.int(pc), cconfig))
}

func RtcDeleteTrack(tr int) int {
	return int(C.rtcDeleteTrack(C.int(tr)))
}

func RtcGetTrackDescription(tr int, buffer []byte, size int) int {
	return int(C.rtcGetTrackDescription(C.int(tr), (*C.char)(unsafe.Pointer(&buffer[0])), C.int(size)))
}

func RtcGetTrackMid(tr int, buffer []byte, size int) int {
	return int(C.rtcGetTrackMid(C.int(tr), (*C.char)(unsafe.Pointer(&buffer[0])), C.int(size)))
}

func RtcGetTrackDirection(tr int, direction *RtcDirection) int {
	return int(C.rtcGetTrackDirection(C.int(tr), (*C.rtcDirection)(unsafe.Pointer(direction))))
}

// #if RTC_ENABLE_MEDIA

// Media

// Define how OBUs are packetizied in a AV1 Sample
type RtcObuPacketization int

const (
	RTC_OBU_PACKETIZED_OBU           RtcObuPacketization = C.RTC_OBU_PACKETIZED_OBU
	RTC_OBU_PACKETIZED_TEMPORAL_UNIT                     = C.RTC_OBU_PACKETIZED_TEMPORAL_UNIT
)

// Define how NAL units are separated in a H264/H265 sample
type RtcNalUnitSeparator int

const (
	RTC_NAL_SEPARATOR_LENGTH               RtcNalUnitSeparator = C.RTC_NAL_SEPARATOR_LENGTH               // first 4 bytes are NAL unit length
	RTC_NAL_SEPARATOR_LONG_START_SEQUENCE                      = C.RTC_NAL_SEPARATOR_LONG_START_SEQUENCE  // 0x00, 0x00, 0x00, 0x01
	RTC_NAL_SEPARATOR_SHORT_START_SEQUENCE                     = C.RTC_NAL_SEPARATOR_SHORT_START_SEQUENCE // 0x00, 0x00, 0x01
	RTC_NAL_SEPARATOR_START_SEQUENCE                           = C.RTC_NAL_SEPARATOR_START_SEQUENCE       // long or short start sequence
)

type RtcPacketizationHandlerInit struct {
	Ssrc            uint32
	Cname           string
	PayloadType     uint8
	ClockRate       uint32
	SequenceNumber  uint16
	Timestamp       uint32
	NalSeparator    RtcNalUnitSeparator
	MaxFragmentSize uint16
}

type RtcSsrcForTypeInit struct {
	Ssrc    uint32
	Name    string
	Msid    string
	TrackId string
}

// Opaque message

// Opaque type used (via rtcMessage*) to reference an rtc::Message
type RtcMessage unsafe.Pointer

// // Allocate a new opaque message.
// // Must be explicitly freed by rtcDeleteOpaqueMessage() unless
// // explicitly returned by a media interceptor callback;

func RtcCreateOpaqueMessage(data unsafe.Pointer, size int) RtcMessage {
	return RtcMessage(C.rtcCreateOpaqueMessage(data, C.int(size)))
}

// RTC_C_EXPORT void rtcDeleteOpaqueMessage(rtcMessage *msg);
func RtcDeleteOpaqueMessage(msg RtcMessage) {
	C.rtcDeleteOpaqueMessage((*C.rtcMessage)(msg))
}

// // Set MediaInterceptor for peer connection
func RtcSetMediaInterceptorCallback(id int, cb RtcTrackCallbackFunc) int {
	rtcTrackCallbackFuncLock.Lock()
	rtcTrackCallbackFuncMap[id] = cb
	rtcTrackCallbackFuncLock.Unlock()

	return int(C.rtcSetMediaInterceptorCallback(C.int(id), C.rtcInterceptorCallbackFunc(C.go_rtcInterceptorCallbackFunc)))
}

func convertRtcPacketizationHandlerInit(init RtcPacketizationHandlerInit) C.rtcPacketizationHandlerInit {
	return C.rtcPacketizationHandlerInit{
		ssrc:            C.uint32_t(init.Ssrc),
		cname:           C.CString(init.Cname),
		payloadType:     C.uint8_t(init.PayloadType),
		clockRate:       C.uint32_t(init.ClockRate),
		sequenceNumber:  C.uint16_t(init.SequenceNumber),
		timestamp:       C.uint32_t(init.Timestamp),
		nalSeparator:    C.rtcNalUnitSeparator(init.NalSeparator),
		maxFragmentSize: C.uint16_t(init.MaxFragmentSize),
	}
}

func convertCSsrcForTypeInit(goInit RtcSsrcForTypeInit) C.rtcSsrcForTypeInit {
	cInit := C.rtcSsrcForTypeInit{
		ssrc: C.uint32_t(goInit.Ssrc),
	}

	if goInit.Name != "" {
		cInit.name = C.CString(goInit.Name)
	}
	if goInit.Msid != "" {
		cInit.msid = C.CString(goInit.Msid)
	}
	if goInit.TrackId != "" {
		cInit.trackId = C.CString(goInit.TrackId)
	}

	return cInit
}

func freeCSsrcForTypeInit(cInit C.rtcSsrcForTypeInit) {
	if cInit.name != nil {
		C.free(unsafe.Pointer(cInit.name))
	}
	if cInit.msid != nil {
		C.free(unsafe.Pointer(cInit.msid))
	}
	if cInit.trackId != nil {
		C.free(unsafe.Pointer(cInit.trackId))
	}
}

// // Set H264PacketizationHandler for track
// RTC_C_EXPORT int rtcSetH264PacketizationHandler(int tr, const rtcPacketizationHandlerInit *init);
func RtcSetH264PacketizationHandler(tr int, init RtcPacketizationHandlerInit) int {
	cInit := convertRtcPacketizationHandlerInit(init)
	defer C.free(unsafe.Pointer(cInit.cname))
	return int(C.rtcSetH264PacketizationHandler(C.int(tr), &cInit))
}

// // Set H265PacketizationHandler for track
// RTC_C_EXPORT int rtcSetH265PacketizationHandler(int tr, const rtcPacketizationHandlerInit *init);
func RtcSetH265PacketizationHandler(tr int, init RtcPacketizationHandlerInit) int {
	cInit := convertRtcPacketizationHandlerInit(init)
	defer C.free(unsafe.Pointer(cInit.cname))
	return int(C.rtcSetH265PacketizationHandler(C.int(tr), &cInit))
}

// // Set AV1PacketizationHandler for track
// RTC_C_EXPORT int rtcSetAV1PacketizationHandler(int tr, const rtcPacketizationHandlerInit *init);
func RtcSetAV1PacketizationHandler(tr int, init RtcPacketizationHandlerInit) int {
	cInit := convertRtcPacketizationHandlerInit(init)
	defer C.free(unsafe.Pointer(cInit.cname))
	return int(C.rtcSetAV1PacketizationHandler(C.int(tr), &cInit))
}

// // Set OpusPacketizationHandler for track
// RTC_C_EXPORT int rtcSetOpusPacketizationHandler(int tr, const rtcPacketizationHandlerInit *init);
func RtcSetOpusPacketizationHandler(tr int, init RtcPacketizationHandlerInit) int {
	cInit := convertRtcPacketizationHandlerInit(init)
	defer C.free(unsafe.Pointer(cInit.cname))
	return int(C.rtcSetOpusPacketizationHandler(C.int(tr), &cInit))
}

// // Set AACPacketizationHandler for track
// RTC_C_EXPORT int rtcSetAACPacketizationHandler(int tr, const rtcPacketizationHandlerInit *init);
func RtcSetAACPacketizationHandler(tr int, init RtcPacketizationHandlerInit) int {
	cInit := convertRtcPacketizationHandlerInit(init)
	defer C.free(unsafe.Pointer(cInit.cname))
	return int(C.rtcSetAACPacketizationHandler(C.int(tr), &cInit))
}

// // Chain RtcpSrReporter to handler chain for given track
// RTC_C_EXPORT int rtcChainRtcpSrReporter(int tr);
func RtcChainRtcpSrReporter(tr int) int {
	return int(C.rtcChainRtcpSrReporter(C.int(tr)))
}

// // Chain RtcpNackResponder to handler chain for given track
// RTC_C_EXPORT int rtcChainRtcpNackResponder(int tr, unsigned int maxStoredPacketsCount);
func RtcChainRtcpNackResponder(tr int, maxStoredPacketsCount uint) int {
	return int(C.rtcChainRtcpNackResponder(C.int(tr), C.uint(maxStoredPacketsCount)))
}

// // Transform seconds to timestamp using track's clock rate, result is written to timestamp
// RTC_C_EXPORT int rtcTransformSecondsToTimestamp(int id, double seconds, uint32_t *timestamp);
func RtcTransformSecondsToTimestamp(id int, seconds float64) int {
	var timestamp C.uint32_t
	return int(C.rtcTransformSecondsToTimestamp(C.int(id), C.double(seconds), &timestamp))
}

// // Transform timestamp to seconds using track's clock rate, result is written to seconds
// RTC_C_EXPORT int rtcTransformTimestampToSeconds(int id, uint32_t timestamp, double *seconds);
func RtcTransformTimestampToSeconds(id int, timestamp uint32) int {
	var seconds C.double
	return int(C.rtcTransformTimestampToSeconds(C.int(id), C.uint32_t(timestamp), &seconds))
}

// // Get current timestamp, result is written to timestamp
// RTC_C_EXPORT int rtcGetCurrentTrackTimestamp(int id, uint32_t *timestamp);
func RtcGetCurrentTrackTimestamp(id int) int {
	var timestamp C.uint32_t
	return int(C.rtcGetCurrentTrackTimestamp(C.int(id), &timestamp))
}

// // Set RTP timestamp for track identified by given id
// RTC_C_EXPORT int rtcSetTrackRtpTimestamp(int id, uint32_t timestamp);
func RtcSetTrackRtpTimestamp(id int, timestamp uint32) int {
	return int(C.rtcSetTrackRtpTimestamp(C.int(id), C.uint32_t(timestamp)))
}

// // Get timestamp of last RTCP SR, result is written to timestamp
// RTC_C_EXPORT int rtcGetLastTrackSenderReportTimestamp(int id, uint32_t *timestamp);
func RtcGetLastTrackSenderReportTimestamp(id int) int {
	var timestamp C.uint32_t
	return int(C.rtcGetLastTrackSenderReportTimestamp(C.int(id), &timestamp))
}

// // Set NeedsToReport flag in RtcpSrReporter handler identified by given track id
// RTC_C_EXPORT int rtcSetNeedsToSendRtcpSr(int id);
func RtcSetNeedsToSendRtcpSr(id int) int {
	return int(C.rtcSetNeedsToSendRtcpSr(C.int(id)))
}

// // Get all available payload types for given codec and stores them in buffer, does nothing if
// // buffer is NULL
// int rtcGetTrackPayloadTypesForCodec(int tr, const char *ccodec, int *buffer, int size);
func RtcGetTrackPayloadTypesForCodec(tr int, codec string, buffer []int, size int) int {
	return int(C.rtcGetTrackPayloadTypesForCodec(C.int(tr), C.CString(codec), (*C.int)(unsafe.Pointer(&buffer[0])), C.int(size)))
}

// // Get all SSRCs for given track
// int rtcGetSsrcsForTrack(int tr, uint32_t *buffer, int count);
func RtcGetSsrcsForTrack(tr int, buffer []uint32, count int) int {
	return int(C.rtcGetSsrcsForTrack(C.int(tr), (*C.uint32_t)(unsafe.Pointer(&buffer[0])), C.int(count)))
}

// // Get CName for SSRC
// int rtcGetCNameForSsrc(int tr, uint32_t ssrc, char *cname, int cnameSize);
func RtcGetCNameForSsrc(tr int, ssrc uint32, cname []byte, cnameSize int) int {
	return int(C.rtcGetCNameForSsrc(C.int(tr), C.uint32_t(ssrc), (*C.char)(unsafe.Pointer(&cname[0])), C.int(cnameSize)))
}

// // Get all SSRCs for given media type in given SDP
// int rtcGetSsrcsForType(const char *mediaType, const char *sdp, uint32_t *buffer, int bufferSize);
func RtcGetSsrcsForType(mediaType string, sdp string, buffer []uint32, bufferSize int) int {
	return int(C.rtcGetSsrcsForType(C.CString(mediaType), C.CString(sdp), (*C.uint32_t)(unsafe.Pointer(&buffer[0])), C.int(bufferSize)))
}

// // Set SSRC for given media type in given SDP
// int rtcSetSsrcForType(const char *mediaType, const char *sdp, char *buffer, const int bufferSize,
//
//	rtcSsrcForTypeInit *init);
func RtcSetSsrcForType(mediaType string, sdp string, buffer []byte, bufferSize int, init RtcSsrcForTypeInit) int {
	cInit := convertCSsrcForTypeInit(init)

	//あやしい
	defer freeCSsrcForTypeInit(cInit)
	return int(C.rtcSetSsrcForType(C.CString(mediaType), C.CString(sdp), (*C.char)(unsafe.Pointer(&buffer[0])), C.int(bufferSize), &cInit))
}

// #endif

// #if RTC_ENABLE_WEBSOCKET

type RtcWsConfiguration struct {
	DisableTlsVerification bool
	ProxyServer            string
	Protocols              []string
	ConnectionTimeoutMs    int
	PingIntervalMs         int
	MaxOutstandingPings    int
}

func RtcCreateWebSocket(url string) int {
	curl := C.CString(url)
	return int(C.rtcCreateWebSocket(curl))
}

func toCWsConfiguration(goConfig RtcWsConfiguration) *C.rtcWsConfiguration {
	cConfig := &C.rtcWsConfiguration{
		disableTlsVerification: C.bool(goConfig.DisableTlsVerification),
		connectionTimeoutMs:    C.int(goConfig.ConnectionTimeoutMs),
		pingIntervalMs:         C.int(goConfig.PingIntervalMs),
		maxOutstandingPings:    C.int(goConfig.MaxOutstandingPings),
	}

	if goConfig.ProxyServer != "" {
		cConfig.proxyServer = C.CString(goConfig.ProxyServer)
	}

	protocolsCount := len(goConfig.Protocols)
	cConfig.protocolsCount = C.int(protocolsCount)
	if protocolsCount > 0 {
		cProtocols := make([]*C.char, protocolsCount)
		for i, proto := range goConfig.Protocols {
			cProtocols[i] = C.CString(proto)
		}
		cConfig.protocols = (**C.char)(unsafe.Pointer(&cProtocols[0]))
	}

	return cConfig
}

func freeCWsConfiguration(cConfig C.rtcWsConfiguration) {
	if cConfig.proxyServer != nil {
		C.free(unsafe.Pointer(cConfig.proxyServer))
	}
	if cConfig.protocolsCount > 0 {
		protocolsSlice := (*[1 << 30]*C.char)(unsafe.Pointer(cConfig.protocols))[:cConfig.protocolsCount:cConfig.protocolsCount]
		for _, proto := range protocolsSlice {
			C.free(unsafe.Pointer(proto))
		}
	}
}

func RtcCreateWebSocketEx(url string, config *RtcWsConfiguration) int {
	cconfig := toCWsConfiguration(*config)
	ret := C.rtcCreateWebSocketEx(C.CString(url), cconfig)

	return int(ret)
}

func RtcDeleteWebSocket(ws int) int {
	return int(C.rtcDeleteWebSocket(C.int(ws)))
}

func RtcGetWebSocketRemoteAddress(ws int, buffer []byte, size int) int {
	return int(C.rtcGetWebSocketRemoteAddress(C.int(ws), (*C.char)(unsafe.Pointer(&buffer[0])), C.int(size)))
}

func RtcGetWebSocketPath(ws int, buffer []byte, size int) int {
	return int(C.rtcGetWebSocketPath(C.int(ws), (*C.char)(unsafe.Pointer(&buffer[0])), C.int(size)))
}

// WebSocketServer

type RtcWsServerConfiguration struct {
	Port                uint16 // 0 means automatic selection
	EnableTls           bool   // if true, enable TLS (WSS)
	CertificatePemFile  string // NULL for autogenerated certificate
	KeyPemFile          string // NULL for autogenerated certificate
	KeyPemPass          string // NULL if no pass
	BindAddress         string // NULL for IP_ANY_ADDR
	ConnectionTimeoutMs int    // in milliseconds, 0 means default, < 0 means disabled
}

func convertWsServerConfiguration(goConfig RtcWsServerConfiguration) *C.rtcWsServerConfiguration {
	cConfig := &C.rtcWsServerConfiguration{
		port:                C.uint16_t(goConfig.Port),
		enableTls:           C.bool(goConfig.EnableTls),
		connectionTimeoutMs: C.int(goConfig.ConnectionTimeoutMs),
	}

	if goConfig.CertificatePemFile != "" {
		cConfig.certificatePemFile = C.CString(goConfig.CertificatePemFile)
	}
	if goConfig.KeyPemFile != "" {
		cConfig.keyPemFile = C.CString(goConfig.KeyPemFile)
	}
	if goConfig.KeyPemPass != "" {
		cConfig.keyPemPass = C.CString(goConfig.KeyPemPass)
	}
	if goConfig.BindAddress != "" {
		cConfig.bindAddress = C.CString(goConfig.BindAddress)
	}

	return cConfig
}

func freeWsServerConfiguration(cConfig *C.rtcWsServerConfiguration) {
	if cConfig.certificatePemFile != nil {
		C.free(unsafe.Pointer(cConfig.certificatePemFile))
	}
	if cConfig.keyPemFile != nil {
		C.free(unsafe.Pointer(cConfig.keyPemFile))
	}
	if cConfig.keyPemPass != nil {
		C.free(unsafe.Pointer(cConfig.keyPemPass))
	}
	if cConfig.bindAddress != nil {
		C.free(unsafe.Pointer(cConfig.bindAddress))
	}
}

type RtcWebSocketClientCallbackFunc func(int, int, unsafe.Pointer)

var RtcCreateWebSocketServerMap = make(map[int]RtcWebSocketClientCallbackFunc)
var RtcCreateWebSocketServerLock = &sync.Mutex{}

//export go_rtcCreateWebSocketServer
func go_rtcCreateWebSocketServer(wsserver C.int, ws C.int, ptr unsafe.Pointer) {
	RtcCreateWebSocketServerLock.Lock()
	cb, ok := RtcCreateWebSocketServerMap[int(wsserver)]
	RtcCreateWebSocketServerLock.Unlock()
	if ok {
		cb(int(wsserver), int(ws), ptr)
	}
}

func RtcCreateWebSocketServer(config RtcWsServerConfiguration, cb RtcWebSocketClientCallbackFunc) int {
	cconfig := convertWsServerConfiguration(config)
	wsserver := C.rtcCreateWebSocketServer(cconfig, C.rtcWebSocketClientCallbackFunc(C.go_rtcCreateWebSocketServer))

	RtcCreateWebSocketServerLock.Lock()
	RtcCreateWebSocketServerMap[int(wsserver)] = cb
	RtcCreateWebSocketServerLock.Unlock()

	return int(wsserver)
}

// RTC_C_EXPORT int rtcDeleteWebSocketServer(int wsserver);
func RtcDeleteWebSocketServer(wsserver int) int {
	return int(C.rtcDeleteWebSocketServer(C.int(wsserver)))
}

// RTC_C_EXPORT int rtcGetWebSocketServerPort(int wsserver);
func RtcGetWebSocketServerPort(wsserver int) int {
	return int(C.rtcGetWebSocketServerPort(C.int(wsserver)))
}

/// #endif

func RtcPreload() {
	C.rtcPreload()
}

func RtcCleanup() {
	C.rtcCleanup()
}

// SCTP global settings

type RtcSctpSettings struct {
	recvBufferSize             int // in bytes, <= 0 means optimized default
	sendBufferSize             int // in bytes, <= 0 means optimized default
	maxChunksOnQueue           int // in chunks, <= 0 means optimized default
	initialCongestionWindow    int // in MTUs, <= 0 means optimized default
	maxBurst                   int // in MTUs, 0 means optimized default, < 0 means disabled
	congestionControlModule    int // 0: RFC2581 (default), 1: HSTCP, 2: H-TCP, 3: RTCC
	delayedSackTimeMs          int // in milliseconds, 0 means optimized default, < 0 means disabled
	minRetransmitTimeoutMs     int // in milliseconds, <= 0 means optimized default
	maxRetransmitTimeoutMs     int // in milliseconds, <= 0 means optimized default
	initialRetransmitTimeoutMs int // in milliseconds, <= 0 means optimized default
	maxRetransmitAttempts      int // number of retransmissions, <= 0 means optimized default
	heartbeatIntervalMs        int // in milliseconds, <= 0 means optimized default
}

// Note: SCTP settings apply to newly-created PeerConnections only
func rtcSetSctpSettings(settings *RtcSctpSettings) int {
	return int(C.rtcSetSctpSettings((*C.rtcSctpSettings)(unsafe.Pointer(settings))))
}
