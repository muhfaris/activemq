package activemq

import (
	"encoding/json"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

const (
	CONST_DOMAIN_KEY               = "CH_DOMAINKEY"
	CONST_PEID                     = "CH_PEID"
	CONST_LOG_TYPE                 = "CH_LogType"
	CONST_LOG_TYPE_INFO            = "INFO"
	CONST_LOG_TYPE_ERROR           = "ERROR"
	CONST_LOG_DIRECTION            = "CH_LogDirection"
	CONST_LOG_DIRECTION_S          = "S"
	CONST_LOG_DIRECTION_T          = "T"
	CONST_LOG_DIRECTION_TR         = "T_R"
	CONST_LOG_DIRECTION_SR         = "S_R"
	CONST_EXECUTION_MODE           = "CH_ExecutionMode"
	CONST_EXECUTION_MODE_LISTENER  = "ListenerExecution"
	CONST_EXECUTION_MODE_MANUAL    = "ManualRetry"
	CONST_EXECUTION_MODE_SCHEDULED = "ScheduledExecution"
	CONST_EXECUTION_USER           = "CH_ExecutionUser"
	CONST_EXECUTION_USER_SYSTEM    = "System"
	CONST_ENDPOINT_TYPE            = "CH_EndpointType"
	CONST_ENDPOINT_TYPE_FTP        = "FTP"
	CONST_ENDPOINT_TYPE_MAIL       = "MAIL"
	CONST_ENDPOINT_TYPE_HTTP       = "HTTP"
	CONST_ENDPOINT_TYPE_JMS        = "JMS"
	CONST_ENDPOINT_TYPE_SOLACE     = "SOLACE"
	CONST_ENDPOINT_META            = "CH_EndpointMeta"
	CONST_FILE_NAME                = "CH_file_name"
	CONST_FTP_USER                 = "CH_ftp_user"
	CONST_FILE_REMOTE_FILE         = "CH_file_remoteFile"
	CONST_FILE_REMOTE_DIRECTORY    = "CH_file_remoteDirectory"
	CONST_FILE_REMOTE_HOST_PORT    = "CH_file_remoteHostPort"
	CONST_JMS_DESTINATION          = "CH_jms_destination"
	CONST_JMS_HOSTPORT             = "CH_jms_hostPort"
	CONST_JMS_USER                 = "CH_jms_user"
	CONST_HTTP_URL                 = "http_url"
	CONST_HTTP_USER                = "http_user"
	CONST_SOURCE_APP               = "CH_SourceApp"
	CONST_TARGET_APP               = "CH_TargetApp"
	CONST_PROCESS_NAME             = "CH_ProcessName"
	CONST_PROCESS_GROUP            = "CH_ProcessGroup"
)

func (am *AmqpClient) getInputHeaders() map[string]string {
	headerMap := make(map[string]string)
	headerMap[CONST_DOMAIN_KEY] = am.DomainKey
	uuid := uuid.NewV4()
	headerMap[CONST_PEID] = uuid.String()
	headerMap[CONST_LOG_TYPE] = CONST_LOG_TYPE_INFO
	headerMap[CONST_LOG_DIRECTION] = CONST_LOG_DIRECTION_S
	headerMap[CONST_EXECUTION_MODE] = CONST_EXECUTION_MODE_LISTENER
	headerMap[CONST_EXECUTION_USER] = CONST_EXECUTION_USER_SYSTEM
	headerMap[CONST_ENDPOINT_TYPE] = CONST_ENDPOINT_TYPE_JMS

	var meta = fmt.Sprintf("%s=%s,%s=%s,%s=%s", CONST_JMS_DESTINATION, am.InputJMSDestination,
		CONST_JMS_HOSTPORT, am.JMS_HostPort, CONST_JMS_USER, am.JMS_User)
	headerMap[CONST_ENDPOINT_META] = meta

	headerMap[CONST_SOURCE_APP] = am.SourceApp
	headerMap[CONST_TARGET_APP] = am.TargetApp
	headerMap[CONST_PROCESS_NAME] = am.ProcessName
	headerMap[CONST_PROCESS_GROUP] = am.ProcessGroup

	return headerMap
}

func (am *AmqpClient) getOutputHeaders(msg string) map[string]string {
	headerMap := make(map[string]string)

	//todo read input headers and add

	headerMap[CONST_ENDPOINT_TYPE] = "JMS"
	headerMap[CONST_LOG_TYPE] = CONST_LOG_TYPE_INFO
	headerMap[CONST_LOG_DIRECTION] = CONST_LOG_DIRECTION_T

	var meta = fmt.Sprintf("%s=%s,%s=%s,%s=%s", CONST_JMS_DESTINATION, am.OutputJMSDestination,
		CONST_JMS_HOSTPORT, am.JMS_HostPort, CONST_JMS_USER, am.JMS_User)
	headerMap[CONST_ENDPOINT_META] = meta

	return headerMap
}

func (am *AmqpClient) LogMessageBeforeProcess(msg, source string) (rsp string) {
	logData.Info(200, "(IN) msg inserted into --> "+am.Channels[0])
	headerMap := am.getInputHeaders()

	am.produce(am.Channels[0], msg, source, headerMap)
	// here channel is in wait until and unless it doesn't get notification from background go routine
	for {
		if am.msgIn != nil {
			var p = new(produceData)
			err := json.Unmarshal(am.msgIn, &p)
			if err != nil {
				fmt.Println(err)
			}
			am.msgIn = nil
			rsp = p.Msg
			logData.WriteToLog()
			return
		}
	}
}

func (am *AmqpClient) LogMessageAfterProcess(msg, source string) (rsp string) {
	logData.Info(200, "(IN) msg inserted into --> "+am.Channels[1])
	headerMap := am.getOutputHeaders(msg)

	am.produce(am.Channels[1], msg, source, headerMap)
	for {
		if am.msgOut != nil {
			var p = new(produceData)
			err := json.Unmarshal(am.msgOut, &p)
			if err != nil {
				fmt.Println(err)
			}
			am.msgOut = nil
			rsp = p.Msg
			logData.WriteToLog()
			return
		}
	}

}

func (am *AmqpClient) Process(inp1 string) (out2 string) {
	logData.Info(200, "(CONVERSION) processing "+am.Channels[0]+" data")
	logData.WriteToLog()
	if am.Custom {
		return am.Processor(inp1)
	}

	//Default Logic
	return strings.ToUpper(inp1)
}

func (am *AmqpClient) Send(inp2 string) (out4 string) {
	fmt.Println(inp2)
	logData.Info(200, "(OUT) calling backend API ")
	logData.WriteToLog()
	if am.Custom {
		return am.Subscriber(inp2)
	}

	//Default Logic
	return "CALLED FROM BACKEND API" + inp2
}

func Close(am *AmqpClient) {
	defer am.conn.Disconnect()
}
