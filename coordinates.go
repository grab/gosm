// Copyright 2021 Grabtaxi Holdings Pte Ltd (GRAB), All rights reserved.

// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

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
