package stamper

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// Stamper allows interaction with the StampStorageContract
type Stamper struct {
	stampStorage *StampStorage
	auth         *bind.TransactOpts
}

// NewStamper creates a Stamper from a reference to the
// StampStorageContract and related authorization
func NewStamper(ss *StampStorage, auth *bind.TransactOpts) *Stamper {
	return &Stamper{ss, auth}
}

// AddStamp submits a new stamp to the StampStorageContract
func (s *Stamper) AddStamp(hash [32]byte) (common.Hash, error) {
	tx, err := s.stampStorage.AddStamp(s.auth, hash)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

func (s *Stamper) LogStamps() {
	sink := make(chan *StampStorageStamped)
	sub, _ := s.stampStorage.WatchStamped(nil, sink)

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case stamp := <-sink:
			hash := hex.EncodeToString(stamp.Hash[:])
			fmt.Printf("%v, %v", stamp.Timestamp, hash)
		}
	}
}
