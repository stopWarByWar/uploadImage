package download

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"uploadImage/utils"
)

const (
	goroutineAmount = 15
)

const (
	Identifier = 0
	ID         = 1
)

func Download(c *http.Client, urlTemplate, canisterID string, supply int, sess *session.Session, bucket *string, standard uint8, fileType string) error {
	path := fmt.Sprintf("./images/%s", canisterID)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(err)
	}
	errUrlChan := make(chan string)
	stopChan := make(chan struct{})
	var errUrls []string
	var goroutinesAmount uint32
	var finishedGoroutinesAmount uint32
	elementPerGoroutines := supply / goroutineAmount
	if supply < goroutineAmount {
		goroutinesAmount++
		go func() {
			DownloadSingleRoutine(errUrlChan, c, urlTemplate, canisterID, 0, supply, sess, bucket, standard, fileType)
			stopChan <- struct{}{}
		}()
	} else {
		for i := 0; i < goroutineAmount; i++ {
			goroutinesAmount++
			var _from int
			var _to int
			if i == goroutineAmount-1 {
				_from = i * elementPerGoroutines
				_to = supply
			} else {
				_from = i * elementPerGoroutines
				_to = (i + 1) * elementPerGoroutines
			}
			infoMSg := fmt.Sprintf("start canisterID %s from %d to %d", canisterID, _from, _to)
			fmt.Println(infoMSg)
			go func(_from, _to int) {
				DownloadSingleRoutine(errUrlChan, c, urlTemplate, canisterID, _from, _to, sess, bucket, standard, fileType)
				stopChan <- struct{}{}
			}(_from, _to)
		}
	}

	for {
		select {
		case errUrl := <-errUrlChan:
			if strings.Contains(errUrl, "fp7fo-2aaaa-aaaaf-qaeea-cai") {
				continue
			}
			errUrls = append(errUrls, errUrl)
		case <-stopChan:
			finishedGoroutinesAmount++
			if finishedGoroutinesAmount == goroutinesAmount {
				data, _ := json.Marshal(&errUrls)
				if len(errUrls) == 0 {
					return nil
				}
				filePath := fmt.Sprintf("errUrl/%s.json", canisterID)
				_ = ioutil.WriteFile(filePath, data, 0644)
				return nil
			}
		}
	}
}

func DownloadSingleRoutine(errUrlChan chan string, client *http.Client, urlTemplate, canisterID string, from, to int, sess *session.Session, bucket *string, standard uint8, fileType string) {
	//uploader := s3manager.NewUploader(sess)
	for i := 0; i < to-from; i++ {
		var url string
		var tokenID uint32
		//fmt.Println(i + from)
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
		if err != nil || resp.StatusCode != 200 {
			errMsg := fmt.Sprintf("Unable to get URL with status code error: %d %s,canisterID:%v,tokenID:%v,error:%v,url:%v", resp.StatusCode, resp.Status, canisterID, i+from, err, url)
			fmt.Println(errMsg)
			errUrlChan <- url
			continue
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errMsg := fmt.Sprintf("Unable read body with error: %v,canisterID:%v,tokenID:%v", err, canisterID, i+from)
			fmt.Println(errMsg)
			errUrlChan <- url
			continue
		}
		resp.Body.Close()
		//buf := bytes.NewReader(data)

		filename := fmt.Sprintf("./images/%s/%d.%s", canisterID, tokenID, fileType)
		fmt.Println(filename)
		//_, err = uploader.Upload(&s3manager.UploadInput{
		//	Bucket: bucket,
		//	Key:    &filename,
		//	Body:   buf,
		//})
		err = ioutil.WriteFile(filename, data, 0766)
		// snippet-end:[s3.go.upload_object.call]
		if err != nil {
			errMsg := fmt.Sprintf("Unable to upload image to s3 with error: %v,canisterID:%v,tokenID:%v", err, canisterID, i+from)
			fmt.Println(errMsg)
			errUrlChan <- url
			continue
		}
	}
	infoMSg := fmt.Sprintf("finish %s from %d to %d", canisterID, from, to)
	fmt.Println(infoMSg)
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
