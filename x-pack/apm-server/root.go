// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package main

import (
	"github.com/spf13/cobra"

	"github.com/elastic/apm-server/internal/beatcmd"
	"github.com/elastic/beats/v7/libbeat/beat"
	_ "github.com/elastic/beats/v7/x-pack/libbeat/management" // Fleet
)

// newXPackRootCommand returns the Elastic licensed "apm-server" root command.
func newXPackRootCommand(newBeat beat.Creator) *cobra.Command {
	return beatcmd.NewRootCommand(beatcmd.BeatParams{
		Create:          newBeat,
		ElasticLicensed: true,
	})
}
