package model

type SearchRequest struct {
	Keyword  string `form:"keyword" binding:"required"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

type SearchResponse struct {
	Total int64         `json:"total"`
	Items []PostListItem `json:"items"`
}

type AISearchRequest struct {
	Question string `json:"question" binding:"required"`
}

type AISearchResponse struct {
	Keywords    []string       `json:"keywords"`
	Intent      string         `json:"intent"`
	Posts       []PostListItem `json:"posts"`
	AISummary   string         `json:"ai_summary"`
	Total       int64          `json:"total"`
}
