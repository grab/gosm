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
