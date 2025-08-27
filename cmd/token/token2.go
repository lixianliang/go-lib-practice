package main

import (
	"context"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/downloader"
	"log"
)

func main() {
	// 源站域名下载方式
	accessKey := "ecrtylqJBfmYCczK41X_J0FfE6VmMcfIrcaadUeO"
	secretKey := "_SWnfLgHJqN3aVsrSo9h4uzls2dqfXb5758qTW5_"
	mac := credentials.NewCredentials(accessKey, secretKey)
	urlsProvider := downloader.SignURLsProvider(downloader.NewDefaultSrcURLsProvider(mac.AccessKey, nil), downloader.NewCredentialsSigner(mac), nil)

	// localFile := "/Users/jemy/Documents/github.png"
	localPath := "./lxl-download.m4a"
	bucket := "lixianliang"
	key := "voices/2025/07/09/191504-1930157190415458306-lxl.m4a"
	downloadManager := downloader.NewDownloadManager(&downloader.DownloadManagerOptions{})
	downloaded, err := downloadManager.DownloadToFile(context.Background(), key, localPath, &downloader.ObjectOptions{
		GenerateOptions:      downloader.GenerateOptions{BucketName: bucket},
		DownloadURLsProvider: urlsProvider,
	})
	if err != nil {
		log.Fatalln("failed download err: ", err)
	}
	log.Println("read n", downloaded)
}
