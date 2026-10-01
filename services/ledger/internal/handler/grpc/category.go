package grpc

import (
	"context"

	ledgerv1 "github.com/moomaideng/bogbank/services/ledger/internal/proto/ledger/v1"
	"github.com/moomaideng/bogbank/services/ledger/internal/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateCategory(ctx context.Context, req *ledgerv1.CreateCategoryRequest) (*ledgerv1.CreateCategoryResponse, error) {
	userID, err := parseUserID(req.GetUserId())
	if err != nil {
		return nil, err
	}

	category, err := s.categories.Create(ctx, userID, req.GetName())
	if err != nil {
		return nil, statusErr(err)
	}

	return &ledgerv1.CreateCategoryResponse{Category: toCategoryProto(category)}, nil
}

func (s *Server) GetCategory(ctx context.Context, req *ledgerv1.GetCategoryRequest) (*ledgerv1.GetCategoryResponse, error) {
	userID, err := parseUserID(req.GetUserId())
	if err != nil {
		return nil, err
	}
	categoryID, err := parseCategoryID(req.GetId())
	if err != nil {
		return nil, err
	}

	category, err := s.categories.Get(ctx, userID, categoryID)
	if err != nil {
		return nil, statusErr(err)
	}

	return &ledgerv1.GetCategoryResponse{Category: toCategoryProto(category)}, nil
}

func (s *Server) ListCategories(ctx context.Context, req *ledgerv1.ListCategoriesRequest) (*ledgerv1.ListCategoriesResponse, error) {
	userID, err := parseUserID(req.GetUserId())
	if err != nil {
		return nil, err
	}
	categories, err := s.categories.List(ctx, userID)
	if err != nil {
		return nil, statusErr(err)
	}

	out := &ledgerv1.ListCategoriesResponse{
		Categories: make([]*ledgerv1.Category, 0, len(categories)),
	}

	for _, category := range categories {
		out.Categories = append(out.Categories, toCategoryProto(category))
	}

	return out, nil
}

func (s *Server) UpdateCategory(ctx context.Context, req *ledgerv1.UpdateCategoryRequest) (*ledgerv1.UpdateCategoryResponse, error) {
	userID, err := parseUserID(req.GetUserId())
	if err != nil {
		return nil, err
	}
	categoryID, err := parseCategoryID(req.GetId())
	if err != nil {
		return nil, err
	}

	category, err := s.categories.Update(ctx, userID, categoryID, req.GetName())
	if err != nil {
		return nil, statusErr(err)
	}

	return &ledgerv1.UpdateCategoryResponse{Category: toCategoryProto(category)}, nil
}

func (s *Server) DeleteCategory(ctx context.Context, req *ledgerv1.DeleteCategoryRequest) (*ledgerv1.DeleteCategoryResponse, error) {
	userID, err := parseUserID(req.GetUserId())
	if err != nil {
		return nil, err
	}
	categoryID, err := parseCategoryID(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := s.categories.Delete(ctx, userID, categoryID); err != nil {
		return nil, statusErr(err)
	}

	return &ledgerv1.DeleteCategoryResponse{}, nil
}

func toCategoryProto(category usecase.Category) *ledgerv1.Category {
	return &ledgerv1.Category{
		Id:        category.ID.String(),
		UserId:    category.UserID.String(),
		Name:      category.Name,
		CreatedAt: timestamppb.New(category.CreatedAt),
	}
}
