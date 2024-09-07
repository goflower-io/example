package main

import (
	"context"
	"fmt"
	"time"

	"github.com/happycrud/xsql"
	"github.com/happycrud/xsql/postgres"

	"github.com/happycrud/example/postgres/crud/user"
)

var (
	db  *xsql.DB
	ctx = context.Background()
)

func main() {
	var err error
	db, err = postgres.NewDB(&xsql.Config{
		DSN:          "postgres://postgres:123456@localhost:5432/postgres",
		ReadDSN:      []string{"postgres://postgres:123456@localhost:5432/postgres"},
		Active:       20,
		Idle:         20,
		IdleTimeout:  time.Hour * 24,
		QueryTimeout: time.Second * 10,
		ExecTimeout:  time.Second * 10,
	})
	if err != nil {
		panic(err)
	}
	debugdb := xsql.Debug(db)
	a := &user.User{
		Id:    3,
		Name:  "sdfs",
		Age:   []int32{11, 12, 22},
		Ctime: time.Now(),
		Mtime: time.Now(),
	}
	b := &user.User{
		Id:    4,
		Name:  "a",
		Age:   []int32{22, 123},
		Ctime: time.Now(),
		Mtime: time.Now(),
	}
	_, err = user.Create(debugdb).SetUser(a, b).Upsert(ctx)
	fmt.Println(a, b, err)

	list, err := user.Find(debugdb).All(ctx)
	for _, v := range list {
		fmt.Printf("%+v ,err:%v\n", v, err)
	}
	user.Update(debugdb).SetAge([]int32{444, 444, 44}).Where(user.IdOp.EQ(1)).Save(ctx)
}
