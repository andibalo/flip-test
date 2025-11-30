package httpresp

import (
	"fmt"

	"net/http"
	"reflect"
	"time"

	"github.com/andibalo/flip-test/pkg/pagination"
	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
)

const (
	StatusCodeCtxKey = "httpStatusCode"
)

type Meta struct {
	Path       string `json:"path"`
	Code       string `json:"code"`
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Error      string `json:"error,omitempty" swaggerignore:"true"`
	Timestamp  string `json:"timestamp"`
}

type Response struct {
	Data       interface{}            `json:"data"`
	Pagination *pagination.Pagination `json:"pagination,omitempty"`
	Success    string                 `json:"success"`
}

type HTTPErrResp struct {
	Meta    Meta   `json:"metadata"`
	Success string `json:"success"`
}

func HttpRespError(c *gin.Context, err error) {

	statusCode := http.StatusInternalServerError
	respCode := ServerError.AsString()

	oopsErr, ok := oops.AsOops(err)

	if ok {
		respCode = oopsErr.Code()
		errCtx := oopsErr.Context()

		sc, exists := errCtx[StatusCodeCtxKey]
		if exists {
			statusCode = sc.(int)
		}
	}

	jsonErrResp := &HTTPErrResp{
		Meta: Meta{
			Path:       c.Request.URL.Path,
			Code:       respCode,
			StatusCode: statusCode,
			Status:     http.StatusText(statusCode),
			Message:    fmt.Sprintf("%s %s [%d] %s", c.Request.Method, c.Request.RequestURI, statusCode, http.StatusText(statusCode)),
			Error:      err.Error(),
			Timestamp:  time.Now().Format(time.RFC3339),
		},
		Success: "false",
	}

	c.Set("status_code", statusCode)
	c.Set("status", http.StatusText(statusCode))
	c.Set("error", fmt.Sprintf("%s %s [%d] %s", c.Request.Method, c.Request.RequestURI, statusCode, http.StatusText(statusCode)))

	c.AbortWithStatusJSON(statusCode, jsonErrResp)
}

func HttpRespSuccess(c *gin.Context, data interface{}, pagination *pagination.Pagination) {

	kind := reflect.ValueOf(data).Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		//check kalo data nya nil / kosong
		if data == nil || reflect.ValueOf(data).IsNil() {
			data = []interface{}{}
		}
	}

	c.JSON(http.StatusOK, Response{
		Data:       data,
		Pagination: pagination,
		Success:    "success",
	})
}
