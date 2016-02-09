package core


type DomainConfiguration struct {
	Email       string               `yaml:"email"`
	ApiKey      string               `yaml:"apikey"`
	Period      int                  `yaml:"period"`
	RecordNames []string             `yaml:"record_names"`
	RecordTypes RecordTypeSlice      `yaml:"record_types"`
}