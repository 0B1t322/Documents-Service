package pgql

import (
	"context"
	"github.com/0B1t322/Documents-Service/document/internal/core/models"
	repository "github.com/0B1t322/Documents-Service/document/internal/repository/styles"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Repository {
	return &Repository{conn: conn}
}

func (r Repository) StoreStyleInDocument(ctx context.Context, documentId uuid.UUID, style *models.Style) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := storeParagraphStyle(ctx, tx, &style.ParagraphStyle); err != nil {
		return err
	}

	if err := storeTextStyle(ctx, tx, &style.TextStyle); err != nil {
		return err
	}

	if err := storeStyle(ctx, tx, style); err != nil {
		return err
	}

	if err := storeStyleInDocument(ctx, tx, documentId, style); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func storeParagraphStyle(ctx context.Context, tx pgx.Tx, style *models.ParagraphStyle) error {
	if err := tx.QueryRow(
		ctx,
		`
			INSERT INTO "ParagraphStyle" 
			    ("Aligment", "LineSpacing")
			VALUES 
			($1,$2)
			RETURNING "Id"`,
		style.Alignment,
		style.LineSpacing,
	).Scan(&style.ID); err != nil {
		return err
	}

	return nil
}

func storeTextStyle(ctx context.Context, tx pgx.Tx, style *models.TextStyle) error {
	if err := tx.QueryRow(
		ctx,
		`
			INSERT INTO "TextStyle"
				("FontFamily", "FontWeight", "FontSize","Bold", "Underline", "Italic", "BackgroundColor", "ForegroundColor")
			VALUES 
				($1,$2,$3,$4,$5,$6,$7,$8)
			RETURNING "Id"`,
		style.FontFamily,
		style.FontWeight,
		style.FontSize,
		style.Bold,
		style.Underline,
		style.Italic,
		style.BackgroundColor,
		style.ForegroundColor,
	).Scan(&style.ID); err != nil {
		return err
	}

	return nil
}

func isStyleWithThisNameExist(err error) bool {
	if pgErr, ok := errors.Into[*pgconn.PgError](err); ok {
		if pgErr.ConstraintName == "Styles_Name_uindex" {
			return true
		}
	}

	return false
}

func storeStyle(ctx context.Context, tx pgx.Tx, style *models.Style) error {
	err := tx.QueryRow(
		ctx,
		`
			INSERT INTO "Styles" 
			    ("Name", "ParagraphStyleId", "TextStyleId") 
			VALUES 
				($1,$2,$3)
			RETURNING "Id"`,
		style.Name,
		style.ParagraphStyle.ID,
		style.TextStyle.ID,
	).Scan(&style.ID)
	if isStyleWithThisNameExist(err) {
		return repository.ErrStyleExist
	} else if err != nil {
		return err
	}

	return nil
}

func storeStyleInDocument(ctx context.Context, tx pgx.Tx, documentId uuid.UUID, style *models.Style) error {
	if _, err := tx.Exec(
		ctx,
		`
			INSERT INTO "DocumentsStyles" 
			    ("DocumentId", "StyleId") 
			VALUES
				($1,$2)`,
		documentId,
		style.ID,
	); err != nil {
		return err
	}

	return nil
}

func (r Repository) FindStyleInDocumentByID(
	ctx context.Context,
	documentId uuid.UUID,
	styleId uuid.UUID,
) (models.Style, error) {
	row := r.conn.QueryRow(
		ctx,
		`
			WITH styleInDocument AS (
				SELECT
					"StyleId"
				FROM "DocumentsStyles"
				WHERE
					"DocumentId" = $1
				  AND "StyleId" = $2
				)
			SELECT
				s."Id",
				"Name",
				TS."Id",
				TS."FontFamily",
				TS."FontWeight",
				TS."FontSize",
				TS."Bold",
				TS."Underline",
				TS."Italic",
				TS."BackgroundColor",
				TS."ForegroundColor",
				PS."Id",
				PS."Aligment",
				PS."LineSpacing"
			FROM "Styles" S
				JOIN "TextStyle" TS on TS."Id" = S."TextStyleId"
				JOIN "ParagraphStyle" PS on PS."Id" = S."ParagraphStyleId"
			WHERE S."Id" = (SELECT "StyleId" from styleInDocument)`,
		documentId,
		styleId,
	)

	style, err := scanStyle(row)
	if err == pgx.ErrNoRows {
		return models.Style{}, repository.ErrStyleNotFound
	} else if err != nil {
		return models.Style{}, err
	}

	return style, nil
}

func scanStyle(s pgx.Row) (style models.Style, err error) {
	err = s.Scan(
		&style.ID, &style.Name,
		&style.TextStyle.ID, &style.TextStyle.FontFamily,
		&style.TextStyle.FontWeight, &style.TextStyle.FontSize,
		&style.TextStyle.Bold, &style.TextStyle.Underline,
		&style.TextStyle.Italic, &style.TextStyle.BackgroundColor,
		&style.TextStyle.BackgroundColor,
		&style.ParagraphStyle.ID, &style.ParagraphStyle.Alignment,
		&style.ParagraphStyle.LineSpacing,
	)
	return
}

func (r Repository) UpdateStyle(ctx context.Context, style models.Style) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update style
	if err := updateStyle(ctx, tx, style); err != nil {
		return err
	}
	// Update paragraph style
	if err := updateParagraphStyle(ctx, tx, style.ParagraphStyle); err != nil {
		return err
	}
	// Update text style
	if err := updateTextStyle(ctx, tx, style.TextStyle); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func updateStyle(ctx context.Context, tx pgx.Tx, style models.Style) error {
	_, err := tx.Exec(
		ctx,
		`
			UPDATE "Styles" 
			SET "Name" = $2
			WHERE "Id" = $1`,
		style.ID,
		style.Name,
	)
	if isStyleWithThisNameExist(err) {
		return repository.ErrStyleExist
	} else if err != nil {
		return err
	}

	return nil
}

func updateParagraphStyle(ctx context.Context, tx pgx.Tx, style models.ParagraphStyle) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE "ParagraphStyle"
			SET 
			    "Aligment" = $2,
				"LineSpacing" = $3
			WHERE "Id" = $1`,
		style.ID,
		style.Alignment,
		style.LineSpacing,
	); err != nil {
		return err
	}

	return nil
}

func updateTextStyle(ctx context.Context, tx pgx.Tx, style models.TextStyle) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE "TextStyle"
			SET 
			    "Bold" = $2,
				"Underline" = $3,
				"BackgroundColor" = $4,
				"ForegroundColor" = $5,
				"FontSize" = $6,
				"Italic" = $7,
				"FontWeight" = $8,
				"FontFamily" = $9
			WHERE "Id" = $1`,
		style.ID,
		style.Bold,
		style.Underline,
		style.BackgroundColor,
		style.ForegroundColor,
		style.FontSize,
		style.Italic,
		style.FontWeight,
		style.FontFamily,
	); err != nil {
		return err
	}

	return nil
}

