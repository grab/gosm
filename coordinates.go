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

// latitude  = .000000001 * (lat_offset + (granularity * lat))
// longitude = .000000001 * (lon_offset + (granularity * lon))
func deltaEncodeCoordinates(lats []float64, lngs []float64) (granularity int32, latOffset int64, lngOffset int64, deltaLats []int64, deltaLngs []int64) {
	l := append(lats, lngs...)
	if len(l) == 0 {
		return 0, 0, 0, nil, nil
	}
	minG := int64(1000000000)
	for _, f := range l {
		a := int64(f * 1e9)
		g := int64(1)
		for {
			if a/10*10 == a {
				g = g * 10
				a = a / 10
			} else {
				break
			}
		}
		if g < minG {
			minG = g
		}
	}
	latOffset, deltaLats = deltaEncodeWithFixedGranularity(lats, int32(minG))
	lngOffset, deltaLngs = deltaEncodeWithFixedGranularity(lngs, int32(minG))
	return int32(minG), latOffset, lngOffset, deltaLats, deltaLngs
}

func deltaEncodeWithFixedGranularity(fn []float64, granularity int32) (offset int64, delta []int64) {
	if len(fn) == 0 {
		return 0, nil
	}
	min := int64(fn[0] * 1e9)
	for _, f := range fn {
		a := int64(f * 1e9)
		if a < min {
			min = a
		}
	}
	delta = make([]int64, len(fn))
	for i, f := range fn {
		delta[i] = (int64(f*1e9) - min) / int64(granularity)
	}
	return min, delta
}
