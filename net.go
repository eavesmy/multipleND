package mnd

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
)

type Ctx struct {
	context context.Context
	client  *http.Client
	cancel  context.CancelFunc
}

func newCtx() *Ctx {
	ctx := &Ctx{}
	ctx.client = &http.Client{}
	ctx.context = context.Background()
	ctx.context, ctx.cancel = context.WithCancel(ctx.context)
	return ctx
}

func (ctx *Ctx) setValue(k, v interface{}) {
	ctx.context = context.WithValue(ctx.context, k, v)
}

func (ctx *Ctx) getValue(k interface{}) interface{} {
	return ctx.context.Value(k)
}

func (ctx *Ctx) post(url string, body []byte) []byte {

	var req *http.Request
	var err error

	if body != nil && len(body) > 0 {
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest("POST", url, nil)
	}

	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}

	if ctx.getValue("token") != nil {
		req.Header.Add("x-token", ctx.getValue("token").(string))
	}

	res, err := ctx.client.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	b, _ := ioutil.ReadAll(res.Body)
	return b
}

func get() {

}

func (ctx *Ctx) put(url string, body []byte) []byte {

	var req *http.Request
	var err error

	if body != nil && len(body) > 0 {
		req, err = http.NewRequest("PUT", url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest("PUT", url, nil)
	}

	if err != nil {
		panic(err)
	}

	if ctx.getValue("token") != nil {
		req.Header.Add("x-token", ctx.getValue("token").(string))
	}

	res, err := ctx.client.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	b, _ := ioutil.ReadAll(res.Body)
	return b
}
