package main

import (
	"encoding/json"
	"github.com/jiyeyuran/mediasoup-go"
	"github.com/jiyeyuran/mediasoup-go/h264"
)

var logger = mediasoup.NewLogger("ExampleApp")

type MediaServer struct {
	worker               *mediasoup.Worker
	router               *mediasoup.Router
	rtcConsumeTransport  *mediasoup.WebRtcTransport
	rtcProduceTransport  *mediasoup.WebRtcTransport
	producer             *mediasoup.Producer
	pipeConsumeTransport *mediasoup.PipeTransport
	pipeProduceTransport *mediasoup.PipeTransport
}

func NewMediaServer() *MediaServer {
	p := &MediaServer{}
	p.CreateWorker()
	return p
}

func (p *MediaServer) CreateWorker() {
	mediasoup.WorkerBin = "./worker/mediasoup-worker"
	var err error
	p.worker, err = mediasoup.NewWorker()
	if err != nil {
		panic(err)
	}
	p.worker.On("died", func(err error) {
		logger.Error("died: %s", err)
	})

	dump, _ := p.worker.Dump()
	logger.Debug("dump: %+v", dump)

	usage, err := p.worker.GetResourceUsage()
	if err != nil {
		panic(err)
	}
	data, _ := json.Marshal(usage)
	logger.Debug("usage: %s", data)

	p.CreateRouter()
}

func (p *MediaServer) CreateRouter() {
	var err error
	p.router, err = p.worker.CreateRouter(mediasoup.RouterOptions{
		MediaCodecs: []*mediasoup.RtpCodecCapability{
			{
				Kind:      "audio",
				MimeType:  "audio/opus",
				ClockRate: 48000,
				Channels:  2,
			},
			{
				Kind:      "video",
				MimeType:  "video/VP8",
				ClockRate: 90000,
				RtcpFeedback: []mediasoup.RtcpFeedback{
					{Type: "nack", Parameter: ""},
					{Type: "nack", Parameter: "pli"},
					{Type: "ccm", Parameter: "fir"},
					{Type: "goog-remb", Parameter: ""},
				},
			},
			{
				Kind:      "video",
				MimeType:  "video/H264",
				ClockRate: 90000,
				Parameters: mediasoup.RtpCodecSpecificParameters{
					RtpParameter: h264.RtpParameter{
						LevelAsymmetryAllowed: 1,
						PacketizationMode:     1,
						ProfileLevelId:        "4d0032",
					},
				},
				RtcpFeedback: []mediasoup.RtcpFeedback{
					{Type: "nack", Parameter: ""},
					{Type: "nack", Parameter: "pli"},
					{Type: "ccm", Parameter: "fir"},
					{Type: "goog-remb", Parameter: ""},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
}

func (p *MediaServer) CreateRtcProduceTransport(ipAddr string) {
	var err error
	p.rtcProduceTransport, err = p.router.CreateWebRtcTransport(mediasoup.WebRtcTransportOptions{
		ListenIps: []mediasoup.TransportListenIp{
			{Ip: ipAddr, AnnouncedIp: ipAddr}, // AnnouncedIp is optional
		},
	})
	if err != nil {
		panic(err)
	}
}

func (p *MediaServer) CreateRtcConsumeTransport(ipAddr string) {
	var err error
	p.rtcConsumeTransport, err = p.router.CreateWebRtcTransport(mediasoup.WebRtcTransportOptions{
		ListenIps: []mediasoup.TransportListenIp{
			{Ip: ipAddr, AnnouncedIp: ipAddr}, // AnnouncedIp is optional
		},
	})
	if err != nil {
		panic(err)
	}
	//	p.rtcConsumeTransport.IceSelectedTuple()
	//	p.pipeConsumeTransport.Tuple()
}

func (p *MediaServer) CreatePipeProduceTransport(ipAddr string) {
	var err error
	p.pipeProduceTransport, err = p.router.CreatePipeTransport(mediasoup.PipeTransportOptions{
		ListenIp:  mediasoup.TransportListenIp{Ip: ipAddr, AnnouncedIp: ipAddr},
		EnableRtx: true,
	})
	if err != nil {
		panic(err)
	}
}

func (p *MediaServer) CreatePipeConsumeTransport(ipAddr string) {
	var err error
	p.pipeConsumeTransport, err = p.router.CreatePipeTransport(mediasoup.PipeTransportOptions{
		ListenIp:  mediasoup.TransportListenIp{Ip: ipAddr, AnnouncedIp: ipAddr},
		EnableRtx: true,
	})
	if err != nil {
		panic(err)
	}
}
