package process

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"uploadImage/utils"
)

const (
	goroutineAmount = 20
)

const (
	Identifier = 0
	ID         = 1
	CCC        = 2
)

func Download(c *http.Client, info utils.EXTNFTImageInfo) ([]utils.ErrUrl, error) {
	var standard uint8
	if info.Standard == "identifier" {
		standard = Identifier
	} else {
		standard = ID
	}

	path := fmt.Sprintf("./images/%s", info.CanisterID)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(err)
	}
	errUrlChan := make(chan utils.ErrUrl)
	stopChan := make(chan struct{})
	var errUrls []utils.ErrUrl
	var goroutinesAmount uint32
	var finishedGoroutinesAmount uint32
	elementPerGoroutines := info.Supply / goroutineAmount
	if info.Supply < goroutineAmount {
		goroutinesAmount++
		go func() {
			DownloadSingleRoutine(errUrlChan, c, info.ImageUrlTemplate, info.CanisterID, info.StorageCanister, 0, info.Supply, standard, info.FileType)
			stopChan <- struct{}{}
		}()
	} else {
		for i := 0; i < goroutineAmount; i++ {
			goroutinesAmount++
			var _from int
			var _to int
			if i == goroutineAmount-1 {
				_from = i * elementPerGoroutines
				_to = info.Supply
			} else {
				_from = i * elementPerGoroutines
				_to = (i + 1) * elementPerGoroutines
			}
			infoMSg := fmt.Sprintf("start canisterID %s from %d to %d", info.CanisterID, _from, _to)
			fmt.Println(infoMSg)
			go func(_from, _to int) {
				DownloadSingleRoutine(errUrlChan, c, info.ImageUrlTemplate, info.CanisterID, info.StorageCanister, _from, _to, standard, info.FileType)
				stopChan <- struct{}{}
			}(_from, _to)
		}
	}

	for {
		select {
		case errUrl := <-errUrlChan:
			if strings.Contains(errUrl.Url, "fp7fo-2aaaa-aaaaf-qaeea-cai") {
				continue
			}
			errUrls = append(errUrls, errUrl)
		case <-stopChan:
			finishedGoroutinesAmount++
			if finishedGoroutinesAmount == goroutinesAmount {
				return errUrls, nil
			}
		}
	}
}

func DownloadSingleRoutine(errUrlChan chan utils.ErrUrl, client *http.Client, urlTemplate, canisterID, storageCanisterID string, from, to int, standard uint8, fileType string) {
	infoMSg := fmt.Sprintf("start canisterID %s from %d to %d", canisterID, from, to)
	fmt.Println(infoMSg)
	for i := 0; i < to-from; i++ {
		var url string
		var tokenID uint32
		//fmt.Println(i + from)
		switch standard {
		case Identifier:
			tokenID = uint32(i + from)
			identifier, _ := utils.TokenId2TokenIdentifier(storageCanisterID, tokenID)
			url = fmt.Sprintf(urlTemplate, identifier)
		case ID:
			tokenID = uint32(i + from + 1)
			url = fmt.Sprintf(urlTemplate, tokenID)
		case CCC:
			tokenID = uint32(i + from)
			url = fmt.Sprintf(urlTemplate, tokenID)
		}
		resp, err := client.Get(url)
		if err != nil {
			errMsg := fmt.Sprintf("Unable to get URL with error:%v", err)
			fmt.Println(errMsg)
			errUrlChan <- utils.ErrUrl{Url: url, TokenID: tokenID, Type: fileType, CanisterID: canisterID}
			continue
		}

		if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 404 {
			errMsg := fmt.Sprintf("Unable to get URL with status code error: %d %s,canisterID:%v,tokenID:%v,url:%v", resp.StatusCode, resp.Status, canisterID, tokenID, url)
			fmt.Println(errMsg)
			errUrlChan <- utils.ErrUrl{Url: url, TokenID: tokenID, Type: fileType, CanisterID: canisterID}
			continue
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errMsg := fmt.Sprintf("Unable read body with error: %v,canisterID:%v,tokenID:%v", err, canisterID, tokenID)
			fmt.Println(errMsg)
			errUrlChan <- utils.ErrUrl{Url: url, TokenID: tokenID, Type: fileType, CanisterID: canisterID}
			continue
		}

		filename := fmt.Sprintf("./images/%s/%d.%s", canisterID, tokenID, fileType)
		fmt.Println(filename)
		err = ioutil.WriteFile(filename, data, 0766)
		// snippet-end:[s3.go.upload_object.call]
		if err != nil {
			errMsg := fmt.Sprintf("Unable to upload image to s3 with error: %v,canisterID:%v,tokenID:%v", err, canisterID, tokenID)
			fmt.Println(errMsg)
			errUrlChan <- utils.ErrUrl{Url: url, TokenID: tokenID, Type: fileType, CanisterID: canisterID}
			continue
		}
		_ = resp.Body.Close()
	}
	infoMSg = fmt.Sprintf("finish %s from %d to %d", canisterID, from, to)
	fmt.Println(infoMSg)
}

