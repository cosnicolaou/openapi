// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type Walker interface {
	Walk(doc *openapi3.T) error
}

type walkerOptions struct {
	followRefs  bool
	visitPrefix []string
	applyPrefix bool
}

type WalkerOption func(o *walkerOptions)

func WalkerFollowRefs(v bool) WalkerOption {
	return func(o *walkerOptions) {
		o.followRefs = v
	}
}

func WalkerVisitPrefix(path ...string) WalkerOption {
	return func(o *walkerOptions) {
		o.visitPrefix = path
		o.applyPrefix = true
	}
}

// Visitor is called for every node in the walk. It returns true for the
// walk to continue, false otherwise. The walk will stop when an error is
// returned.
type Visitor func(path []string, parent, node any) (ok bool, err error)

// NewWalker returns a Walker that will visit every node in an openapi3 document.
func NewWalker(v Visitor, opts ...WalkerOption) Walker {
	w := &nodeWalker{visitor: v}
	for _, opt := range opts {
		opt(&w.opts)
	}
	return w
}

type nodeWalker struct {
	opts    walkerOptions
	visitor Visitor
}

// returns true if b is a prefix of a
func prefixMatch(a, b []string) bool {
	if len(b) > len(a) {
		return false
	}
	for i := range b {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (wn nodeWalker) visit(path []string, parent, node any) (ok bool, err error) {
	if wn.opts.applyPrefix {
		if !prefixMatch(path, wn.opts.visitPrefix) {
			return true, nil
		}
	}
	return wn.visitor(path, parent, node)
}

func (wn nodeWalker) Walk(doc *openapi3.T) error {
	_, err := wn.components([]string{"components"}, doc, doc.Components)
	// info
	// paths
	// SecurityRequirements
	// servers
	// tags
	// externaldocs
	return err
}

func (wn nodeWalker) components(path []string, parent any, c openapi3.Components) (ok bool, err error) {
	if ok, err = wn.visit(path, parent, c); !ok || err != nil {
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
	return false, nil
}

func (wn nodeWalker) schemas(path []string, parent any, schemas openapi3.Schemas) (ok bool, err error) {
	for name, schema := range schemas {
		if ok, err = wn.schemaRef(append(path, name), parent, schema); !ok || err != nil {
			return
		}
	}
	return
}

func (wn nodeWalker) schemaRefs(path []string, parent any, srefs openapi3.SchemaRefs) (ok bool, err error) {
	for _, sref := range srefs {
		if ok, err = wn.schemaRef(path, parent, sref); !ok || err != nil {
			return ok, err
		}
	}
	return
}

func (wn nodeWalker) schemaRef(path []string, parent any, sref *openapi3.SchemaRef) (ok bool, err error) {
	if sref == nil {
		return
	}
	if ok, err = wn.visit(path, parent, sref); !ok || err != nil {
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
	return true, nil
}

func (wn nodeWalker) extensions(path []string, parent any, exts map[string]interface{}) (ok bool, err error) {
	if ok, err = wn.visit(path, parent, exts); !ok || err != nil {
		return
	}
	for name, ext := range exts {
		if ok, err = wn.visit(append(path, name), parent, ext); !ok || err != nil {
			return
		}
	}
	return
}

func (wn nodeWalker) discriminator(path []string, parent any, disc *openapi3.Discriminator) (ok bool, err error) {
	if disc == nil {
		return
	}
	if ok, err = wn.visit(path, parent, disc); !ok || err != nil {
		return
	}
	if ok, err = wn.visit(append(path, "mapping"), parent, disc.Mapping); !ok || err != nil {
		return
	}
	wn.extensions(path, parent, disc.Extensions)
	for name, mapping := range disc.Mapping {
		if ok, err = wn.visit(append(path, "mapping", name), parent, mapping); !ok || err != nil {
			return
		}
	}
	return true, nil
}
