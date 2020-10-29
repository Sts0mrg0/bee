// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package push

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethersphere/bee/pkg/file/pipeline"
	"github.com/ethersphere/bee/pkg/pushsync"
)

var errInvalidData = errors.New("store: invalid data")

type pushWriter struct {
	p    pushsync.PushSyncer
	ctx  context.Context
	next pipeline.ChainWriter
}

// NewPushSyncWriter returns a pushWriter. It writes the given data to the network
// using the PushSyncer.
func NewPushSyncWriter(ctx context.Context, p pushsync.PushSyncer, next pipeline.ChainWriter) pipeline.ChainWriter {
	return &storeWriter{ctx: ctx, p: p, next: next}
}

func (w *storeWriter) ChainWrite(p *pipeline.PipeWriteArgs) error {
	var err error
PUSH:
	_, err = w.p.PushChunkToClosest()
	if err != nil {
		fmt.Println("push err", err)
		goto PUSH
	}
	if w.next == nil {
		return nil
	}

	return w.next.ChainWrite(p)
}

func (w *storeWriter) Sum() ([]byte, error) {
	return w.next.Sum()
}