package vo

type Pagination struct {
	Total       int64       `json:"total"`
	PageSize    int64       `json:"ps"`
	CurrentPage int64       `json:"ps"`
	Data        interface{} `json:"contentList"`
}

func NewPagination(total, pageSize, currentPage int64, data interface{}) Pagination {
	return Pagination{Total: total, PageSize: pageSize, CurrentPage: currentPage, Data: data}
}
