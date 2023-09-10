package objects

/*UpdateResult represents the response of getUpdates method*/
// type UpdateResult struct {
// 	Ok     bool      `json:"ok"`
// 	Result []*Update `json:"result"`
// }

// FailureResult represents a failure response that has "ok : false" field.
type FailureResult struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

// Result is generic struct conataining results on success
type Result[k any] struct {
	Ok     bool `json:"ok"`
	Result k    `json:"result"`
}
