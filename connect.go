package activemq

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	stomp "github.com/go-stomp/stomp"
	"github.com/go-stomp/stomp/frame"
)

// Connect to activeMQ
var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
	// stomp.ConnOpt.Login("userid", "userpassword"),
	// stomp.ConnOpt.Host("localhost"),
	stomp.ConnOpt.HeartBeat(360*time.Second, 360*time.Second),
	stomp.ConnOpt.HeartBeatError(360 * time.Second),
}

func (am *AmqpClient) Connect() {
	options = append(options, stomp.ConnOpt.Login(am.Username, am.Password))
	conn, err := stomp.Dial("tcp", am.Address, options...)
	if err != nil {
		logData.AppendMiscError("Error occurred while connecting ActiveMQ: " + err.Error())
		logData.WriteToLog()
		os.Exit(1)
	}
	logData.AppendMiscLogger("Connncted to Active-MQ")
	logData.WriteToInfoLog()
	am.conn = conn
	am.sub = make(map[string]*stomp.Subscription)

	if am.Custom {
		if am.Processor == nil {
			logData.AppendMiscError("Processor Function Couldn't found for Custom Functions")
			logData.WriteToLog()
			os.Exit(1)
		}
		if am.Subscriber == nil {
			logData.AppendMiscError("Subscriber Function Couldn't found for Custom Functions")
			logData.WriteToLog()
			os.Exit(1)
		}
	}

	// starting all channels
	go am.listner()
}

func (am *AmqpClient) listner() {

	go func() {
		channelName := am.Channels[0]
		sub, err := am.conn.Subscribe(channelName,
			stomp.AckAuto,
		// stomp.SubscribeOpt.Header("", ""),
		)
		if err != nil {
			logData.AppendMiscError("Error occurred while receiving: " + err.Error())
			logData.WriteToLog()
		}
		am.sub[channelName] = sub
		for {
			msg := am.subscribeIn(sub, channelName)
			if msg != nil {
				var p = new(produceData)
				err := json.Unmarshal(am.msgIn, &p)
				if err != nil {
					am.processAndSend(string(am.msgIn), "ASYNC")
					am.msgIn = nil
				}
				// else if p.Source == "API" {
				// 	// channel notify
				// 	// chan <- p.Msg
				// }
			}
		}
	}()

	go func() {
		channelName := am.Channels[1]

		sub, err := am.conn.Subscribe(channelName, stomp.AckAuto)
		if err != nil {
			logData.AppendMiscError("Error occurred while receiving: " + err.Error())
			logData.WriteToLog()
		}
		am.sub[channelName] = sub

		for {
			msg := am.subscribeOut(sub, channelName)
			if msg != nil {

				var p = new(produceData)
				err := json.Unmarshal(am.msgOut, &p)
				if err != nil {
					fmt.Println(err)
				}
				if p.Source != "API" {
					am.Send(p.Msg)
					am.msgOut = nil
				}
				// if p.Source != "API" {
				// 	am.Send(p.Msg, "ASYNC")
				// 	am.msgOut = nil
				// }
			}
		}
	}()
}

func (am *AmqpClient) subscribeIn(s *stomp.Subscription, channel string) (msg interface{}) {
	logData.Info(200, "(OUT) Getting msg from "+channel)
	m := <-s.C
	if m.Body != nil {
		am.msgIn = m.Body
		msg = m.Body
		return
	}
	return
}

func (am *AmqpClient) subscribeOut(s *stomp.Subscription, channel string) (msg interface{}) {
	logData.Info(200, "(OUT) Getting msg from "+channel)
	m := <-s.C
	if m.Body != nil {
		am.msgOut = m.Body
		msg = m.Body
		return
	}
	return
}

func (am *AmqpClient) produce(channelName, message, source string, headerMap map[string]string) (err error) {
	var p = new(produceData)
	p.Msg = message
	p.Source = source
	b, _ := json.Marshal(p)

	headers := make([]func(*frame.Frame) error, len(headerMap))
	i := 0
	for key, value := range headerMap {
		fmt.Println(key, ":", value)
		headers[i] = stomp.SendOpt.Header(key, value)
		i = i + 1
	}

retry:
	serr := am.conn.Send(
		channelName,  // destination
		"text/plain", // content-type
		[]byte(b),
		// stomp.SendOpt.Receipt,
		// stomp.SendOpt.Header("", ""),
	) // body
	if serr != nil {
		if serr == stomp.ErrAlreadyClosed {
			logData.Error(500, "(IN) Active MQ Server Closed retry...")
			am.Connect()
			goto retry
		}
		logData.Error(500, "(IN) Error occurred while sending to producer: "+err.Error())
		logData.WriteToLog()
		err = serr
		return
	}
	logData.Info(200, "(IN) Msg submitted successfully to "+channelName)
	logData.WriteToInfoLog()

	return nil
}

func (am *AmqpClient) processAndSend(out1, source string) (out4 string) {
	//do spme process or conversion
	out2 := am.Process(out1)
	fmt.Println("OUT2", out2)
	// logType=INFO, direction=OUT
	out3 := am.LogMessageAfterProcess(out2, source)
	fmt.Println("OUT3", out3)

	//come here
	out4 = am.Send(out3)
	return
}
