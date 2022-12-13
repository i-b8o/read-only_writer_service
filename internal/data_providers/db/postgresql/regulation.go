package postgressql

import (
	"context"
	"errors"
	"fmt"
	client "read-only_writer_service/pkg/client/postgresql"
	"time"

	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
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

	const sql = `INSERT INTO regulation ("name", "abbreviation", "title", "created_at") VALUES ($1, $2, $3, $4) RETURNING "id"`

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
	const sql = `delete from regulation where id=$1`
	_, err := rs.client.Exec(ctx, sql, regulationID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return fmt.Errorf("message: %s, code: %s", pgErr.Message, pgErr.Code)
	}

	return err
}

// GetAll
func (rs *regulationStorage) GetAll(ctx context.Context) (regulations []*pb.WriterRegulation, err error) {
	const sql = `SELECT id, name, abbreviation, title FROM "regulation"`

	rows, err := rs.client.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var regulation pb.WriterRegulation
		if err = rows.Scan(
			&regulation.ID, &regulation.Name, &regulation.Abbreviation, &regulation.Title,
		); err != nil {
			return nil, err
		}

		regulations = append(regulations, &regulation)
	}

	return regulations, nil
}
