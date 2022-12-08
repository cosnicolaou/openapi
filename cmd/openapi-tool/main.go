// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"

	"cloudeng.io/cmdutil/subcmd"
)

const spec = `
name: openapi
summary: command line for manipulating openapi specs
commands:
  - name: download
    summary: download an openapi spec
    arguments:
      - url
  - name: format
    summary: format an openapi spec, optionally inspecting to YAML if the schema was originally in json format.
    arguments:
      - filename
  - name: transform
    summary: transform an openapi spec using a
    arguments:
      - filename
  - name: validate
    summary: validate an openapi v3 spec.
    arguments:
      - filename
  - name: convert
    summary: convert an openapi v2 spec to v3.
    arguments:
      - filename
  - name: inspect
    summary: display the element at a path in the spec
    arguments:
      - filename
`

var cmdSet *subcmd.CommandSetYAML

func init() {
	cmdSet = subcmd.MustFromYAML(spec)
	cmdSet.Set("download").RunnerAndFlags(downloadCmd,
		subcmd.MustRegisteredFlagSet(&DownloadFlags{}))
	cmdSet.Set("format").RunnerAndFlags(formatCmd,
		subcmd.MustRegisteredFlagSet(&FormatFlags{}))
	cmdSet.Set("transform").RunnerAndFlags(transformCmd,
		subcmd.MustRegisteredFlagSet(&TransformFlags{}))
	cmdSet.Set("validate").RunnerAndFlags(validateCmd,
		subcmd.MustRegisteredFlagSet(&struct{}{}))
	cmdSet.Set("convert").RunnerAndFlags(convertCmd,
		subcmd.MustRegisteredFlagSet(&ConvertFlags{}))
	cmdSet.Set("inspect").RunnerAndFlags(inspectCmd,
		subcmd.MustRegisteredFlagSet(&InspectFlags{}))
}

func main() {
	cmdSet.MustDispatch(context.Background())
}
