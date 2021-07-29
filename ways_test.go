// Copyright 2021 Grabtaxi Holdings Pte Ltd (GRAB), All rights reserved.

// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package gosm

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"testing"

	"github.com/qedus/osmpbf"
	"github.com/stretchr/testify/assert"
)

func TestWriteWays(t *testing.T) {
	ways := []*Way{
		{
			ID: 9650669,
			Tags: map[string]string{
				"k1": "v1",
				"k2": "v2",
			},
			NodeIDs: []int64{75385503, 75390364, 75390426, 1116369070},
		},
		{
			ID: 9650692,
			Tags: map[string]string{
				"k3": "v3",
				"k4": "v4",
				"k5": "v5",
			},
			NodeIDs: []int64{603386705, 75444477, 760790597, 760790382, 760790527, 75444457},
		},
		{
			ID: 11750310,
			Tags: map[string]string{
				"k6":         "v6",
				"emptyValue": "",
			},
			NodeIDs: []int64{105207733, 105207737, 105207726},
		},
		{
			ID: 11750349,
			Tags: map[string]string{
				"k7":  "v7",
				"k8":  "v8",
				"k9":  "v9",
				"k10": "v10",
			},
			NodeIDs: []int64{105208152, 105208155, 105208157, 105208163, 105208165, 105208167, 105208168, 105208174, 2363909540},
		},
	}
	for _, enableZlip := range []bool{false, true} {
		wc := &myWriter{}

		encoder := NewEncoder(&NewEncoderRequiredInput{
			RequiredFeatures: []string{"OsmSchema-V0.6"},
			Writer:           wc,
		},
			WithWritingProgram("wayTesting"),
			WithZlipEnabled(enableZlip),
		)
		errChan, err := encoder.Start()
		assert.Nil(t, err)

		var errs []error
		go func() {
			for {
				e, ok := <-errChan
				if ok {
					errs = append(errs, e)
				}
			}
		}()
		encoder.AppendWays(ways)
		assert.Nil(t, encoder.Close())
		assert.Nil(t, errs)
		decoder := osmpbf.NewDecoder(wc)
		decodeErr := decoder.Start(runtime.GOMAXPROCS(-1))
		if decodeErr != nil {
			log.Fatalln(decodeErr)
		}

		var resultWays []*osmpbf.Way
		for {
			w, err := decoder.Decode()
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("decode ways err: %+v", err)
			}
			switch w := w.(type) {
			case *osmpbf.Way:
				resultWays = append(resultWays, w)
			}
		}
		assert.Equal(t, 4, len(resultWays))
		for i, w := range resultWays {
			assert.Equal(t, w.ID, ways[i].ID)
			for k, v := range w.Tags {
				assert.Equal(t, v, w.Tags[k])
			}
			for j, id := range w.NodeIDs {
				assert.Equal(t, id, ways[i].NodeIDs[j])
			}
		}
	}
}
