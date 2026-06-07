package handler

import (
	"strconv"

	"community-server/internal/model"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService CategoryService
}

func NewCategoryHandler(categoryService CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// Create 创建分类
// @Summary 创建分类
// @Tags 分类
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.CreateCategoryRequest true "分类信息"
// @Success 200 {object} response.Response
// @Router /admin/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	catID, err := h.categoryService.Create(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, gin.H{"category_id": catID})
}

// Update 编辑分类
// @Summary 编辑分类
// @Tags 分类
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path uint true "分类ID"
// @Param body body model.UpdateCategoryRequest true "更新内容"
// @Success 200 {object} response.Response
// @Router /admin/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	catIDStr := c.Param("id")
	catID, err := strconv.ParseUint(catIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的分类ID")
		return
	}

	var req model.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	if err := h.categoryService.Update(uint(catID), &req); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除分类
// @Summary 删除分类
// @Tags 分类
// @Security BearerAuth
// @Produce json
// @Param id path uint true "分类ID"
// @Success 200 {object} response.Response
// @Router /admin/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	catIDStr := c.Param("id")
	catID, err := strconv.ParseUint(catIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的分类ID")
		return
	}

	if err := h.categoryService.Delete(uint(catID)); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetList 分类列表
// @Summary 分类列表
// @Tags 分类
// @Produce json
// @Success 200 {object} response.Response
// @Router /categories [get]
func (h *CategoryHandler) GetList(c *gin.Context) {
	items, err := h.categoryService.GetList()
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, items)
}

// GetByID 分类详情
// @Summary 分类详情
// @Tags 分类
// @Produce json
// @Param id path uint true "分类ID"
// @Success 200 {object} response.Response
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	catIDStr := c.Param("id")
	catID, err := strconv.ParseUint(catIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的分类ID")
		return
	}

	cat, err := h.categoryService.GetByID(uint(catID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeNotFound, err.Error())
		return
	}
	response.Success(c, cat)
}
