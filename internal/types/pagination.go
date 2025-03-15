package types

type PaginationFilter struct {
	Query        *string
	Offset       *int
	Count        *int
	SortProperty *string
	SortOrder    *string
}
