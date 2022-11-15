package postgressql

import (
	"context"
	"errors"
	"fmt"
	client "read-only_writer_service/pkg/client/postgresql"

	pb "github.com/i-b8o/regulations_contracts/pb/writer/v1"
	"github.com/jackc/pgconn"
)

type paragraphStorage struct {
	client client.PostgreSQLClient
}

func NewParagraphStorage(client client.PostgreSQLClient) *paragraphStorage {
	return &paragraphStorage{client: client}
}

// CreateAll
func (ps *paragraphStorage) CreateAll(ctx context.Context, paragraphs []*pb.WriterParagraph) error {
	vals := []interface{}{}
	sql := `INSERT INTO paragraphs ("paragraph_id","order_num","is_table","is_nft","has_links","class","content","c_id") VALUES `
	i := 1
	for _, p := range paragraphs {
		sql += fmt.Sprintf("($%d, $%d, $%d , $%d, $%d, $%d, $%d, $%d),", i, i+1, i+2, i+3, i+4, i+5, i+6, i+7)
		i = i + 8
		vals = append(vals, p.ID, p.Num, p.IsTable, p.IsNFT, p.HasLinks, p.Class, p.Content, p.ChapterID)
	}
	sql = sql[:len(sql)-1]

	if _, err := ps.client.Exec(ctx, sql, vals...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("message: %s, code: %s", pgErr.Message, pgErr.Code)
		}

		return err
	}

	return nil
}

func (ps *paragraphStorage) UpdateOne(ctx context.Context, content string, paragraphID uint64) error {
	sql := `UPDATE "paragraphs" SET content = $1 WHERE paragraph_id = $2 RETURNING "id"`
	row := ps.client.QueryRow(ctx, sql, content, paragraphID)
	var ID uint64

	err := row.Scan(&ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("message: %s, code: %s", pgErr.Message, pgErr.Code)
		}

		return err
	}

	return nil
}

// Delete
func (ps *paragraphStorage) DeleteForChapter(ctx context.Context, chapterID uint64) error {
	sql := `delete from paragraphs where c_id=$1`
	_, err := ps.client.Exec(ctx, sql, chapterID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return fmt.Errorf("message: %s, code: %s", pgErr.Message, pgErr.Code)
	}

	return err
}

func (ps *paragraphStorage) GetWithHrefs(ctx context.Context, chapterID uint64) ([]*pb.WriterParagraph, error) {
	const sql = `SELECT paragraph_id, content FROM "paragraphs" WHERE c_id = $1 AND has_links=true`

	var paragraphs []*pb.WriterParagraph

	rows, err := ps.client.Query(ctx, sql, chapterID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		paragraph := pb.WriterParagraph{}
		if err = rows.Scan(
			&paragraph.ID, &paragraph.Content,
		); err != nil {
			return nil, err
		}

		paragraphs = append(paragraphs, &paragraph)
	}

	return paragraphs, nil

}
