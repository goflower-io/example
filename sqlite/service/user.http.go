package service

import (
	"net/http"

	"github.com/goflower-io/golib/net/app"

	"github.com/goflower-io/example/sqlite/api"
	"github.com/goflower-io/example/sqlite/views"
)

// UserHandler wraps UserServiceServer with HTTP handlers.
// It supports middleware chains and flexible response types (JSON, HTML, HTMX)
// via github.com/goflower-io/golib/net/app.
//
// Quick start:
//
//	h := NewUserHandler(svc)
//	h.Middlewares = append(h.Middlewares, app.RecoveryMiddle, app.LogMidddle)
//	h.AddPath(func(method, path string, hf http.HandlerFunc) {
//	    mux.HandleFunc(method+" "+path, hf)
//	})
//
// Response format is negotiated per-request via Accept / HX-Request headers:
//
//	Accept: application/json  → JSON envelope  (app.ResponseJSON)
//	HX-Request: true          → HTMX fragment  (app.ResponseHTMX)
//	(default)                 → full HTML page (app.ResponseHTML)
type UserHandler struct {
	api.UserServiceServer

	// Middlewares are applied to every route in registration order
	// (first element is outermost / executes first).
	// Compatible with app.RecoveryMiddle and app.LogMidddle.
	Middlewares []func(http.HandlerFunc) http.HandlerFunc

	// OnError writes an error response. Override to customise the format.
	// Default: JSON error envelope for API callers, plain text otherwise.
	OnError func(w http.ResponseWriter, r *http.Request, statusCode int, err error)

	// OnSuccess writes a success response for mutation operations (Create/Update/Delete).
	// Default: JSON {"message":"..."} for API callers, plain text for HTML callers.
	OnSuccess func(w http.ResponseWriter, r *http.Request, msg string)
}

// NewUserHandler creates a UserHandler with sensible defaults.
func NewUserHandler(s api.UserServiceServer) *UserHandler {
	h := &UserHandler{
		UserServiceServer: s,
	}
	h.OnError = func(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
		wr := app.NewWriter(w, r)
		switch app.ResponseConentType(r) {
		case app.ResponseJSON:
			wr.JSONErr(app.NewError(statusCode, err.Error()))
		default:
			wr.Text(statusCode, err.Error())
		}
	}
	h.OnSuccess = func(w http.ResponseWriter, r *http.Request, msg string) {
		wr := app.NewWriter(w, r)
		switch app.ResponseConentType(r) {
		case app.ResponseJSON:
			wr.JSONOk(map[string]string{"message": msg})
		default:
			wr.Text(http.StatusOK, msg)
		}
	}
	return h
}

// applyMiddlewares wraps fn with all configured middlewares.
// The first element of Middlewares is the outermost wrapper.
func (h *UserHandler) applyMiddlewares(fn http.HandlerFunc) http.HandlerFunc {
	for i := len(h.Middlewares) - 1; i >= 0; i-- {
		fn = h.Middlewares[i](fn)
	}
	return fn
}

// AddPath registers all CRUD routes using the provided addPathFn.
// Configured middlewares are applied to each route.
//
// Example with stdlib mux:
//
//	h.AddPath(func(method, path string, hf http.HandlerFunc) {
//	    mux.HandleFunc(method+" "+path, hf)
//	})
func (h *UserHandler) AddPath(addPathFn func(method, path string, hf http.HandlerFunc)) {
	addPathFn(http.MethodGet, api.UserService_CreateUser_FullMethodName, h.applyMiddlewares(h.CreateUserHandle))
	addPathFn(http.MethodPut, api.UserService_CreateUser_FullMethodName, h.applyMiddlewares(h.CreateUserHandle))

	addPathFn(http.MethodGet, api.UserService_GetUser_FullMethodName, h.applyMiddlewares(h.GetUserHandle))
	addPathFn(http.MethodDelete, api.UserService_DeleteUser_FullMethodName, h.applyMiddlewares(h.DeleteUserHandle))

	addPathFn(http.MethodGet, api.UserService_UpdateUser_FullMethodName, h.applyMiddlewares(h.UpdateUserHandle))
	addPathFn(http.MethodPost, api.UserService_UpdateUser_FullMethodName, h.applyMiddlewares(h.UpdateUserHandle))

	addPathFn(http.MethodGet, api.UserService_ListUsers_FullMethodName, h.applyMiddlewares(h.ListUsersHandle))
	addPathFn(http.MethodGet, api.UserService_ListUsersMore_FullMethodName, h.applyMiddlewares(h.ListUsersMoreHandle))
}

// ─────────────────────────────────────────────
// ListMore (cursor-based infinite scroll)
// ─────────────────────────────────────────────

func (h *UserHandler) ListUsersMoreHandle(w http.ResponseWriter, req *http.Request) {
	reqb := new(api.ListUsersMoreReq)
	if err := app.GetRequestParams(reqb, req); err != nil {
		h.OnError(w, req, http.StatusBadRequest, err)
		return
	}
	// First page: ensure a default cursor so the service can establish ordering.
	if reqb.GetPageToken() == "" && reqb.GetCursor() == nil {
		reqb.Cursor = &api.UserCursor{}
	}
	if reqb.GetPageSize() <= 0 {
		reqb.PageSize = 20
	}
	resp, err := h.UserServiceServer.ListUsersMore(req.Context(), reqb)
	if err != nil {
		h.OnError(w, req, http.StatusInternalServerError, err)
		return
	}
	wr := app.NewWriter(w, req)
	switch app.ResponseConentType(req) {
	case app.ResponseJSON:
		wr.JSONOk(resp)
	case app.ResponseHTMX:
		wr.TemplFragment(views.UserListMoreView(resp, req.URL.Path, reqb.GetPageSize()))
	default:
		wr.TemplOk(views.UserListMorePage(resp, req.URL.Path, reqb.GetPageSize()))
	}
}

