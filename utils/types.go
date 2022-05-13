package utils

import (
	"github.com/mix-labs/IC-Go/utils/principal"
)

const EntrepotNFTInfoRowLen = 8

type EXTNFTImageInfo struct {
	CanisterID       string `json:"canister_id"`
	Name             string `json:"name"`
	ImageUrlTemplate string `json:"image_url_template"`
	Supply           int    `json:"supply"`
	Standard         string `json:"standard"`
	FileType         string `json:"file_type"`
}

type CCCNFTImagesInfo struct {
	CanisterID        string                `json:"canister_id"`
	Name              string                `json:"name"`
	ImageUrlTemplates []CCCImageUrlTemplate `json:"image_url_templates"`
	Supply            int                   `json:"supply"`
	FileType          string                `json:"file_type"`
	Type              string                `json:"type"`
	ImageUrlTemplate  string                `json:"image_url_template"`
}

type CCCImageUrlTemplate struct {
	ImageUrlTemplate string `json:"image_url_template"`
	From             int    `json:"from"`
	To               int    `json:"to"`
}

type ErrUrl struct {
	CanisterID string `json:"canister_id"`
	Url        string `json:"url"`
	TokenID    uint32 `json:"token_id"`
	Type       string `json:"type"`
}

type CCCNFTInfo struct {
	CanisterID    string `json:"canister_id"`
	TokenID       uint64 `json:"token_id"`
	ImageUrl      string `json:"image_url"`
	VideoUrl      string `json:"video_url"`
	ImageFileType string `json:"file_type"`
	VideoFileType string `json:"video_file_type"`
}

type textOp struct {
	Some string `ic:"some"`
	None uint8  `ic:"none"`
}

type nftStoreInfo struct {
	Index     uint64 `ic:"index"`
	PhotoLink textOp `ic:"photoLink"`
	VideoLink textOp `ic:"videoLink"`
}

type registry struct {
	Principal principal.Principal `ic:"0"`
	Tokens    []nftStoreInfo      `ic:"1"`
}

type TokenImageUrl struct {
	TokenIndex  uint64      `ic:"0"`
	NFTLinkInfo NFTLinkInfo `ic:"1"`
}

type NFTLinkInfo struct {
	ID        uint64 `ic:"id"`
	PhotoLink string `ic:"photoLink"`
	VideoLink string `ic:"videoLink"`
}
