package process

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nfnt/resize"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"image/gif"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"uploadImage/s3"
	"uploadImage/utils"
)

func GetCCCImagesURL(filePath string, dns string, autoMigrate bool) {
	_db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: _db}), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
	})
	if err != nil {
		panic(err)
	}
	if autoMigrate {
		err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&utils.NFTUrl{})
		if err != nil {
			panic(err)
		}
	}

	start := time.Now()
	var cccInfos map[string]utils.CCCNFTImagesInfo
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
			err := GetCCCUrlsFromIC(db, info)
			if err != nil {
				errMsg := fmt.Sprintf("can not download images from ic error:%v,url:%v", err, info.ImageUrlTemplate)
				fmt.Println(errMsg)
				panic(err)
			}
		} else {
			singleNFTInfos, err := utils.GetCCCNFTImageURL(info.CanisterID, info.FileType, info.ImageUrlTemplate, info.Type)
			if err != nil {
				errMsg := fmt.Sprintf("can not get images url from ic with error:%v,url:%v", err, info.ImageUrlTemplate)
				fmt.Println(errMsg)
				continue
			}
			var urls []utils.NFTUrl
			for i, nftInfo := range singleNFTInfos {
				//fmt.Println(nftInfo.CanisterID, nftInfo.TokenID, nftInfo.ImageUrl)
				urls = append(urls, utils.NFTUrl{CanisterID: nftInfo.CanisterID, TokenID: uint32(nftInfo.TokenID), ImageUrl: nftInfo.ImageUrl, VideoUrl: nftInfo.VideoUrl})
				if len(urls) >= 500 || i == len(singleNFTInfos)-1 {
					fmt.Println("set urls into db with length", len(urls))
					if err := db.Save(&urls).Error; err != nil {
						fmt.Printf("can not save images urls canister:%s,tokenid:%d,url:%s,err:%v\n", nftInfo.CanisterID, nftInfo.TokenID, nftInfo.ImageUrl, err)
						panic(err)
					}
					urls = []utils.NFTUrl{}
				}
			}
		}
		fmt.Printf("successfully download ccc images, canisterID:%s\n", info.CanisterID)
	}

	fmt.Printf("successfully download ccc images\n")
	fmt.Println(time.Now().Sub(start))

}

func GetDip721ImageUrl(filePath string, dns string, autoMigrate bool) {
	_db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: _db}), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
	})
	if err != nil {
		panic(err)
	}
	if autoMigrate {
		err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&utils.NFTUrl{})
		if err != nil {
			panic(err)
		}
	}
	var infos map[string]utils.BasicNFTInfo
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}
	if err = json.Unmarshal(data, &infos); err != nil {
		panic(err.Error())
	}
	for _, info := range infos {
		if err = GetDIPUrlsFromIC(db, info); err != nil {
			panic(err.Error())
		}
		fmt.Println("successfully get url canister", info.CanisterID)
	}
}

func GetYumiImagesUrl(filePath string, dns string, autoMigrate bool) {
	_db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: _db}), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
	})
	if err != nil {
		panic(err)
	}
	if autoMigrate {
		err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&utils.NFTUrl{})
		if err != nil {
			panic(err)
		}
	}

	var yumiInfos map[string]utils.BasicNFTInfo
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}
	if err = json.Unmarshal(data, &yumiInfos); err != nil {
		panic(err.Error())
	}
	for _, yumiNFT := range yumiInfos {
		if err = GetYumiUrlsFromIC(db, yumiNFT.CanisterID); err != nil {
			fmt.Printf("can not update yumi %s urls with error:%v\n", yumiNFT.Name, err)
			panic(err)
		}
		fmt.Println("successfully get url canister:", yumiNFT.Name)
	}
}

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
