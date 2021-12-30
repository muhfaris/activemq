package activemq

import (
	stomp "github.com/go-stomp/stomp"
	"github.com/muhfaris/activemq/pkg/logger"
)

var logData = logger.LogModel{}

type produceData struct {
	Msg    string `json:"msg"`
	Source string `json:"source"`
}

type AmqpClient struct {
	Username string
	Password string
	Address  string
	Channels []string

	Custom     bool
	Processor  func(inp string) string
	Subscriber func(inp string) string

	conn   *stomp.Conn
	sub    map[string]*stomp.Subscription
	msgIn  []byte
	msgOut []byte

	//check
	DomainKey            string
	SourceApp            string
	TargetApp            string
	ProcessName          string
	ProcessGroup         string
	InputJMSDestination  string
	OutputJMSDestination string
	JMS_HostPort         string
	JMS_User             string
}

type Process struct {
	IAmqpClient
}

type IAmqpClient interface {
	Connect()
	LogMessageBeforeProcess(msg, source string) (rsp string)
	LogMessageAfterProcess(msg, source string) (rsp string)
	Process(inp1 string) (out2 string)
	Send(msg string) string
}
