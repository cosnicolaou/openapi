// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package transforms

import (
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

type T interface {
	Name() string
	Describe(node yaml.Node) string
	Configure(node yaml.Node) error
	Transform(*openapi3.T) (*openapi3.T, error)
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
