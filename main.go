package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/http"
	"time"
	"uploadImage/download"
)

func main() {
	c := new(http.Client)
	start := time.Now()
	var bucket = "pyd-nft"
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	download.Download(c, "https://bzsui-sqaaa-aaaah-qce2a-cai.raw.ic0.app/?type=thumbnail&tokenid=%s", "zzk67-giaaa-aaaaj-qaujq-cai", 10000, sess, &bucket, download.Identifier)
	fmt.Println(time.Now().Sub(start))
}
