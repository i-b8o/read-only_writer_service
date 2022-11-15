package postgressql

import (
	"context"
	"errors"
	"fmt"
	client "read-only_writer_service/pkg/client/postgresql"

	pb "github.com/i-b8o/regulations_contracts/pb/writer/v1"
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

// GetAllById returns all chapters associated with the given ID
func (cs *chapterStorage) GetAllById(ctx context.Context, regulationID uint64) ([]*pb.WriterChapter, error) {
	const sql = `SELECT id,name,num,order_num FROM "chapters" WHERE r_id = $1 ORDER BY order_num`

	var chapters []*pb.WriterChapter

	rows, err := cs.client.Query(ctx, sql, regulationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		chapter := pb.WriterChapter{}
		if err = rows.Scan(
			&chapter.ID, &chapter.Name, &chapter.Num, &chapter.OrderNum,
		); err != nil {
			return nil, err
		}

		chapters = append(chapters, &chapter)
	}

	return chapters, nil

}

// GetOneById returns an chapter associated with the given ID
func (cs *chapterStorage) GetRegulationId(ctx context.Context, chapterID uint64) (uint64, error) {
	const sql = `SELECT r_id FROM "chapters" WHERE id = $1`
	row := cs.client.QueryRow(ctx, sql, chapterID)
	var ID uint64
	err := row.Scan(&ID)
	if err != nil {
		return ID, err
	}

	return ID, nil
}
