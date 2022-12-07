package paragraph_controller

import (
	"context"
	"fmt"
	"log"
	"read-only_writer_service/pkg/client/postgresql"
	"testing"
	"time"

	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

const (
	dbHost     = "0.0.0.0"
	dbPort     = "5436"
	dbUser     = "reader"
	dbPassword = "postgres"
	dbName     = "reader"
)

func setupDB() *pgxpool.Pool {
	pgConfig := postgresql.NewPgConfig(
		dbUser, dbPassword,
		dbHost, dbPort, dbName,
	)

	pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	if err != nil {
		log.Fatal(err)
	}

	return pgClient
}

func TestCreateAll(t *testing.T) {
	assert := assert.New(t)
	pgClient := setupDB()
	defer pgClient.Close()
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", "0.0.0.0", "30001"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewWriterParagraphGRPCClient(conn)
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		input *pb.CreateAllParagraphsRequest
		err   error
	}{
		{
			input: &pb.CreateAllParagraphsRequest{Paragraphs: []*pb.WriterParagraph{&pb.WriterParagraph{ID: 4, Num: 4, HasLinks: false, IsTable: false, IsNFT: false, Class: "class", Content: "Содержимое четвертого параграфа", ChapterID: 3}, &pb.WriterParagraph{ID: 5, Num: 5, HasLinks: true, IsTable: true, IsNFT: true, Class: "class", Content: "Содержимое пятого параграфа", ChapterID: 3}}},
			err:   nil,
		},
	}

	for _, test := range tests {
		_, err := client.CreateAll(ctx, test.input)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(test.err, err, err)
		for _, p := range test.input.Paragraphs {
			sql := fmt.Sprintf("select paragraph_id,order_num,is_table,is_nft,has_links,class,content,c_id from paragraph where id=%d", p.ID)
			rows, err := pgClient.Query(ctx, sql)
			if err != nil {
				t.Log(err)
			}
			defer rows.Close()

			var pId, cId uint64
			var orderNum uint32
			var isTable, isNft, hasLinks bool
			var class, content string
			for rows.Next() {
				if err = rows.Scan(
					&pId, &orderNum, &isTable, &isNft, &hasLinks, &class, &content, &cId,
				); err != nil {
					t.Log(err)
				}
			}
			assert.Equal(p.ID, pId)
			assert.Equal(p.Num, orderNum)
			assert.Equal(p.ChapterID, cId)
			assert.Equal(p.Class, class)
			assert.Equal(p.Content, content)
			assert.Equal(p.HasLinks, hasLinks)
			assert.Equal(p.IsTable, isTable)
			assert.Equal(p.IsNFT, isNft)
		}

	}

	_, err = pgClient.Exec(ctx, resetDB)
	if err != nil {
		t.Log(err)
	}
}

func TestGetOne(t *testing.T) {
	assert := assert.New(t)
	pgClient := setupDB()
	defer pgClient.Close()
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", "0.0.0.0", "30001"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewWriterParagraphGRPCClient(conn)
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		input    *pb.GetOneParagraphRequest
		expected *pb.GetOneParagraphResponse
		err      error
	}{
		{
			input:    &pb.GetOneParagraphRequest{ID: 1},
			expected: &pb.GetOneParagraphResponse{Content: "Содержимое <a id=\"dst101675\"></a> первого <a href='11111/a3a3a3/111'>параграфа</a>"},
			err:      nil,
		},
	}

	for _, test := range tests {
		resp, err := client.GetOne(ctx, test.input)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(test.err, err)
		assert.True(proto.Equal(test.expected, resp))

	}
	_, err = pgClient.Exec(ctx, resetDB)
	if err != nil {
		t.Log(err)
	}
}