// ─────────────────────────────────────────────
// List
// ─────────────────────────────────────────────

func (h *UserHandler) ListUsersHandle(w http.ResponseWriter, req *http.Request) {
	reqb := new(api.ListUsersReq)
	if err := app.GetRequestParams(reqb, req); err != nil {
		h.OnError(w, req, http.StatusBadRequest, err)
		return
	}
	resp, err := h.UserServiceServer.ListUsers(req.Context(), reqb)
	if err != nil {
		h.OnError(w, req, http.StatusInternalServerError, err)
		return
	}
	wr := app.NewWriter(w, req)
	switch app.ResponseConentType(req) {
	case app.ResponseJSON:
		wr.JSONOk(resp)
	case app.ResponseHTMX:
		wr.TemplFragment(views.UserListView(resp, req.URL.Path, req.Form))
	default:
		wr.TemplOk(views.UserListPage(resp, req.URL.Path, req.Form))
	}
}

// ─────────────────────────────────────────────
// Get
// ─────────────────────────────────────────────

func (h *UserHandler) GetUserHandle(w http.ResponseWriter, req *http.Request) {
	reqb := new(api.UserId)
	if err := app.GetRequestParams(reqb, req); err != nil {
		h.OnError(w, req, http.StatusBadRequest, err)
		return
	}
	resp, err := h.UserServiceServer.GetUser(req.Context(), reqb)
	if err != nil {
		h.OnError(w, req, http.StatusInternalServerError, err)
		return
	}
	wr := app.NewWriter(w, req)
	switch app.ResponseConentType(req) {
	case app.ResponseJSON:
		wr.JSONOk(resp)
	case app.ResponseHTMX:
		wr.TemplFragment(views.UserDetailView(resp))
	default:
		wr.TemplOk(views.UserDetailPage(resp))
	}
}

// ─────────────────────────────────────────────
// Create  (GET → show form, PUT → create record)
// ─────────────────────────────────────────────

func (h *UserHandler) CreateUserHandle(w http.ResponseWriter, req *http.Request) {
	wr := app.NewWriter(w, req)
	if req.Method == http.MethodGet {
		switch app.ResponseConentType(req) {
		case app.ResponseHTMX:
			wr.TemplFragment(views.UserCreateView())
		default:
			wr.TemplOk(views.UserCreatePage())
		}
		return
	}
	reqb := new(api.User)
	if err := app.GetRequestParams(reqb, req); err != nil {
		h.OnError(w, req, http.StatusBadRequest, err)
		return
	}
	resp, err := h.UserServiceServer.CreateUser(req.Context(), reqb)
	if err != nil {
		h.OnError(w, req, http.StatusInternalServerError, err)
		return
	}
	switch app.ResponseConentType(req) {
	case app.ResponseJSON:
		wr.JSON(http.StatusCreated, resp)
	default:
		h.OnSuccess(w, req, "User Created")
	}
}

// ─────────────────────────────────────────────
// Update  (GET → show form, POST → update record)
// ─────────────────────────────────────────────

func (h *UserHandler) UpdateUserHandle(w http.ResponseWriter, req *http.Request) {
	wr := app.NewWriter(w, req)
	if req.Method == http.MethodGet {
		reqb := new(api.UserId)
		if err := app.GetRequestParams(reqb, req); err != nil {
			h.OnError(w, req, http.StatusBadRequest, err)
			return
		}
		resp, err := h.UserServiceServer.GetUser(req.Context(), reqb)
		if err != nil {
			h.OnError(w, req, http.StatusInternalServerError, err)
			return
		}
		switch app.ResponseConentType(req) {
		case app.ResponseHTMX:
			wr.TemplFragment(views.UserUpdateView(resp))
		default:
			wr.TemplOk(views.UserUpdatePage(resp))
		}
		return
	}
	reqb := new(api.UpdateUserReq)
	if err := app.GetRequestParams(reqb, req); err != nil {
		h.OnError(w, req, http.StatusBadRequest, err)
		return
	}
	resp, err := h.UserServiceServer.UpdateUser(req.Context(), reqb)
	if err != nil {
		h.OnError(w, req, http.StatusInternalServerError, err)
		return
	}
	switch app.ResponseConentType(req) {
	case app.ResponseJSON:
		wr.JSONOk(resp)
	default:
		h.OnSuccess(w, req, "User Updated")
	}
}

// ─────────────────────────────────────────────
// Delete
// ─────────────────────────────────────────────

func (h *UserHandler) DeleteUserHandle(w http.ResponseWriter, req *http.Request) {
	reqb := new(api.UserId)
	if err := app.GetRequestParams(reqb, req); err != nil {
		h.OnError(w, req, http.StatusBadRequest, err)
		return
	}
	_, err := h.UserServiceServer.DeleteUser(req.Context(), reqb)
	if err != nil {
		h.OnError(w, req, http.StatusInternalServerError, err)
		return
	}
	h.OnSuccess(w, req, "User Deleted")
}
