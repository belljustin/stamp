package db

import (
	"bytes"
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type StampRequestDAO struct {
	db *sql.DB
}

func NewStampRequestDAO(conn *sql.DB) *StampRequestDAO {
	return &StampRequestDAO{conn}
}

type StampRequest struct {
	Id    uuid.UUID `json:"id"`
	Hash  string    `json:"hash"`
	State string    `json:"state"`
}

// Insert a new stamp request and return the corresponding UUID.
func (dao *StampRequestDAO) Insert(hash string) (*uuid.UUID, error) {
	return insertStampRequest(dao.db, hash)
}

// GetStampRequest returns a stamp request as identified by th provided UUID.
func (dao *StampRequestDAO) GetStampRequest(id uuid.UUID) (*StampRequest, error) {
	return getStampRequest(dao.db, id)
}

// Batch all pending stamp requests with the batchId.
// Returns the set of batched stamp requests.
func (dao *StampRequestDAO) Batch(batchId uuid.UUID) ([]StampRequest, error) {
	return batchStampRequests(dao.db, batchId, "pending", "batched")
}

// FailRequests fails all stamp_requests where their id is in ids.
// Returns the number of rows affected.
func (dao *StampRequestDAO) FailRequests(ids []uuid.UUID) (int64, error) {
	// TODO: get rid of this? May not be neccessary with input side validation.
	return batchUpdateStampRequestState(dao.db, ids, "batched", "failed")
}

func insertStampRequest(h Handle, hash string) (*uuid.UUID, error) {
	sqlStatement := `
	INSERT INTO stamp_requests (id, hash, state)
	VALUES ($1, $2, $3)
	RETURNING id`

	var id uuid.UUID
	err := h.QueryRow(sqlStatement, uuid.New(), hash, "pending").Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func getStampRequest(h Handle, id uuid.UUID) (*StampRequest, error) {
	sqlStatement := `
	SELECT id, hash, state
	FROM stamp_requests
	WHERE id = $1`

	var sr StampRequest
	err := h.QueryRow(sqlStatement, id).Scan(&sr.Id, &sr.Hash, &sr.State)
	if err != nil {
		return nil, err
	}
	return &sr, nil
}

func batchStampRequests(h Handle, batchId uuid.UUID, oldState, newState string) ([]StampRequest, error) {
	// TODO: should there be a limit?
	sqlStatement := `
	UPDATE stamp_requests
	SET state = $2, stampId = $3 
	WHERE state = $1 
	RETURNING id, hash, state`

	rows, err := h.Query(sqlStatement, oldState, newState, batchId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stampRequests []StampRequest
	for rows.Next() {
		var sr StampRequest
		err = rows.Scan(&sr.Id, &sr.Hash, &sr.State)
		if err != nil {
			return nil, err
		}
		stampRequests = append(stampRequests, sr)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return stampRequests, nil
}

// batchUpdateState changes requets with the provided ids from oldState to newState
func batchUpdateStampRequestState(h Handle, ids []uuid.UUID, oldState, newState string) (int64, error) {
	sqlStatement := `
	UPDATE stamp_requests
	SET state = $2 
	WHERE state = $1 `

	sqlStatement = sqlStatement + generateWhereIdsInClause(ids)
	result, err := h.Exec(sqlStatement, oldState, newState)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()

}

func generateWhereIdsInClause(ids []uuid.UUID) string {
	// TODO: there must be a cleaner way to do this....
	var whereClause bytes.Buffer
	whereClause.WriteString("AND id in ('" + ids[0].String() + "'")
	for _, id := range ids[1:] {
		whereClause.WriteString(",'" + id.String() + "'")
	}
	whereClause.WriteString(")")
	return whereClause.String()
}
