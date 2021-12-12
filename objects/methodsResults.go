package objects

/*This object represents the response of getUpdates method*/
type UpdateResult struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type FailureResult struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

type SendMethodsResult struct {
	Ok     bool    `json:"ok"`
	Result Message `json:"result"`
}
