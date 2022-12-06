// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type Visitor func(parent any, child *openapi3.SchemaRef)

type Walker interface {
	Walk(doc *openapi3.T)
}

// NewSchemaWalker returns a Walker that will visit every instance of
// a SchemaRef in an openapi3 document.
func NewSchemaWalker(v Visitor) Walker {
	return schemaWalker{v}
}

type schemaWalker struct{ visit Visitor }

func (ws schemaWalker) Walk(doc *openapi3.T) {
	ws.paths(doc.Paths)
	ws.responses(doc.Components.Responses)
	for n, sr := range doc.Components.Schemas {
		ws.schemaRef(n, sr)
	}
}

func (ws schemaWalker) paths(paths openapi3.Paths) {
	for n, path := range paths {
		ws.parameters(n, path.Parameters)
		ws.operations(path.Connect, path.Delete, path.Get, path.Head,
			path.Options, path.Patch, path.Post, path.Put, path.Trace)
	}
}

func (ws schemaWalker) operations(operations ...*openapi3.Operation) {
	for _, op := range operations {
		if op == nil {
			continue
		}
		ws.parameters(op, op.Parameters)
		ws.responses(op.Responses)
	}
}

func (ws schemaWalker) parameters(parent any, parameters openapi3.Parameters) {
	for _, parameter := range parameters {
		if parameter.Value == nil {
			continue
		}
		ws.schemaRef(parent, parameter.Value.Schema)
	}
}

func (ws schemaWalker) parametersMap(parameters openapi3.ParametersMap) {
	for n, parameter := range parameters {
		if parameter.Value == nil {
			continue
		}
		ws.schemaRef(n, parameter.Value.Schema)
	}
}

func (ws schemaWalker) responses(responses openapi3.Responses) {
	for _, response := range responses {
		if response.Value == nil {
			continue
		}
		for _, media := range response.Value.Content {
			ws.schemaRef(response.Value, media.Schema)
		}
	}
}

func (ws schemaWalker) schemaRef(parent any, ref *openapi3.SchemaRef) {
	ws.visit(parent, ref)
	for _, sr := range ref.Value.AnyOf {
		ws.schemaRef(ref, sr)
	}
	for _, sr := range ref.Value.AllOf {
		ws.schemaRef(ref, sr)
	}
	for _, sr := range ref.Value.OneOf {
		ws.schemaRef(ref, sr)
	}
	for id, sr := range ref.Value.Properties {
		ws.schemaRef(id, sr)
	}
	if sr := ref.Value.Items; sr != nil {
		ws.schemaRef(ref, sr)
	}
}
