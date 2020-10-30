package post

import "github.com/ethersphere/bee/pkg/swarm"

type Stamper interface {
	Stamp(swarm.Chunk) Stamp
}

type Stamp interface {
	BatchId() []byte
	Signature() []byte
}
