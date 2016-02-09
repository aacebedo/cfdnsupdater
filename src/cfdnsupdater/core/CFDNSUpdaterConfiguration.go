package core

type CFDNSUpdaterConfiguration struct {
	Verbose       bool           
	Quiet         bool           
	Syslog        bool           
	DomainConfigs map[string]DomainConfiguration
}