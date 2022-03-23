// Copyright 2021 Grabtaxi Holdings Pte Ltd (GRAB), All rights reserved.

// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package gosm

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"
	"testing"

	"github.com/qedus/osmpbf"
	"github.com/stretchr/testify/assert"
)

type myWriter struct {
	buf  bytes.Buffer
	lock sync.Mutex
}

func (w *myWriter) Write(p []byte) (n int, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.buf.Write(p)
}

func (w *myWriter) Close() error {
	return nil
}

func (w *myWriter) Read(p []byte) (n int, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.buf.Read(p)
}

func TestWriteNodes(t *testing.T) {
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
				"node":       "node3",
				"align":      "way",
				"emptyValue": "",
				"erp":        "yes",
			},
		},
		{
			ID: 6978510774,
			Latitude:  -7.2383445,
			Longitude: 112.6734548,
			Tags: map[string]string{
				"node":       "node3",
				"align":      "way",
				"emptyValue": "",
				"ref":        "0",
			},
		},
	}
	for _, enableZlip := range []bool{false, true} {
		wc := &myWriter{}
		encoder := NewEncoder(&NewEncoderRequiredInput{
			RequiredFeatures: []string{"OsmSchema-V0.6"},
			Writer:           wc,
		},
			WithWritingProgram("nodeTesting"),
			WithZlipEnabled(enableZlip),
		)
		errChan, err := encoder.Start()
		assert.Nil(t, err)
		var errs []error
		var errsLock sync.Mutex
		go func() {
			for {
				e, ok := <-errChan
				if ok {
					errsLock.Lock()
					errs = append(errs, e)
					errsLock.Unlock()
				}
			}
		}()
		encoder.AppendNodes(nodes)
		assert.Nil(t, encoder.Close())
		errsLock.Lock()
		assert.Nil(t, errs)
		errsLock.Unlock()

		decoder := osmpbf.NewDecoder(wc)
		decoder.SetBufferSize(osmpbf.MaxBlobSize)
		decodeErr := decoder.Start(runtime.GOMAXPROCS(-1))
		if decodeErr != nil {
			log.Fatalln(decodeErr)
		}

		var resultNodes []*osmpbf.Node
		for {
			v, err := decoder.Decode()
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("decode err :%+v\n", err)
			}

			switch v := v.(type) {
			case *osmpbf.Node:
				resultNodes = append(resultNodes, v)
			}
		}
		assert.Equal(t, 4, len(resultNodes))
		for i, n := range resultNodes {
			assert.Equal(t, n.ID, nodes[i].ID)
			assert.True(t, n.Lat-nodes[i].Latitude > -0.001 && n.Lat-nodes[i].Latitude < 0.001)
			assert.True(t, n.Lon-nodes[i].Longitude > -0.001 && n.Lon-nodes[i].Longitude < 0.001)
			for k, v := range n.Tags {
				assert.Equal(t, v, nodes[i].Tags[k])
			}
		}
	}
}
