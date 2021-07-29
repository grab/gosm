// Copyright 2021 Grabtaxi Holdings Pte Ltd (GRAB), All rights reserved.

// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package gosm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeltaEncodeWithFixedGranularity(t *testing.T) {
	cases := []struct {
		inputFn          []float64
		inputGranularity int32
		expetectOffset   int64
		expectedInts     []int64
	}{
		{
			inputFn:          []float64{0.002345, 0.002346},
			inputGranularity: int32(1000),
			expetectOffset:   int64(2345000),
			expectedInts:     []int64{0, 1},
		},
		{
			inputFn:          []float64{0.002345, 0.002346, 1.0023456},
			inputGranularity: int32(100),
			expetectOffset:   int64(2345000),
			expectedInts:     []int64{0, 10, 10000006},
		},
		{
			inputFn:          []float64{-0.002345, -0.002346, 1.01},
			inputGranularity: int32(1000),
			expetectOffset:   int64(-2346000),
			expectedInts:     []int64{1, 0, 1012346},
		},
		{
			inputFn:          []float64{-0.002345, -0.002346, -1.01000007},
			inputGranularity: int32(10),
			expetectOffset:   int64(-1010000070),
			expectedInts:     []int64{100765507, 100765407, 0},
		},
		{
			inputFn:          []float64{0.002345, 0.002346},
			inputGranularity: int32(100),
			expetectOffset:   int64(2345000),
			expectedInts:     []int64{0, 10},
		},
	}
	for _, c := range cases {
		o, d := deltaEncodeWithFixedGranularity(c.inputFn, c.inputGranularity)
		assert.Equal(t, c.expetectOffset, o)
		assert.Equal(t, c.expectedInts, d)
		for i, in := range c.inputFn {
			assert.True(t, in-0.000000001*float64(o+int64(c.inputGranularity)*d[i]) > -0.001 && in-0.000000001*float64(o+int64(c.inputGranularity)*d[i]) < 0.01)
		}
	}
}

func TestDeltaEncodeCoordinates(t *testing.T) {
	cases := []struct {
		lats                []float64
		lngs                []float64
		expectedGranularity int32
		expectedLatOffset   int64
		expectedLngOffset   int64
		expectedDeltaLats   []int64
		expectedDeltaLngs   []int64
	}{
		{
			lats: []float64{1.3334369543955960, 1.3307062250574480, 1.3306006657420268, 1.3309222648360457, 1.3436963896736138, 1.3159375599949470, 1.3158660081205550},
			lngs:                []float64{103.9152950162963500, 103.8889361238400400, 103.8883285932482200, 103.8892271364421200, 103.7352140139293800, 103.7315496138722800, 103.7318148345315200},
			expectedGranularity: 1,
			expectedLatOffset:   1315866008,
			expectedLngOffset:   103731549613,
			expectedDeltaLats:   []int64{17570946, 14840217, 14734657, 15056256, 27830381, 71551, 0},
			expectedDeltaLngs:   []int64{183745403, 157386510, 156778980, 157677523, 3664400, 0, 265221},
		},
	}

	for _, c := range cases {
		g, latOffset, lngOffset, deltaLats, deltaLngs := deltaEncodeCoordinates(c.lats, c.lngs)
		assert.Equal(t, c.expectedGranularity, g)
		assert.Equal(t, c.expectedLatOffset, latOffset)
		assert.Equal(t, c.expectedLngOffset, lngOffset)
		assert.Equal(t, c.expectedDeltaLats, deltaLats)
		assert.Equal(t, c.expectedDeltaLngs, deltaLngs)
	}
}
