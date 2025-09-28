package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mibrgmv/document-service/internal/domain"
	"github.com/mibrgmv/document-service/internal/repository"
)

type documentRepository struct {
	pool *pgxpool.Pool
}

func NewDocumentRepository(pool *pgxpool.Pool) repository.DocumentRepository {
	return &documentRepository{pool: pool}
}

func (r *documentRepository) CreateDocument(ctx context.Context, doc *domain.Document) error {
	sql := `
	insert into documents (id, name, mime, file, public, created, grant_list, owner, data, json)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.pool.Exec(ctx, sql, doc.ID, doc.Name, doc.Mime, doc.File, doc.Public,
		doc.Created, doc.Grant, doc.Owner, doc.Data, doc.JSON)
	return err
}

func (r *documentRepository) GetDocumentByID(ctx context.Context, id string) (*domain.Document, error) {
	sql := `
	select
		id, name, mime, file, public,
	    created, grant_list, owner, data, json
	from documents 
	where id = $1
	`

	row := r.pool.QueryRow(ctx, sql, id)

	var doc domain.Document
	err := row.Scan(&doc.ID, &doc.Name, &doc.Mime, &doc.File, &doc.Public, &doc.Created,
		&doc.Grant, &doc.Owner, &doc.Data, &doc.JSON)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (r *documentRepository) GetUserDocuments(ctx context.Context, login string, limit int) ([]domain.Document, error) {
	sql := `
	select id, name, mime, file, public, created, grant_list
	from documents
	where (owner = $1 or $1 = any(grant_list) or public = true)
	order by name, created limit $2
	`

	rows, err := r.pool.Query(ctx, sql, login, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []domain.Document
	for rows.Next() {
		var doc domain.Document
		err := rows.Scan(&doc.ID, &doc.Name, &doc.Mime, &doc.File, &doc.Public, &doc.Created, &doc.Grant)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

func (r *documentRepository) DeleteDocument(ctx context.Context, id, owner string) error {
	sql := `
	delete from documents
    where id = $1 and owner = $2
	`

	_, err := r.pool.Exec(ctx, sql, id, owner)
	return err
}

func (r *documentRepository) DocumentExists(ctx context.Context, id string) (bool, error) {
	sql := `
	select exists(select 1 from documents where id = $1)
	`

	var exists bool
	err := r.pool.QueryRow(ctx, sql, id).Scan(&exists)
	return exists, err
}
