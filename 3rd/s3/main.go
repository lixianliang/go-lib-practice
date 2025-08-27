package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/url"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	ctx := context.Background()

	// 加载 AWS 配置（从环境变量或 ~/.aws/config）
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"), // 指定 S3 存储桶区域
	)
	if err != nil {
		panic(err)
	}

	// 原始请求参数
	bucket := "your-bucket"
	key := "path/to/object.txt"
	expires := time.Now().Add(3600 * time.Second) // URL 有效期

	// 构建原始请求
	req, err := http.NewRequest("GET", "", nil)
	req.URL = &url.URL{
		Scheme:   "https",
		Host:     fmt.Sprintf("%s.s3.%s.amazonaws.com", bucket, cfg.Region),
		Path:     key,
		RawQuery: "imageMogr2/thumbnail/10x&versionId=xyz&response-content-disposition=attachment", // 原始查询参数
	}
	if err != nil {
		panic(err)
	}

	// 创建 V4 签名器
	credentials, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		panic(err)
	}
	signer := v4.NewSigner()

	// 计算 Body 的 SHA256（GET 请求 Body 为空）
	hasher := sha256.New()
	hasher.Write([]byte{}) // 空 Body
	sha256Hash := fmt.Sprintf("%x", hasher.Sum(nil))

	// 生成预签名 URL
	presignedURL, _, err := signer.PresignHTTP(ctx, credentials, req, sha256Hash, "s3", cfg.Region, expires)
	if err != nil {
		panic(err)
	}

	fmt.Println("预签名 URL:", presignedURL)
}
