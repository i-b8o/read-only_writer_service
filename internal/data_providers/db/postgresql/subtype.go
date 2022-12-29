package postgressql

import (
	"context"
	"errors"
	"fmt"
	client "read-only_writer_service/pkg/client/postgresql"

	"github.com/jackc/pgconn"
)

type subTypeStorage struct {
	client client.PostgreSQLClient
}

func NewSubTypeStorage(client client.PostgreSQLClient) *subTypeStorage {
	return &subTypeStorage{client: client}
}

// Create returns the ID of the inserted type
func (rs *subTypeStorage) Create(ctx context.Context, name string, typeID uint64) (uint64, error) {

	const sql = `INSERT INTO subtype ("name", "type_id") VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name RETURNING "id"`

	row := rs.client.QueryRow(ctx, sql, name, typeID)
	var subTypeID uint64

	err := row.Scan(&subTypeID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return subTypeID, fmt.Errorf("message: %s, code: %s", pgErr.Message, pgErr.Code)
	}

	return subTypeID, err
}
