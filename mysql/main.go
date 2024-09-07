package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/happycrud/xsql"

	"github.com/happycrud/example/mysql/api"
	"github.com/happycrud/example/mysql/crud"
	"github.com/happycrud/example/mysql/service"
	"github.com/happycrud/example/mysql/views"
)

var (
	db  *crud.Client
	ctx = context.Background()
)

func main() {
	var err error
	db, err = crud.NewClient(&xsql.Config{
		DSN:          "root:123456@tcp(127.0.0.1:3306)/test?parseTime=true&loc=Local",
		ReadDSN:      []string{"root:123456@tcp(127.0.0.1:3306)/test?parseTime=true&loc=Local"},
		Active:       20,
		Idle:         20,
		IdleTimeout:  time.Hour * 24,
		QueryTimeout: time.Second * 10,
		ExecTimeout:  time.Second * 10,
	})
	if err != nil {
		panic(err)
	}
	s := &service.UserServiceImpl{Client: db}
	hs := service.NewUserHandler(s)
	mux := http.NewServeMux()
	hs.AddPath(func(method, path string, hf http.HandlerFunc) {
		fmt.Println(method + " " + path)
		mux.HandleFunc(method+" "+path, hf)
	})
	mux.HandleFunc("GET /index", func(w http.ResponseWriter, r *http.Request) {
		a := &api.User{
			Id:   1,
			Name: "ddd",
			Age:  100,
		}
		views.UserUpdateView(a).Render(r.Context(), w)
	})
	http.ListenAndServe("0.0.0.0:8088", mux)
}
