package download

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io/ioutil"
	"net/http"
	"sync"
	"uploadImage/utils"
)

const SingleGoroutineNftAmount = 100

const (
	Identifier = 0
	ID         = 1
)

func Download(c *http.Client, urlTemplate, canisterID string, supply int, sess *session.Session, bucket *string, standard uint8) error {
	var wg sync.WaitGroup
	from := 0
	for from < supply {
		wg.Add(1)
		var _from int
		var _to int
		if from+SingleGoroutineNftAmount < supply {
			_from = from
			_to = from + SingleGoroutineNftAmount
			from = from + SingleGoroutineNftAmount
		} else {
			_from = from
			_to = supply
			from = supply
		}
		go func(_from, _to int) {
			if err := DownloadSingleRoutine(c, urlTemplate, canisterID, _from, _to, sess, bucket, standard); err != nil {
				fmt.Printf("can not download identifier imgage with canisterID:%v,from:%v,to:%v\n", canisterID, _from, _to)
			}
			wg.Done()
		}(_from, _to)
	}

	wg.Wait()
	return nil
}

func DownloadSingleRoutine(client *http.Client, urlTemplate, canisterID string, from, to int, sess *session.Session, bucket *string, standard uint8) error {
	uploader := s3manager.NewUploader(sess)
	for i := 0; i < to-from; i++ {
		var url string
		var tokenID uint32
		fmt.Println(i + from)
		switch standard {
		case Identifier:
			tokenID = uint32(i + from)
			identifier, _ := utils.TokenId2TokenIdentifier(canisterID, tokenID)
			url = fmt.Sprintf(urlTemplate, identifier)
		case ID:
			tokenID = uint32(i + from + 1)
			url = fmt.Sprintf(urlTemplate, tokenID)
		}
		resp, err := client.Get(url)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			errMsg := fmt.Sprintf("Unable to get URL with status code error: %d %s,canisterID:%v,tokenID:%v", resp.StatusCode, resp.Status, canisterID, i+from)
			return errors.New(errMsg)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		buf := bytes.NewReader(data)

		filename := fmt.Sprintf("%s-%d", canisterID, tokenID)
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: bucket,
			Key:    &filename,
			Body:   buf,
		})
		// snippet-end:[s3.go.upload_object.call]
		if err != nil {
			errMsg := fmt.Sprintf("Unable to upload image to s3 with error: %v,canisterID:%v,tokenID:%v", err, canisterID, i+from)
			return errors.New(errMsg)
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
