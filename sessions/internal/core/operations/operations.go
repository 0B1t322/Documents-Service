package operations

type (
	Insert interface {
		GetIndex() int
		SetIndex(int) Insert
		InsertLength() int
		IncludeInsert(Insert) Insert
		IncludeDelete(Delete) Insert
	}

	Delete interface {
		GetIndex() int
		SetIndex(int) Delete
		DeleteLength() int
		IncludeInsert(Insert) Delete
		IncludeDelete(Delete) Delete
	}
)
