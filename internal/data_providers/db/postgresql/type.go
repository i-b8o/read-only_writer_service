package postgressql

import (
	"context"
	"errors"
	"fmt"
	client "read-only_writer_service/pkg/client/postgresql"

	"github.com/jackc/pgconn"
)

type typeStorage struct {
	client client.PostgreSQLClient
}

func NewTypeStorage(client client.PostgreSQLClient) *typeStorage {
	return &typeStorage{client: client}
}

// Create returns the ID of the inserted type
func (rs *typeStorage) Create(ctx context.Context, name string) (uint64, error) {

	const sql = `INSERT INTO type ("name") VALUES ($1) ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name RETURNING "id";`

	row := rs.client.QueryRow(ctx, sql, name)
	var typeID uint64

	err := row.Scan(&typeID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return typeID, fmt.Errorf("message: %s, code: %s", pgErr.Message, pgErr.Code)
	}

	return typeID, err
}
