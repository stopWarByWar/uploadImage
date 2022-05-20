package s3

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/http"
	"sync"
	"testing"
	"uploadImage/utils"
)

func TestListBucket(t *testing.T) {
	listBuckets()
}

func TestCreateBucket(t *testing.T) {
	var bucket string
	bucket = "pyd-nft"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	err = MakeBucket(sess, &bucket)
	if err != nil {
		t.Log("Could not create bucket "+bucket, "error:", err)
	}
}

func TestRemoveBucket(t *testing.T) {
	var bucket string
	bucket = "pyd-nft"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	err = RemoveBucket(sess, &bucket)
	if err != nil {
		t.Log("Could not create bucket "+bucket, "error:", err)
	}
}

func TestPutFile(t *testing.T) {
	bucket := "pyd-nft"
	file := "ccc_nft_archive.json"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	err = PutFile(sess, &bucket, &file, []byte{})
	if err != nil {
		t.Log("Could not delete bucket "+bucket, "error:", err)
	}
}

func TestGetObjects(t *testing.T) {
	var bucket string
	bucket = "pyd-nft"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	result, err := GetObjects(sess, &bucket)
	if err != nil {
		t.Log("Could not get bucket "+bucket, "error:", err)
	}

	for _, o := range result.Contents {
		t.Log(o.String())
	}
}

func TestDeleteItem(t *testing.T) {
	bucket := "pyd-nft"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	result, err := GetObjects(sess, &bucket)

	for _, item := range result.Contents {
		fmt.Printf("item: %v\n", *item.Key)
		err = DeleteItem(sess, &bucket, item.Key)
		if err != nil {
			t.Log("Could not create bucket "+bucket, "error:", err)
		}
	}
}
func TestDeleteItems(t *testing.T) {
	bucket := "pyd-nft"
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	t.Log(DeleteItems(sess, &bucket))
}

func TestPutDirectory(t *testing.T) {
	bucket := "pyd-nft"
	directory := "../download/images/ctt6t-faaaa-aaaah-qcpbq-cai"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	err = UploadDirectory(sess, &bucket, &directory)
	if err != nil {
		fmt.Println("Got an error uploading directory " + directory + " to bucket " + bucket)
		return
	}
}

func TestSetBucket(t *testing.T) {
	bucket := "pyd-nft"
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err := SetBucketPublic(sess, &bucket); err != nil {
		t.Log(err)
	}

	t.Log("Bucket " + bucket + " is now public")
}

func TestGetFile(t *testing.T) {
	filename := "ccc_nft_archive.json"
	bucket := "pyd-nft"
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	data, err := GetFile(sess, &bucket, &filename)
	if err != nil {
		t.Fatalf(err.Error())
	}
	var cccInfos map[string]utils.CCCNFTImagesInfo
	t.Log(json.Unmarshal(data, &cccInfos))
	t.Log(len(cccInfos))

}

func TestDownloadObject(t *testing.T) {
	filename := "ccc_nft_archive.json"
	bucket := "pyd-nft"
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	t.Log(DownloadObject(sess, &filename, &bucket))
}

func TestName(t *testing.T) {
	for {
		var wg sync.WaitGroup
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				http.Get("https://linda.mirror.xyz/df649d61efb92c910464a4e74ae213c4cab150b9cbcc4b7fb6090fc77881a95d")
			}()
			wg.Wait()
		}
	}
}
