package pgql

import (
	"context"
	"fmt"
	"github.com/0B1t322/Documents-Service/document/internal/core/models"
	repository "github.com/0B1t322/Documents-Service/document/internal/repository/elements"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
)

type Repository struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Repository {
	return &Repository{conn: conn}
}

func (r Repository) StoreStructuralElement(
	ctx context.Context,
	bodyId uuid.UUID,
	element *models.StructuralElement,
) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Store child element
	switch element.GetType() {
	case models.SEParagraph:
		if err := storeParagraph(ctx, tx, element); err != nil {
			return err
		}
	case models.SESectionBreak:
		if err := storeSectionBreak(ctx, tx, element); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown structural element type")
	}

	// Store structural element
	if err := storeStructuralElement(ctx, tx, element, bodyId); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func storeParagraph(ctx context.Context, tx pgx.Tx, element *models.StructuralElement) error {
	row := tx.QueryRow(
		ctx,
		`INSERT INTO "Paragraphs" DEFAULT VALUES RETURNING "Id"`,
	)
	if err := row.Scan(&element.Paragraph.ID); err != nil {
		return err
	}

	return nil
}

func storeSectionBreak(ctx context.Context, tx pgx.Tx, element *models.StructuralElement) error {
	row := tx.QueryRow(
		ctx,
		`INSERT INTO "SectionBreak" DEFAULT VALUES RETURNING "Id"`,
	)
	if err := row.Scan(&element.SectionBreak.ID); err != nil {
		return err
	}

	return nil
}

// isNeedReindex check is need reindex elements and can we insert them into this position
// if err != nil: we can't insert
func isNeedReindex(insertPosition int, maxIndex int, elementsCount int) (bool, error) {
	// Check that can create element with this index
	insertAsLastElement := elementsCount != 0 && insertPosition-maxIndex == 1
	insertAsFirstElement := insertPosition == 0
	insertBetween := elementsCount > 0 && maxIndex > insertPosition

	if !insertAsLastElement && !insertAsFirstElement && !insertBetween {
		return false, repository.ErrBadIndex
	}

	needReindex := (elementsCount > 0 && insertAsFirstElement) || insertBetween

	return needReindex, nil
}

func storeStructuralElement(
	ctx context.Context, tx pgx.Tx, element *models.StructuralElement,
	bodyId uuid.UUID,
) error {
	if element.Index < 0 {
		return repository.ErrBadIndex
	}

	// get max index and counts of elements
	var (
		// Can be nil
		maxIndex      = new(int)
		elementsCount int
	)
	if err := tx.QueryRow(
		ctx,
		`SELECT max("Index"), count(1) FROM "StructuralElements" WHERE "BodyId" = $1`,
		bodyId,
	).Scan(maxIndex, &elementsCount); err != nil {
		return err
	}

	needReindex, err := isNeedReindex(element.Index, lo.FromPtr(maxIndex), elementsCount)
	if err != nil {
		return err
	}

	// Reindex
	if needReindex {
		if _, err := tx.Exec(
			ctx,
			`
			UPDATE "StructuralElements" 
			SET "Index" = "Index" + 1 
			WHERE 
			    "BodyId" = $1 AND
			    "Index" >= $2`,
			bodyId,
			element.Index,
		); err != nil {
			return err
		}
	}

	var query string
	if t := element.GetType(); t == models.SEParagraph {
		query = `
				INSERT INTO "StructuralElements" ("BodyId", "ParagraphId", "Index") 
				VALUES ($1, $2, $3) 
				RETURNING "Id"`
	} else if t == models.SESectionBreak {
		query = `
				INSERT INTO "StructuralElements" ("BodyId", "SectionBreakId", "Index") 
				VALUES ($1, $2, $3) 
				RETURNING "Id"`
	}

	// insert as element with next index
	if err := tx.QueryRow(ctx, query, bodyId, element.GetElementID(), element.Index).Scan(&element.ID); err != nil {
		return err
	}

	return nil
}

