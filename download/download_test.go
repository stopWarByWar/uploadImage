package download

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
	"uploadImage/utils"
)

func TestDownload(t *testing.T) {
	c := new(http.Client)
	start := time.Now()
	var bucket = "pyd-nft"
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	Download(c, "https://bzsui-sqaaa-aaaah-qce2a-cai.raw.ic0.app/?type=thumbnail&tokenid=%s", "zzk67-giaaa-aaaaj-qaujq-cai", 10000, sess, &bucket, Identifier)
	t.Log(time.Now().Sub(start))
}

func TestDownloadIdentifier(t *testing.T) {
	c := new(http.Client)
	start := time.Now()
	var wg sync.WaitGroup
	var bucket = "pyd-nft"
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	for i := 0; i < 10; i++ {
		from := i * 10
		to := i*10 + 10
		go func(from, to int) {
			wg.Add(1)
			defer wg.Done()
			t.Log("from: ", from, "to: ", to)
			t.Log(DownloadSingleRoutine(c, "https://bzsui-sqaaa-aaaah-qce2a-cai.raw.ic0.app/?type=thumbnail&tokenid=%s", "zzk67-giaaa-aaaaj-qaujq-cai", from, to, sess, &bucket, Identifier))
		}(from, to)
	}
	wg.Wait()
	t.Log(time.Now().Sub(start))
}

func TestDownload1(t *testing.T) {
	c := new(http.Client)
	start := time.Now()
	var bucket = "pyd-nft"
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	t.Log(DownloadSingleRoutine(c, "https://zzk67-giaaa-aaaaj-qaujq-cai.raw.ic0.app/?type=thumbnail&tokenid=%s", "zzk67-giaaa-aaaaj-qaujq-cai", 0, 10, sess, &bucket, Identifier))
	t.Log(time.Now().Sub(start))
}

func TestDownloadAll(t *testing.T) {
	start := time.Now()
	c := new(http.Client)
	nftInfos, err := utils.GetNFTImageInfos("./nft_imge.xlsx")
	if err != nil {
		t.Fatalf(err.Error())
	}

	for name, info := range nftInfos {
		var url string
		if info.Standard == "identifier" {
			identifier, _ := utils.TokenId2TokenIdentifier(info.CanisterID, 12)
			url = fmt.Sprintf(info.ImageUrlTemplate, identifier)
		} else {
			url = fmt.Sprintf(info.ImageUrlTemplate, 1)
		}
		t.Log(name + ":" + url)
		data, err := DownloadSingle(c, url)
		if err != nil {
			errMsg := fmt.Sprintf("error:%v,url:%v", err, url)
			t.Log(errMsg)
			continue
		}
		path := fmt.Sprintf("images/%s.%s", name, info.FileType)
		ioutil.WriteFile(path, data, 0644)
	}
	t.Log(time.Now().Sub(start))
}

func TestDownloadSingle(t *testing.T) {
	c := new(http.Client)
	var url = "https://d3ttm-qaaaa-aaaai-qam4a-cai.raw.ic0.app/?tokenId=4"

	data, err := DownloadSingle(c, url)
	if err != nil {
		errMsg := fmt.Sprintf("error:%v,url:%v", err, url)
		t.Log(errMsg)
	}
	ioutil.WriteFile("image.jpg", data, 0644)

}
