package process

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nfnt/resize"
	"image/gif"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"uploadImage/s3"
	"uploadImage/utils"
)

func DownloadImages(filePath string) {
	var errUrls []utils.ErrUrl
	c := new(http.Client)
	start := time.Now()
	nftInfos, err := utils.GetEXTNFTImageInfos(filePath)
	if err != nil {
		panic(err.Error())
	}

	for name, info := range nftInfos {
		fmt.Println(name, info.Supply, info.ImageUrlTemplate)
		errUrl, err := Download(c, info)
		if err != nil {
			errMsg := fmt.Sprintf("error:%v,url:%v", err, info.ImageUrlTemplate)
			fmt.Println(errMsg)
			continue
		}
		errUrls = append(errUrls, errUrl...)
	}

	if len(errUrls) != 0 {
		errURLData, _ := json.Marshal(errUrls)
		_ = ioutil.WriteFile("err.json", errURLData, 0644)
	}
	fmt.Println("successfully download ext images")
	fmt.Println(time.Now().Sub(start))
}

func DownloadCCCImages(filePath string) {
	start := time.Now()
	var errUrls []utils.ErrUrl
	var cccInfos map[string]utils.CCCNFTImagesInfo
	c := new(http.Client)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}
	if err = json.Unmarshal(data, &cccInfos); err != nil {
		panic(err.Error())
	}

	for name, info := range cccInfos {
		fmt.Println(name, info.Supply, info.Type)
		if info.Type == "ic" {
			errUrl, err := DownloadCCCFromIC(c, info)
			if err != nil {
				errMsg := fmt.Sprintf("can not download images from ic error:%v,url:%v", err, info.ImageUrlTemplate)
				fmt.Println(errMsg)
				continue
			}
			errUrls = append(errUrls, errUrl...)
		} else {
			singleNFTInfos, err := utils.GetCCCNFTImageURL(info.CanisterID, info.FileType, info.ImageUrlTemplate, info.Type)
			if err != nil {
				errMsg := fmt.Sprintf("can not get images url from ic with error:%v,url:%v", err, info.ImageUrlTemplate)
				fmt.Println(errMsg)
				continue
			}
			errUrl, err := DownloadCCCFromIPFS(c, singleNFTInfos)
			if err != nil {
				errMsg := fmt.Sprintf("can not down load images from ipfs error:%v,url:%v", err, info.ImageUrlTemplate)
				fmt.Println(errMsg)
				continue
			}
			errUrls = append(errUrls, errUrl...)
		}
		fmt.Printf("successfully download ccc images, canisterID:%s\n", info.CanisterID)
	}

	if len(errUrls) != 0 {
		errURLData, _ := json.Marshal(errUrls)
		_ = ioutil.WriteFile("err.json", errURLData, 0644)
	}
	fmt.Printf("successfully download ccc images\n")
	fmt.Println(time.Now().Sub(start))
}

func Retry(errPath string) {
	var errUrls []utils.ErrUrl
	data, err := ioutil.ReadFile(errPath)
	if err != nil {
		panic(err)
	}
	_ = json.Unmarshal(data, &errUrls)
	newErrUrls, _ := ReRequestUrl(&http.Client{}, errUrls)
	_data, _ := json.Marshal(&newErrUrls)
	ioutil.WriteFile(errPath, _data, 0644)
}

func Upload(bucket string, directory string) error {
	start := time.Now()
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		panic(err)
	}

	err = s3.UploadDirectory(sess, &bucket, &directory)
	if err != nil {
		errMsg := fmt.Sprintf("can not upload directory:%s,error:%v", directory, err)
		fmt.Println(errMsg)
	}
	msg := fmt.Sprintf("upload directory %s successfully", directory)
	fmt.Println(msg)
	fmt.Println(time.Now().Sub(start))
	return nil
}

//func CompassGifs(dir string) {
//	var paths []string
//	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
//		if !info.IsDir() {
//			paths = append(paths, path)
//		}
//		return nil
//	})
//
//	for _, path := range paths {
//		if err := compassGif(path); err != nil {
//			panic(err)
//		}
//	}
//
//}

func compassGif(path string) error {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	image, err := gif.Decode(file)
	if err != nil {
		panic(err)
	}
	m := resize.Resize(400, 400, image, resize.Bilinear)
	out, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	return gif.Encode(out, m, nil)
}
