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
