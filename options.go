// Copyright 2021 Grabtaxi Holdings Pte Ltd (GRAB), All rights reserved.

// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package gosm

import (
	"gosm/gosmpb"
)

type logger interface {
	Printf(format string, v ...interface{})
}

// Option type provides different options when initializing the encoder.
type Option func(e *Encoder)

// WithLogger sets a logger.
func WithLogger(l logger) Option {
	return func(e *Encoder) {
		e.logger = l
	}
}

// WithOptionalFeatures sets optional features for the encoder
func WithOptionalFeatures(optionalFeatures []string) Option {
	return func(e *Encoder) {
		e.optionalFeatures = optionalFeatures
	}
}

// WithZlipEnabled enabled zlib when writing data to pbf file.
func WithZlipEnabled(v bool) Option {
	return func(e *Encoder) {
		e.enableZlip = v
	}
}

// WithBbox provides the pbf bbox information in the header block, just an information, does not limit the writing if the data is out of the bbox.
func WithBbox(b *gosmpb.HeaderBBox) Option {
	return func(e *Encoder) {
		e.bbox = b
	}
}

// WithWritingProgram sets the writing program, does not impact the encoding behaviour.
func WithWritingProgram(w string) Option {
	return func(e *Encoder) {
		e.writingProgram = w
	}
}
