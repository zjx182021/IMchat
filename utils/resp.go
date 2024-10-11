package utils

import (
	"encoding/json"
	"net/http"
)

type H struct {
	Code  int
	Msg   string
	Data  any
	Rows  any
	Total any
}

func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, 0, nil, msg)
}
func RespOK(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg)
}

func Resp(w http.ResponseWriter, code int, data any, msg string) {
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(ret)
}

func RespOKList(w http.ResponseWriter, data any, total any) {
	RespList(w, 0, data, total)
}
func RespList(w http.ResponseWriter, code int, data any, total any) {
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code: code,
		Data: data,
		Rows: total,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(ret)
}
