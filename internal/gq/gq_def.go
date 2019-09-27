package gq

type IGQ interface {
	// Parse 解析文档
	Parse(ParseReq) ParseResp
}

// Response defines the API response format.
type ParseResp struct {
	// Data is the carrier for returning data.
	Data interface{} `json:"data"`
	// Error records the errors in this request.
	Error error `json:"error"`
	// TimeCost recorded the time wastage of the request.
	TimeCost int64 `json:"timecost"`
}

type ParseReq struct {
	doc, expr string
}
