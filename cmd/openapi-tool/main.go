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
    summary: download an openapi specification
    arguments:
      - url
  - name: format
    summary: format an openapi spec.
    arguments:
      - filename
  - name: transform
    summary: edit an openapi v3 specification using a set of built in transformations
             and rewrites.
    arguments:
      - filename
  - name: validate
    summary: validate an openapi v3 spec.
    arguments:
      - filename
  - name: convert
    summary: convert an openapi v2 specification to v3.
    arguments:
      - filename
  - name: inspect
    summary: display the element at a path in an openapi v3 specification
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
