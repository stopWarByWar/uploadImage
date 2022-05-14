package process

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"io/ioutil"
	"net/http"
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

	errURLData, _ := json.Marshal(errUrls)
	_ = ioutil.WriteFile("err.json", errURLData, 0644)
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

	errURLData, _ := json.Marshal(errUrls)
	_ = ioutil.WriteFile("err.json", errURLData, 0644)
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

//func RetryFromLogs(filePath, logPath string) {
//	nftInfos, err := utils.GetEXTNFTImageInfos(filePath)
//	if err != nil {
//		panic(err)
//	}
//	c := new(http.Client)
//	var errUrls []utils.ErrUrl
//	fi, err := os.Open(logPath)
//	if err != nil {
//		fmt.Printf("Error: %s\n", err)
//		return
//	}
//	defer fi.Close()
//
//	br := bufio.NewReader(fi)
//	for {
//		a, _, c := br.ReadLine()
//		if c == io.EOF {
//			break
//		}
//
//		result := strings.Split(string(a), ",")
//		if len(result) < 5 {
//			fmt.Println(result)
//			continue
//		}
//		canisterID := strings.Split(result[1], ":")[1]
//		tokenID := strings.Split(result[2], ":")[1]
//		url := strings.Split(result[4], ":")[1] + ":" + strings.Split(result[4], ":")[2]
//		errCode := strings.Split(strings.Split(result[0], ":")[1], " ")[1]
//
//		if strings.Contains(errCode, "20") || strings.Contains(errCode, "404") {
//			continue
//		}
//
//		if canisterID == "po6n2-uiaaa-aaaaj-qaiua-cai" {
//			url = "https://po6n2-uiaaa-aaaaj-qaiua-cai.raw.ic0.app/?tokenid="
//			id, _ := strconv.Atoi(tokenID)
//			identifier, _ := utils.TokenId2TokenIdentifier("po6n2-uiaaa-aaaaj-qaiua-cai", uint32(id))
//			url = url + identifier
//		}
//		//fmt.Println(url, tokenID, canisterID, nftInfos[canisterID].FileType, errCode)
//		id, _ := strconv.Atoi(tokenID)
//		errUrls = append(errUrls, utils.ErrUrl{Url: url, TokenID: uint32(id), CanisterID: canisterID, Type: nftInfos[canisterID].FileType})
//	}
//	fmt.Println(len(errUrls))
//	ReRequestUrlOld(c, errUrls)
//}

func Upload(bucket string, directory string) error {
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
	return nil
}
