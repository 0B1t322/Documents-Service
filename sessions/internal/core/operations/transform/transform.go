package transform

import "github.com/0B1t322/Documents-Service/sessions/internal/core/operations"

type (
	InsertOperation = operations.Insert
	DeleteOperation = operations.Delete
)

func InsertDelete(op1 InsertOperation, op2 DeleteOperation) InsertOperation {
	return op1.IncludeDelete(op2)
}

func InsertInsert(op1, op2 InsertOperation) InsertOperation {
	return op1.IncludeInsert(op2)
}

func DeleteInsert(op1 DeleteOperation, op2 InsertOperation) DeleteOperation {
	return op1.IncludeInsert(op2)
}

func DeleteDelete(op1, op2 DeleteOperation) DeleteOperation {
	return op1.IncludeDelete(op2)
}
