// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mock

import (
	"context"
	"sync"

	"github.com/ethersphere/bee/pkg/swarm"
	"github.com/ethersphere/bee/pkg/topology"
)

type mock struct {
	peers           []swarm.Address
	closestPeer     swarm.Address
	closestPeerErr  error
	addPeersErr     error
	marshalJSONFunc func() ([]byte, error)
	mtx             sync.Mutex
}

func WithAddPeersErr(err error) Option {
	return optionFunc(func(d *mock) {
		d.addPeersErr = err
	})
}

func WithClosestPeer(addr swarm.Address) Option {
	return optionFunc(func(d *mock) {
		d.closestPeer = addr
	})
}

func WithClosestPeerErr(err error) Option {
	return optionFunc(func(d *mock) {
		d.closestPeerErr = err
	})
}

func WithMarshalJSONFunc(f func() ([]byte, error)) Option {
	return optionFunc(func(d *mock) {
		d.marshalJSONFunc = f
	})
}

func NewTopologyDriver(opts ...Option) topology.Driver {
	d := new(mock)
	for _, o := range opts {
		o.apply(d)
	}
	return d
}

func (d *mock) AddPeers(_ context.Context, addrs ...swarm.Address) error {
	if d.addPeersErr != nil {
		return d.addPeersErr
	}

	for _, addr := range addrs {
		d.mtx.Lock()
		d.peers = append(d.peers, addr)
		d.mtx.Unlock()
	}

	return nil
}
func (d *mock) Connected(ctx context.Context, addr swarm.Address) error {
	return d.AddPeers(ctx, addr)
}

func (d *mock) Disconnected(swarm.Address) {
	panic("todo")
}

func (d *mock) Peers() []swarm.Address {
	return d.peers
}

func (d *mock) ClosestPeer(addr swarm.Address) (peerAddr swarm.Address, err error) {
	return d.closestPeer, d.closestPeerErr
}

func (d *mock) SubscribePeersChange() (c <-chan struct{}, unsubscribe func()) {
	return c, unsubscribe
}

func (_ *mock) NeighborhoodDepth() uint8 {
	return 0
}

// EachPeer iterates from closest bin to farthest
func (_ *mock) EachPeer(_ topology.EachPeerFunc) error {
	panic("not implemented") // TODO: Implement
}

// EachPeerRev iterates from farthest bin to closest
func (_ *mock) EachPeerRev(_ topology.EachPeerFunc) error {
	panic("not implemented") // TODO: Implement
}

func (d *mock) MarshalJSON() ([]byte, error) {
	return d.marshalJSONFunc()
}

func (d *mock) Close() error {
	return nil
}

type Option interface {
	apply(*mock)
}

type optionFunc func(*mock)

func (f optionFunc) apply(r *mock) { f(r) }