func DownloadSingle(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 404 {
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

func DownloadCCCFromIC(client *http.Client, info utils.CCCNFTImagesInfo) ([]utils.ErrUrl, error) {
	if len(info.ImageUrlTemplates) == 0 {
		return nil, nil
	}

	path := fmt.Sprintf("./images/%s", info.CanisterID)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(err)
	}
	errUrlChan := make(chan utils.ErrUrl)
	stopChan := make(chan struct{})
	var errUrls []utils.ErrUrl
	for _, urlsTemplate := range info.ImageUrlTemplates {
		var goroutinesAmount uint32
		var finishedGoroutinesAmount uint32
		singeCanisterSupply := urlsTemplate.To - urlsTemplate.From
		elementPerGoroutines := singeCanisterSupply / goroutineAmount
		if singeCanisterSupply < goroutineAmount {
			goroutinesAmount++
			go func() {
				DownloadSingleRoutine(errUrlChan, client, urlsTemplate.ImageUrlTemplate, info.CanisterID, info.StorageCanister, urlsTemplate.From, urlsTemplate.To, CCC, info.FileType)
				stopChan <- struct{}{}
			}()
		} else {
			for i := 0; i < goroutineAmount; i++ {
				goroutinesAmount++
				var _from int
				var _to int
				if i == goroutineAmount-1 {
					_from = i*elementPerGoroutines + urlsTemplate.From
					_to = urlsTemplate.To
				} else {
					_from = i*elementPerGoroutines + urlsTemplate.From
					_to = (i+1)*elementPerGoroutines + urlsTemplate.From
				}
				go func(_from, _to int) {
					DownloadSingleRoutine(errUrlChan, client, urlsTemplate.ImageUrlTemplate, info.CanisterID, info.StorageCanister, _from, _to, CCC, info.FileType)
					stopChan <- struct{}{}
				}(_from, _to)
			}
		}
	Loop:
		for {
			select {
			case errUrl := <-errUrlChan:
				if strings.Contains(errUrl.Url, "fp7fo-2aaaa-aaaaf-qaeea-cai") {
					continue
				}
				errUrls = append(errUrls, errUrl)
			case <-stopChan:
				finishedGoroutinesAmount++
				if finishedGoroutinesAmount == goroutinesAmount {
					break Loop
				}
			}
		}
	}
	return errUrls, nil

}

func DownloadCCCFromIPFS(client *http.Client, infos []utils.CCCNFTInfo) ([]utils.ErrUrl, error) {
	supply := len(infos)
	if supply == 0 {
		return nil, nil
	}
	path := fmt.Sprintf("./images/%s", infos[0].CanisterID)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(err)
	}
	errUrlChan := make(chan utils.ErrUrl)
	stopChan := make(chan struct{})
	var errUrls []utils.ErrUrl
	var goroutinesAmount uint32
	var finishedGoroutinesAmount uint32

	elementPerGoroutines := supply / goroutineAmount
	if supply < goroutineAmount {
		goroutinesAmount++
		go func() {
			DownloadCCCFromIPFSSingleRoutine(errUrlChan, client, infos)
			stopChan <- struct{}{}
		}()
	} else {
		var from int
		var to int
		for i := 0; i < goroutineAmount; i++ {
			goroutinesAmount++

			if i < goroutineAmount-1 {
				from = i * elementPerGoroutines
				to = (i + 1) * elementPerGoroutines
			} else {
				from = i * elementPerGoroutines
				to = supply
			}
			infoMSg := fmt.Sprintf("start canisterID %s from %d to %d", infos[0].CanisterID, from, to)
			fmt.Println(infoMSg)
			go func(infos []utils.CCCNFTInfo) {
				DownloadCCCFromIPFSSingleRoutine(errUrlChan, client, infos)
				stopChan <- struct{}{}
			}(infos[from:to])
		}
	}
Loop:
	for {
		select {
		case errUrl := <-errUrlChan:
			if strings.Contains(errUrl.Url, "fp7fo-2aaaa-aaaaf-qaeea-cai") {
				continue
			}
			errUrls = append(errUrls, errUrl)
			//fmt.Printf("error: %v\n", errUrls)
			//if len(errUrls) >= 1 {
			//	break Loop
			//}
		case <-stopChan:
			finishedGoroutinesAmount++
			if finishedGoroutinesAmount == goroutinesAmount {
				break Loop
			}
		}
	}

	return errUrls, nil
}

