package utils

import (
	"bytes"
	"fmt"
	"yunshu_Forum/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var ossClient *oss.Client
var ossBucket *oss.Bucket

func OssInit() {
	clientOptions := []oss.ClientOption{}
	clientOptions = append(clientOptions, oss.Region("cn-shenzhen"))
	clientOptions = append(clientOptions, oss.AuthVersion(oss.AuthV4))
	var err error
	ossClient, err = oss.New(config.Oss.Endpoint, config.Oss.AccessKeyId, config.Oss.AccessKeySecret, clientOptions...)
	if err != nil {
		fmt.Println("oss用户实例创建失败:", err)
		return
	}
	fmt.Println(config.Oss.Bucket)
	ossBucket, err = ossClient.Bucket(config.Oss.Bucket)
	if err != nil {
		fmt.Println("oss存储同创建失败: ", err)
		return
	}
}

func UploadFileToOss(objectName string, data string) (string, error) {
	//定义路径
	objectName = "postImage/" + objectName
	err := ossBucket.PutObject(objectName, bytes.NewReader([]byte(data)))
	if err != nil {
		return "", err
	}
	return objectName, nil
}

func DeleteFileFromOss(objectName string) error {
	objectName = "postImage/" + objectName
	err := ossBucket.DeleteObject(objectName)
	if err != nil {
		return err
	}
	return nil
}

func GetUrl(objectName string) string {
	objectName = "postImage/" + objectName
	signedURL, err := ossBucket.SignURL(objectName, oss.HTTPGet, 60)
	if err != nil {
		return ""
	}
	return signedURL
}
