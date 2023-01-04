package postgressql

import (
	"context"
	"errors"
	"fmt"
	client "read-only_writer_service/pkg/client/postgresql"

	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
	"github.com/jackc/pgconn"
)

type paragraphStorage struct {
	client client.PostgreSQLClient
}

func NewParagraphStorage(client client.PostgreSQLClient) *paragraphStorage {
	return &paragraphStorage{client: client}
}

func (ps *paragraphStorage) Check(ctx context.Context, id uint64) (bool, error) {
	const sql = `select id from paragraph where id = $1 limit 1`
	var IDs []uint64
	rows, err := ps.client.Query(ctx, sql, id)
	if err != nil {
		return true, err
	}

	defer rows.Close()

	for rows.Next() {
		var ID uint64
		if err = rows.Scan(
			&ID,
		); err != nil {
			return true, err
		}

		IDs = append(IDs, ID)
	}
	return len(IDs) > 0, nil
}

// CreateAll
func (ps *paragraphStorage) CreateAll(ctx context.Context, paragraphs []*pb.WriterParagraph) error {
	vals := []interface{}{}
	sql := `INSERT INTO paragraph ("id","order_num","is_table","is_nft","has_links","class","content","c_id") VALUES `
	i := 1
	for _, p := range paragraphs {
		sql += fmt.Sprintf("($%d, $%d, $%d , $%d, $%d, $%d, $%d, $%d),", i, i+1, i+2, i+3, i+4, i+5, i+6, i+7)
		i = i + 8
		// CreateAll will be used only for internal purposes by admin
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

func (ps *paragraphStorage) UpdateOne(ctx context.Context, content string, paragraphID, chapterID uint64) error {
	if chapterID <= 0 {
		return fmt.Errorf("chapter id error")
	}
	const sql = `UPDATE "paragraph" SET content = $1 WHERE id = $2 AND c_id = $3 RETURNING "id"`
	row := ps.client.QueryRow(ctx, sql, content, paragraphID, chapterID)
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
	const sql = `delete from paragraph where c_id=$1`
	_, err := ps.client.Exec(ctx, sql, chapterID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return fmt.Errorf("message: %s, code: %s", pgErr.Message, pgErr.Code)
	}

	return err
}

// Delete
func (ps *paragraphStorage) GetOne(ctx context.Context, paragraphID, chapterID uint64) (*pb.WriterParagraph, error) {
	const sql = `SELECT content FROM paragraph WHERE id=$1 AND c_id=$2`
	row := ps.client.QueryRow(ctx, sql, paragraphID, chapterID)
	paragraph := &pb.WriterParagraph{}
	err := row.Scan(&paragraph.Content)
	if err != nil {
		return paragraph, err
	}
	return paragraph, nil

}

func (ps *paragraphStorage) GetWithHrefs(ctx context.Context, chapterID uint64) ([]*pb.WriterParagraph, error) {
	const sql = `SELECT id, c_id, content FROM "paragraph" WHERE c_id = $1 AND has_links=true`

	var paragraphs []*pb.WriterParagraph

	rows, err := ps.client.Query(ctx, sql, chapterID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		paragraph := pb.WriterParagraph{}
		if err = rows.Scan(
			&paragraph.ID, &paragraph.ChapterID, &paragraph.Content,
		); err != nil {
			return nil, err
		}

		paragraphs = append(paragraphs, &paragraph)
	}

	return paragraphs, nil

}
