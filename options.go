// Copyright 2021 Grabtaxi Holdings Pte Ltd (GRAB), All rights reserved.

// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package gosm

import (
	"gosm/gosmpb"
)

type logger interface {
	Printf(format string, v ...interface{})
}

// Option ...
type Option func(e *Encoder)

// WithLogger will set a logger
func WithLogger(l logger) Option {
	return func(e *Encoder) {
		e.logger = l
	}
}

// WithOptionalFeatures ...
func WithOptionalFeatures(optionalFeatures []string) Option {
	return func(e *Encoder) {
		e.optionalFeatures = optionalFeatures
	}
}

// WithZlipEnabled ...
func WithZlipEnabled(v bool) Option {
	return func(e *Encoder) {
		e.enableZlip = v
	}
}

// WithBbox ...
func WithBbox(b *gosmpb.HeaderBBox) Option {
	return func(e *Encoder) {
		e.bbox = b
	}
}

// WithWritingProgram ...
func WithWritingProgram(w string) Option {
	return func(e *Encoder) {
		e.writingProgram = w
	}
}