func (r Repository) FindStructuralElementByID(
	ctx context.Context,
	seId int,
) (models.StructuralElement, error) {
	var (
		id          int
		elementType models.SEType
		elementId   int
		styleId     int
		index       int
	)
	{
		if err := r.conn.QueryRow(
			ctx,
			`
			SELECT SE."Id",
					CASE
						WHEN "SectionBreakId" IS NULL THEN 1
						ELSE 2
					END                                       "ElementType",
					coalesce("ParagraphId", "SectionBreakId") "ElementId",
					coalesce("ParagraphStyleId", "SectionStyleId", 0) "StyleId",
					"Index" 
			FROM "StructuralElements" SE
			LEFT OUTER JOIN "Paragraphs" P on SE."ParagraphId" = P."Id"
			LEFT OUTER JOIN "SectionBreak" SB on SB."Id" = SE."SectionBreakId"
			WHERE SE."Id" = $1`,
			seId,
		).Scan(&id, &elementType, &elementId, &styleId, &index); err == pgx.ErrNoRows {
			return models.StructuralElement{}, repository.ErrStructuralElementNotFound
		} else if err != nil {
			return models.StructuralElement{}, err
		}
	}

	element := scanStructuralElement(id, elementType, elementId, styleId, index)

	return element, nil
}

func (r Repository) FindStructuralElementByIDAndBodyID(
	ctx context.Context,
	seId int,
	bodyId uuid.UUID,
) (models.StructuralElement, error) {
	var (
		id          int
		elementType models.SEType
		elementId   int
		styleId     int
		index       int
	)
	{
		if err := r.conn.QueryRow(
			ctx,
			`
			SELECT SE."Id",
					CASE
						WHEN "SectionBreakId" IS NULL THEN 1
						ELSE 2
					END                                       "ElementType",
					coalesce("ParagraphId", "SectionBreakId") "ElementId",
					coalesce("ParagraphStyleId", "SectionStyleId", 0) "StyleId",
					"Index" 
			FROM "StructuralElements" SE
			LEFT OUTER JOIN "Paragraphs" P on SE."ParagraphId" = P."Id"
			LEFT OUTER JOIN "SectionBreak" SB on SB."Id" = SE."SectionBreakId"
			WHERE SE."Id" = $1 AND 
			      SE."BodyId" = $2`,
			seId,
			bodyId,
		).Scan(&id, &elementType, &elementId, &styleId, &index); err == pgx.ErrNoRows {
			return models.StructuralElement{}, repository.ErrStructuralElementNotFound
		} else if err != nil {
			return models.StructuralElement{}, err
		}
	}

	element := scanStructuralElement(id, elementType, elementId, styleId, index)

	return element, nil
}

func (r Repository) GetStructuralElements(
	ctx context.Context,
	bodyId uuid.UUID,
	cursor int,
	limit uint,
) (slice []models.StructuralElement, next int, err error) {
	rows, err := r.conn.Query(
		ctx,
		`
			SELECT  SE."Id",
					CASE
						WHEN "SectionBreakId" IS NULL THEN 1
						ELSE 2
					END                                       "ElementType",
					coalesce("ParagraphId", "SectionBreakId") "ElementId",
					coalesce("ParagraphStyleId", "SectionStyleId", 0) "StyleId",
					"Index"
			FROM "StructuralElements" as SE
			LEFT OUTER JOIN "Paragraphs" P on SE."ParagraphId" = P."Id"
			LEFT OUTER JOIN "SectionBreak" SB on SB."Id" = SE."SectionBreakId"
			WHERE 
			    "BodyId" = $1 AND
			    "Index" > $2
			ORDER BY "Index"
			LIMIT $3`,
		bodyId,
		cursor,
		Limit(limit),
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id          int
			elementType models.SEType
			elementId   int
			styleId     int
			index       int
		)
		if err := rows.Scan(&id, &elementType, &elementId, &styleId, &index); err != nil {
			return nil, 0, err
		}

		element := scanStructuralElement(id, elementType, elementId, styleId, index)

		slice = append(slice, element)
	}
	if count := len(slice); count > 0 {
		next = slice[count-1].Index
	}
	return
}

