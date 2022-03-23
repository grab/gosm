// Copyright 2021 Grabtaxi Holdings Pte Ltd (GRAB), All rights reserved.

// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package gosm

type stringTable struct {
	indexMap map[string]int
	table    []string
	all      []string
}

func newStringTable() *stringTable {
	return &stringTable{
		indexMap: map[string]int{},
		table:    []string{""},
	}
}

//  --- DenseNodes ---

func (t *stringTable) add(str string) {
	if _, ok := t.indexMap[str]; !ok {
		t.table = append(t.table, str)
		t.indexMap[str] = len(t.table) - 1
	}
	t.all = append(t.all, str)
}

func (t *stringTable) index(str string) int {
	if v, ok := t.indexMap[str]; ok {
		return v
	}
	return -1
}

func (t *stringTable) endOne() {
	t.all = append(t.all, "##")
}

// toKeysVals is used in DenseNodes.
func (t *stringTable) toKeysVals() []int32 {
	var result []int32
	for _, s := range t.all {
		if s == "##" {
			result = append(result, 0)
			continue
		}
		result = append(result, int32(t.index(s)))
	}
	return result
}

// --- Node / Way / Relation ---

func (t *stringTable) addTags(tags map[string]string) ([]uint32, []uint32) {
	keyIDs := make([]uint32, len(tags))
	valIDs := make([]uint32, len(tags))
	cnt := 0
	for k, v := range tags {
		kIdx := t.append(k)
		keyIDs[cnt] = kIdx

		vIdx := t.append(v)
		valIDs[cnt] = vIdx

		cnt++
	}
	return keyIDs, valIDs
}

func (t *stringTable) append(str string) uint32 {
	idx, ok := t.indexMap[str]
	if !ok {
		t.table = append(t.table, str)
		idx = len(t.table) - 1
	}
	return uint32(idx)
}

// used in relations.
func (t *stringTable) addRoles(roles []string) []int32 {
	result := make([]int32, len(roles))
	for i, role := range roles {
		id := t.append(role)
		result[i] = int32(id)
	}
	return result
}
