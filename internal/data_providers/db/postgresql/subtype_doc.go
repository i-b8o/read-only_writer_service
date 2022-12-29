package postgressql

import (
	"context"
	"errors"
	"fmt"
	client "read-only_writer_service/pkg/client/postgresql"

	"github.com/jackc/pgconn"
)

type subTypeDocStorage struct {
	client client.PostgreSQLClient
}

func NewSubTypeDocStorage(client client.PostgreSQLClient) *subTypeDocStorage {
	return &subTypeDocStorage{client: client}
}

// Create returns the ID of the inserted type
func (s *subTypeDocStorage) Create(ctx context.Context, subTypeID, docID uint64) error {
	const sql = `INSERT INTO subtype_doc ("subtype_id", "doc_id") VALUES ($1, $2)`
	_, err := s.client.Exec(ctx, sql, subTypeID, docID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return fmt.Errorf("message: %s, code: %s", pgErr.Message, pgErr.Code)
	}

	return err
}
