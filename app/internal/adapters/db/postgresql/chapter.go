package postgressql

import (
	"context"
	"errors"
	"fmt"
	client "regulations_writable_service/pkg/client/postgresql"

	"github.com/i-b8o/regulations_contracts/pb"
	"github.com/jackc/pgconn"
)

type chapterStorage struct {
	client client.PostgreSQLClient
}

func NewChapterStorage(client client.PostgreSQLClient) *chapterStorage {
	return &chapterStorage{client: client}
}

// Create returns the ID of the inserted chapter
func (cs *chapterStorage) Create(ctx context.Context, chapter *pb.CreateChapterRequest) (uint64, error) {
	sql := `INSERT INTO chapters ("name", "num", "order_num","r_id") VALUES ($1,$2,$3,$4) RETURNING "id"`

	row := cs.client.QueryRow(ctx, sql, chapter.Name, chapter.Num, chapter.OrderNum, chapter.RegulationID)

	var chapterID uint64

	err := row.Scan(&chapterID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return chapterID, fmt.Errorf("message: %s, code: %s", pgErr.Message, pgErr.Code)
	}

	return chapterID, err
}

// Delete
func (cs *chapterStorage) DeleteAllForRegulation(ctx context.Context, regulationID uint64) error {
	const sql1 = `delete from chapters where r_id=$1`
	_, err := cs.client.Exec(ctx, sql1, regulationID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return fmt.Errorf("message: %s, code: %s", pgErr.Message, pgErr.Code)
	}

	return err
}
