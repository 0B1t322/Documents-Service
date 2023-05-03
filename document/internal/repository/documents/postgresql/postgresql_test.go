package postgresql_test

import (
	"context"
	"github.com/0B1t322/Documents-Service/document/internal/core/models"
	"github.com/0B1t322/Documents-Service/document/internal/repository/documents/postgresql"
	"github.com/0B1t322/Documents-Service/internal/database/pgql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFunc_DocumentsPostgresRepository(t *testing.T) {
	url, ok := os.LookupEnv("DOCUMENT_APP_PGSQL_URL")
	if !ok {
		url = "postgres://postgres:password@localhost:5432/Documents"
	}

	pool, err := pgql.NewConnectionPool(context.Background(), url)
	require.NoError(t, err)

	repo := postgresql.New(pool)

	t.Run(
		"CreateDocument",
		func(t *testing.T) {
			var document = models.Document{
				Title: "New Document",
				Body:  models.Body{},
				Style: models.DocumentStyle{
					PageSize: models.Size{
						Height: models.Dimension{
							Magnitude: 16,
							Unit:      models.UNIT_PT,
						},
						Width: models.Dimension{
							Magnitude: 16,
							Unit:      models.UNIT_PT,
						},
					},
				},
			}
			err = repo.Store(context.Background(), &document)
			require.NoError(t, err)

			t.Logf("%+v", document)
		},
	)

	t.Run(
		"Get Document", func(t *testing.T) {
			d, err := repo.FindByID(context.Background(), uuid.MustParse("bbfc39f4-dea1-41f1-8443-3c6818d728e6"))
			require.NoError(t, err)

			t.Logf("%+v", d)
		},
	)

	t.Run(
		"Update document", func(t *testing.T) {
			d, err := repo.FindByID(context.Background(), uuid.MustParse("bbfc39f4-dea1-41f1-8443-3c6818d728e6"))
			require.NoError(t, err)
			d.Title = "Updated title"
			d.Style.PageSize.Height.Magnitude = 100

			require.NoError(t, repo.Update(context.Background(), d))
		},
	)

	t.Run(
		"Get documents", func(t *testing.T) {
			ds, nextCursor, err := repo.Get(
				context.Background(),
				uuid.MustParse("bbfc39f4-dea1-41f1-8443-3c6818d728e6"),
				10,
			)
			require.NoError(t, err)

			t.Logf("%+v", ds)
			t.Logf("%+v", nextCursor)
		},
	)
}
