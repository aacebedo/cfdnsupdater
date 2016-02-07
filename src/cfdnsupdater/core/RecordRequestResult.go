package core

type RecordRequestResult struct {
	Result []struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Content string `json:"content"`
		Id      string `json:"id"`
	}
}