func DownloadCCCFromIPFSSingleRoutine(errUrlChan chan utils.ErrUrl, client *http.Client, infos []utils.CCCNFTInfo) {
	for _, info := range infos {
		resp, err := client.Get(info.ImageUrl)
		if err != nil {
			errMsg := fmt.Sprintf("Unable to get URL with error:%v", err)
			fmt.Println(errMsg)
			errUrlChan <- utils.ErrUrl{Url: info.ImageUrl, TokenID: uint32(info.TokenID), Type: info.ImageFileType, CanisterID: info.CanisterID}
			continue
		}

		if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 404 {
			errMsg := fmt.Sprintf("Unable to get URL with status code error: %d %s,canisterID:%v,tokenID:%v,url:%v", resp.StatusCode, resp.Status, info.CanisterID, info.TokenID, info.ImageUrl)
			fmt.Println(errMsg)
			errUrlChan <- utils.ErrUrl{Url: info.ImageUrl, TokenID: uint32(info.TokenID), Type: info.ImageFileType, CanisterID: info.CanisterID}
			continue
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errMsg := fmt.Sprintf("Unable read body with error:%v ,canisterID:%v,tokenID:%v,error:%v,url:%v", err, info.CanisterID, info.TokenID, err, info.ImageUrl)
			fmt.Println(errMsg)
			errUrlChan <- utils.ErrUrl{Url: info.ImageUrl, TokenID: uint32(info.TokenID), Type: info.ImageFileType, CanisterID: info.CanisterID}
			continue
		}
		_ = resp.Body.Close()

		filename := fmt.Sprintf("./images/%s/%d.%s", info.CanisterID, info.TokenID, info.ImageFileType)
		fmt.Println(filename)
		err = ioutil.WriteFile(filename, data, 0766)
		// snippet-end:[s3.go.upload_object.call]
		if err != nil {
			errMsg := fmt.Sprintf("Unable to save image file with error: %v,canisterID:%v,tokenID:%v", err, info.CanisterID, info.TokenID)
			fmt.Println(errMsg)
			errUrlChan <- utils.ErrUrl{Url: info.ImageUrl, TokenID: uint32(info.TokenID), Type: info.ImageFileType, CanisterID: info.CanisterID}
			continue
		}
	}
	//infoMSg := fmt.Sprintf("finish nft %s ", infos[0].CanisterID)
	//fmt.Println(infoMSg)
}

func ReRequestUrl(client *http.Client, urls []utils.ErrUrl) ([]utils.ErrUrl, error) {
	var newErrUrls []utils.ErrUrl
	if len(urls) == 0 {
		return nil, nil
	}
	for _, url := range urls {
		data, err := DownloadSingle(client, url.Url)
		if err != nil || data == nil {
			fmt.Printf("can not download url %s with error: %v", url.Url, err)
			newErrUrls = append(newErrUrls, url)
		}
		filename := fmt.Sprintf("./images/%s/%d.%s", url.CanisterID, url.TokenID, url.Type)
		fmt.Println(filename)
		err = ioutil.WriteFile(filename, data, 0766)
		if err != nil {
			fmt.Printf("can not save images %s with error: %v", url.Url, err)
			newErrUrls = append(newErrUrls, url)
		}
	}

	return newErrUrls, nil
}

//func ReRequestUrlOld(client *http.Client, urls []utils.ErrUrl) error {
//	var data []byte
//	var err error
//	for _, url := range urls {
//		if data, err = DownloadSingle(client, url.Url); err != nil || data == nil {
//			//fmt.Printf("can not download, url:%s, canisterID:%s,tokenID:%s,error: %v\n", url.Url, url.CanisterID, url.TokenID, err)
//			id := url.TokenID
//			identifier, _ := utils.TokenId2TokenIdentifier(url.CanisterID, uint32(id))
//			entrepotUrl := fmt.Sprintf("https://images.entrepot.app/t/%s/%s", url.CanisterID, identifier)
//			if data, err = DownloadSingle(client, entrepotUrl); err != nil {
//				fmt.Printf("can not download, url:%s, canisterID:%s,tokenID:%d,error: %v\n", entrepotUrl, url.CanisterID, url.TokenID, err)
//			}
//		}
//		filename := fmt.Sprintf("./images/%s/%d.%s", url.CanisterID, url.TokenID, url.Type)
//		//fmt.Println(filename)
//		err = ioutil.WriteFile(filename, data, 0766)
//		if err != nil {
//			path := fmt.Sprintf("./images/%s", url.CanisterID)
//			if err := os.MkdirAll(path, os.ModePerm); err != nil {
//				panic(err)
//			}
//			if err = ioutil.WriteFile(filename, data, 0766); err != nil {
//				fmt.Printf("can not download, url:%s, canisterID:%s,tokenID:%d,error: %v\n", url.Url, url.CanisterID, url.TokenID, err)
//			}
//		}
//	}
//	return nil
//}
