package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/goflower-io/golib/net/app"
	"github.com/goflower-io/xsql"

	"github.com/goflower-io/example/sqlite/api"
	"github.com/goflower-io/example/sqlite/crud"
	"github.com/goflower-io/example/sqlite/service"
	"github.com/goflower-io/example/sqlite/views"
)

var (
	db  *crud.Client
	ctx = context.Background()
)

func main() {
	var err error
	db, err = crud.NewClient(&xsql.Config{
		DSN:          "./sqlite.db",
		ReadDSN:      []string{"./sqlite.db"},
		Active:       20,
		Idle:         20,
		IdleTimeout:  time.Hour * 24,
		QueryTimeout: time.Second * 10,
		ExecTimeout:  time.Second * 10,
	}, true)
	if err != nil {
		panic(err)
	}

	s := &service.UserServiceImpl{Client: db}
	hs := service.NewUserHandler(s)
	// Add middleware: recovery (panic → 500) + structured request logging.
	hs.Middlewares = append(hs.Middlewares, app.RecoveryMiddle, app.LogMidddle)
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
	// db.User.Update().SetName("xxx").Where(user.IdOp.EQ(4005)).Save(ctx)

	// list, err := db.User.Find().Select().Where(user.AgeOp.EQ(11)).All(ctx)
	// b, _ := json.Marshal(list)
	// fmt.Println(string(b), err)

	// db.User.Delete().Where(user.IdOp.EQ(a.Id)).Exec(ctx)
}
