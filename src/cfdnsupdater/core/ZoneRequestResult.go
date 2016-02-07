package core

type ZoneRequestResult struct {
	Result []struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	}
	Result_info struct {
		Total_count int `json:"total_count"`
	}
}



