package main

import (
	"context"
	"fmt"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

func main() {
	// 配置信息
	accessKey := "xxx"
	secretKey := "xxx"
	bucket := "lixianliang"
	filePath := "/Users/lixianliang/Documents/audio/lxl.m4a"

	// 1. 生成上传凭证
	upToken := generateUploadToken(accessKey, secretKey, bucket)
	upToken = "ecrtylqJBfmYCczK41X_J0FfE6VmMcfIrcaadUeO:vYB5-4fl5AKYicaJfa6RLCuuxrM=:eyJkZWFkbGluZSI6MTc1MjA2Njg3NiwiZm9yY2VTYXZlS2V5Ijp0cnVlLCJyZXR1cm5Cb2R5Ijoie1wia2V5XCI6XCIkKGtleSlcIixcImhhc2hcIjpcIiQoZXRhZylcIixcImZzaXplXCI6JChmc2l6ZSksXCJ0eXBlXCI6JChtaW1lVHlwZSl9Iiwic2F2ZUtleSI6InZvaWNlcy8ke3llYXJ9LyR7bW9ufS8ke2RheX0vJHtob3VyfSR7bWlufSR7c2VjfS0xOTMwMTU3MTkwNDE1NDU4MzA2LSR7Zm5hbWV9Iiwic2NvcGUiOiJsaXhpYW5saWFuZyJ9"
	println("token ", upToken)

	// 2. 配置上传参数
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan, // 华东区域
		UseHTTPS:      true,
		UseCdnDomains: true,
	}

	// 3. 创建表单上传对象
	formUploader := storage.NewFormUploader(&cfg)

	// 4. 执行文件上传
	//ret := storage.PutRet{}
	type XPutRet struct {
		Hash  string `json:"hash"`
		Key   string `json:"key"`
		Fsize int64  `json:"fsize"`
		Type  string `json:"type"`
	}
	ret := XPutRet{}
	err := formUploader.PutFile(
		context.Background(),
		&ret,
		upToken,
		"lxl.m4a",
		filePath,
		nil,
	)

	if err != nil {
		fmt.Println("文件上传失败:", err)
		return
	}

	fmt.Printf("文件上传成功!\nKey: %s\nHash: %s %d %s\n", ret.Key, ret.Hash, ret.Fsize, ret.Type)
}

func generateUploadToken(accessKey, secretKey, bucket string) string {
	mac := auth.New(accessKey, secretKey)
	policy := storage.PutPolicy{
		Scope:        bucket,
		Expires:      3600,
		ReturnBody:   `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"type":$(mimeType)}`,
		ForceSaveKey: true,
		SaveKey:      "voices/${year}/${mon}/${day}/${fname}",
	}
	return policy.UploadToken(mac)
}
