// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package openapi

import (
	"bytes"
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

func FormatV3(doc *openapi3.T, isYAML bool) ([]byte, error) {
	data, err := doc.MarshalJSON()
	if !isYAML {
		return data, err
	}
	var tmp any
	if err := json.Unmarshal(data, &tmp); err != nil {
		return nil, err
	}
	out := &bytes.Buffer{}
	enc := yaml.NewEncoder(out)
	enc.SetIndent(2)
	err = enc.Encode(tmp)
	return out.Bytes(), err
}
