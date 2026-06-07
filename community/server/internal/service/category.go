package service

import (
	"errors"

	"community-server/internal/db/mysql"
	"community-server/internal/model"
	"community-server/internal/repository"
)

type CategoryService struct {
	catRepo repository.CategoryRepository
}

func NewCategoryService(catRepo repository.CategoryRepository) *CategoryService {
	return &CategoryService{catRepo: catRepo}
}

func (s *CategoryService) Create(req *model.CreateCategoryRequest) (uint, error) {
	cat := mysql.Category{
		Name: req.Name, Description: req.Description, SortOrder: req.SortOrder, Status: 1,
	}
	if err := s.catRepo.Create(&cat); err != nil {
		return 0, errors.New("创建分类失败")
	}
	return cat.ID, nil
}

func (s *CategoryService) Update(catID uint, req *model.UpdateCategoryRequest) error {
	if _, err := s.catRepo.FindByID(catID); err != nil {
		return errors.New("分类不存在")
	}
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}
	return s.catRepo.Update(catID, updates)
}

func (s *CategoryService) Delete(catID uint) error {
	if _, err := s.catRepo.FindByID(catID); err != nil {
		return errors.New("分类不存在")
	}
	return s.catRepo.Delete(catID)
}

func (s *CategoryService) GetList() ([]model.CategoryResponse, error) {
	cats, err := s.catRepo.List()
	if err != nil {
		return nil, err
	}
	items := make([]model.CategoryResponse, 0, len(cats))
	for _, c := range cats {
		items = append(items, model.CategoryResponse{
			ID: c.ID, Name: c.Name, Description: c.Description,
			SortOrder: c.SortOrder, PostCount: c.PostCount,
		})
	}
	return items, nil
}

func (s *CategoryService) GetByID(catID uint) (*model.CategoryResponse, error) {
	cat, err := s.catRepo.FindByID(catID)
	if err != nil {
		return nil, errors.New("分类不存在")
	}
	return &model.CategoryResponse{
		ID: cat.ID, Name: cat.Name, Description: cat.Description,
		SortOrder: cat.SortOrder, PostCount: cat.PostCount,
	}, nil
}
