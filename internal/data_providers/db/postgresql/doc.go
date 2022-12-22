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

type docStorage struct {
	client client.PostgreSQLClient
}

func NewDocStorage(client client.PostgreSQLClient) *docStorage {
	return &docStorage{client: client}
}

// Create returns the ID of the inserted chapter
func (rs *docStorage) Create(ctx context.Context, doc *pb.CreateDocRequest) (uint64, error) {
	t := time.Now()

	const sql = `INSERT INTO doc ("name", "abbreviation", "header", "title", "description", "keywords", "created_at") VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING "id"`

	row := rs.client.QueryRow(ctx, sql, doc.Name, doc.Abbreviation, doc.Header, doc.Title, doc.Description, doc.Keywords, t)
	var docID uint64

	err := row.Scan(&docID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return docID, fmt.Errorf("message: %s, code: %s", pgErr.Message, pgErr.Code)
	}

	return docID, err
}

// Delete
func (rs *docStorage) Delete(ctx context.Context, docID uint64) error {
	const sql = `delete from doc where id=$1`
	_, err := rs.client.Exec(ctx, sql, docID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return fmt.Errorf("message: %s, code: %s", pgErr.Message, pgErr.Code)
	}

	return err
}

// GetAll
func (rs *docStorage) GetAll(ctx context.Context) (docs []*pb.WriterDoc, err error) {
	const sql = `SELECT id, name, abbreviation, title FROM "doc"`

	rows, err := rs.client.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var doc pb.WriterDoc
		if err = rows.Scan(
			&doc.ID, &doc.Name, &doc.Abbreviation, &doc.Title,
		); err != nil {
			return nil, err
		}

		docs = append(docs, &doc)
	}

	return docs, nil
}
