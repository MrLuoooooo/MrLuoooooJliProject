package model

type CreateTagRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description" binding:"max=200"`
}

type UpdateTagRequest struct {
	Name        string `json:"name" binding:"max=50"`
	Description string `json:"description" binding:"max=200"`
}

type TagListRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

type TagResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PostCount   int    `json:"post_count"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type TagListResponse struct {
	Total int64        `json:"total"`
	Items []TagResponse `json:"items"`
}

type AddPostTagsRequest struct {
	TagIDs []uint `json:"tag_ids" binding:"required"`
}
