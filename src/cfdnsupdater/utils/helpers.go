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

package utils

func StringInSlice(valueToCheck string, array []string) (res bool) {
	res = false
	for _, v := range array {
		if v == valueToCheck {
			res = true
		}
	}
	return
}
