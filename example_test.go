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
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"
	"testing"

	"github.com/qedus/osmpbf"
	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
	nodes := []*Node{
		{
			ID: 7278995748,
			Latitude:  -7.2380901,
			Longitude: 112.6773289,
			Tags: map[string]string{
				"node": "node1",
				"erp":  "no",
			},
		},
		{
			ID: 6978510772,
			Latitude:  -7.2381273,
			Longitude: 112.6775354,
		},
		{
			ID: 6978510773,
			Latitude:  -7.2383685,
			Longitude: 112.6782548,
			Tags: map[string]string{
				"node":  "node3",
				"align": "way",
				"erp":   "yes",
			},
		},
	}
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
				"k6": "v6",
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
	relations := []*Relation{
		{
			ID: 437710,
			Tags: map[string]string{
				"currency":         "IDR",
				"city_code":        "SUB",
				"gantry:price:car": "8000.00",
			},
			Members: []Member{
				{
					ID:   260114973,
					Type: NodeType,
					Role: "from",
				},
				{
					ID:   1780865796,
					Type: NodeType,
					Role: "to",
				},
			},
		},
		{
			ID: 436226,
			Members: []Member{
				{
					ID:   2418135255,
					Type: NodeType,
					Role: "through",
				},
			},
		},
		{
			ID: 2000143,
			Tags: map[string]string{
				"restriction": "no_left_turn",
				"start_date":  "2017-01-01",
				"end_date":    "9999-12-31",
				"city_code":   "KNO",
			},
			Members: []Member{
				{
					ID:   3354491584,
					Type: NodeType,
					Role: "to",
				},
				{
					ID:   3354491587,
					Type: NodeType,
					Role: "via",
				},
				{
					ID:   5416035667,
					Type: NodeType,
					Role: "from",
				},
			},
		},
	}

	for _, enableZlip := range []bool{false, true} {
		wc := &myWriter{}
		encoder := NewEncoder(&NewEncoderRequiredInput{
			RequiredFeatures: []string{"OsmSchema-V0.6"},
			Writer:           wc,
		},
			WithWritingProgram("example"),
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
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			encoder.AppendNodes(nodes)
			encoder.Flush(NodeType)
			wg.Done()
		}()
		go func() {
			encoder.AppendWays(ways)
			encoder.Flush(WayType)
			wg.Done()
		}()
		go func() {
			encoder.AppendRelations(relations)
			encoder.Flush(RelationType)
			wg.Done()
		}()
		wg.Wait()
		assert.Nil(t, encoder.Close())
		assert.Nil(t, errs)
		decoder := osmpbf.NewDecoder(wc)
		decodeErr := decoder.Start(runtime.GOMAXPROCS(-1))
		if decodeErr != nil {
			log.Fatalln(decodeErr)
		}
		var resultNodes []*osmpbf.Node
		var resultRelations []*osmpbf.Relation
		var resultWays []*osmpbf.Way
		for {
			v, err := decoder.Decode()
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("decode ways err: %+v", err)
			}
			switch v := v.(type) {
			case *osmpbf.Node:
				resultNodes = append(resultNodes, v)
			case *osmpbf.Relation:
				resultRelations = append(resultRelations, v)
			case *osmpbf.Way:
				resultWays = append(resultWays, v)
			}
		}
		assert.Equal(t, 3, len(resultNodes))
		for i, n := range resultNodes {
			assert.Equal(t, n.ID, nodes[i].ID)
			assert.True(t, n.Lat-nodes[i].Latitude > -0.001 && n.Lat-nodes[i].Latitude < 0.001)
			assert.True(t, n.Lon-nodes[i].Longitude > -0.001 && n.Lon-nodes[i].Longitude < 0.001)
			for k, v := range n.Tags {
				assert.Equal(t, v, nodes[i].Tags[k])
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
		assert.Equal(t, 3, len(resultRelations))
		for i, r := range resultRelations {
			assert.Equal(t, r.ID, relations[i].ID)
			for k, v := range r.Tags {
				assert.Equal(t, v, r.Tags[k])
			}
			for j, mem := range r.Members {
				assert.Equal(t, relations[i].Members[j].ID, mem.ID)
				assert.Equal(t, relations[i].Members[j].Role, mem.Role)
				assert.Equal(t, int(relations[i].Members[j].Type), int(mem.Type))
			}
		}
	}
}
