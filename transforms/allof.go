// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package transforms

import (
	"strings"

	"cloudeng.io/text/linewrap"
	"github.com/cosnicolaou/openapi"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

func init() {
	Register(&fixAllOf{})
}

type readOnly struct {
	Path []string
}

type merge struct {
}

type fixAllOf struct {
	IgnoreReadOnly  []string `yaml:"ignoreReadOnly"`
	MergeProperties []string `yaml:"mergeProperties"`
}

func (t *fixAllOf) Name() string {
	return "allOf"
}

func (t *fixAllOf) Configure(node yaml.Node) error {
	return node.Decode(t)
}

func (t *fixAllOf) Describe(node yaml.Node) string {
	out := &strings.Builder{}
	out.WriteString(linewrap.Block(0, 80, `
The allOf transform handles the case where an allOff list is incorrectly specified with a -properties, -required, -description block that is intended to be merged with the preceeding allOf item. For example:
`))

	out.WriteString(`
  allOff:
    - $ref: "#/components/schemas/something"
    - properties:
      url:
        example: http://example
        type string

  allOff:
    - $ref: "#/components/schemas/something"
    - required:
      - type-in-something
`)
	tmp := &fixAllOf{}
	node.Decode(tmp)
	out.WriteString("\noptions:\n")
	out.WriteString(formatYAML(2, tmp))
	return out.String()
}

func (t *fixAllOf) TransformV2(doc *openapi2.T) (*openapi2.T, error) {
	return nil, ErrTransformNotImplementedForV2
}

func mergeProperties(base, overlay map[string]*openapi3.SchemaRef) {
	for k, v := range overlay {
		base[k] = v
	}
}

func (t *fixAllOf) mergeProperties(parent string, schema *openapi3.SchemaRef) {
	for _, tm := range t.MergeProperties {
		if parent == tm {
			var prev *openapi3.SchemaRef
			na := []*openapi3.SchemaRef{}
			for _, s := range schema.Value.AllOf {
				if len(s.Value.Type) != 0 {
					na = append(na, s)
					prev = s
					continue
				}
				if prev != nil {
					mergeProperties(prev.Value.Properties, s.Value.Properties)
					// Remove any use of $ref to allow the merged properties
					// to be appear directly in the spec.
					prev.Ref = ""
				}
			}
			schema.Value.AllOf = na
		}
	}
}

func (t *fixAllOf) ignoreReadOnly(parent string, schema *openapi3.SchemaRef) {
	for _, ro := range t.IgnoreReadOnly {
		if parent == ro {
			na := []*openapi3.SchemaRef{}
			for _, s := range schema.Value.AllOf {
				if len(s.Value.Type) == 0 && s.Value.ReadOnly {
					continue
				}
				na = append(na, s)
			}
			schema.Value.AllOf = na
		}
	}
}

func (t *fixAllOf) visitor(parent any, schema *openapi3.SchemaRef) {
	if len(schema.Value.AllOf) == 0 {
		return
	}
	p, ok := parent.(string)
	if !ok {
		return
	}
	t.ignoreReadOnly(p, schema)
	t.mergeProperties(p, schema)
	return
}

func (t *fixAllOf) TransformV3(doc *openapi3.T) (*openapi3.T, error) {
	walker := openapi.NewSchemaWalker(t.visitor)
	walker.Walk(doc)
	return doc, nil
}
