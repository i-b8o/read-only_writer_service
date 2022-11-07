package postgressql

import (
	"context"
	"regulations_writable_service/internal/pb"
	client "regulations_writable_service/pkg/client/postgresql"
)

type chapterStorage struct {
	client client.PostgreSQLClient
}

func NewChapterStorage(client client.PostgreSQLClient) *chapterStorage {
	return &chapterStorage{client: client}
}

// Create returns the ID of the inserted chapter
func (cs *chapterStorage) Create(ctx context.Context, chapter *pb.Chapter) (uint64, error) {
	sql := `INSERT INTO chapters ("name", "num", "order_num","r_id") VALUES ($1,$2,$3,$4) RETURNING "id"`

	row := cs.client.QueryRow(ctx, sql, chapter.Name, chapter.Num, chapter.OrderNum, chapter.RegulationID)

	var chapterID uint64

	err := row.Scan(&chapterID)

	return chapterID, err
}

// Delete
func (cs *chapterStorage) DeleteAllForRegulation(ctx context.Context, regulationID uint64) error {
	const sql1 = `delete from chapters where r_id=$1`
	_, err := cs.client.Exec(ctx, sql1, regulationID)
	return err
}
