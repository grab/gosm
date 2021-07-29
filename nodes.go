// Copyright 2021 Grabtaxi Holdings Pte Ltd (GRAB), All rights reserved.

// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package gosm

import (
	"gosm/gosmpb"
)

// Node ...
type Node struct {
	ID        int64
	Latitude  float64
	Longitude float64
	Tags      map[string]string
}

// AppendNodes will append nodes to the buffer, when it meets the limit(8000 entities or 32MB)
// it will convert the nodes to dense nodes and write to the writer
func (e *Encoder) AppendNodes(nodes []*Node) {
	e.nodesBuf <- &nodeMembers{
		ns: nodes,
	}
}

// nodeMembers implements members
type nodeMembers struct {
	ns []*Node
}

// we will simplify the design
func (n *nodeMembers) toPrimitiveBlock() (*gosmpb.PrimitiveBlock, error) {
	st := newStringTable()
	lats := make([]float64, len(n.ns))
	lngs := make([]float64, len(n.ns))
	ids := make([]int64, len(n.ns))
	for i, n := range n.ns {
		for k, v := range n.Tags {
			st.add(k)
			st.add(v)
		}
		st.endOne()
		lats[i] = n.Latitude
		lngs[i] = n.Longitude
		ids[i] = n.ID
	}
	granularity, latOffset, lngOffset, deltaLats, deltaLngs := deltaEncodeCoordinates(lats, lngs)
	pb := &gosmpb.PrimitiveBlock{
		Stringtable: &gosmpb.StringTable{
			S: st.table,
		},
		Granularity: int32ToPointer(granularity),
		LatOffset:   int64ToPointer(latOffset),
		LonOffset:   int64ToPointer(lngOffset),
		Primitivegroup: []*gosmpb.PrimitiveGroup{
			{
				Dense: &gosmpb.DenseNodes{
					Id:       deltaEncodeInt64s(ids),
					Lat:      deltaEncodeInt64s(deltaLats),
					Lon:      deltaEncodeInt64s(deltaLngs),
					KeysVals: st.toKeysVals(),
				},
			},
		},
	}

	return pb, nil
}

func (n *nodeMembers) len() int {
	return len(n.ns)
}

func (n *nodeMembers) clear() {
	n.ns = nil
}

func (n *nodeMembers) appendMembers(m members) {
	n1, ok := m.(*nodeMembers)
	if ok {
		n.ns = append(n.ns, n1.ns...)
	}
}

func deltaEncodeInt64s(nums []int64) []int64 {
	if len(nums) == 0 {
		return nil
	}
	result := make([]int64, len(nums))
	result[0] = nums[0]
	for i := 1; i < len(result); i++ {
		result[i] = nums[i] - nums[i-1]
	}
	return result
}
