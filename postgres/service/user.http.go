package service

import (
	"encoding/json"
	"net/http"

	h1 "github.com/happycrud/golib/net/http"

	"github.com/happycrud/example/postgres/api"
	"github.com/happycrud/example/postgres/views"
)

type UserHandler struct {
	api.UserServiceServer
	Error   func(w http.ResponseWriter, statusCode int, err error)
	Success func(w http.ResponseWriter, msg string)
}

func NewUserHandler(s api.UserServiceServer) *UserHandler {
	h := &UserHandler{
		UserServiceServer: s,
		Error: func(w http.ResponseWriter, statusCode int, err error) {
			w.WriteHeader(statusCode)
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
	if err := h1.GetRequestParams(reqb, req); err != nil {
		h.Error(w, http.StatusBadRequest, err)
		return
	}
	resp, err := h.UserServiceServer.ListUsers(req.Context(), reqb)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err)
		return
	}
	switch h1.ResponseConentType(req) {
	case h1.ResponseJSON:
		d, _ := json.Marshal(resp)
		w.Write(d)
	case h1.ResponseHTMX:
		views.UserListView(resp, req.URL.Path, req.Form).Render(req.Context(), w)
	case h1.ResponseHTML:
		views.UserListPage(resp, req.URL.Path, req.Form).Render(req.Context(), w)
	}
}

func (h *UserHandler) GetUserHandle(w http.ResponseWriter, req *http.Request) {
	reqb := new(api.UserId)
	if err := h1.GetRequestParams(reqb, req); err != nil {
		h.Error(w, http.StatusBadRequest, err)
		return
	}
	resp, err := h.UserServiceServer.GetUser(req.Context(), reqb)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err)
		return
	}
	switch h1.ResponseConentType(req) {
	case h1.ResponseJSON:
		d, _ := json.Marshal(resp)
		w.Write(d)
	case h1.ResponseHTMX:
		views.UserDetailView(resp).Render(req.Context(), w)
	case h1.ResponseHTML:
		views.UserDetailPage(resp).Render(req.Context(), w)
	}
}

func (h *UserHandler) UpdateUserHandle(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		reqb := new(api.UserId)
		if err := h1.GetRequestParams(reqb, req); err != nil {
			h.Error(w, http.StatusBadRequest, err)
			return
		}
		resp, err := h.UserServiceServer.GetUser(req.Context(), reqb)
		if err != nil {
			h.Error(w, http.StatusInternalServerError, err)
			return
		}
		switch h1.ResponseConentType(req) {
		case h1.ResponseHTMX:
			views.UserUpdateView(resp).Render(req.Context(), w)
		default:
			views.UserUpdatePage(resp).Render(req.Context(), w)
		}
		return
	}
	reqb := new(api.UpdateUserReq)
	if err := h1.GetRequestParams(reqb, req); err != nil {
		h.Error(w, http.StatusBadRequest, err)
		return
	}
	_, err := h.UserServiceServer.UpdateUser(req.Context(), reqb)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err)
		return
	}
	h.Success(w, "User Updated")
}

func (h *UserHandler) DeleteUserHandle(w http.ResponseWriter, req *http.Request) {
	reqb := new(api.UserId)
	if err := h1.GetRequestParams(reqb, req); err != nil {
		h.Error(w, http.StatusBadRequest, err)
		return
	}
	_, err := h.UserServiceServer.DeleteUser(req.Context(), reqb)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err)
		return
	}
	h.Success(w, "User Deleted")
}

func (h *UserHandler) CreateUserHandle(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		switch h1.ResponseConentType(req) {
		case h1.ResponseHTMX:
			views.UserCreateView().Render(req.Context(), w)
		default:
			views.UserCreatePage().Render(req.Context(), w)
		}
		return
	}
	reqb := new(api.User)
	if err := h1.GetRequestParams(reqb, req); err != nil {
		h.Error(w, http.StatusBadRequest, err)
		return
	}
	_, err := h.UserServiceServer.CreateUser(req.Context(), reqb)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err)
		return
	}
	h.Success(w, "User Created")
}
