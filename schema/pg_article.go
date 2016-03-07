package schema

import (
	"bytes"
	"database/sql"
)

var schema_v1 = map[string][]string{
	"article": []string{
		"root_uuid UUID",
		"document_id INTEGER",
		"meps_language CHARACTER VARYING(4) NOT NULL",
		"file_signature CHARACTER VARYING(64)",
		"ml_signature CHARACTER VARYING(64)",
	},
}

func write_create_table(table_name string, column_defs []string) string {
	var q bytes.Buffer

	q.WriteString("CREATE TABLE ")
	q.WriteString(table_name)
	q.WriteString("(")
	for i, c := range column_defs {
		q.WriteString(c)
		if i < len(column_defs)-1 {
			q.WriteString(",")
		}
	}
	q.WriteString(")")
	return q.String()
}

func create_schema_v1(db *sql.DB) error {
	for k, v := range schema_v1 {
		column_defs := append(v, k+"_uuid UUID NOT NULL")

		// Create normal table
		q := write_create_table(k, column_defs)
		_, err := db.Exec(q)
		if err != nil {
			return err
		}

		// Create conflict table
		cq := write_create_table(k+"_conflict", column_defs)
		_, err = db.Exec(cq)
		if err != nil {
			return err
		}
	}
	return nil
}

func table_exists(name string, db *sql.DB) bool {
	exists := false
	err := db.QueryRow(`
	SELECT EXISTS (
		SELECT 1
		FROM   information_schema.tables 
		WHERE  table_name=?
	)`, name).Scan(&exists)

	if err != nil {
		panic(err)
	}

	return exists
}

func get_schema_version(db *sql.DB) int {
	version := 0
	err := db.QueryRow(`
		SELECT version 
		FROM   article_schema
		LIMIT 1
	`).Scan(&version)

	if err != nil {
		panic(err)
	}

	return version
}

func article_init(db *sql.DB) {

}

var article_object_query = `
SELECT uuid,
		 root_uuid,
		 document_id,
		 meps_language,
		 file_signature,
		 ml_signature
FROM article
WHERE uuid=?
`

func (a *Article) Get(uuid string, db *sql.DB) {
	a = &Article{}
	db.QueryRow(article_object_query, uuid).Scan(
		&a.Uuid,
		&a.RootUuid,
		&a.DocumentId,
		&a.MepsLanguage,
		&a.FileSignature,
		&a.MlSignature,
	)
}
