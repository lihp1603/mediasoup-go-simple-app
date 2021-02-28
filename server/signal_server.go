package main

import (
	"fmt"
	"net/http"

	"github.com/jiyeyuran/mediasoup-go"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type WebrtcTranParam struct {
	Id             string                   `json:"id"`
	IceCandidates  []mediasoup.IceCandidate `json:"iceCandidates"`
	IceParameters  mediasoup.IceParameters  `json:"iceParameters"`
	DtlsParameters mediasoup.DtlsParameters `json:"dtlsParameters"`
}

type WebrtcTranDtlsParam struct {
	DtlsParameters mediasoup.DtlsParameters `json:"dtlsParameters"`
}

type PipeTranParam struct {
	Ip   string `json:"ip"`
	Port uint16 `json:"port"`
}

type ProduceParam struct {
	Id string `json:"id"`
}

type ConsumeParam struct {
	ProducerId    string                  `json:"producerId"`
	Id            string                  `json:"id"`
	Kind          string                  `json:"kind"`
	RtpParameters mediasoup.RtpParameters `json:"rtpParameters"`
}

type SignalServer struct {
	server      *gosocketio.Server
	mediaServer *MediaServer
}

func NewSignalServer(mediaServer *MediaServer) *SignalServer {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	p := &SignalServer{
		server:      server,
		mediaServer: mediaServer,
	}
	return p
}

func (p *SignalServer) Start(ip, port string, cert, key string) {
	p.server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		fmt.Println("connected:")
	})

	p.server.On(gosocketio.OnError, func(c *gosocketio.Channel) {
		fmt.Println("meet error:")
	})

	p.server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		fmt.Println("closed:")
	})

	p.server.On("getRouterRtpCapabilities", func(c *gosocketio.Channel) interface{} {
		fmt.Println("getRouterRtpCapabilities:")
		return p.mediaServer.router.RtpCapabilities()
	})

	p.server.On("createProducerTransport", func(c *gosocketio.Channel) interface{} {
		fmt.Println("createProducerTransport")
		p.mediaServer.CreateRtcProduceTransport(ip)
		param := &WebrtcTranParam{
			Id:             p.mediaServer.rtcProduceTransport.Id(),
			IceCandidates:  p.mediaServer.rtcProduceTransport.IceCandidates(),
			IceParameters:  p.mediaServer.rtcProduceTransport.IceParameters(),
			DtlsParameters: p.mediaServer.rtcProduceTransport.DtlsParameters(),
		}
		return param
	})

	p.server.On("connectProducerTransport", func(c *gosocketio.Channel, msg WebrtcTranDtlsParam) interface{} {
		fmt.Println("connectProducerTransport:", msg)
		opt := mediasoup.TransportConnectOptions{
			DtlsParameters: &msg.DtlsParameters,
		}
		err := p.mediaServer.rtcProduceTransport.Connect(opt)
		if err != nil {
			fmt.Println("connectProducerTransport err:", err)
		}
		return nil
	})

	p.server.On("produce", func(c *gosocketio.Channel, opt mediasoup.ProducerOptions) interface{} {
		fmt.Printf("produce: %+v\n", opt)
		for _, codec := range opt.RtpParameters.Codecs {
			fmt.Printf("codec: %+v\n", codec)
		}

		var err error
		p.mediaServer.producer, err = p.mediaServer.rtcProduceTransport.Produce(opt)
		if err != nil {
			fmt.Println("produce err:", err)
		}
		fmt.Println("produce end")
		param := &ProduceParam{
			Id: p.mediaServer.producer.Id(),
		}

		return param
	})

	//////////////////////////////////////////////////////////////////////////////////////////////
	p.server.On("createConsumerTransport", func(c *gosocketio.Channel) interface{} {
		fmt.Println("createConsumerTransport:")
		p.mediaServer.CreateRtcConsumeTransport(ip)
		param := &WebrtcTranParam{
			Id:             p.mediaServer.rtcConsumeTransport.Id(),
			IceCandidates:  p.mediaServer.rtcConsumeTransport.IceCandidates(),
			IceParameters:  p.mediaServer.rtcConsumeTransport.IceParameters(),
			DtlsParameters: p.mediaServer.rtcConsumeTransport.DtlsParameters(),
		}
		return param
	})

	p.server.On("connectConsumerTransport", func(c *gosocketio.Channel, msg WebrtcTranDtlsParam) interface{} {
		fmt.Println("connectConsumerTransport:", msg)
		opt := mediasoup.TransportConnectOptions{
			DtlsParameters: &msg.DtlsParameters,
		}
		err := p.mediaServer.rtcConsumeTransport.Connect(opt)
		if err != nil {
			fmt.Println("connectConsumerTransport err:", err)
		}
		return nil
	})

	p.server.On("consume", func(c *gosocketio.Channel, opt mediasoup.ConsumerOptions) interface{} {
		fmt.Println("consume:", opt)
		opt.ProducerId = p.mediaServer.producer.Id()
		consumer, err := p.mediaServer.rtcConsumeTransport.Consume(opt)
		if err != nil {
			fmt.Println("consume err:", err)
		}
		param := &ConsumeParam{
			ProducerId:    consumer.ProducerId(),
			Id:            consumer.Id(),
			Kind:          string(consumer.Kind()),
			RtpParameters: consumer.RtpParameters(),
		}
		fmt.Println("consume end")
		return param
	})

	p.server.On("resume", func(c *gosocketio.Channel) interface{} {
		fmt.Println("resume:")
		return nil
	})
	//////////////////////////////////////////////////////////////////////////////////////////////
	p.server.On("createPipeTransport", func(c *gosocketio.Channel, msg PipeTranParam) interface{} {
		fmt.Println("createPipeTransport:")
		p.mediaServer.CreatePipeConsumeTransport(ip)
		opt := mediasoup.TransportConnectOptions{
			Ip:   msg.Ip,
			Port: msg.Port,
		}
		err := p.mediaServer.pipeConsumeTransport.Connect(opt)
		if err != nil {
			fmt.Println("createPipeTransport err:", err)
		}

		param := &PipeTranParam{
			Ip:   p.mediaServer.pipeConsumeTransport.Tuple().LocalIp,
			Port: p.mediaServer.pipeConsumeTransport.Tuple().LocalPort,
		}
		return param
	})

	p.server.On("pipeConsume", func(c *gosocketio.Channel, opt mediasoup.ConsumerOptions) interface{} {
		fmt.Println("pipeConsume:", opt)
		opt.ProducerId = p.mediaServer.producer.Id()
		consumer, err := p.mediaServer.pipeConsumeTransport.Consume(opt)
		if err != nil {
			fmt.Println("pipeConsume err:", err)
		}
		param := &ConsumeParam{
			ProducerId:    consumer.ProducerId(),
			Id:            consumer.Id(),
			Kind:          string(consumer.Kind()),
			RtpParameters: consumer.RtpParameters(),
		}
		fmt.Println("pipeConsume end")
		return param
	})

	//go p.server.Serve()

	//http.HandleFunc("/server/", func(w http.ResponseWriter, r *http.Request) {
	//	p.server.ServeHTTP(w, r)
	//})
	//go http.ListenAndServe(ip +":"+ port, nil)

	//
	//setup http server
	serveMux := http.NewServeMux()
	serveMux.Handle("/server/", p.server)
	addr := ip + ":" + port
	fmt.Println("http listen addr:" + addr)

	//	go http.ListenAndServe(addr, serveMux)
	//使用https
	go http.ListenAndServeTLS(addr, cert, key, serveMux)

}
