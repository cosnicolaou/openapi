// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package transforms_test

import (
	"testing"

	"github.com/cosnicolaou/openapi/transforms"
)

const allOfConfig = `
transforms:
  - name: allOf
    ignoreReadOnly:
      - readonlyEg
    mergeProperties:
      - MergeEg

`

func TestAllOf(t *testing.T) {
	doc, cfg := loadForTest("allof-eg.yaml", allOfConfig)
	tr := transforms.Get("allOf")
	if err := cfg.Configure(tr); err != nil {
		t.Fatal(err)
	}
	doc, err := tr.TransformV3(doc)
	if err != nil {
		t.Fatal(err)
	}
	txt := asYAML(t, doc)
	contains(t, txt, `type: object
properties:
type:
type: string
format: date
example: "2019-09-12"
`)
}
