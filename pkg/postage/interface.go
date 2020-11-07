package postage

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// BatchStore stores batches by their ID as well as maintains a value-based priority index
// Its interface definitions reflect the updates triggered by events emitted by
// the postage contract on the blockchain
type BatchStore interface {
	Create(id []byte, owner []byte, amount *big.Int, depth uint8) error
	TopUp(id []byte, amount *big.Int) error
	UpdateDepth(id []byte, depth uint8) error
	UpdatePrice(price *big.Int) error
}

// Event is the interface subsuming all postage contract blockchain events
//
// postage contract event  | golang Event              | Update call on BatchStore
// ------------------------+---------------------------+---------------------------
// BatchCreated            | batchCreatedEvent         | Create
// BatchTopUp              | batchTopUpEvent           | TopUp
// BatchDepthIncrease      | batchDepthIncreaseEvent   | UpdateDepth
// PriceUpdate             | priceUpdateEvent          | UpdatePrice
type Event interface {
	Update(s BatchStore) error
}

// Listener is an event iterator
type Listener interface {
	Listen(from uint64, quit chan struct{}, update func(uint64, Event) error, addrs ...common.Address) error
}
