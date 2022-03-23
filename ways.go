// Copyright 2021 Grabtaxi Holdings Pte Ltd (GRAB), All rights reserved.

// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

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

// wayMembers implements members.
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
