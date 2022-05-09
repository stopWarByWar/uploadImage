package utils

const EntrepotNFTInfoRowLen = 7

type NFTImageInfo struct {
	CanisterID       string `json:"canister_id"`
	Name             string `json:"name"`
	ImageUrlTemplate string `json:"image_url_template"`
	Supply           uint32 `json:"supply"`
	Standard         string `json:"standard"`
	FileType         string `json:"file_type"`
}
