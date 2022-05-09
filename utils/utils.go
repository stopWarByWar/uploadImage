package utils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/mix-labs/IC-Go/utils/principal"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func TokenId2TokenIdentifier(canisterID string, tokenID uint32) (string, error) {
	var tokenIdentifier string
	canisterBlob, err := principal.Decode(canisterID)
	if err != nil {
		return "", err
	}
	idBlob := make([]byte, 4)
	binary.BigEndian.PutUint32(idBlob, tokenID)

	data := []uint8{10, 116, 105, 100}
	data = append(data, canisterBlob...)
	data = append(data, idBlob...)
	tokenIdentifier = principal.New(data).Encode()
	return tokenIdentifier, err
}

func GetNFTImageInfos(path string) (map[string]NFTImageInfo, error) {
	var nftInfos = make(map[string]NFTImageInfo)
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	rows, err := f.GetRows("Sheet1")

	for i := 1; i < len(rows); i++ {
		if len(rows[i]) < EntrepotNFTInfoRowLen {
			errMsg := fmt.Sprintf("the amount of element of nft in entrepot should be %d", EntrepotNFTInfoRowLen)
			fmt.Println(len(rows[i]))
			fmt.Println(rows[i])
			return nil, errors.New(errMsg)
		}
		//if rows[i][0] == "" {
		//	fmt.Println(rows[i][1])
		//}
		if rows[i][3] == "N" {
			continue
		}
		supply, _ := strconv.Atoi(rows[i][6])
		nftInfos[rows[i][1]] = NFTImageInfo{
			CanisterID:       rows[i][0],
			Name:             rows[i][1],
			ImageUrlTemplate: rows[i][2],
			Standard:         rows[i][5],
			Supply:           uint32(supply),
			FileType:         rows[i][7],
		}
	}
	return nftInfos, nil
}
