package utils

const EntrepotNFTInfoRowLen = 7

type EXTNFTImageInfo struct {
	CanisterID       string `json:"canister_id"`
	Name             string `json:"name"`
	ImageUrlTemplate string `json:"image_url_template"`
	Supply           int    `json:"supply"`
	Standard         string `json:"standard"`
	FileType         string `json:"file_type"`
}

type CCCNFTImageInfo struct {
	CanisterID        string                `json:"canister_id"`
	Name              string                `json:"name"`
	ImageUrlTemplates []CCCImageUrlTemplate `json:"image_url_templates"`
	Supply            int                   `json:"supply"`
	Standard          string                `json:"standard"`
	FileType          string                `json:"file_type"`
}

type CCCImageUrlTemplate struct {
	ImageUrlTemplate string `json:"image_url_template"`
	From             uint32 `json:"from"`
	To               uint32 `json:"to"`
}