// UpdateOneParagraph
func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	pgClient := setupDB()
	defer pgClient.Close()
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", "0.0.0.0", "30001"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewWriterParagraphGRPCClient(conn)
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		input *pb.UpdateOneParagraphRequest
		err   error
	}{
		{
			input: &pb.UpdateOneParagraphRequest{ID: 3, Content: "Измененное содержимое третьего параграфа"},
			err:   nil,
		},
	}

	for _, test := range tests {
		_, err := client.Update(ctx, test.input)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(test.err, err, err)
		sql := fmt.Sprintf("select content from paragraph where id=%d", test.input.ID)
		rows, err := pgClient.Query(ctx, sql)
		if err != nil {
			t.Log(err)
		}
		defer rows.Close()
		var content string
		for rows.Next() {
			if err = rows.Scan(
				&content,
			); err != nil {
				t.Log(err)
			}
		}
		assert.Equal(test.input.Content, content)

	}

	_, err = pgClient.Exec(ctx, resetDB)
	if err != nil {
		t.Log(err)
	}
}

// GetParagraphsWithHrefs
func TestGetWithHrefs(t *testing.T) {
	assert := assert.New(t)
	pgClient := setupDB()
	defer pgClient.Close()
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", "0.0.0.0", "30001"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewWriterParagraphGRPCClient(conn)
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		input    *pb.GetParagraphsWithHrefsRequest
		expected *pb.GetParagraphsWithHrefsResponse
		err      error
	}{
		{
			input:    &pb.GetParagraphsWithHrefsRequest{ID: 1},
			expected: &pb.GetParagraphsWithHrefsResponse{Paragraphs: []*pb.WriterParagraph{&pb.WriterParagraph{ID: 1, Content: "Содержимое <a id=\"dst101675\"></a> первого <a href='11111/a3a3a3/111'>параграфа</a>"}, &pb.WriterParagraph{ID: 2, Content: "Содержимое второго <a href='372952/4e92c731969781306ebd1095867d2385f83ac7af/335104'>пункта 5.14</a> параграфа"}, &pb.WriterParagraph{ID: 3, Content: "<a id='335050'></a>Содержимое третьего параграфа<a href='/document/cons_doc_LAW_2875/'>таблицей N 2</a>."}}},
			err:      nil,
		},
	}

	for _, test := range tests {
		e, err := client.GetWithHrefs(ctx, test.input)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(test.err, err, err)

		for i, t := range test.expected.Paragraphs {
			assert.Equal(t.ID, e.Paragraphs[i].ID, i)
			assert.Equal(t.Content, e.Paragraphs[i].Content, i)
			assert.Equal(t.HasLinks, e.Paragraphs[i].HasLinks, i)

		}

	}

	_, err = pgClient.Exec(ctx, resetDB)
	if err != nil {
		t.Log(err)
	}
}

