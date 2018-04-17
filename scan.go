// Copy from redigo
// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package simplecache

import (
	"errors"
	"fmt"
	"reflect"
)

func convertAssign(d interface{}, s interface{}) (err error) {
	v_dst := reflect.ValueOf(d)
	v_src := reflect.ValueOf(s)

	v_dst.Elem().Set(v_src)

	return nil
}

func Scan(src []interface{}, dest ...interface{}) ([]interface{}, error) {
	if len(src) < len(dest) {
		return nil, errors.New("simplecacheScan: array short")
	}
	var err error
	for i, d := range dest {
		err = convertAssign(d, src[i])
		if err != nil {
			err = fmt.Errorf("simplecacheScan: cannot assign to dest %d: %v", i, err)
			break
		}
	}
	return src[len(dest):], err
}
