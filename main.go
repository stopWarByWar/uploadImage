package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/http"
	"time"
	"uploadImage/download"
	"uploadImage/utils"
)

func main() {
	c := new(http.Client)
	start := time.Now()
	var bucket = "pyd-nft"
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	nftInfos, err := utils.GetEXTNFTImageInfos("./download/nft_image1.xlsx")
	if err != nil {
		panic(err.Error())
	}

	for name, info := range nftInfos {
		var standard uint8
		if info.Standard == "identifier" {
			standard = download.Identifier
		} else {
			standard = download.ID
		}
		fmt.Println(name, info.Supply, info.ImageUrlTemplate)
		err = download.Download(c, info.ImageUrlTemplate, info.CanisterID, info.Supply, sess, &bucket, standard, info.FileType)
		if err != nil {
			errMsg := fmt.Sprintf("error:%v,url:%v", err, info.ImageUrlTemplate)
			fmt.Println(errMsg)
			continue
		}
	}
	fmt.Println(time.Now().Sub(start))
}
