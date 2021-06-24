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

import "gosm/gosmpb"

// Way ...
type Way struct {
	ID      int64
	Tags    map[string]string
	NodeIDs []int64
}

// AppendWays ...
func (e *Encoder) AppendWays(ways []*Way) {
	e.waysBuf <- &wayMembers{
		ws: ways,
	}
}

// wayMembers implements members
type wayMembers struct {
	ws []*Way
}

func (w *wayMembers) toPrimitiveBlock() (*gosmpb.PrimitiveBlock, error) {
	st := newStringTable()
	var ways []*gosmpb.Way
	for _, w := range w.ws {
		keyIDs, valIDs := st.addTags(w.Tags)
		ways = append(ways, &gosmpb.Way{
			Id:   int64ToPointer(w.ID),
			Keys: keyIDs,
			Vals: valIDs,
			Refs: deltaEncodeInt64s(w.NodeIDs),
		})
	}

	pb := &gosmpb.PrimitiveBlock{
		Stringtable: &gosmpb.StringTable{
			S: st.table,
		},
		Primitivegroup: []*gosmpb.PrimitiveGroup{
			{
				Ways: ways,
			},
		},
	}
	return pb, nil
}

func (w *wayMembers) len() int {
	return len(w.ws)
}

func (w *wayMembers) clear() {
	w.ws = nil
}

func (w *wayMembers) appendMembers(m members) {
	w1, ok := m.(*wayMembers)
	if ok {
		w.ws = append(w.ws, w1.ws...)
	}
}
