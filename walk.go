// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type Walker interface {
	Walk(doc *openapi3.T)
}

type walkerOptions struct {
	followRefs bool
}

type WalkerOption func(o *walkerOptions)

func WalkerFollowRefs(v bool) WalkerOption {
	return func(o *walkerOptions) {
		o.followRefs = v
	}
}

type Visitor func(path []string, parent, node any) bool

// NewWalker returns a Walker that will visit every node in an openapi3 document.
func NewWalker(v Visitor, opts ...WalkerOption) Walker {
	w := &nodeWalker{visit: v}
	for _, opt := range opts {
		opt(&w.opts)
	}
	return w
}

type nodeWalker struct {
	opts  walkerOptions
	visit Visitor
}

func (wn nodeWalker) Walk(doc *openapi3.T) {
	wn.components([]string{"components"}, doc, doc.Components)
	// info
	// paths
	// SecurityRequirements
	// servers
	// tags
	// externaldocs
}

func (wn nodeWalker) components(path []string, parent any, c openapi3.Components) {
	if !wn.visit(path, parent, c) {
		return
	}
	wn.schemas(append(path, "schemas"), parent, c.Schemas)
	// parameters
	// headers
	// requestbodies
	// responses
	// securityschemas
	// examples
	// links
	// callbacks
}

func (wn nodeWalker) schemas(path []string, parent any, schemas openapi3.Schemas) {
	for name, schema := range schemas {
		wn.schemaRef(append(path, name), parent, schema)
	}
}

func (wn nodeWalker) schemaRefs(path []string, parent any, srefs openapi3.SchemaRefs) {
	for _, sref := range srefs {
		wn.schemaRef(path, parent, sref)
	}
}

func (wn nodeWalker) schemaRef(path []string, parent any, sref *openapi3.SchemaRef) {
	if sref == nil {
		return
	}
	if !wn.visit(path, parent, sref) {
		return
	}
	if !wn.opts.followRefs && len(sref.Ref) > 0 {
		return
	}
	wn.schemaRefs(append(path, "oneOf"), sref, sref.Value.OneOf)
	wn.schemaRefs(append(path, "anyOf"), sref, sref.Value.AnyOf)
	wn.schemaRefs(append(path, "allOf"), sref, sref.Value.AllOf)
	wn.schemaRef(append(path, "not"), sref, sref.Value.Not)
	wn.schemas(append(path, "properties"), sref, sref.Value.Properties)
	wn.schemaRef(append(path, "items"), sref, sref.Value.Items)
	wn.schemaRef(append(path, "additionalProperties"), sref, sref.Value.AdditionalProperties)
	wn.extensions(append(path, "extensions"), sref, sref.Value.Extensions)
	wn.discriminator(append(path, "discriminator"), sref, sref.Value.Discriminator)
}

func (wn nodeWalker) extensions(path []string, parent any, exts map[string]interface{}) {
	if !wn.visit(path, parent, exts) {
		return
	}
	for name, ext := range exts {
		wn.visit(append(path, name), parent, ext)
	}
}

func (wn nodeWalker) discriminator(path []string, parent any, disc *openapi3.Discriminator) {
	if disc == nil {
		return
	}
	if !wn.visit(path, parent, disc) {
		return
	}
	if !wn.visit(append(path, "mapping"), parent, disc.Mapping) {
		return
	}
	wn.extensions(path, parent, disc.Extensions)
	for name, mapping := range disc.Mapping {
		wn.visit(append(path, "mapping", name), parent, mapping)
	}
}
