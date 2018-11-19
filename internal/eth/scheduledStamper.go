package eth

import (
	"encoding/hex"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/belljustin/stamp/internal/db"
	"github.com/belljustin/stamp/pkg/merkle"
	"github.com/belljustin/stamp/pkg/stamper"
)

type ScheduledStamper struct {
	stamper         *stamper.Stamper
	stampRequestDAO *db.StampRequestDAO
	stampDAO        *db.StampDAO
}

func NewScheduledStamper(s *stamper.Stamper, srDAO *db.StampRequestDAO, sDAO *db.StampDAO) *ScheduledStamper {
	return &ScheduledStamper{s, srDAO, sDAO}
}

// Start beings stamping at scheduled intervals of d duration.
// The return value is a channel for sending a stop signal
func (s *ScheduledStamper) Start(d time.Duration) chan interface{} {
	ticker := time.NewTicker(d)

	c := make(chan interface{})

	// add stamp at d intervals, exit if recieve a signal on the channel
	go func() {
		for {
			select {
			case <-ticker.C:
				s.AddStamp()
			case <-c:
				return
			}
		}
	}()

	return c
}

func (s *ScheduledStamper) AddStamp() (*merkle.Tree, error) {
	stampId := uuid.New()
	stampRequests, err := s.stampDAO.New(stampId)
	if err != nil {
		log.Fatalf("Could not add new stamp: %v", err)
	}
	valHashes, invalidReqIds := validateHashes(stampRequests)

	// If there are invalid requests, fail them
	if len(invalidReqIds) > 0 {
		n, err := s.stampRequestDAO.FailRequests(invalidReqIds)
		log.Printf("Failed %v requests", n)
		if err != nil {
			log.Fatalf("Error: couldn't fail invalid requests: %v", err)
		}
	}

	// If none of the request ids validated, just return
	if len(valHashes.requestIds) == 0 {
		log.Println("No valid hashes were found")
		return nil, nil
	}

	// Build a merkle tree with the valid data
	mt := merkle.NewTree(valHashes.hashes)
	rootHash, err := hex.DecodeString(mt.Root.Hash)
	if err != nil {
		log.Fatalf("Error: merkle tree root did not contain a valid hash")
	}

	err = s.stampDAO.AddTree(stampId, mt)
	if err != nil {
		log.Fatalf("Could not add tree to db for stampId %v: %v", stampId, err)
	}

	// Add the stamp to the contract
	// stamper.AddStamp takes a fixed sized array
	var stampHash [32]byte
	copy(stampHash[:], rootHash)
	txhash, err := s.stamper.AddStamp(stampHash)
	if err != nil {
		log.Fatalf("Could not submit stampId %v to blockhain: %v", stampId, err)
	}

	err = s.stampDAO.MarkSent(stampId, txhash)
	if err != nil {
		log.Fatalf("Could not mark stampId %v as sent: %v", stampId, err)
	}

	log.Printf("Stamped %d requests", len(valHashes.requestIds))
	return mt, nil
}

type validatedHashes struct {
	requestIds []uuid.UUID
	hashes     []string
}

func (v *validatedHashes) add(requestId uuid.UUID, hash string) {
	v.requestIds = append(v.requestIds, requestId)
	v.hashes = append(v.hashes, hash)
}

func validateHashes(stampRequests []db.StampRequest) (*validatedHashes, []uuid.UUID) {
	var invalidReqIds []uuid.UUID
	valHashes := new(validatedHashes)
	for _, sr := range stampRequests {
		hash, err := hex.DecodeString(sr.Hash)
		if err != nil || len(hash) > 32 {
			log.Printf("Invalid hash for stampReqId %v", sr.Id)
			invalidReqIds = append(invalidReqIds, sr.Id)
		} else {
			valHashes.add(sr.Id, sr.Hash)
		}
	}
	return valHashes, invalidReqIds
}