func (r Repository) GetAllStylesInDocument(ctx context.Context, documentId uuid.UUID) ([]models.Style, error) {
	rows, err := r.conn.Query(
		ctx,
		`
			WITH styleInDocument AS (
				SELECT
					"StyleId"
				FROM "DocumentsStyles"
				WHERE
					"DocumentId" = $1
				)
			SELECT
				s."Id",
				"Name",
				TS."Id",
				TS."FontFamily",
				TS."FontWeight",
				TS."FontSize",
				TS."Bold",
				TS."Underline",
				TS."Italic",
				TS."BackgroundColor",
				TS."ForegroundColor",
				PS."Id",
				PS."Aligment",
				PS."LineSpacing"
			FROM "Styles" S
				JOIN "TextStyle" TS on TS."Id" = S."TextStyleId"
				JOIN "ParagraphStyle" PS on PS."Id" = S."ParagraphStyleId"
			WHERE S."Id" IN (SELECT "StyleId" from styleInDocument)`,
		documentId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var styles []models.Style
	for rows.Next() {
		if style, err := scanStyle(rows); err == nil {
			styles = append(styles, style)
		} else {
			return nil, err
		}
	}

	return styles, nil
}

func (r Repository) DeleteStyleInDocument(ctx context.Context, documentId uuid.UUID, style models.Style) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := deleteStyleInDocument(ctx, tx, style, documentId); err != nil {
		return err
	}

	if err := deleteStyle(ctx, tx, style); err != nil {
		return err
	}

	if err := deleteTextStyle(ctx, tx, style.TextStyle); err != nil {
		return err
	}

	if err := deleteParagraphStyle(ctx, tx, style.ParagraphStyle); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func deleteTextStyle(ctx context.Context, tx pgx.Tx, style models.TextStyle) error {
	if _, err := tx.Exec(
		ctx,
		`
			DELETE FROM "TextStyle" WHERE "Id" = $1`,
		style.ID,
	); err != nil {
		return err
	}

	return nil
}

func deleteParagraphStyle(ctx context.Context, tx pgx.Tx, style models.ParagraphStyle) error {
	if _, err := tx.Exec(
		ctx,
		`
			DELETE FROM "ParagraphStyle" WHERE "Id" = $1`,
		style.ID,
	); err != nil {
		return err
	}

	return nil
}

func deleteStyle(ctx context.Context, tx pgx.Tx, style models.Style) error {
	if _, err := tx.Exec(
		ctx,
		`
			DELETE FROM "Styles" WHERE "Id" = $1`,
		style.ID,
	); err != nil {
		return err
	}

	return nil
}

func deleteStyleInDocument(ctx context.Context, tx pgx.Tx, style models.Style, documentId uuid.UUID) error {
	if _, err := tx.Exec(
		ctx,
		`
			DELETE FROM "DocumentsStyles" 
			       WHERE 
			           "Id" = $1 AND
			           "DocumentId" = $2`,
		style.ID,
		documentId,
	); err != nil {
		return err
	}

	return nil
}
