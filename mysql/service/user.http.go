package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	form "github.com/go-playground/form/v4"

	"github.com/happycrud/example/mysql/api"
)

type UserHandler struct {
	api.UserServiceServer
}

func (h *UserHandler) AddPath(addPathFn func(method, path string, hf http.HandlerFunc)) {
	addPathFn(http.MethodPut, api.UserService_CreateUser_FullMethodName, h.CreateUserHandle)
	addPathFn(http.MethodGet, api.UserService_GetUser_FullMethodName, h.GetUserHandle)
	addPathFn(http.MethodDelete, api.UserService_DeleteUser_FullMethodName, h.DeleteUserHandle)
	addPathFn(http.MethodPost, api.UserService_UpdateUser_FullMethodName, h.UpdateUserHandle)
	addPathFn(http.MethodGet, api.UserService_ListUsers_FullMethodName, h.ListUsersHandle)
}

func (h *UserHandler) ListUsersHandle(w http.ResponseWriter, req *http.Request) {
	reqb := new(api.ListUsersReq)
	if err := GetRequestBody(reqb, req); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	resp, err := h.UserServiceServer.ListUsers(req.Context(), reqb)
	if err != nil {
		fmt.Fprint(w, err.Error())

		return
	}

	fmt.Fprintf(w, "%v", resp)
}

func (h *UserHandler) GetUserHandle(w http.ResponseWriter, req *http.Request) {
	reqb := new(api.UserId)
	if err := GetRequestBody(reqb, req); err != nil {
		fmt.Fprint(w, err.Error())

		return
	}
	resp, err := h.UserServiceServer.GetUser(req.Context(), reqb)
	if err != nil {
		fmt.Fprint(w, err.Error())

		return
	}
	fmt.Fprintf(w, "%v", resp)
}

func (h *UserHandler) UpdateUserHandle(w http.ResponseWriter, req *http.Request) {
	reqb := new(api.UpdateUserReq)
	if err := GetRequestBody(reqb, req); err != nil {
		fmt.Fprint(w, err.Error())

		return
	}
	resp, err := h.UserServiceServer.UpdateUser(req.Context(), reqb)
	if err != nil {
		fmt.Fprint(w, err.Error())

		return
	}
	fmt.Fprintf(w, "%v", resp)
}

func (h *UserHandler) DeleteUserHandle(w http.ResponseWriter, req *http.Request) {
	reqb := new(api.UserId)
	if err := GetRequestBody(reqb, req); err != nil {
		fmt.Fprint(w, err.Error())

		return
	}
	resp, err := h.UserServiceServer.GetUser(req.Context(), reqb)
	if err != nil {
		fmt.Fprint(w, err.Error())

		return
	}
	fmt.Fprintf(w, "%v", resp)
}

func (h *UserHandler) CreateUserHandle(w http.ResponseWriter, req *http.Request) {
	reqb := new(api.User)
	if err := GetRequestBody(reqb, req); err != nil {
		fmt.Fprint(w, err.Error())

		return
	}
	resp, err := h.UserServiceServer.CreateUser(req.Context(), reqb)
	if err != nil {
		fmt.Fprint(w, err.Error())

		return
	}
	fmt.Fprintf(w, "%v", resp)
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
	default:
		// "application/x-www-form-urlencoded":
		if err := req.ParseForm(); err != nil {
			return err
		}
		if err := formParser().Decode(reqb, req.Form); err != nil {
			return err
		}
	}
	return nil
}
