package pgql_test

import (
	"context"
	"github.com/0B1t322/Documents-Service/document/internal/core/models"
	"github.com/0B1t322/Documents-Service/document/internal/repository/elements/pgql"
	database "github.com/0B1t322/Documents-Service/internal/database/pgql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFunc_PGQLElementRepository(t *testing.T) {
	url, ok := os.LookupEnv("DOCUMENT_APP_PGSQL_URL")
	if !ok {
		url = "postgres://postgres:password@localhost:5432/Documents"
	}

	pool, err := database.NewConnectionPool(context.Background(), url)
	require.NoError(t, err)

	repo := pgql.New(pool)

	t.Run(
		"Store", func(t *testing.T) {
			ctx := context.Background()
			bodyId := uuid.MustParse("f32da28e-bdb0-49d8-9687-71070c98e49c")

			err := repo.StoreStructuralElement(
				ctx, bodyId, &models.StructuralElement{
					Index:        2,
					SectionBreak: &models.SectionBreak{},
				},
			)
			require.NoError(t, err)
		},
	)

	t.Run(
		"Get", func(t *testing.T) {
			bodyId := uuid.MustParse("f32da28e-bdb0-49d8-9687-71070c98e49c")
			ses, next, err := repo.GetStructuralElements(context.Background(), bodyId, 1, 2)
			require.NoError(t, err)
			t.Log("Next", next)
			t.Logf("%+v", ses)
		},
	)

	t.Run(
		"Delete", func(t *testing.T) {
			se, err := repo.FindStructuralElementByID(context.Background(), 14)
			require.NoError(t, err)

			t.Logf("%+v", se)
			t.Logf("%+v", se.SectionBreak)

			err = repo.DeleteStructuralElement(context.Background(), se)
			require.NoError(t, err)
		},
	)
}
