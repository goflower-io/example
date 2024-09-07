package service

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	form "github.com/go-playground/form/v4"

	"github.com/happycrud/example/mysql/api"
	"github.com/happycrud/example/mysql/views"
)

type UserHandler struct {
	api.UserServiceServer
	Error   func(w http.ResponseWriter, err error)
	Success func(w http.ResponseWriter, msg string)
}

func NewUserHandler(s api.UserServiceServer) *UserHandler {
	h := &UserHandler{
		UserServiceServer: s,
		Error: func(w http.ResponseWriter, err error) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		},
		Success: func(w http.ResponseWriter, msg string) {
			w.Write([]byte(msg))
		},
	}
	return h
}

func (h *UserHandler) AddPath(addPathFn func(method, path string, hf http.HandlerFunc)) {
	addPathFn(http.MethodGet, api.UserService_CreateUser_FullMethodName, h.CreateUserHandle)
	addPathFn(http.MethodPut, api.UserService_CreateUser_FullMethodName, h.CreateUserHandle)

	addPathFn(http.MethodGet, api.UserService_GetUser_FullMethodName, h.GetUserHandle)
	addPathFn(http.MethodDelete, api.UserService_DeleteUser_FullMethodName, h.DeleteUserHandle)

	addPathFn(http.MethodGet, api.UserService_UpdateUser_FullMethodName, h.UpdateUserHandle)
	addPathFn(http.MethodPost, api.UserService_UpdateUser_FullMethodName, h.UpdateUserHandle)

	addPathFn(http.MethodGet, api.UserService_ListUsers_FullMethodName, h.ListUsersHandle)
}

func (h *UserHandler) ListUsersHandle(w http.ResponseWriter, req *http.Request) {
	reqb := new(api.ListUsersReq)
	if err := GetRequestParams(reqb, req); err != nil {
		h.Error(w, err)
		return
	}
	resp, err := h.UserServiceServer.ListUsers(req.Context(), reqb)
	if err != nil {
		h.Error(w, err)
		return
	}
	switch ResponseConentType(req) {
	case ResponseJSON:
		d, _ := json.Marshal(resp)
		w.Write(d)
	case ResponseHTMX:
		views.UserListView(resp).Render(req.Context(), w)
	case ResponseHTML:
		views.UserListPage(resp).Render(req.Context(), w)
	}
}

func (h *UserHandler) GetUserHandle(w http.ResponseWriter, req *http.Request) {
	reqb := new(api.UserId)
	if err := GetRequestParams(reqb, req); err != nil {
		h.Error(w, err)
		return
	}
	resp, err := h.UserServiceServer.GetUser(req.Context(), reqb)
	if err != nil {
		h.Error(w, err)
		return
	}
	switch ResponseConentType(req) {
	case ResponseJSON:
		d, _ := json.Marshal(resp)
		w.Write(d)
	case ResponseHTMX:
		views.UserDetailView(resp).Render(req.Context(), w)
	case ResponseHTML:
		views.UserDetailPage(resp).Render(req.Context(), w)
	}
}

func (h *UserHandler) UpdateUserHandle(w http.ResponseWriter, req *http.Request) {
	// 区分是GET还是POST方法
	if req.Method == http.MethodGet {
		reqb := new(api.UserId)
		if err := GetRequestParams(reqb, req); err != nil {
			h.Error(w, err)
			return
		}
		resp, err := h.UserServiceServer.GetUser(req.Context(), reqb)
		if err != nil {
			h.Error(w, err)
			return
		}
		switch ResponseConentType(req) {
		case ResponseHTMX:
			views.UserUpdateView(resp).Render(req.Context(), w)
		default:
			views.UserUpdatePage(resp).Render(req.Context(), w)
		}
		return
	}
	reqb := new(api.UpdateUserReq)
	if err := GetRequestParams(reqb, req); err != nil {
		h.Error(w, err)
		return
	}
	_, err := h.UserServiceServer.UpdateUser(req.Context(), reqb)
	if err != nil {
		h.Error(w, err)
		return
	}
	h.Success(w, "User Updated")
}

func (h *UserHandler) DeleteUserHandle(w http.ResponseWriter, req *http.Request) {
	reqb := new(api.UserId)
	if err := GetRequestParams(reqb, req); err != nil {
		h.Error(w, err)
		return
	}
	_, err := h.UserServiceServer.DeleteUser(req.Context(), reqb)
	if err != nil {
		h.Error(w, err)
		return
	}
	h.Success(w, "User Deleted")
}

func (h *UserHandler) CreateUserHandle(w http.ResponseWriter, req *http.Request) {
	// 区分是 GET还是PUT方法
	if req.Method == http.MethodGet {
		switch ResponseConentType(req) {
		case ResponseHTMX:
			views.UserCreateView().Render(req.Context(), w)
		default:
			views.UserCreateView().Render(req.Context(), w)
		}
		return
	}
	reqb := new(api.User)
	if err := GetRequestParams(reqb, req); err != nil {
		h.Error(w, err)
		return
	}
	_, err := h.UserServiceServer.CreateUser(req.Context(), reqb)
	if err != nil {
		h.Error(w, err)
		return
	}
	h.Success(w, "User Created")
}

var formParser = sync.OnceValue(func() *form.Decoder {
	return form.NewDecoder()
})

func ResponseConentType(req *http.Request) ResponseType {
	if req.Header.Get("Accept") == "application/json" {
		return ResponseJSON
	}
	if req.Header.Get("HX-Request") == "true" {
		return ResponseHTMX
	}
	return ResponseHTML
}

type ResponseType string

const (
	ResponseJSON ResponseType = "json"
	ResponseHTML ResponseType = "html"
	ResponseHTMX ResponseType = "htmx"
)

func GetRequestParams(reqb any, req *http.Request) error {
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
