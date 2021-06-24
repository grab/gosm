// Copyright (c) 2012-2021 Grabtaxi Holdings PTE LTD (GRAB), All Rights Reserved. NOTICE: All information contained herein
// is, and remains the property of GRAB. The intellectual and technical concepts contained herein are confidential, proprietary
// and controlled by GRAB and may be covered by patents, patents in process, and are protected by trade secret or copyright law.
//
// You are strictly forbidden to copy, download, store (in any medium), transmit, disseminate, adapt or change this material
// in any way unless prior written permission is obtained from GRAB. Access to the source code contained herein is hereby
// forbidden to anyone except current GRAB employees or contractors with binding Confidentiality and Non-disclosure agreements
// explicitly covering such access.
//
// The copyright notice above does not evidence any actual or intended publication or disclosure of this source code,
// which includes information that is confidential and/or proprietary, and is a trade secret, of GRAB.
//
// ANY REPRODUCTION, MODIFICATION, DISTRIBUTION, PUBLIC PERFORMANCE, OR PUBLIC DISPLAY OF OR THROUGH USE OF THIS SOURCE
// CODE WITHOUT THE EXPRESS WRITTEN CONSENT OF GRAB IS STRICTLY PROHIBITED, AND IN VIOLATION OF APPLICABLE LAWS AND
// INTERNATIONAL TREATIES. THE RECEIPT OR POSSESSION OF THIS SOURCE CODE AND/OR RELATED INFORMATION DOES NOT CONVEY
// OR IMPLY ANY RIGHTS TO REPRODUCE, DISCLOSE OR DISTRIBUTE ITS CONTENTS, OR TO MANUFACTURE, USE, OR SELL ANYTHING
// THAT IT MAY DESCRIBE, IN WHOLE OR IN PART.

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

// toKeysVals is used in DenseNodes
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

// used in relations
func (t *stringTable) addRoles(roles []string) []int32 {
	result := make([]int32, len(roles))
	for i, role := range roles {
		id := t.append(role)
		result[i] = int32(id)
	}
	return result
}
