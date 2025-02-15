// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

//go:build !integration
// +build !integration

package beatcmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/elastic-agent-libs/config"
)

func TestLocker(t *testing.T) {
	initCfgfile(t, `output.console.enabled: true`)

	nopBeater := &nopBeater{
		running: make(chan struct{}),
		done:    make(chan struct{}),
	}
	beat1, err := NewBeat(BeatParams{
		Create: func(*beat.Beat, *config.C) (beat.Beater, error) {
			return nopBeater, nil
		},
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var g errgroup.Group
	g.Go(func() error { return beat1.Run(ctx) })

	// Wait for the first beater to be running, at which point
	// the lock should be held.
	<-nopBeater.running

	beat2, err := NewBeat(BeatParams{
		Create: func(*beat.Beat, *config.C) (beat.Beater, error) {
			panic("should not be called")
		},
	})
	require.NoError(t, err)
	err = beat2.Run(ctx)
	require.ErrorIs(t, err, ErrAlreadyLocked)

	cancel()
	assert.NoError(t, g.Wait())
}

type nopBeater struct {
	running chan struct{}
	done    chan struct{}
}

func (b *nopBeater) Run(*beat.Beat) error {
	close(b.running)
	<-b.done
	return nil
}

func (b *nopBeater) Stop() {
	close(b.done)
}
