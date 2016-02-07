package core


type DomainConfiguration struct {
	Email       string               `json:"email"`
	ApiKey      string               `json:"apikey"`
	Period      int                  `json:"period"`
	Domain      string               `json:"domain"`
	RecordNames []string             `json:"record_names"`
	RecordTypes RecordTypeSlice `json:"record_types,RecordType"`
}