func scanStructuralElement(
	id int, elementType models.SEType, elementId int, styleId int,
	index int,
) models.StructuralElement {
	element := models.StructuralElement{
		ID:    id,
		Index: index,
	}
	switch elementType {
	case models.SEParagraph:
		element.Paragraph = &models.Paragraph{
			ID:               elementId,
			ParagraphStyleId: styleId,
		}
	case models.SESectionBreak:
		element.SectionBreak = &models.SectionBreak{
			ID:                  elementId,
			SectionBreakStyleId: styleId,
		}
	}

	return element
}

func (r Repository) DeleteStructuralElement(
	ctx context.Context,
	element models.StructuralElement,
) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := deleteStructuralElement(ctx, tx, element); err != nil {
		return err
	}

	if err := deleteChildElement(ctx, tx, element); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func deleteChildElement(ctx context.Context, tx pgx.Tx, from models.StructuralElement) error {
	switch from.GetType() {
	case models.SEParagraph:
		return deleteParagraph(ctx, tx, from.Paragraph)
	case models.SESectionBreak:
		return deleteSectionBreak(ctx, tx, from.SectionBreak)
	}

	return fmt.Errorf("Unknown child element type")
}

func deleteParagraph(ctx context.Context, tx pgx.Tx, paragraph *models.Paragraph) error {
	// TODO: Delete paragraphs elements
	if _, err := tx.Exec(
		ctx,
		`
			DELETE FROM "Paragraphs" WHERE "Id" = $1`,
		paragraph.ID,
	); err != nil {
		return err
	}

	return nil
}

func deleteSectionBreak(ctx context.Context, tx pgx.Tx, sectionBreak *models.SectionBreak) error {
	if _, err := tx.Exec(
		ctx,
		`
		DELETE FROM "SectionBreak" WHERE "Id" = $1`,
		sectionBreak.ID,
	); err != nil {
		return err
	}

	return nil
}

func deleteStructuralElement(ctx context.Context, tx pgx.Tx, element models.StructuralElement) error {
	var bodyId uuid.UUID
	{
		if err := tx.QueryRow(
			ctx,
			`
			DELETE FROM "StructuralElements" 
		   	WHERE "Id" = $1
		   	RETURNING "BodyId"`,
			element.ID,
		).Scan(&bodyId); err != nil {
			return err
		}
	}

	// Reindex structural elements
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE "StructuralElements" 
			SET "Index" = "Index" - 1
			WHERE 
			    "BodyId" = $1 AND
			    "Index" > $2`,
		bodyId,
		element.Index,
	); err != nil {
		return err
	}

	return nil
}

func (r Repository) UpdateStructuralElement(ctx context.Context, element models.StructuralElement) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := updateStyleOfChildElement(ctx, tx, element); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func updateStyleOfChildElement(ctx context.Context, tx pgx.Tx, element models.StructuralElement) error {
	switch element.GetType() {
	case models.SEParagraph:
		return updateStyleOfParagraph(ctx, tx, element.Paragraph)
	case models.SESectionBreak:
		return updateStyleOfSectionBreak(ctx, tx, element.SectionBreak)
	}

	return repository.ErrSEBadType
}

func updateStyleOfParagraph(ctx context.Context, tx pgx.Tx, paragraph *models.Paragraph) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE "Paragraphs" 
			SET "ParagraphStyleId" = $1 
			WHERE "Id" = $2`,
		paragraph.ParagraphStyleId,
		paragraph.ID,
	); err != nil {
		return err
	}

	return nil
}

