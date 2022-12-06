// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package transforms_test

import (
	"testing"

	"github.com/cosnicolaou/openapi/transforms"
)

const enumConfig = `
transforms:
  - name: enums
    flatten_single_enum:
      - match: date
        type: string
        format: date
        example: "2019-09-12"
`

func TestEnums(t *testing.T) {
	doc, cfg := loadForTest("enum-eg.yaml", enumConfig)
	tr := transforms.Get("enums")
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
