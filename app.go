package mnd

import (
	"io"
	"path/filepath"
)

const _maxfilesize = 5368709120
const _maxbuffer = 2097152

// 获取下载直链
func GetDownloadUrl() string {
	return ""
}

// 调用下载
func Download(url string) io.Reader {

	ctx := newCtx()
	defer ctx.cancel()

	return nil
}

// 上传单个文件, 5G 以内
func Upload(filePath string) string {

	ctx := newCtx()
	defer ctx.cancel()

	ctx.setValue("u", filePath)

	name := filepath.Base(filePath)
	ctx.setValue("filename", name)

	return upload(ctx)
}
