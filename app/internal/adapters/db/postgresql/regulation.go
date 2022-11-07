package postgressql

import (
	"context"
	"regulations_writable_service/internal/pb"
	client "regulations_writable_service/pkg/client/postgresql"
	"time"
)

type regulationStorage struct {
	client client.PostgreSQLClient
}

func NewRegulationStorage(client client.PostgreSQLClient) *regulationStorage {
	return &regulationStorage{client: client}
}

// Create returns the ID of the inserted chapter
func (rs *regulationStorage) Create(ctx context.Context, regulation *pb.Regulation) (uint64, error) {
	t := time.Now()

	const sql = `INSERT INTO regulations ("name", "abbreviation", "title", "created_at") VALUES ($1, $2, $3, $4) RETURNING "id"`

	row := rs.client.QueryRow(ctx, sql, regulation.Name, regulation.Abbreviation, regulation.Title, t)
	var regulationID uint64

	err := row.Scan(&regulationID)

	return regulationID, err
}

// Delete
func (rs *regulationStorage) Delete(ctx context.Context, regulationID uint64) error {
	sql := `delete from regulations where id=$1`
	_, err := rs.client.Exec(ctx, sql, regulationID)

	return err
}