func updateStyleOfSectionBreak(ctx context.Context, tx pgx.Tx, sectionBreak *models.SectionBreak) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE "SectionBreak" 
			SET "SectionStyleId" = $1 
			WHERE "Id" = $2`,
		sectionBreak.SectionBreakStyleId,
		sectionBreak.ID,
	); err != nil {
		return err
	}

	return nil
}

func (r Repository) StoreParagraphElement(
	ctx context.Context,
	paragraphID int,
	element *models.ParagraphElement,
) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	switch element.GetType() {
	case models.PETextRune:
		if err := storeTextRun(ctx, tx, element); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown paragraph element type")
	}

	if err := storeParagraphElement(ctx, tx, element, paragraphID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func storeTextRun(ctx context.Context, tx pgx.Tx, element *models.ParagraphElement) error {
	row := tx.QueryRow(
		ctx,
		`
			INSERT INTO "TextRun" 
			    ("Content", "TextStyleId") 
			VALUES ($1,$2)
			RETURNING "Id"`,
		[]byte(element.TextRune.Content),
		element.TextRune.TextStyleID,
	)

	if err := row.Scan(&element.TextRune.ID); err != nil {
		return err
	}

	return nil
}

func storeParagraphElement(
	ctx context.Context, tx pgx.Tx, element *models.ParagraphElement,
	paragraphID int,
) error {
	if element.Index < 0 {
		return repository.ErrBadIndex
	}

	// get max index and counts of elements
	var (
		// Can be nil
		maxIndex      int
		elementsCount int
	)
	{
		if err := tx.QueryRow(
			ctx,
			`SELECT coalesce(max("Index"), 0), count(1) FROM "ParagraphElements" WHERE "ParagraphId" = $1`,
			paragraphID,
		).Scan(&maxIndex, &elementsCount); err != nil {
			return err
		}
	}

	needReindex, err := isNeedReindex(element.Index, maxIndex, elementsCount)
	if err != nil {
		return err
	}

	if needReindex {
		if _, err := tx.Exec(
			ctx,
			`
			UPDATE "ParagraphElements" 
			SET "Index" = "Index" + 1 
			WHERE 
			    "ParagraphId" = $1 AND
			    "Index" >= $2`,
			paragraphID,
			element.Index,
		); err != nil {
			return err
		}
	}

	var childField string
	{
		switch element.GetType() {
		case models.PETextRune:
			childField = `"TextRunId"`
		default:
			return fmt.Errorf("Can't store paragraph element: uknown element type")
		}
	}
	sql := fmt.Sprintf(
		`
			INSERT INTO "ParagraphElements" 
			    ("ParagraphId", "Index", %s)
			    VALUES ($1,$2, $3)
			    RETURNING "Id"`,
		childField,
	)

	row := tx.QueryRow(
		ctx,
		sql,
		paragraphID,
		element.Index,
		element.GetChildElementID(),
	)
	if err := row.Scan(&element.ID); err != nil {
		return err
	}

	return nil
}

func (r Repository) GetParagraphElements(
	ctx context.Context,
	paragraphId int,
	cursor int,
	limit uint,
) (slice []models.ParagraphElement, nextCursor int, err error) {
	rows, err := r.conn.Query(
		ctx,
		`
			SELECT
				PE."Id",
				coalesce("TextRunId", "InlineObjectElementId", "PageBreakId", "EquationId") "ElementId",
				CASE
					WHEN "TextRunId" IS NOT NULL THEN 1
					WHEN "InlineObjectElementId" IS NOT NULL THEN 2
					WHEN "PageBreakId" IS NOT NULL THEN 3
					WHEN "EquationId" IS NOT NULL THEN 4
					ELSE 0
				END "ElementType",
				"Index",
				coalesce(TR."TextStyleId", IOE."TextStyleId", PB."TextStyleId", E."TextStyleId", 0) "TextStyleId",
				coalesce(TR."Content", E."Content", '') "Content",
				coalesce("InlineObjectId", '00000000-0000-0000-0000-000000000000') "InlineObjectId"
			FROM "ParagraphElements" PE
			LEFT OUTER JOIN "TextRun" TR on PE."TextRunId" = TR."Id"
			LEFT OUTER JOIN "InlineObjectsElements" IOE on PE."InlineObjectElementId" = IOE."Id"
			LEFT OUTER JOIN "PageBreak" PB on PB."Id" = PE."PageBreakId"
			LEFT OUTER JOIN "Equation" E on E."Id" = PE."EquationId"
			WHERE "ParagraphId" = $1 AND
			      "Index" > $2
			ORDER BY "Index"
			LIMIT $3`,
		paragraphId,
		cursor,
		Limit(limit),
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id             int
			elementId      int
			elementType    models.PEType
			index          int
			textStyleId    int
			content        string
			inlineObjectID uuid.UUID
		)

		if err := rows.Scan(
			&id, &elementId, &elementType, &index, &textStyleId, &content,
			&inlineObjectID,
		); err != nil {
			return nil, 0, err
		}

		element, err := scanParagraphElement(id, elementId, elementType, index, textStyleId, content, inlineObjectID)
		if err != nil {
			return nil, 0, err
		}

		slice = append(slice, element)
	}

	count := len(slice)
	if count > 0 {
		nextCursor = slice[count-1].Index
	}

	return
}

func scanParagraphElement(
	id int, elementId int, elementType models.PEType, index int, textStyleId int, content string,
	inlineObjectId uuid.UUID,
) (models.ParagraphElement, error) {
	element := models.ParagraphElement{
		ID:    id,
		Index: index,
	}

	switch elementType {
	case models.PETextRune:
		element.TextRune = &models.TextRune{
			ID:          elementId,
			Content:     content,
			TextStyleID: textStyleId,
		}
	default:
		return models.ParagraphElement{}, fmt.Errorf("Can't scan paragraph element: unsuported element type")
	}

	return element, nil
}

func (r Repository) GetParagraphElement(
	ctx context.Context,
	paragraphId int,
	paragraphElementId int,
) (models.ParagraphElement, error) {
	var (
		id             int
		elementId      int
		elementType    models.PEType
		index          int
		textStyleId    int
		content        string
		inlineObjectID uuid.UUID
	)

	if err := r.conn.QueryRow(
		ctx,
		`
			SELECT
				PE."Id",
				coalesce("TextRunId", "InlineObjectElementId", "PageBreakId", "EquationId") "ElementId",
				CASE
					WHEN "TextRunId" IS NOT NULL THEN 1
					WHEN "InlineObjectElementId" IS NOT NULL THEN 2
					WHEN "PageBreakId" IS NOT NULL THEN 3
					WHEN "EquationId" IS NOT NULL THEN 4
					ELSE 0
				END "ElementType",
				"Index",
				coalesce(TR."TextStyleId", IOE."TextStyleId", PB."TextStyleId", E."TextStyleId", 0) "TextStyleId",
				coalesce(TR."Content", E."Content", '') "Content",
				coalesce("InlineObjectId", '00000000-0000-0000-0000-000000000000') "InlineObjectId"
			FROM "ParagraphElements" PE
			LEFT OUTER JOIN "TextRun" TR on PE."TextRunId" = TR."Id"
			LEFT OUTER JOIN "InlineObjectsElements" IOE on PE."InlineObjectElementId" = IOE."Id"
			LEFT OUTER JOIN "PageBreak" PB on PB."Id" = PE."PageBreakId"
			LEFT OUTER JOIN "Equation" E on E."Id" = PE."EquationId"
			WHERE "ParagraphId" = $1 AND
			      PE."Id" = $2`,
		paragraphId,
		paragraphElementId,
	).Scan(
		&id, &elementId, &elementType, &index, &textStyleId, &content,
		&inlineObjectID,
	); err == pgx.ErrNoRows {
		return models.ParagraphElement{}, repository.ErrParagraphElementNotFound
	} else if err != nil {
		return models.ParagraphElement{}, err
	}

	return scanParagraphElement(id, elementId, elementType, index, textStyleId, content, inlineObjectID)
}

func (r Repository) DeleteParagraphElement(
	ctx context.Context,
	element models.ParagraphElement,
) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := deleteParagraphElement(ctx, tx, element); err != nil {
		return err
	}

	if err := deleteParagraphElementChild(ctx, tx, element); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func deleteParagraphElement(ctx context.Context, tx pgx.Tx, element models.ParagraphElement) error {
	var peId int
	{
		if err := tx.QueryRow(
			ctx,
			`
			DELETE FROM "ParagraphElements" 
		   	WHERE "Id" = $1
			RETURNING "ParagraphId"`,
			element.ID,
		).Scan(&peId); err != nil {
			return err
		}
	}

	// Reindex
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE "ParagraphElements" 
			SET "Index" = "Index" - 1
			WHERE 
			    "ParagraphId" = $1 AND
			    "Index" > $2`,
		peId,
		element.Index,
	); err != nil {
		return err
	}

	return nil
}

