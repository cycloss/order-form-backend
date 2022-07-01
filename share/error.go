package share

import "log"

// if internalMessage is nil, it will not be logged by the buildResponse function
func NewApiErr(statusCode int, externalMessage string, internalMessage string) *ApiErr {
	return &ApiErr{StatusCode: statusCode, ExternalMessage: externalMessage, InternalMessage: internalMessage}
}

type ApiErr struct {
	StatusCode      int
	ExternalMessage string
	InternalMessage string
}

func (ae *ApiErr) Error() string {
	if ae.InternalMessage != "" {
		return ae.InternalMessage
	}
	return ae.ExternalMessage
}

func (ae *ApiErr) LogError() {
	if ae.InternalMessage == "" {
		return
	}
	log.Printf("api error internal message: %s", ae.Error())
}
