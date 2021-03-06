/*
 * Copyright 2014-2015 Jason Woods.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import "encoding/json"

// Event holds a key-value map that represents a single log event
type Event map[string]interface{}

// EventDescriptor describes an Event, such as it's source and offset, which can
// be used in order to resume log files
type EventDescriptor struct {
	Stream Stream
	Offset int64
	Event  []byte
}

// Encode returns the Event in JSON format
func (e Event) Encode() ([]byte, error) {
	return json.Marshal(e)
}
