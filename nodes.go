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
