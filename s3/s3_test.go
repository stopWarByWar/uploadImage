package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"testing"
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
	bucket := "ic-nft-thumbnails"
	file := "s3.go"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	err = PutFile(sess, &bucket, &file, []byte{})
	if err != nil {
		t.Log("Could not create bucket "+bucket, "error:", err)
	}
}

func TestGetObjects(t *testing.T) {
	var bucket string
	bucket = "ic-nft-thumbnails"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	result, err := GetObjects(sess, &bucket)
	if err != nil {
		t.Log("Could not create bucket "+bucket, "error:", err)
	}

	for _, o := range result.Contents {
		t.Log(o.String())
	}
}

func TestDeleteItem(t *testing.T) {
	bucket := "ic-nft-thumbnails"
	file := "s3.go"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	err = DeleteItem(sess, &bucket, &file)
	if err != nil {
		t.Log("Could not create bucket "+bucket, "error:", err)
	}
}
