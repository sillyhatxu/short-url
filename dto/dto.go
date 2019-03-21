package dto

type Health struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

type ErrorResp struct {
	Msg string `json:"msg"`
}

type ShortReq struct {
	LongURL string `json:"long-url"`
}

type ShortResp struct {
	ShortURL string `json:"short-url"`
}

type ExpandReq struct {
	ShortURL string `json:"short-url"`
}

type ExpandResp struct {
	LongURL string `json:"long-url"`
}

type CompressDTO struct {
	ShortURL  string `json:"short-url"`
	LongURL   string `json:"long-url"`
	ClickTime int64  `json:"click-time"`
}
