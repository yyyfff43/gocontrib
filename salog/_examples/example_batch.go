/*
 * Created by dengshiwei on 2020/01/06.
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

package main

import (
	"fmt"
	sdk "git.zhwenxue.com/zhgo/sheikah/salog"
)

func main() {
	c, err := sdk.InitBatchConsumer("http://localhost:8106/sa?project=production", 3, 1000)
	if err != nil {
		fmt.Println(err)
		return
	}

	sa := sdk.InitSensorsAnalytics(c, "default", false)
	defer sa.Close()

	distinctId := "ABCDEF123456"
	event := "ViewProduct"
	properties := map[string]interface{}{
		"$ip":            "2.2.2.2",
		"ProductId":      "123456",
		"ProductCatalog": "Laptop Computer",
		"IsAddedToFav":   true,
	}

	err = sa.Track(distinctId, event, properties, true)
	if err != nil {
		fmt.Println("track failed", err)
		return
	}

	err = sa.Track(distinctId, event, properties, true)
	if err != nil {
		fmt.Println("track failed", err)
		return
	}

	err = sa.Track(distinctId, event, properties, true)
	if err != nil {
		fmt.Println("track failed", err)
		return
	}

	fmt.Println("track done")
}
