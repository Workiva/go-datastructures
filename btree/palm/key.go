/*
Copyright 2014 Workiva, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package palm

import "github.com/Workiva/go-datastructures/common"

func reverseKeys(cmps common.Comparators) common.Comparators {
	reversed := make(common.Comparators, len(cmps))
	for i := len(cmps) - 1; i >= 0; i-- {
		reversed[len(cmps)-1-i] = cmps[i]
	}

	return reversed
}

func chunkKeys(keys common.Comparators, numParts int64) []common.Comparators {
	parts := make([]common.Comparators, numParts)
	for i := int64(0); i < numParts; i++ {
		parts[i] = keys[i*int64(len(keys))/numParts : (i+1)*int64(len(keys))/numParts]
	}
	return parts
}
