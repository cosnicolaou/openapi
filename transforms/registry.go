// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package transforms

import (
	"errors"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

var (
	ErrTransformNotImplementedForV2 = errors.New("transform not implemented for v2 schemas")
	ErrTransformNotImplementedForV3 = errors.New("transform not implemented for v3 schemas")
)

type T interface {
	Name() string
	Describe(node yaml.Node) string
	Configure(node yaml.Node) error
	TransformV2(*openapi2.T) (*openapi2.T, error)
	TransformV3(*openapi3.T) (*openapi3.T, error)
}

var installed = map[string]T{}

func Register(t T) {
	installed[t.Name()] = t
}

func List() []string {
	var r []string
	for k := range installed {
		r = append(r, k)
	}
	return r
}

func Get(name string) T {
	return installed[name]
}
