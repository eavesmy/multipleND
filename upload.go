package mnd

import (
	"encoding/json"
	"fmt"
	"github.com/eavesmy/golang-lib"
	"os"
)

func upload(ctx *Ctx) string {

	fi, err := os.Stat(ctx.getValue("u").(string))
	if err != nil {
		panic(err)
	}
	size := fi.Size()
	if size > _maxfilesize {
		panic(error_oversize)
	}

	ispart := false
	if size > _maxbuffer {
		ispart = true
		ctx.setValue("needSplit", ispart)
	}

	ctx.setValue("filesize", size)

	token := getToken(ctx)

	ctx.setValue("token", token)

	userinfo(ctx)
	go storageinfo(ctx)
	addSend(ctx)
	getUpId(ctx)

	if !ispart {
		url := psurl(ctx)

	}

	return ""
}

func userinfo(ctx *Ctx) {
	ret := ctx.post("https://www.wenshushu.cn/ap/user/userinfo", nil)
	fmt.Println(string(ret))
}

func storageinfo(ctx *Ctx) {
	ret := ctx.post("https://www.wenshushu.cn/ap/user/storage", nil)

	j := map[string]string{}
	json.Unmarshal(ret, &j)

	fmt.Println(j)
}

func getToken(ctx *Ctx) string {
	body := map[string]string{"dev_info": "{}"}
	b, _ := json.Marshal(body)
	ret := ctx.post("https://www.wenshushu.cn/ap/login/anonymous", b)

	r := &Ret{}

	if err := json.Unmarshal(ret, &r); err != nil {
		panic(err)
	}

	return r.Data["token"]
}

func addSend(ctx *Ctx) {
	body := &m_addSend{
		Expire: 2,
		Recvs: []string{
			"social",
			"public",
		},
		FileSize:  ctx.getValue("filesize").(int64),
		FileCount: 1,
	}
	b, _ := json.Marshal(body)
	ret := ctx.post("https://www.wenshushu.cn/ap/task/addsend", b)

	j := &Ret{}
	json.Unmarshal(ret, &j)

	ctx.setValue("bid", j.Data["bid"])
	ctx.setValue("social_token", j.Data["social_token"])
	ctx.setValue("tid", j.Data["tid"])
	ctx.setValue("ufileid", j.Data["ufileid"])
}

func getUpId(ctx *Ctx) {
	body := &m_getUpId{
		Preid:  ctx.getValue("ufileid").(string),
		Boxid:  ctx.getValue("bid").(string),
		Linkid: ctx.getValue("tid").(string),
		Utype:  "sendcopy",
		Length: ctx.getValue("filesize").(int64),
		Count:  1,
	}
	b, _ := json.Marshal(body)
	ret := ctx.post("https://www.wenshushu.cn/ap/uploadv2/getupid", b)

	j := &Ret{}

	json.Unmarshal(ret, &j)

	ctx.setValue("upId", j.Data["upId"])
}

func psurl(ctx *Ctx) string {
	body := &m_upload{
		IsPart: ctx.getValue("needSplit").(bool),
		Fname:  ctx.getValue("filename").(string),
		Fsize:  ctx.getValue("filesize").(int64),
		Upid:   ctx.getValue("upId").(string),
	}
	b, _ := json.Marshal(body)
	ret := ctx.post("https://www.wenshushu.cn/ap/uploadv2/psurl", b)
	return string(ret)
}

func file_put(ctx *Ctx) {
	ctx.put()
}
