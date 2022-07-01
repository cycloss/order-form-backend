package share

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Aborts the handler chain. data interface must be json serializable.
func RespondJson(c *gin.Context, data interface{}, err error) {
	var resp = buildResponse(data, err)
	c.JSON(resp.Status, resp)
	c.Abort()
}

// If error is not nil or data is nil then a non 200 response will be returned
// If error if ServerError or SfAuthError, the error is automatically logged
func buildResponse(data interface{}, err error) *Response {

	// error logging should only occur in this block
	if err != nil {
		if ae, ok := err.(*ApiErr); ok {
			ae.LogError()
			return newResponse(ae.StatusCode, ae.ExternalMessage, nil)
		}
		log.Printf("api error internal message: %s", err.Error())
		return newResponse(http.StatusInternalServerError, "the server encountered an unknown internal error", nil)

	}
	if data == nil {
		return newResponse(http.StatusInternalServerError, "No error occured but the server failed to generate any data", nil)
	}

	return newResponse(http.StatusOK, "Success", data)

}

type Response struct {
	Meta `json:"meta"`
	// allows any struct to be passed in
	Response interface{} `json:"response"`
}

type Meta struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func newResponse(status int, message string, response interface{}) *Response {
	return &Response{
		Meta{
			status,
			message,
		},
		response,
	}
}
