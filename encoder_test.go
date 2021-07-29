// Copyright 2021 Grabtaxi Holdings Pte Ltd (GRAB), All rights reserved.

// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package gosm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewEncoder ...
func TestNewEncoder(t *testing.T) {
	wc := &myWriter{}
	encoder := NewEncoder(&NewEncoderRequiredInput{
		RequiredFeatures: []string{"OsmSchema-V0.6"},
		Writer:           wc,
	},
		WithWritingProgram("example"),
		WithZlipEnabled(false),
	)

	assert.Equal(t, "example", encoder.writingProgram)
	assert.Equal(t, []string{"OsmSchema-V0.6"}, encoder.requiredFeatures)
	assert.Equal(t, false, encoder.enableZlip)
}
