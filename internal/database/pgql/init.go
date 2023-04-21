package pgql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"sync"
)

// All types
var typesArray = []string{
	// Unit enum
	"unit",
	// Dimensions composite types
	"dimensions",
	// Size composite types
	"size",
	// Alignment enum
	"alignment",
}

var pgsqlTypes []*pgtype.Type

var initType sync.Once

// init types from postgresql
// should be call after database connection
func initTypes(conn *pgx.Conn) error {
	var initErr error = nil
	initType.Do(
		func() {
			for _, pgType := range typesArray {
				dt, err := conn.LoadType(context.Background(), pgType)
				if err != nil {
					initErr = err
					return
				}
				// Need to register types because they can be related
				conn.TypeMap().RegisterType(dt)

				pgsqlTypes = append(pgsqlTypes, dt)
			}
		},
	)

	typeMap := conn.TypeMap()

	for _, t := range pgsqlTypes {
		typeMap.RegisterType(t)
	}

	return initErr
}
