package logger

import (
	"io/ioutil"
	"net/http"

	xid "github.com/rs/xid"
	logrus "github.com/sirupsen/logrus"
)

type LogModel struct {
	Request          HttpRequest       `json:"api"`
	DataServices     []DataService     `json:"dataservices"`
	ExternalServices []ExternalService `json:"externalService"`
	Message          string            `json:"message"`
	ID               string            `json:"id"`
	Misc             []string          `json:"misc"`
}

func (log *LogModel) Initialize(r *http.Request) {
	log.Request = HttpRequest{
		Method: r.Method,
		Url:    r.RequestURI,
		IP:     r.RemoteAddr,
	}
	log.ID = xid.New().String()
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err == nil {
		bodyString := string(bodyBytes)
		log.Request.Body = bodyString
	}
	log.Request.Headers = map[string]string{}
	for key, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			log.Request.Headers[key] = value
		}
	}
	log.DataServices = []DataService{}
	log.ExternalServices = []ExternalService{}
	log.Misc = []string{}
}

func (log *LogModel) InitializeInternalLog() {
	log.ID = xid.New().String()
	log.DataServices = []DataService{}
	log.ExternalServices = []ExternalService{}
	log.Misc = []string{}
}

func (log *LogModel) AppendDataService(dataService DataService) {
	log.DataServices = append(log.DataServices, dataService)
}

func (log *LogModel) AddUserIdentity(user interface{}) {
	if user == nil {
		return
	}
	// log.UserIdentity = UserIdentity{
	// 	UserID: user.(accounts.UserAccount).UserID,
	// 	Role:   user.(accounts.UserAccount).Role,
	// }
}

func (log *LogModel) AddExternalService(service ExternalService) {
	log.ExternalServices = append(log.ExternalServices, service)
}

func (log *LogModel) AppendMiscError(misc string) {
	log.Misc = append(log.Misc, misc)
}

func (log *LogModel) AppendMiscLogger(misc string) {
	log.Misc = append(log.Misc, misc)
}

func (log *LogModel) Info(status int, message string) {
	log.Message = message
	log.Request.Status = status
}

func (log *LogModel) Error(status int, message string) {
	log.Message = message
	log.Request.Status = status
}

func (log *LogModel) WriteToLog() {
	logrus.WithFields(logrus.Fields{"request": log.Request, "dataservices": log.DataServices, "externalServices": log.ExternalServices, "logId": log.ID, "misc": log.Misc}).Error(log.Message)
}

func (log *LogModel) WriteToInfoLog() {
	logrus.WithFields(logrus.Fields{"request": log.Request, "dataservices": log.DataServices, "externalServices": log.ExternalServices, "logId": log.ID, "misc": log.Misc}).Info(log.Message)
}
