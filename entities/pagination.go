package entities

//
// Paginator describes pagination of a collection
//
type Paginator struct {
	Page        int   `json:"page"`
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
	HasNextPage bool  `json:"-"`
}

// TotalPages calculate total number of pages based on Page and PerPage attributes
func (p Paginator) TotalPages() int64 {
	pages := p.Total / int64(p.PerPage)
	if p.Total%int64(p.PerPage) > 0 {
		pages++
	}
	return pages
}

//
// Pager simply defines pagination request for collections
//
type Pager struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// Offset calculate offset based on Page and PerPage attributes
func (p Pager) Offset() int64 {
	return int64(p.Page-1) * int64(p.PerPage)
}

// CreatePaginator create paginator with Total and HasNextPage filled in
func (p Pager) CreatePaginator(resultsOnCurrentPage int, totalCount int64) *Paginator {
	res := new(Paginator)
	res.Page = p.Page
	res.PerPage = p.PerPage
	res.HasNextPage = (p.Offset() + int64(p.PerPage)) < totalCount
	res.Total = totalCount
	return res
}

//
// Sorter simply fines the soring options for collections
//
type Sorter struct {
	Field string `json:"field"`
	Asc   bool   `json:"asc"`
}
