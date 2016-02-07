package core

type CFDNSUpdaterConfiguration struct {
	Verbose       bool           `json:"verbose"`
	Quiet         bool           `json:"quiet"`
	Syslog        bool           `json:"syslog"`
	DomainConfigs []DomainConfiguration `json:"domain_configs"`
}