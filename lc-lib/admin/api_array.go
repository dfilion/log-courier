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

package admin

import (
	"bytes"
	"encoding/json"
	"strconv"
)

type apiArrayEntry struct {
	row   int
	entry APIEntry
}

// APIArray represents an array of entries in the API accessible through a
// primary key
type APIArray struct {
	entryMap map[string]apiArrayEntry
	entries  []APIEntry
}

// AddEntry a new array entry
func (a *APIArray) AddEntry(key string, entry APIEntry) {
	if a.entryMap == nil {
		a.entryMap = make(map[string]apiArrayEntry)
	} else {
		if _, ok := a.entryMap[key]; ok {
			panic("Key already exists")
		}
	}

	a.entryMap[key] = apiArrayEntry{
		row:   len(a.entries),
		entry: entry,
	}

	a.entries = append(a.entries, entry)
}

// RemoveEntry removes an array entry
func (a *APIArray) RemoveEntry(key string) {
	if a.entryMap == nil {
		panic("Array has no entries")
	}

	entry, ok := a.entryMap[key]
	if !ok {
		panic("Entry not found")
	}

	delete(a.entryMap, key)
	for i := entry.row; i+1 < len(a.entries); i++ {
		a.entries[i] = a.entries[i+1]
	}
	a.entries = a.entries[:len(a.entries)-1]
}

// Get returns an entry using it's primary key name or row number
func (a *APIArray) Get(path string) (APIEntry, error) {
	if a.entryMap == nil {
		return nil, nil
	}

	if entry, ok := a.entryMap[path]; ok {
		return entry.entry, nil
	}

	// Try to parse as a number
	entryNum, err := strconv.ParseInt(path, 10, 0)
	if err != nil {
		return nil, err
	}

	if entryNum < 0 || entryNum >= int64(len(a.entries)) {
		return nil, nil
	}

	return a.entries[entryNum], nil
}

// MarshalJSON returns the APIArray in JSON form
func (a *APIArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.entries)
}

// HumanReadable returns the APIArray as a string
func (a *APIArray) HumanReadable(indent string) ([]byte, error) {
	if a.entryMap == nil {
		return nil, nil
	}

	var result bytes.Buffer
	newIndent := indent + APIIndentation

	for primaryKey, entry := range a.entryMap {
		subResult, err := entry.entry.HumanReadable(newIndent)
		if err != nil {
			return nil, err
		}

		if bytes.IndexRune(subResult, '\n') != -1 {
			result.WriteString(primaryKey)
			result.WriteString(":\n")
			result.Write(subResult)
			continue
		}

		result.WriteString(primaryKey)
		result.WriteString(": ")
		result.Write(subResult)
	}

	return result.Bytes(), nil
}

// Update ensures the data we have is up to date - should be overriden by users
// if required to keep the contents up to date on each request
// Default behaviour is to update each of the array entries
func (a *APIArray) Update() error {
	for _, entry := range a.entries {
		if err := entry.Update(); err != nil {
			return err
		}
	}

	return nil
}