func deleteParagraphElementChild(ctx context.Context, tx pgx.Tx, element models.ParagraphElement) error {
	var tableName string
	{
		switch element.GetType() {
		case models.PETextRune:
			tableName = "TextRun"
		default:
			return fmt.Errorf("Can't delete paragraph element child: unsuported child type")
		}
	}
	sql := fmt.Sprintf(
		`
		DELETE FROM "%s" WHERE "Id" = $1`,
		tableName,
	)

	if _, err := tx.Exec(ctx, sql, element.GetChildElementID()); err != nil {
		return err
	}

	return nil
}

func (r Repository) UpdateParagraphElement(
	ctx context.Context,
	element models.ParagraphElement,
) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := updateParagraphElementChild(ctx, tx, element); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func updateParagraphElementChild(ctx context.Context, tx pgx.Tx, element models.ParagraphElement) error {
	switch element.GetType() {
	case models.PETextRune:
		if _, err := tx.Exec(
			ctx,
			`UPDATE "TextRun" 
				SET 
				    "Content" = $1, 
				    "TextStyleId" = $2 
				WHERE "Id" = $3`,
			[]byte(element.TextRune.Content),
			element.TextRune.TextStyleID,
			element.TextRune.ID,
		); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Cant update paragraph element child: unsupported type")
	}

	return nil
}

