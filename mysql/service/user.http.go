package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/go-playground/form"

	"github.com/happycrud/example/mysql/api"
)

type UserHander struct {
	api.UserServiceServer
}

func (h *UserHander) AddRouter(r *http.ServeMux) {
	r.HandleFunc("PUT "+api.UserService_CreateUser_FullMethodName, h.CreateUserHandle)
	r.HandleFunc("GET "+api.UserService_GetUser_FullMethodName, nil)
	r.HandleFunc("DELETE "+api.UserService_DeleteUser_FullMethodName, nil)
	r.HandleFunc("GET "+api.UserService_ListUsers_FullMethodName, nil)
	r.HandleFunc("POST "+api.UserService_UpdateUser_FullMethodName, nil)
}

func (h *UserHander) CreateUserHandle(w http.ResponseWriter, req *http.Request) {
	var reqb api.User
	if err := GetRequestBody(&reqb, req); err != nil {
		return
	}
	resp, err := h.UserServiceServer.CreateUser(req.Context(), &reqb)
	if err != nil {
		return
	}
	fmt.Println(resp)
}

var formParser = sync.OnceValue(func() *form.Decoder {
	return form.NewDecoder()
})

func MakeResponse() {
}

func GetRequestBody(reqb any, req *http.Request) error {
	contentType := req.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		var body []byte
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(body, reqb); err != nil {
			return err
		}

	case "application/x-www-form-urlencoded":
		if err := req.ParseForm(); err != nil {
			return err
		}
		if err := formParser().Decode(reqb, req.Form); err != nil {
			return err
		}
	default:
		return errors.New("not support contentType")
	}
	return nil
}
