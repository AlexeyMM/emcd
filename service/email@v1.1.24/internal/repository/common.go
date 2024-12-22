package repository

type Pagination struct {
	Page int
	Size int
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Size
}

func (p Pagination) Limit() int {
	return p.Size
}

func (p Pagination) Next() Pagination {
	return Pagination{
		Page: p.Page + 1,
		Size: p.Size,
	}
}

func NewPagination(size int) Pagination {
	return Pagination{
		Page: 1,
		Size: size,
	}
}
