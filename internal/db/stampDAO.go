package db

import (
	"database/sql"
	"encoding/json"

	"github.com/belljustin/stamp/pkg/merkle"

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

// AddTree transitions the stamp to the processing state and adds a merkle tree.
func (dao *StampDAO) AddTree(id uuid.UUID, merkleTree *merkle.Tree) error {
	return addStampTree(dao.db, id, merkleTree)
}

// MarkSent transitions the stamp to the sent state.
func (dao *StampDAO) MarkSent(id uuid.UUID) error {
	return updateStampState(dao.db, id, "sent")
}

func newStamp(h Handle, id uuid.UUID) error {
	sqlStatement := `
	INSERT INTO stamps (id, state)
	VALUES ($1, $2)`

	_, err := h.Exec(sqlStatement, id, "pending")
	return err
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

func updateStampState(h Handle, id uuid.UUID, newState string) error {
	sqlStatement := `
	UPDATE stamps
	SET state = $1
	WHERE id = $2`

	_, err := h.Exec(sqlStatement, newState, id)
	return err
}
