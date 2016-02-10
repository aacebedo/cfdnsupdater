package configuration

import (
    "cfdnsupdater/core"
  )


type DomainConfiguration struct {
	Email       string               `yaml:"email"`
	ApiKey      string               `yaml:"apikey"`
	Period      int                  `yaml:"period"`
	RecordNames []string             `yaml:"record_names"`
	RecordTypes core.RecordTypeSlice      `yaml:"record_types"`
}

type CFDNSUpdaterConfiguration struct {
	Verbose       bool           
	Quiet         bool           
	Syslog        bool           
	DomainConfigs map[string]DomainConfiguration
}