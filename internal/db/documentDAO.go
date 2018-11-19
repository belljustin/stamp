package db

import (
	"bytes"
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type DocumentDAO struct {
	db *sql.DB
}

func NewDocumentDAO(conn *sql.DB) *DocumentDAO {
	return &DocumentDAO{conn}
}

type Document struct {
	Id      uuid.UUID `json:"id"`
	Hash    string    `json:"hash"`
	State   string    `json:"state"`
	StampId uuid.UUID `json:"stampId"`
}

// Insert a new document and return the corresponding UUID.
func (dao *DocumentDAO) Insert(hash string) (*uuid.UUID, error) {
	return insertDocument(dao.db, hash)
}

// Get returns a document as identified by the provided UUID.
func (dao *DocumentDAO) Get(id uuid.UUID) (*Document, error) {
	return getDocument(dao.db, id)
}

// Stamp all pending documents with the stampId.
// Returns the set of stamped documents.
func (dao *DocumentDAO) Stamp(stampId uuid.UUID) ([]Document, error) {
	return stampDocuments(dao.db, stampId, "pending", "batched")
}

// Fail fails all documents where their id is in ids.
// Returns the number of rows affected.
func (dao *DocumentDAO) Fail(ids []uuid.UUID) (int64, error) {
	// TODO: get rid of this? May not be neccessary with input side validation.
	return batchUpdateDocumentState(dao.db, ids, "batched", "failed")
}

func insertDocument(h Handle, hash string) (*uuid.UUID, error) {
	sqlStatement := `
	INSERT INTO documents (id, hash, state)
	VALUES ($1, $2, $3)
	RETURNING id`

	var id uuid.UUID
	err := h.QueryRow(sqlStatement, uuid.New(), hash, "pending").Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func getDocument(h Handle, id uuid.UUID) (*Document, error) {
	sqlStatement := `
	SELECT id, hash, state, stampId
	FROM documents
	WHERE id = $1`

	var d Document
	err := h.QueryRow(sqlStatement, id).Scan(&d.Id, &d.Hash, &d.State, &d.StampId)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func stampDocuments(h Handle, batchId uuid.UUID, oldState, newState string) ([]Document, error) {
	// TODO: should there be a limit?
	sqlStatement := `
	UPDATE documents
	SET state = $2, stampId = $3 
	WHERE state = $1 
	RETURNING id, hash, state`

	rows, err := h.Query(sqlStatement, oldState, newState, batchId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []Document
	for rows.Next() {
		var sr Document
		err = rows.Scan(&sr.Id, &sr.Hash, &sr.State)
		if err != nil {
			return nil, err
		}
		documents = append(documents, sr)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return documents, nil
}

// batchUpdateDocumentState changes documents with the provided ids from oldState to newState
func batchUpdateDocumentState(h Handle, ids []uuid.UUID, oldState, newState string) (int64, error) {
	sqlStatement := `
	UPDATE documents
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
