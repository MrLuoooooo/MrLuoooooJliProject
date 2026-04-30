package model

type SearchRequest struct {
	Keyword  string `form:"keyword" binding:"required"`
	Type     string `form:"type,default=post"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

type SearchResponse struct {
	Total int64          `json:"total"`
	Type  string         `json:"type"`
	Posts []PostListItem `json:"posts,omitempty"`
	Users []SearchUser   `json:"users,omitempty"`
}

type SearchUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

type AISearchRequest struct {
	Question string `json:"question" binding:"required"`
}

type AISearchResponse struct {
	Keywords  []string       `json:"keywords"`
	Intent    string         `json:"intent"`
	Posts     []PostListItem `json:"posts"`
	AISummary string         `json:"ai_summary"`
	Total     int64          `json:"total"`
}
