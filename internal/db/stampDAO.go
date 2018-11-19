package db

import (
	"database/sql"
	"encoding/json"

	"github.com/belljustin/stamp/pkg/merkle"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type StampDAO struct {
	db *sql.DB
}

func NewStampDAO(conn *sql.DB) *StampDAO {
	return &StampDAO{conn}
}

type Stamp struct {
	Id         uuid.UUID    `json:"id"`
	TxHash     string       `json:"txhash"`
	MerkleTree *merkle.Tree `json:"merkleTree"`
	State      string       `json:"state"`
}

// New batches all pending stamp request into a new stamp with the provided stampId.
func (dao *StampDAO) New(stampId uuid.UUID) ([]StampRequest, error) {
	tx, err := dao.db.Begin()
	if err != nil {
		return nil, err
	}

	err = newStamp(tx, stampId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	ids, err := batchStampRequests(tx, stampId, "pending", "batched")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(ids) == 0 {
		tx.Rollback()
		return nil, nil
	}
	err = tx.Commit()
	return ids, err
}

// Get a stamp with the provided id
func (dao *StampDAO) Get(id uuid.UUID) (*Stamp, error) {
	return getStamp(dao.db, id)
}

// AddTree transitions the stamp to the processing state and adds a merkle tree.
func (dao *StampDAO) AddTree(id uuid.UUID, merkleTree *merkle.Tree) error {
	return addStampTree(dao.db, id, merkleTree)
}

// MarkSent transitions the stamp to the sent state and set the txhash.
func (dao *StampDAO) MarkSent(id uuid.UUID, txhash common.Hash) error {
	return updateStampState(dao.db, id, txhash, "sent")
}

func newStamp(h Handle, id uuid.UUID) error {
	sqlStatement := `
	INSERT INTO stamps (id, state)
	VALUES ($1, $2)`

	_, err := h.Exec(sqlStatement, id, "pending")
	return err
}

func getStamp(h Handle, id uuid.UUID) (*Stamp, error) {
	sqlStatement := `
	SELECT id, txhash, merkletree, state FROM stamps
	WHERE id = $1`

	var s Stamp
	var bMT []byte
	err := h.QueryRow(sqlStatement, id).Scan(&s.Id, &s.TxHash, &bMT, &s.State)
	if err != nil {
		return nil, err
	}

	var mt merkle.Tree
	if err = json.Unmarshal(bMT, &mt); err != nil {
		return nil, err
	}
	s.MerkleTree = &mt
	return &s, nil
}

func addStampTree(h Handle, id uuid.UUID, merkleTree *merkle.Tree) error {
	sqlStatement := `
	UPDATE stamps
	SET merkleTree = $1, state = $2
	WHERE id = $3`

	jMerkleTree, err := json.Marshal(merkleTree)
	if err != nil {
		return err
	}

	_, err = h.Exec(sqlStatement, jMerkleTree, "processing", id)
	return err
}

func updateStampState(h Handle, id uuid.UUID, txhash common.Hash, newState string) error {
	sqlStatement := `
	UPDATE stamps
	SET state = $1, txhash = $2
	WHERE id = $3`

	_, err := h.Exec(sqlStatement, newState, txhash.Hex(), id)
	return err
}
