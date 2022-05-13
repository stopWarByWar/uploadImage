package utils

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
