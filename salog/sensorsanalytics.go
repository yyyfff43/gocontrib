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

package sensorsanalytics

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"

	"git.zhwenxue.com/zhgo/gocontrib/salog/consumers"
	"git.zhwenxue.com/zhgo/gocontrib/salog/structs"
	"git.zhwenxue.com/zhgo/gocontrib/salog/utils"
)

const (
	TRACK            = "track"
	TrackSignup      = "track_signup"
	ProfileSet       = "profile_set"
	ProfileSetOnce   = "profile_set_once"
	ProfileIncrement = "profile_increment"
	ProfileAppend    = "profile_append"
	ProfileUnset     = "profile_unset"
	ProfileDelete    = "profile_delete"
	ItemSet          = "item_set"
	ItemDelete       = "item_delete"

	SdkVersion = "2.0.3"
	LibName    = "Golang"

	MaxIdLen = 255
)

type SensorsAnalytics struct {
	C           consumers.Consumer
	ProjectName string
	TimeFree    bool
}

func InitSensorsAnalytics(c consumers.Consumer, projectName string, timeFree bool) SensorsAnalytics {
	return SensorsAnalytics{C: c, ProjectName: projectName, TimeFree: timeFree}
}

func (sa *SensorsAnalytics) track(etype, event, distinctId, originId string, properties map[string]interface{}, isLoginId bool) error {
	eventTime := utils.NowMs()
	if et := extractUserTime(properties); et > 0 {
		eventTime = et
	}

	data := structs.EventData{
		Type:          etype,
		Time:          eventTime,
		DistinctId:    distinctId,
		Properties:    properties,
		LibProperties: getLibProperties(),
	}

	if sa.ProjectName != "" {
		data.Project = sa.ProjectName
	}

	if etype == TRACK || etype == TrackSignup {
		data.Event = event
	}

	if etype == TrackSignup {
		data.OriginId = originId
	}

	if sa.TimeFree {
		data.TimeFree = true
	}

	if isLoginId {
		properties["$is_login_id"] = true
	}

	err := data.NormalizeData()
	if err != nil {
		return err
	}

	return sa.C.Send(data)
}

func (sa *SensorsAnalytics) Flush() {
	_ = sa.C.Flush()
}

func (sa *SensorsAnalytics) Close() {
	_ = sa.C.Close()
}

func (sa *SensorsAnalytics) Track(distinctId, event string, properties map[string]interface{}, isLoginId bool) error {
	var nproperties map[string]interface{}

	// merge properties
	if properties == nil {
		nproperties = make(map[string]interface{})
	} else {
		nproperties = utils.DeepCopy(properties)
	}

	nproperties["$lib"] = LibName
	nproperties["$lib_version"] = SdkVersion

	return sa.track(TRACK, event, distinctId, "", nproperties, isLoginId)
}

func (sa *SensorsAnalytics) TrackSignup(distinctId, originId string) error {
	// check originId and merge properties
	if originId == "" {
		return errors.New("property [original_id] must not be empty")
	}
	if len(originId) > MaxIdLen {
		return errors.New("the max length of property [original_id] is 255")
	}

	properties := make(map[string]interface{})

	properties["$lib"] = LibName
	properties["$lib_version"] = SdkVersion

	return sa.track(TrackSignup, "$SignUp", distinctId, originId, properties, false)
}

func (sa *SensorsAnalytics) ProfileSet(distinctId string, properties map[string]interface{}, isLoginId bool) error {
	var nproperties map[string]interface{}

	if properties == nil {
		return errors.New("property should not be nil")
	} else {
		nproperties = utils.DeepCopy(properties)
	}

	return sa.track(ProfileSet, "", distinctId, "", nproperties, isLoginId)
}

func (sa *SensorsAnalytics) ProfileSetOnce(distinctId string, properties map[string]interface{}, isLoginId bool) error {
	var nproperties map[string]interface{}

	if properties == nil {
		return errors.New("property should not be nil")
	} else {
		nproperties = utils.DeepCopy(properties)
	}

	return sa.track(ProfileSetOnce, "", distinctId, "", nproperties, isLoginId)
}

