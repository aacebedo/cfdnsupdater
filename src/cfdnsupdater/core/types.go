//This file is part of CFDNSUpdater.
//
//    CFDNSUpdater is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    CFDNSUpdater is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with CFDNSUpdater.  If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"errors"
	"fmt"
)

type RecordType int

type RecordTypeSlice []RecordType

const (
	A     RecordType = 0 << iota
	AAAA  RecordType = 1
	CNAME RecordType = 2
	MX    RecordType = 3
	LOC   RecordType = 4
	SRV   RecordType = 5
	SPF   RecordType = 6
	TXT   RecordType = 7
	NS    RecordType = 8
)

func (self *RecordType) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
  var val string
  if err = unmarshal(&val); err == nil {
      *self, err = FromString(val)
  }
	return
}

func (self RecordType) String() (res string, err error) {
	conversionArray := []string{"A", "AAAA", "CNAME", "MX", "LOC", "SRV", "SPF", "TXT", "NS"}
	if int(self) < len(conversionArray) {
		res = conversionArray[self]
	} else {
		err = errors.New("Invalid type of record")
	}
	return
}

func (self RecordTypeSlice) Contains(valueToSearch RecordType) (res bool) {
	res = false
	for _, v := range self {
		if v == valueToSearch {
			res = true
			break
		}
	}
	return
}

func FromString(valToConvert string) (res RecordType, err error) {
	conversionArray := []string{"A", "AAAA", "CNAME", "MX", "LOC", "SRV", "SPF", "TXT", "NS"}
	for p, v := range conversionArray {
		if v == valToConvert {
			res = RecordType(p)
			return
		}
	}
	err = errors.New(fmt.Sprintf("'%s' is not a valid record type", valToConvert))
	return
}


type RecordRequestResult struct {
	Result []struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Content string `json:"content"`
		Id      string `json:"id"`
	}
}

type ZoneRequestResult struct {
	Result []struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	}
	Result_info struct {
		Total_count int `json:"total_count"`
	}
}
