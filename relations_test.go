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
	"testing"

	"github.com/qedus/osmpbf"
	"github.com/stretchr/testify/assert"
)

func TestRelations(t *testing.T) {
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
				"emptyValue":  "",
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
			WithWritingProgram("relationTesting"),
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
		encoder.AppendRelations(relations)
		assert.Nil(t, encoder.Close())
		assert.Nil(t, errs)
		decoder := osmpbf.NewDecoder(wc)
		decodeErr := decoder.Start(runtime.GOMAXPROCS(-1))
		if decodeErr != nil {
			log.Fatalln(decodeErr)
		}

		var resultRelations []*osmpbf.Relation
		for {
			r, err := decoder.Decode()
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("decode ways err: %+v", err)
			}
			switch r := r.(type) {
			case *osmpbf.Relation:
				resultRelations = append(resultRelations, r)
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
