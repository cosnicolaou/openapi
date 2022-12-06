// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package transforms

import (
	"encoding/json"
	"fmt"
	"strings"

	"cloudeng.io/errors"
	"cloudeng.io/text/linewrap"
	"github.com/cosnicolaou/openapi"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

func init() {
	Register(&rewrites{})
}

type rewrite struct {
	Replace    string
	Rewrite    string
	Match      map[string]string
	Properties []string
	repl       Replacement
}

type rewrites struct {
	Rewrites []rewrite
}

func (t *rewrites) configure(rewrites []rewrite) ([]rewrite, error) {
	for i, rw := range rewrites {
		repl, err := NewReplacement(rw.Rewrite)
		if err != nil {
			return nil, err
		}
		rewrites[i].repl = repl
	}
	return rewrites, nil
}

func (t *rewrites) Name() string {
	return "rewrites"
}

func (t *rewrites) Configure(node yaml.Node) error {
	if err := node.Decode(t); err != nil {
		return err
	}
	var errs errors.M
	var err error
	t.Rewrites, err = t.configure(t.Rewrites)
	errs.Append(err)
	return errs.Err()
}

func (t *rewrites) Describe(node yaml.Node) string {
	out := &strings.Builder{}
	fmt.Fprintf(out, linewrap.Block(0, 80, `
The rewrites transform rewriting of fields using expressions of the form "/regexp/replacement/"

The rewrites may be contrained by specifying a context .....

`))
	tmp := &fixEnums{}
	node.Decode(tmp)
	out.WriteString("\noptions:\n")
	out.WriteString(formatYAML(2, tmp))
	return out.String()
}

func (t *rewrites) TransformV2(doc *openapi2.T) (*openapi2.T, error) {
	return nil, ErrTransformNotImplementedForV2
}

func (t *rewrites) TransformV3(doc *openapi3.T) (*openapi3.T, error) {
	walker := openapi.NewSchemaWalker(t.visitor)
	walker.Walk(doc)
	return doc, nil
}

func jsonMap(v any) map[string]any {
	var r map[string]any
	buf, _ := json.Marshal(v)
	json.Unmarshal(buf, &r)
	return r
}

func marshalMap(in map[string]any, out any) error {
	buf, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, out)
}

func jsonSlice(v any) []any {
	var r []any
	buf, _ := json.Marshal(v)
	json.Unmarshal(buf, &r)
	return r
}

func (t *rewrites) matchFields(rw rewrite, srv map[string]any) bool {
	for k, v := range rw.Match {
		m, ok := srv[k]
		if !ok || m != v {
			return false
		}
	}
	return true
}

func (t *rewrites) matchProperties(rw rewrite, props []any) bool {
	for _, v := range rw.Properties {
		for _, prop := range props {
			p, ok := prop.(string)
			if !ok {
				return false
			}
			if p != v {
				return false
			}
		}
	}
	return true
}

func (t *rewrites) visitor(parent any, sr *openapi3.SchemaRef) {
	v := sr.Value

	properties := jsonSlice(parent)
	schemaFields := jsonMap(v)

	for _, rw := range t.Rewrites {
		if len(rw.Replace) == 0 {
			continue
		}
		// Check that the field to be rewritten exists in the current
		// Schema.
		if _, ok := schemaFields[rw.Replace]; !ok {
			continue
		}
		if _, ok := schemaFields[rw.Replace].(string); !ok {
			continue
		}
		if !t.matchFields(rw, schemaFields) {
			continue
		}
		if !t.matchProperties(rw, properties) {
			continue
		}
		v := schemaFields[rw.Replace].(string)
		if !rw.repl.MatchString(v) {
			continue
		}
		var nsv openapi3.Schema
		ov := schemaFields[rw.Replace]
		schemaFields[rw.Replace] = rw.repl.ReplaceAllString(v)
		if err := marshalMap(schemaFields, &nsv); err != nil {
			schemaFields[rw.Replace] = ov
			continue
		}
		sr.Value = &nsv
	}

}

/*
func findField(v any, field string) (string, bool) {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return "", false
	}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		fmt.Printf("%v == %v\n", field, f.Tag)
		if f.Name == field {
			val := reflect.ValueOf(v).Field(i)
			if val.Kind() == reflect.String {
				return val.String(), true
			}
		}
	}
	return "", false
}
*/

/*
	for _, rw := range t.Rewrites {
		if fv, ok := findField(v, rw.Field); ok {
			fmt.Printf("%v: %v\n", rw.Field, fv)
		}
		eg, ok := v.Example.(string)
		if !ok {
			return
		}
		if rw.repl.MatchString(eg) && t.match(parent, v, rw.Context) {
			v.Example = rw.repl.ReplaceAllString(eg)
		}
	}
*/
/*
func (t *rewrites) matchFields(sv *openapi3.Schema, key string, val any) bool {
	switch v := val.(type) {
	case string:
		switch key {
		case "type":
			return sv.Type == v
		case "format":
			return sv.Format == v
		}
	}
	return false
}

func (t *rewrites) matchProperties(parent any, val any) bool {
	id, ok := parent.(string)
	if !ok {
		return false
	}
	vals, ok := val.([]interface{})
	if ok {
		for _, av := range vals {
			v, ok := av.(string)
			if !ok {
				return false
			}
			if id == v {
				return true
			}
		}
	}
	return false
}

func (t *rewrites) match(parent any, sv *openapi3.Schema, context map[string]any) bool {
	psr, ok := parent.(string)
	for k, v := range context {
		switch k {
		case "type", "format":
			if !t.matchFields(sv, k, v) {
				return false
			}
		case "properties":
			if !ok {
				return false
			}
			if !t.matchProperties(psr, v) {
				return false
			}
		}
	}
	return true
}
*/
