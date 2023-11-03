package service

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
	"github.com/thedevsaddam/renderer"
)

type CoreInterface interface {
	ErrorResponse(http.ResponseWriter, error, int, string)
	DecodeQueryString(req interface{}, url *url.URL) error
	DecodeBody(req interface{}, body io.ReadCloser) error
}

type CoreService struct{}

func (c CoreService) ErrorResponse(w http.ResponseWriter, err error, code int, message string) {
	rnd := renderer.New()
	errMsg := err.Error()
	if message != "" {
		errMsg = message
	}

	rnd.JSON(w, code, map[string]interface{}{
		"error": errMsg,
	})
}

func (c CoreService) SuccessResponse(w http.ResponseWriter, data interface{}) {
	rnd := renderer.New()
	rnd.JSON(w, http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

func (c CoreService) DecodeQueryString(req interface{}, url *url.URL) error {
	var decoder = schema.NewDecoder()
	errDecode := decoder.Decode(req, url.Query())
	if errDecode != nil {
		return errDecode
	}
	return nil
}

func (c CoreService) DecodeBody(req interface{}, body io.ReadCloser) error {
	errParse := json.NewDecoder(body).Decode(&req)
	if errParse != nil {
		return errParse
	}
	return nil
}
