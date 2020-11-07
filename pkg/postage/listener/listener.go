package listener

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethersphere/bee/pkg/postage"
)

var _ postage.Listener = (*ethListener)(nil)

type ethListener struct {
	client         *ethclient.Client
	logToEventFunc func(types.Log) postage.Event
}

func (lis *ethListener) Listen(from uint64, quit chan struct{}, f func(uint64, postage.Event) error, addrs ...common.Address) error {
	// need to cancel context even if terminate without error with quit channel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// using an event channel, we subscribe to filter logs emmitted by the postage contract(s)
	// since the last block recorded that had a relevant event
	events := make(chan types.Log)
	sub, err := lis.client.SubscribeFilterLogs(ctx, query(from, addrs...), events)
	if err != nil {
		return err
	}
LOOP:
	for {
		select {
		case err = <-sub.Err(): // subscription error
			break LOOP
		case <-quit: // normal quit
			break LOOP
		case ev := <-events:
			// read and parse log into event and call the function
			// supplying the blocknumber of the event and the event as arguments
			// if this call returns an error the listen loop terminates
			err = f(ev.BlockNumber, lis.logToEventFunc(ev))
			if err != nil {
				break LOOP
			}
		}
	}
	return err
}

func query(from uint64, addrs ...common.Address) ethereum.FilterQuery {
	return ethereum.FilterQuery{
		Addresses: addrs,
		FromBlock: big.NewInt(0).SetUint64(from),
	}
}
