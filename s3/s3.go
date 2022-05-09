package s3

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

func listBuckets() {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create S3 service client
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func MakeBucket(sess *session.Session, bucket *string) error {
	// snippet-start:[s3.go.create_bucket.call]
	svc := s3.New(sess)

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: bucket,
	})
	// snippet-end:[s3.go.create_bucket.call]
	if err != nil {
		return err
	}

	// snippet-start:[s3.go.create_bucket.wait]
	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: bucket,
	})
	// snippet-end:[s3.go.create_bucket.wait]
	if err != nil {
		return err
	}

	return nil
}

func RemoveBucket(sess *session.Session, bucket *string) error {
	// snippet-start:[s3.go.delete_bucket.call]
	svc := s3.New(sess)

	_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: bucket,
	})
	// snippet-end:[s3.go.delete_bucket.call]
	if err != nil {
		return err
	}

	// snippet-start:[s3.go.delete_bucket.wait]
	err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: bucket,
	})
	// snippet-end:[s3.go.delete_bucket.wait]
	if err != nil {
		return err
	}

	return nil
}

func PutFile(sess *session.Session, bucket *string, filename *string, data []byte) error {
	// snippet-start:[s3.go.upload_object.open]
	// snippet-start:[s3.go.upload_object.call]
	buf := bytes.NewReader(data)
	uploader := s3manager.NewUploader(sess)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: bucket,
		Key:    filename,
		Body:   buf,
	})
	// snippet-end:[s3.go.upload_object.call]
	if err != nil {
		return err
	}
	return nil
}

func GetObjects(sess *session.Session, bucket *string) (*s3.ListObjectsV2Output, error) {
	// snippet-start:[s3.go.list_objects.call]
	svc := s3.New(sess)

	// Get the list of items
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: bucket})
	// snippet-end:[s3.go.list_objects.call]
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func DeleteItem(sess *session.Session, bucket *string, item *string) error {
	// snippet-start:[s3.go.delete_object.call]
	svc := s3.New(sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: bucket,
		Key:    item,
	})
	// snippet-end:[s3.go.delete_object.call]
	if err != nil {
		return err
	}

	// snippet-start:[s3.go.delete_object.wait]
	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: bucket,
		Key:    item,
	})
	// snippet-end:[s3.go.delete_object.wait]
	if err != nil {
		return err
	}

	return nil
}
