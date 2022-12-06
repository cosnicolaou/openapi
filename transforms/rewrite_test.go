// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package transforms_test

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/cosnicolaou/openapi/transforms"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

func loadForTest(filename string, config string) (*openapi3.T, transforms.Config) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	doc, err := loader.LoadFromFile(filepath.Join("testdata", filename))
	if err != nil {
		panic(err)
	}
	cfg, err := transforms.ParseConfig([]byte(config))
	if err != nil {
		panic(err)
	}
	return doc, cfg
}

const rewriteConfig = `
transforms:
  - name: rewrites
    rewrites:
      - rewrite: "/^example_replacement$/something-new/"
        replace: example
        match:
          type: string
      - rewrite: "/^integer_error$/integer/"
        replace: type
        properties:
          - end
`

func asYAML(t *testing.T, v any) string {
	out := &strings.Builder{}
	enc := yaml.NewEncoder(out)
	enc.SetIndent(1)
	if err := enc.Encode(v); err != nil {
		t.Fatal(err)
	}
	return out.String()
}

func trimSpace(s string) string {
	out := strings.Builder{}
	for _, l := range strings.Split(s, "\n") {
		out.WriteString(strings.TrimSpace(l))
		out.WriteRune('\n')
	}
	return strings.TrimSpace(out.String())
}

func contains(t *testing.T, got, want string) {
	got, want = trimSpace(got), trimSpace(want)
	if !strings.Contains(got, want) {
		_, file, line, _ := runtime.Caller(1)
		file = filepath.Base(file)
		t.Errorf("%s:%v: %v\n\ndoes not contain\n\n%v", file, line, got, want)
	}
}

func TestRewrite(t *testing.T) {
	doc, cfg := loadForTest("rewrite-eg.yaml", rewriteConfig)
	tr := transforms.Get("rewrites")
	if err := cfg.Configure(tr); err != nil {
		t.Fatal(err)
	}
	doc, err := tr.TransformV3(doc)
	if err != nil {
		t.Fatal(err)
	}
	txt := asYAML(t, doc)
	contains(t, txt, `end:
type: integer
example: example_replacement
id:
type: string
example: something-new
name:
type: string
maxLength: 255
start:
type: integer
example: example_replacement
`)
}
