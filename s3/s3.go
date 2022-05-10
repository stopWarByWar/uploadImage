package s3

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"path/filepath"
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

func DeleteItems(sess *session.Session, bucket *string) error {
	// snippet-start:[s3.go.delete_objects.call]
	svc := s3.New(sess)

	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: bucket,
	})

	err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter)
	// snippet-end:[s3.go.delete_objects.call]
	if err != nil {
		return err
	}

	return nil
}

func SetBucketPublic(sess *session.Session, bucket *string) error {
	// snippet-start:[s3.go.make_bucket_public.call]
	svc := s3.New(sess)

	params := &s3.PutBucketAclInput{
		Bucket: bucket,
		ACL:    aws.String("public-read"),
	}

	_, err := svc.PutBucketAcl(params)
	// snippet-end:[s3.go.make_bucket_public.call]
	if err != nil {
		return err
	}

	return nil
}

type DirectoryIterator struct {
	filePaths []string
	bucket    string
	next      struct {
		path string
		f    *os.File
	}
	err error
}

// const exitError = 1

// UploadDirectory uploads the files in a directory to a bucket
// Inputs:
//     sess is the current session, which provides configuration for the SDK's service clients
//     bucket is the name of the bucket
//     path is the path to the directory to upload
// Output:
//     If success, nil
//     Otherwise, an error from the call to UploadWithIterator
func UploadDirectory(sess *session.Session, bucket *string, path *string) error {
	di := NewDirectoryIterator(bucket, path)
	uploader := s3manager.NewUploader(sess)

	err := uploader.UploadWithIterator(aws.BackgroundContext(), di)
	if err != nil {
		return err
	}

	return nil
}

// NewDirectoryIterator builds a new DirectoryIterator
func NewDirectoryIterator(bucket *string, dir *string) s3manager.BatchUploadIterator {
	var paths []string
	filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	return &DirectoryIterator{
		filePaths: paths,
		bucket:    *bucket,
	}
}

// Next returns whether next file exists
func (di *DirectoryIterator) Next() bool {
	if len(di.filePaths) == 0 {
		di.next.f = nil
		return false
	}

	f, err := os.Open(di.filePaths[0])
	di.err = err
	di.next.f = f
	di.next.path = di.filePaths[0]
	di.filePaths = di.filePaths[1:]

	return true && di.Err() == nil
}

// Err returns error of DirectoryIterator
func (di *DirectoryIterator) Err() error {
	return di.err
}

// UploadObject uploads a file
func (di *DirectoryIterator) UploadObject() s3manager.BatchUploadObject {
	f := di.next.f
	return s3manager.BatchUploadObject{
		Object: &s3manager.UploadInput{
			Bucket: &di.bucket,
			Key:    &di.next.path,
			Body:   f,
		},
		After: func() error {
			return f.Close()
		},
	}
}
