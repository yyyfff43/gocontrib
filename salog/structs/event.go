/*
 * Created by deng shi wei on 2020/01/06.
 * Copyright 2015－2020 Sensors Data Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package structs

import (
	"errors"
	"regexp"
	"time"
)

const (
	KeyMax   = 256
	ValueMax = 8192

	NamePatternBad = "^(^distinct_id$|^original_id$|^time$|^properties$|^id$|^first_id$|^second_id$|^users$|^events$|^event$|^user_id$|^date$|^datetime$)$"
	NamePatternOk  = "^[a-zA-Z_$][a-zA-Z\\d_$]{0,99}$"
)

var patternBad, patternOk *regexp.Regexp

type EventData struct {
	Type          string                 `json:"type"`
	Time          int64                  `json:"time"`
	DistinctId    string                 `json:"distinct_id"`
	Properties    map[string]interface{} `json:"properties"`
	LibProperties LibProperties          `json:"lib"`
	Project       string                 `json:"project"`
	Event         string                 `json:"event"`
	OriginId      string                 `json:"original_id,omitempty"`
	TimeFree      bool                   `json:"time_free,omitempty"`
}

func init() {
	patternBad, _ = regexp.Compile(NamePatternBad)
	patternOk, _ = regexp.Compile(NamePatternOk)
}

func (e *EventData) NormalizeData() error {
	//check distinct id
	if e.DistinctId == "" || len(e.DistinctId) == 0 {
		return errors.New("property [distinct_id] can't be empty")
	}

	if len(e.DistinctId) > 255 {
		return errors.New("the max length of property [distinct_id] is 255")
	}

	//check event
	if e.Event != "" {
		isMatch := checkPattern([]byte(e.Event))
		if !isMatch {
			return errors.New("event name must be a valid variable name")
		}
	}

	//check project
	if e.Project != "" {
		isMatch := checkPattern([]byte(e.Project))
		if !isMatch {
			return errors.New("project name must be a valid variable name")
		}
	}

	//check properties
	if e.Properties != nil {
		for k, v := range e.Properties {
			//check key
			if len(k) > KeyMax {
				return errors.New("the max length of property key is 256")
			}
			isMatch := checkPattern([]byte(k))
			if !isMatch {
				return errors.New("property key must be a valid variable name")
			}

			//check value
			switch vv := v.(type) {
			case int:
			case bool:
			case float64:
			case string:
				if len(vv) > ValueMax {
					return errors.New("the max length of property key is 8192")
				}
			case []string: //value in properties list MUST be string
			case time.Time: //only support time.Time
				e.Properties[k] = v.(time.Time).Format("2006-01-02 15:04:05.999")

			default:
				return errors.New("property value must be a string/int/float64/bool/time.Time/[]string")
			}
		}
	}

	return nil
}

func checkPattern(name []byte) bool {
	return !patternBad.Match(name) && patternOk.Match(name)
}
