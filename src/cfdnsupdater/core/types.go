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

type RecordRequestResult struct {
	Result []struct {
		Name        string `json:"name"`
		Type string `json:"type"`
		Content     string `json:"content"`
		Id          string `json:"id"`
	}
}

func stringInSlice(valueToCheck string, 
  array []string) bool {
	for _, v := range array {
		if v == valueToCheck {
			return true
		}
	}
	return false
}