package download

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestDownload(t *testing.T) {
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
			t.Log(DownloadIdentifier(c, "https://zzk67-giaaa-aaaaj-qaujq-cai.raw.ic0.app/?type=thumbnail&tokenid=%s", "zzk67-giaaa-aaaaj-qaujq-cai", from, to, sess, &bucket))
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

	t.Log(DownloadIdentifier(c, "https://zzk67-giaaa-aaaaj-qaujq-cai.raw.ic0.app/?type=thumbnail&tokenid=%s", "zzk67-giaaa-aaaaj-qaujq-cai", 0, 10, sess, &bucket))
	t.Log(time.Now().Sub(start))
}

func TestDownloadSingle(t *testing.T) {
	c := new(http.Client)
	data, err := DownloadSingle(c, "https://zzk67-giaaa-aaaaj-qaujq-cai.raw.ic0.app/?type=thumbnail&tokenid=amudd-kikor-uwiaa-aaaaa-cmafc-maqca-aaaam-q")
	if err != nil {
		t.Log(err)
	}
	ioutil.WriteFile("result", data, 0644)
}
