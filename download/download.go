package download

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io/ioutil"
	"net/http"
	"uploadImage/utils"
)

func DownloadIdentifier(client *http.Client, urlTemplate, canisterID string, from, to int, sess *session.Session, bucket *string) error {
	uploader := s3manager.NewUploader(sess)

	for i := 0; i < to-from; i++ {
		identifier, err := utils.TokenId2TokenIdentifier(canisterID, uint32(i+from))
		url := fmt.Sprintf(urlTemplate, identifier)
		resp, err := client.Get(url)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			errMsg := fmt.Sprintf("Unable to get URL with status code error: %d %s", resp.StatusCode, resp.Status)
			return errors.New(errMsg)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		buf := bytes.NewReader(data)

		filename := "token-" + identifier
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: bucket,
			Key:    &filename,
			Body:   buf,
		})
		// snippet-end:[s3.go.upload_object.call]
		if err != nil {
			return err
		}
	}
	return nil
}

func DownloadSingle(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		errMsg := fmt.Sprintf("Unable to get URL with status code error: %d %s", resp.StatusCode, resp.Status)
		return nil, errors.New(errMsg)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
