package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/bradenrayhorn/beans/beans"
)

var codeToHTTPStatus = map[string]int{
	beans.EINTERNAL: http.StatusInternalServerError,
	beans.EINVALID:  http.StatusUnprocessableEntity,
}

func Error(w http.ResponseWriter, err error) {
	type errorResponse struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}

	var code = beans.EINTERNAL
	var msg = beans.ErrorInternal.Error()
	var beansError beans.Error
	if errors.As(err, &beansError) {
		code, msg = beansError.BeansError()
	}

	if code == beans.EINTERNAL {
		log.Println(err)
	}

	w.WriteHeader(codeToHTTPStatus[code])
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(&errorResponse{Error: msg, Code: code})
}
