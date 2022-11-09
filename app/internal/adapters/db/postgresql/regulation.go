package postgressql

import (
	"context"
	"errors"
	"fmt"
	client "regulations_writable_service/pkg/client/postgresql"
	"time"

	"github.com/i-b8o/regulations_contracts/pb"
	"github.com/jackc/pgconn"
)

type regulationStorage struct {
	client client.PostgreSQLClient
}

func NewRegulationStorage(client client.PostgreSQLClient) *regulationStorage {
	return &regulationStorage{client: client}
}

// Create returns the ID of the inserted chapter
func (rs *regulationStorage) Create(ctx context.Context, regulation *pb.CreateRegulationRequest) (uint64, error) {
	t := time.Now()

	const sql = `INSERT INTO regulations ("name", "abbreviation", "title", "created_at") VALUES ($1, $2, $3, $4) RETURNING "id"`

	row := rs.client.QueryRow(ctx, sql, regulation.Name, regulation.Abbreviation, regulation.Title, t)
	var regulationID uint64

	err := row.Scan(&regulationID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return regulationID, fmt.Errorf("message: %s, code: %s", pgErr.Message, pgErr.Code)
	}

	return regulationID, err
}

// Delete
func (rs *regulationStorage) Delete(ctx context.Context, regulationID uint64) error {
	sql := `delete from regulations where id=$1`
	_, err := rs.client.Exec(ctx, sql, regulationID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return fmt.Errorf("message: %s, code: %s", pgErr.Message, pgErr.Code)
	}

	return err
}
