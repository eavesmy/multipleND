package mnd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/eavesmy/golang-lib/crypto"
	"io/ioutil"
	"os"
	"sync"
)

var lock = sync.WaitGroup{}

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
	}
	ctx.setValue("needSplit", ispart)
	ctx.setValue("filesize", size)

	token := getToken(ctx)

	ctx.setValue("token", token)

	userinfo(ctx)
	go storageinfo(ctx)
	addSend(ctx)
	getUpId(ctx)

	/*
		if _fast := fast(ctx); _fast != "" {
			return _fast
		}
	*/

	if !ispart {
		url := psurl(ctx)
		ctx.setValue("uploadlocation", url)
		b, _ := ioutil.ReadFile(ctx.getValue("u").(string))
		file_put(ctx, b)
	} else {

		var offset int64 = 0
		_fi, _ := os.Open(ctx.getValue("u").(string))
		reader := bufio.NewReader(_fi)
		index := 1

		for ; offset < ctx.getValue("filesize").(int64); offset += _maxbuffer {
			_buffer := make([]byte, _maxbuffer)
			ctx.setValue("filesize", len(_buffer))
			reader.Read(_buffer)
			index++

			// 最多10个协程，机子跟不上,开多了没用

			lock.Add(1)
			go func() {
				url := psurl(ctx, index)
				ctx.setValue("uploadlocation", url)
				file_put(ctx, _buffer)
				lock.Done()
			}()
		}

		lock.Wait()
	}

	return complete(ctx)
}

func userinfo(ctx *Ctx) {
	ctx.post("https://www.wenshushu.cn/ap/user/userinfo", nil)
}

func storageinfo(ctx *Ctx) {
	ret := ctx.post("https://www.wenshushu.cn/ap/user/storage", nil)

	r := &Ret{}
	json.Unmarshal(ret, &r)
}

func getToken(ctx *Ctx) string {
	body := &m_token{DevInfo: "{}"}
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

	fmt.Println(string(b))

	ret := ctx.post("https://www.wenshushu.cn/ap/uploadv2/getupid", b)

	j := &Ret{}

	json.Unmarshal(ret, &j)

	ctx.setValue("upId", j.Data["upId"])
}

func psurl(ctx *Ctx, nus ...int) string {
	body := &m_upload{
		IsPart: ctx.getValue("needSplit").(bool),
		Fname:  ctx.getValue("filename").(string),
		Fsize:  ctx.getValue("filesize").(int64),
		Upid:   ctx.getValue("upId").(string),
	}

	if len(nus) > 0 {
		body.PartNu = nus[0]
	}

	b, _ := json.Marshal(body)

	r := &Ret{}

	ret := ctx.post("https://www.wenshushu.cn/ap/uploadv2/psurl", b)
	json.Unmarshal(ret, &r)

	return r.Data["url"]
}

func file_put(ctx *Ctx, body []byte) {

	ctx.put(ctx.getValue("uploadlocation").(string), body)
	return
}

func fast(ctx *Ctx) string {

	cm1, cs1 := fileHash(ctx)

	m := &m_fast{
		Hash: map[string]string{
			"cm1": cm1,
			"cs1": cs1,
		},
		Uf: map[string]string{
			"name":  ctx.getValue("filename").(string),
			"boxid": ctx.getValue("bid").(string),
			"preid": ctx.getValue("ufileid").(string),
		},
		UpId: ctx.getValue("upId").(string),
	}

	if !ctx.getValue("needSplit").(bool) {
		m.Hash["cm"] = crypto.Sha1(cm1)
	}

	d, _ := json.Marshal(m)

	for i := 0; i < 2; i++ {
		ret := ctx.post("https://www.wenshushu.cn/ap/uploadv2/fast", d)
		r := &RetFast{}
		json.Unmarshal(ret, &r)

		can_fast := r.Data.Status
		ufile := r.Data.Ufile

		if can_fast != 0 && ufile == nil {
			hash_code := ""
			// 遍历整个文件获取每段的 md5 hash

			var offset int64 = 0
			fi, _ := os.Open(ctx.getValue("u").(string))
			reader := bufio.NewReader(fi)

			for ; offset < ctx.getValue("filesize").(int64); offset += _maxbuffer {
				_buffer := make([]byte, _maxbuffer)
				reader.Read(_buffer)
				hash_code += crypto.Md5Bytes(_buffer)
			}

			m.Hash["cm"] = hash_code

		} else if can_fast != 0 && ufile != nil {
			return copySend(ctx)
		}
	}

	// 拿到签名

	return ""
}

func fileHash(ctx *Ctx, blocks ...[]byte) (cm, cs string) {

	block := []byte{}

	if len(blocks) > 0 {
		block = blocks[0]
	}

	if len(block) == 0 {
		filepath := ctx.getValue("u").(string)
		fi, _ := os.Open(filepath)
		block, _ = ioutil.ReadAll(fi)
	}

	cm = crypto.Md5Bytes(block)
	cs = crypto.Sha1Bytes(block)

	return
}

func copySend(ctx *Ctx) string {
	body := &m_send{
		Bid:     ctx.getValue("bid").(string),
		Tid:     ctx.getValue("tid").(string),
		UfileId: ctx.getValue("ufileid").(string),
	}

	b, _ := json.Marshal(body)

	r := ctx.post("https://www.wenshushu.cn/ap/task/copysend", b)
	ret := &Ret{}
	json.Unmarshal(r, &ret)

	fmt.Println("m", ret.Data["mgr_url"])
	fmt.Println("p", ret.Data["public_url"])
	fmt.Println("----------")

	return ret.Data["public_url"]
}

func complete(ctx *Ctx) string {
	m := &m_complete{
		IsPart: ctx.getValue("needSplit").(bool),
		Fname:  ctx.getValue("filename").(string),
		UpId:   ctx.getValue("upId").(string),
		Location: map[string]string{
			"boxid": ctx.getValue("bid").(string),
			"preid": ctx.getValue("ufileid").(string),
		},
	}
	b, _ := json.Marshal(m)
	ctx.post("https://www.wenshushu.cn/ap/uploadv2/complete", b)

	return copySend(ctx)
}
