package logger

type DataService struct {
	Name    string      `json:"name"`
	Query   string      `json:"query"`
	Message string      `json:"message"`
	Args    interface{} `json:"args"`
}

type ExternalService struct {
	URL          string `json:"url"`
	Method       string `json:"method"`
	RequestBody  string `json:"requestBody"`
	ResponseCode int    `json:"responseCode"`
	ResponseBody string `json:"responseBody"`
	Message      string `json:"message"`
}

type HttpRequest struct {
	Method      string            `json:"method"`
	Url         string            `json:"url"`
	Agent       string            `json:"agent"`
	IP          string            `json:"ip"`
	Status      int               `json:"status"`
	Headers     map[string]string `json:"string"`
	QueryParams map[string]string `json:"string"`
	Body        string            `json:"body"`
}
