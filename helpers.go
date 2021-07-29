// Copyright 2021 Grabtaxi Holdings Pte Ltd (GRAB), All rights reserved.

// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package gosm

func countInt32LenOfBytes(b []byte) *int32 {
	l := int32(len(b))
	return &l
}

func stringToPointer(s string) *string {
	ss := s
	return &ss
}

func int32ToPointer(n int32) *int32 {
	a := n
	return &a
}

func int64ToPointer(n int64) *int64 {
	a := n
	return &a
}
