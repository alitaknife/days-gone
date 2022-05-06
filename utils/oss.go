package utils

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/frame/g"
	"os"
)

var Client = (*oss.Client)(nil)

func init() {
	err := (error)(nil)
	endpoint := g.Config().GetString("gitBed.Ht") + g.Config().GetString("gitBed.Endpoint")
	accessKeyId := g.Config().GetString("gitBed.AccessKeyId")
	accessKeySecret := g.Config().GetString("gitBed.AccessKeySecret")

	Client, err = oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
}
