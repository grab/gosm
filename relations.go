// Copyright 2021 Grabtaxi Holdings Pte Ltd (GRAB), All rights reserved.

// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package gosm

import "gosm/gosmpb"

// Relation ...
type Relation struct {
	ID      int64
	Tags    map[string]string
	Members []Member
	Info    gosmpb.Info
}

// MemberType ...
type MemberType int

const (
	// NodeType ...
	NodeType MemberType = iota
	// WayType ...
	WayType
	// RelationType ...
	RelationType
)

// Member ...
type Member struct {
	ID   int64
	Type MemberType
	Role string
}

// AppendRelations ...
func (e *Encoder) AppendRelations(relations []*Relation) {
	e.relationsBuf <- &relationMembers{
		rs: relations,
	}
}

// relationMembers implements members
type relationMembers struct {
	rs []*Relation
}

func (r *relationMembers) toPrimitiveBlock() (*gosmpb.PrimitiveBlock, error) {
	st := newStringTable()
	var relations []*gosmpb.Relation
	for _, m := range r.rs {
		keys, values := st.addTags(m.Tags)
		memberIDs := make([]int64, len(m.Members))
		memberTypes := make([]gosmpb.Relation_MemberType, len(m.Members))
		roles := make([]string, len(m.Members))
		for i, mem := range m.Members {
			memberIDs[i] = mem.ID
			memberTypes[i] = gosmpb.Relation_MemberType(mem.Type)
			roles[i] = mem.Role
		}
		roleIDs := st.addRoles(roles)
		relationPb := &gosmpb.Relation{
			Id:       int64ToPointer(m.ID),
			Keys:     keys,
			Vals:     values,
			RolesSid: roleIDs,
			Memids:   deltaEncodeInt64s(memberIDs),
			Types:    memberTypes,
		}
		relations = append(relations, relationPb)
	}

	pb := &gosmpb.PrimitiveBlock{
		Stringtable: &gosmpb.StringTable{
			S: st.table,
		},
		Primitivegroup: []*gosmpb.PrimitiveGroup{
			{
				Relations: relations,
			},
		},
	}
	return pb, nil
}

func (r *relationMembers) len() int {
	return len(r.rs)
}

func (r *relationMembers) clear() {
	r.rs = nil
}

func (r *relationMembers) appendMembers(m members) {
	r1, ok := m.(*relationMembers)
	if ok {
		r.rs = append(r.rs, r1.rs...)
	}
}
