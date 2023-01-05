package utils

import (
	"github.com/mix-labs/IC-Go/utils/principal"
	"time"
)

const EntrepotNFTInfoRowLen = 8

type EXTNFTImageInfo struct {
	CanisterID       string `json:"canister_id"`
	Name             string `json:"name"`
	ImageUrlTemplate string `json:"image_url_template"`
	Supply           int    `json:"supply"`
	Standard         string `json:"standard"`
	FileType         string `json:"file_type"`
	StorageCanister  string `json:"storage_canister"`
}

type CCCNFTImagesInfo struct {
	CanisterID        string                `json:"canister_id"`
	Name              string                `json:"name"`
	ImageUrlTemplates []CCCImageUrlTemplate `json:"image_url_templates"`
	Supply            int                   `json:"supply"`
	FileType          string                `json:"file_type"`
	Type              string                `json:"type"`
	ImageUrlTemplate  string                `json:"image_url_template"`
	StorageCanister   string                `json:"storage_canister"`
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
	CanisterID      string `json:"canister_id"`
	TokenID         uint64 `json:"token_id"`
	ImageUrl        string `json:"image_url"`
	VideoUrl        string `json:"video_url"`
	ImageFileType   string `json:"file_type"`
	VideoFileType   string `json:"video_file_type"`
	StorageCanister string `json:"storage_canister"`
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

type NFTPhotoLink struct {
	TokenIndex  uint64 `ic:"0"`
	NFTLinkInfo string `ic:"1"`
}

type StringOp struct {
	Some string `ic:"some"`
	None uint8  `ic:"none"`
}

type IPFS3ImageLink struct {
	ImageLink StringOp `ic:"0"`
	VideoLink StringOp `ic:"1"`
}

type BasicNFTInfo struct {
	CanisterID string `json:"canister_id"`
	Name       string `json:"name"`
}

type GenericValue struct {
	TextContent string `ic:"TextContent"`
	EnumIndex   string `ic:"EnumIndex"`
}

type record struct {
	Text         string       `ic:"0"`
	GenericValue GenericValue `ic:"1"`
}
type TokenMetadata struct {
	Properties []record `ic:"properties"`
}

type ManualReply_2 struct {
	Ok        TokenMetadata `ic:"Ok"`
	EnumIndex string        `ic:"EnumIndex"`
}

type BlobOpt struct {
	Some []byte `ic:"some"`
	None uint8  `ic:"none"`
}
type NonFungibleType struct {
	MetaData BlobOpt `ic:"metadata"`
}
type Metadata__1 struct {
	NonFungible NonFungibleType `ic:"nonfungible"`
	Index       string          `ic:"EnumIndex"`
}
type MetaData struct {
	TokenIndex uint32      `ic:"0"`
	MetaData   Metadata__1 `ic:"1"`
}

type YumUrl struct {
	Thumb    string `json:"thumb"`
	Url      string `json:"url"`
	MimeType string `json:"mime_type"`
}

type NFTUrl struct {
	CanisterID string `gorm:"column:canister_id;type:char(27) not null;primaryKey:token_instance,priority 1" json:"canister_id"`
	TokenID    uint32 `gorm:"column:token_id;primaryKey:token_instance,priority 2" json:"token_id"`
	ImageUrl   string `gorm:"column:image_url" json:"image_url"`
	VideoUrl   string `gorm:"column:video_url" json:"video_url"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (NFTUrl) TableName() string {
	return "nft_urls"
}

type BasicEntrepotImageInfo struct {
	Name       string `json:"name"`
	CanisterId string `json:"canister_id"`
	Supply     int    `json:"supply"`
	Type       string `json:"type"`
}

//
//type VecNat8Opt struct {
//	Some []uint8 `ic:"some"`
//	None uint8   `ic:"none"`
//}
//type Fungible struct {
//	Decimals uint8      `ic:"decimals"`
//	Metadata VecNat8Opt `ic:"metadata"`
//	Name     string     `ic:"name"`
//	Symbol   string     `ic:"symbol"`
//}
//
//type Nonfungible struct {
//	Metadata VecNat8Opt `ic:"metadata"`
//}
//type Metadata struct {
//	Fungible    Fungible    `ic:"fungible"`
//	Nonfungible Nonfungible `ic:"nonfungible"`
//	Index       string      `ic:"EnumIndex"`
//}
//
//type result_2 struct {
//	Ok    Metadata `ic:"Ok"`
//	Index string   `ic:"EnumIndex"`
//}

type GenericValue_1 struct {
	Nat64Content uint64  `ic:"Nat64Content"`
	Nat32Content uint32  `ic:"Nat32Content"`
	BoolContent  bool    `ic:"BoolContent"`
	Nat8Content  uint8   `ic:"Nat8Content"`
	Int64Content int64   `ic:"Int64Content"`
	IntContent   int     `ic:"IntContent"`
	NatContent   uint    `ic:"NatContent"`
	Nat16Content uint16  `ic:"Nat16Content"`
	Int32Content int32   `ic:"Int32Content"`
	Int8Content  int8    `ic:"Int8Content"`
	FloatContent float64 `ic:"FloatContent"`
	Int16Content int16   `ic:"Int16Content"`
	BlobContent  []uint8 `ic:"BlobContent"`
	Principal    principal.Principal
	TextContent  string `ic:"TextContent"`
	Index        string `ic:"EnumIndex"`
}

type TokenPropertie struct {
	Name  string         `ic:"0"`
	Value GenericValue_1 `ic:"1"`
}

type TokenMetaData struct {
	Properties []TokenPropertie `ic:"properties"`
}

type NftError struct{}
type ManualReply_3 struct {
	Ok    TokenMetaData `ic:"Ok"`
	Err   NftError      `ic:"Err"`
	Index string        `ic:"EnumIndex"`
}