func (r Repository) FindParagraphElementByIndexes(
	ctx context.Context,
	bodyId uuid.UUID,
	seId,
	peId int,
) (models.ParagraphElement, error) {
	var (
		id             int
		elementId      int
		elementType    models.PEType
		index          int
		textStyleId    int
		content        string
		inlineObjectID uuid.UUID
	)

	if err := r.conn.QueryRow(
		ctx,
		`
			SELECT
				PE."Id",
				coalesce("TextRunId", "InlineObjectElementId", "PageBreakId", "EquationId") "ElementId",
				CASE
					WHEN "TextRunId" IS NOT NULL THEN 1
					WHEN "InlineObjectElementId" IS NOT NULL THEN 2
					WHEN "PageBreakId" IS NOT NULL THEN 3
					WHEN "EquationId" IS NOT NULL THEN 4
					ELSE 0
				END "ElementType",
				PE."Index",
				coalesce(TR."TextStyleId", IOE."TextStyleId", PB."TextStyleId", E."TextStyleId", 0) "TextStyleId",
				coalesce(TR."Content", E."Content", '') "Content",
				coalesce("InlineObjectId", '00000000-0000-0000-0000-000000000000') "InlineObjectId"
			FROM "ParagraphElements" PE
			LEFT OUTER JOIN "TextRun" TR on PE."TextRunId" = TR."Id"
			LEFT OUTER JOIN "InlineObjectsElements" IOE on PE."InlineObjectElementId" = IOE."Id"
			LEFT OUTER JOIN "PageBreak" PB on PB."Id" = PE."PageBreakId"
			LEFT OUTER JOIN "Equation" E on E."Id" = PE."EquationId"
			LEFT OUTER JOIN "StructuralElements" SE on PE."ParagraphId" = SE."ParagraphId"
			WHERE 
				SE."BodyId" = $1 AND
				SE."Index" = $2 AND
				PE."Index" = $3`,
		bodyId,
		seId,
		peId,
	).Scan(
		&id, &elementId, &elementType, &index, &textStyleId, &content,
		&inlineObjectID,
	); err == pgx.ErrNoRows {
		return models.ParagraphElement{}, repository.ErrParagraphElementNotFound
	} else if err != nil {
		return models.ParagraphElement{}, repository.ErrParagraphElementNotFound
	}

	element, err := scanParagraphElement(id, elementId, elementType, index, textStyleId, content, inlineObjectID)
	if err != nil {
		return models.ParagraphElement{}, err
	}

	return element, nil
}