const resetDB = `
DROP TABLE IF EXISTS absent_reg;
DROP TABLE IF EXISTS pseudo_chapter;
DROP TABLE IF EXISTS pseudo_regulation;
DROP TABLE IF EXISTS link;
DROP MATERIALIZED VIEW IF EXISTS reg_search;
DROP INDEX IF EXISTS idx_search;
DROP TABLE IF EXISTS paragraph;
DROP TABLE IF EXISTS chapter;
DROP TABLE IF EXISTS regulation;


CREATE TABLE regulation (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL CHECK (NAME != '') UNIQUE,
    abbreviation TEXT,
    title TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE regulation ADD COLUMN ts tsvector GENERATED ALWAYS AS (setweight(to_tsvector('russian', coalesce(name, '')), 'A') || setweight(to_tsvector('russian', coalesce(title, '')), 'B')) STORED;
CREATE INDEX reg_ts_idx ON regulation USING GIN (ts);



CREATE TABLE chapter (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL CHECK (name != ''),
    order_num SMALLINT NOT NULL CHECK (order_num >= 0),
    num TEXT,
    r_id integer REFERENCES regulation,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE chapter ADD COLUMN ts tsvector GENERATED ALWAYS AS (to_tsvector('russian', name)) STORED;
CREATE INDEX ch_ts_idx ON chapter USING GIN (ts);

CREATE TABLE paragraph (
    id SERIAL PRIMARY KEY,
    paragraph_id INT NOT NULL CHECK (paragraph_id >= 0),
    order_num INT NOT NULL CHECK (order_num >= 0),
    is_table BOOLEAN NOT NULL,
    is_nft BOOLEAN NOT NULL,
    has_links BOOLEAN NOT NULL,
    class TEXT,
    content TEXT NOT NULL,
    c_id integer REFERENCES chapter
);

ALTER TABLE paragraph ADD COLUMN ts tsvector GENERATED ALWAYS AS (to_tsvector('russian', content)) STORED;
CREATE INDEX p_ts_idx ON paragraph USING GIN (ts);


CREATE MATERIALIZED VIEW reg_search 
AS SELECT 
r.id AS "r_id", r.name AS "r_name", NULL AS "c_id", NULL AS "c_name", CAST(NULL AS integer) AS "p_id", NULL AS "p_text", r.name AS "text",
to_tsvector('russian', r.name) AS ts FROM regulation AS r UNION SELECT 
NULL AS "r_id", r.name AS "r_name", c.id AS "c_id", c.name AS "c_name", NULL AS "p_id", NULL AS "p_text", c.name AS "text",
to_tsvector('russian', c.name) AS ts FROM chapter AS c INNER JOIN regulation AS r ON r.id= c.r_id
UNION SELECT 
NULL AS "r_id", r.name AS "r_name", c.id AS "c_id", c.name AS "c_name", p.paragraph_id AS "p_id", p.content AS "p_text", p.content AS "text",
to_tsvector('russian', content) AS ts 
FROM paragraph AS p INNER JOIN chapter AS c ON p.c_id= c.id INNER JOIN regulation AS r ON c.r_id = r.id;

create index idx_search on reg_search using GIN(ts);

CREATE TABLE pseudo_regulation (
    r_id integer,
    pseudo TEXT NOT NULL CHECK (pseudo != '')
);

CREATE TABLE pseudo_chapter (
    c_id integer,
    pseudo TEXT NOT NULL CHECK (pseudo != '')
);

CREATE TABLE absent_reg (
    id SERIAL PRIMARY KEY,
    pseudo TEXT NOT NULL CHECK (pseudo != ''),
    done BOOLEAN NOT NULL DEFAULT false,
    paragraph_id integer  
);

CREATE TABLE link (
    id INT NOT NULL UNIQUE,
    paragraph_num INT NOT NULL CHECK (paragraph_num >= 0),
    c_id integer,
    r_id integer
);

INSERT INTO regulation ("name", "abbreviation", "title", "created_at") VALUES ('Имя первой записи', 'Аббревиатура первой записи', 'Заголовок первой записи', '2023-01-01 00:00:00');
INSERT INTO chapter ("name", "num", "order_num","r_id", "updated_at") VALUES ('Имя первой записи','I',1,1, '2023-01-01 00:00:00'), ('Имя второй записи','II',2,1, '2023-01-01 00:00:00'), ('Имя третьей записи','III',3,1, '2023-01-01 00:00:00');
INSERT INTO paragraph ("paragraph_id","order_num","is_table","is_nft","has_links","class","content","c_id") VALUES (1,1,false,false,true,'any-class','Содержимое <a id="dst101675"></a> первого <a href=''11111/a3a3a3/111''>параграфа</a>', 1), (2,2,true,true,true,'any-class','Содержимое второго <a href=''372952/4e92c731969781306ebd1095867d2385f83ac7af/335104''>пункта 5.14</a> параграфа', 1), (3,3,false,false,true,'any-class','<a id=''335050''></a>Содержимое третьего параграфа<a href=''/document/cons_doc_LAW_2875/''>таблицей N 2</a>.', 1);
INSERT INTO pseudo_regulation ("r_id", "pseudo") VALUES (1, 11111);
INSERT INTO pseudo_chapter ("c_id", "pseudo") VALUES (3, 'a3a3a3');
INSERT INTO absent_reg ("pseudo", "done", "paragraph_id") VALUES ('aaaaa', false, 1), ('bbbbb', true, 2), ('ccccc', false, 3);
`
