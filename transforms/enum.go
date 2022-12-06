// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package transforms

import (
	"fmt"
	"strings"

	"cloudeng.io/text/linewrap"
	"github.com/cosnicolaou/openapi"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

func init() {
	Register(&fixEnums{})
}

type singleEnumRewrite struct {
	Match   string
	Type    string
	Format  string
	Example string
}

type fixEnums struct {
	FlattenSingleEnum []singleEnumRewrite `yaml:"flatten_single_enum"`
}

func (t *fixEnums) Name() string {
	return "enums"
}

func (t *fixEnums) Configure(node yaml.Node) error {
	return node.Decode(t)
}

func (t *fixEnums) Describe(node yaml.Node) string {
	out := &strings.Builder{}
	fmt.Fprintf(out, linewrap.Block(0, 80, `
The enums transform handles various cases where enums are incorrectly specified`))
	tmp := &fixEnums{}
	node.Decode(tmp)
	out.WriteString("\noptions:\n")
	out.WriteString(formatYAML(2, tmp))
	return out.String()
}

func (t *fixEnums) TransformV2(doc *openapi2.T) (*openapi2.T, error) {
	return nil, ErrTransformNotImplementedForV2
}

func (t *fixEnums) TransformV3(doc *openapi3.T) (*openapi3.T, error) {
	walker := openapi.NewSchemaWalker(t.visitor)
	walker.Walk(doc)
	return doc, nil
}

func (t *fixEnums) visitor(parent any, sr *openapi3.SchemaRef) {
	v := sr.Value
	if len(v.Enum) == 1 && len(t.FlattenSingleEnum) > 0 {
		ev := v.Enum[0]
		for _, rw := range t.FlattenSingleEnum {
			str, ok := ev.(string)
			if !ok {
				continue
			}
			if str == rw.Match {
				v.Enum = nil
				v.Type = rw.Type
				v.Format = rw.Format
				v.Example = rw.Example
				return
			}
		}
	}
}
