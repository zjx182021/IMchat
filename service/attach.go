package service

import "net/http"

func upload(writer http.ResponseWriter, request *http.Request){
	request.FormFile("file")
}