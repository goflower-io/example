package main

import (
	"fmt"
	"net/url"

	"github.com/go-playground/form"

	"github.com/goflower-io/example/mysql/api"
)

func main() {
	a := api.UpdateUserReq{}
	decode := form.NewDecoder()
	decode.Decode(&a,
		url.Values{
			"User.Id":       []string{"1"},
			"UpdateMask[0]": []string{"id"},
			"UpdateMask[1]": []string{"name"},
		},
	)
	fmt.Printf("%+v", &a)

	b := api.UserFilter{}
	decode.Decode(&b, url.Values{
		"Field": []string{"1"},
		"Op":    []string{"="},
		"Value": []string{"1"},
	})
	fmt.Printf("%+v", &b)
}
