package postgresql

import (
	"context"
	"github.com/0B1t322/Online-Document-Redactor/document/internal/core/models"
	"github.com/0B1t322/Online-Document-Redactor/document/internal/repository/documents"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Repository {
	return &Repository{conn: conn}
}

func (r Repository) Store(
	ctx context.Context,
	document *models.Document,
) error {
	// Create documents
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}

	// Create document style
	if err := r.storeDocumentStyle(ctx, tx, &document.Style); err != nil {
		tx.Rollback(ctx)
		return err
	}

	// create body of document
	if err := r.storeBody(ctx, tx, &document.Body); err != nil {
		tx.Rollback(ctx)
		return err
	}

	// Create documents
	if err := r.storeDocument(ctx, tx, document); err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r Repository) storeDocumentStyle(ctx context.Context, tx pgx.Tx, style *models.DocumentStyle) error {
	row := tx.QueryRow(
		ctx,
		`INSERT INTO "DocumentStyles" ("PageSize") VALUES ($1) RETURNING "Id"`,
		style.PageSize,
	)
	if err := row.Scan(&style.ID); err != nil {
		return err
	}

	return nil
}

func (r Repository) storeBody(ctx context.Context, tx pgx.Tx, body *models.Body) error {
	row := tx.QueryRow(
		ctx,
		`INSERT INTO "Body" DEFAULT VALUES RETURNING "Id"`,
	)

	if err := row.Scan(&body.ID); err != nil {
		return err
	}

	return nil
}

func (r Repository) storeDocument(ctx context.Context, tx pgx.Tx, document *models.Document) error {
	row := tx.QueryRow(
		ctx,
		`INSERT INTO "Documents" ("Title", "BodyId", "DocumentStyleId") 
				VALUES ($1,$2,$3) RETURNING "Id"`,
		document.Title,
		document.Body.ID,
		document.Style.ID,
	)

	if err := row.Scan(&document.ID); err != nil {
		return err
	}

	return nil
}

func (r Repository) FindByID(ctx context.Context, documentId uuid.UUID) (models.Document, error) {
	row := r.conn.QueryRow(
		ctx,
		`SELECT d."Id", d."Title", b."Id", ds."Id", ds."PageSize"
			FROM "Documents" d
    		JOIN "Body" b on b."Id" = d."BodyId"
    		JOIN "DocumentStyles" ds on ds."Id" = d."DocumentStyleId"
    		WHERE d."Id" = $1`,
		documentId,
	)

	var document models.Document

	if err := row.Scan(
		&document.ID,
		&document.Title,
		&document.Body.ID,
		&document.Style.ID,
		&document.Style.PageSize,
	); err == pgx.ErrNoRows {
		return document, repository.ErrNotFound
	} else if err != nil {
		return document, err
	}

	return document, nil
}

func (r Repository) Update(ctx context.Context, document models.Document) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}

	if err := r.updateDocument(ctx, tx, document); err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err := r.updateDocumentStyle(ctx, tx, document.Style); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (r Repository) updateDocument(ctx context.Context, tx pgx.Tx, document models.Document) error {
	_, err := tx.Exec(
		ctx,
		`UPDATE "Documents" SET "Title" = $2
			WHERE "Id" = $1`,
		document.ID,
		document.Title,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) updateDocumentStyle(ctx context.Context, tx pgx.Tx, style models.DocumentStyle) error {
	_, err := tx.Exec(
		ctx,
		`Update "DocumentStyles" SET "PageSize" = $2
			WHERE "Id" = $1`,
		style.ID,
		style.PageSize,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) Get(
	ctx context.Context,
	cursor uuid.UUID,
	limit uint64,
) ([]models.Document, uuid.UUID, error) {
	rows, err := r.conn.Query(
		ctx,
		`SELECT d."Id", d."Title" 
			FROM "Documents" d
			WHERE d."Id" > $1
			ORDER BY d."Id" ASC 
			LIMIT $2`,
		cursor,
		limit,
	)
	if err != nil {
		return nil, uuid.Nil, err
	}
	defer rows.Conn()

	var documents []models.Document

	for rows.Next() {
		d := models.Document{}
		if err := rows.Scan(&d.ID, &d.Title); err != nil {
			return nil, uuid.Nil, err
		}

		documents = append(documents, d)
	}
	var nextCursor = uuid.Nil
	if len(documents) > 0 {
		nextCursor = documents[len(documents)-1].ID
	}

	return documents, nextCursor, nil
}