func (sa *SensorsAnalytics) ProfileIncrement(distinctId string, properties map[string]interface{}, isLoginId bool) error {
	var nproperties map[string]interface{}

	if properties == nil {
		return errors.New("property should not be nil")
	} else {
		nproperties = utils.DeepCopy(properties)
	}

	return sa.track(ProfileIncrement, "", distinctId, "", nproperties, isLoginId)
}

func (sa *SensorsAnalytics) ProfileAppend(distinctId string, properties map[string]interface{}, isLoginId bool) error {
	var nproperties map[string]interface{}

	if properties == nil {
		return errors.New("property should not be nil")
	} else {
		nproperties = utils.DeepCopy(properties)
	}

	return sa.track(ProfileAppend, "", distinctId, "", nproperties, isLoginId)
}

func (sa *SensorsAnalytics) ProfileUnset(distinctId string, properties map[string]interface{}, isLoginId bool) error {
	var nproperties map[string]interface{}

	if properties == nil {
		return errors.New("property should not be nil")
	} else {
		nproperties = utils.DeepCopy(properties)
	}

	return sa.track(ProfileUnset, "", distinctId, "", nproperties, isLoginId)
}

func (sa *SensorsAnalytics) ProfileDelete(distinctId string, isLoginId bool) error {
	nproperties := make(map[string]interface{})

	return sa.track(ProfileDelete, "", distinctId, "", nproperties, isLoginId)
}

func (sa *SensorsAnalytics) ItemSet(itemType string, itemId string, properties map[string]interface{}) error {
	libProperties := getLibProperties()
	time := utils.NowMs()
	if properties == nil {
		properties = map[string]interface{}{}
	}

	itemData := structs.Item{
		Type:          ItemSet,
		ItemId:        itemId,
		Time:          time,
		ItemType:      itemType,
		Properties:    properties,
		LibProperties: libProperties,
	}

	err := itemData.NormalizeItem()
	if err != nil {
		return err
	}

	return sa.C.ItemSend(itemData)
}

func (sa *SensorsAnalytics) ItemDelete(itemType string, itemId string) error {
	libProperties := getLibProperties()
	time := utils.NowMs()

	itemData := structs.Item{
		Type:          ItemDelete,
		ItemId:        itemId,
		Time:          time,
		ItemType:      itemType,
		Properties:    map[string]interface{}{},
		LibProperties: libProperties,
	}

	err := itemData.NormalizeItem()
	if err != nil {
		return err
	}

	return sa.C.ItemSend(itemData)
}

func getLibProperties() structs.LibProperties {
	lp := structs.LibProperties{}
	lp.Lib = LibName
	lp.LibVersion = SdkVersion
	lp.LibMethod = "code"
	if pc, file, line, ok := runtime.Caller(3); ok { //3 means sdk's caller
		f := runtime.FuncForPC(pc)
		lp.LibDetail = fmt.Sprintf("##%s##%s##%d", f.Name(), file, line)
	}

	return lp
}

func extractUserTime(p map[string]interface{}) int64 {
	if t, ok := p["$time"]; ok {
		v, ok := t.(int64)
		if !ok {
			fmt.Fprintln(os.Stderr, "It's not ok for type string")
			return 0
		}
		delete(p, "$time")

		return v
	}

	return 0
}

func InitDefaultConsumer(url string, timeout int) (*consumers.DefaultConsumer, error) {
	return consumers.InitDefaultConsumer(url, timeout)
}

func InitBatchConsumer(url string, max, timeout int) (*consumers.BatchConsumer, error) {
	return consumers.InitBatchConsumer(url, max, timeout)
}

func InitLoggingConsumer(filename string, hourRotate bool) (*consumers.LoggingConsumer, error) {
	return consumers.InitLoggingConsumer(filename, hourRotate)
}

func InitDebugConsumer(url string, writeData bool, timeout int) (*consumers.DebugConsumer, error) {
	return consumers.InitDebugConsumer(url, writeData, timeout)
}

func InitIOConsumer(w io.Writer) (*consumers.IOConsumer, error) {
	return consumers.InitIOConsumer(w)
}

func NewNopSensorsAnalytics() SensorsAnalytics {
	consumer := consumers.InitNopConsumer()
	s := InitSensorsAnalytics(consumer, "default", false)
	return s
}
