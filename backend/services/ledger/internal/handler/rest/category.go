package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/moomaideng/bogbank/services/ledger/internal/usecase"
)

type categoryHandler struct {
	categories *usecase.CategoryUsecase
}

func registerCategories(api huma.API, categories *usecase.CategoryUsecase) {
	h := &categoryHandler{categories: categories}
	huma.Register(api, huma.Operation{
		OperationID: "createCategory",
		Tags:        []string{"categories"},
		Method:      http.MethodPost,
		Path:        "/users/{user_id}/categories",
	}, h.createCategory)
	huma.Register(api, huma.Operation{
		OperationID: "listCategories",
		Tags:        []string{"categories"},
		Method:      http.MethodGet,
		Path:        "/users/{user_id}/categories",
	}, h.listCategories)
	huma.Register(api, huma.Operation{
		OperationID: "getCategory",
		Tags:        []string{"categories"},
		Method:      http.MethodGet,
		Path:        "/users/{user_id}/categories/{id}",
	}, h.getCategory)
	huma.Register(api, huma.Operation{
		OperationID: "updateCategory",
		Tags:        []string{"categories"},
		Method:      http.MethodPut,
		Path:        "/users/{user_id}/categories/{id}",
	}, h.updateCategory)
	huma.Register(api, huma.Operation{
		OperationID: "deleteCategory",
		Tags:        []string{"categories"},
		Method:      http.MethodDelete,
		Path:        "/users/{user_id}/categories/{id}",
	}, h.deleteCategory)
}

type categoryBody struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type categoryNameBody struct {
	Name string `json:"name"`
}

type categoryCreateInput struct {
	UserID string `path:"user_id"`
	Body   categoryNameBody
}

type categoryUpdateInput struct {
	UserID string `path:"user_id"`
	ID     string `path:"id"`
	Body   categoryNameBody
}

type categoryOutput struct {
	Body categoryBody
}

type categoryListOutput struct {
	Body []categoryBody
}

func (h *categoryHandler) createCategory(ctx context.Context, in *categoryCreateInput) (*categoryOutput, error) {
	userID, err := parseUserID(in.UserID)
	if err != nil {
		return nil, err
	}

	category, err := h.categories.Create(ctx, userID, in.Body.Name)
	if err != nil {
		return nil, writeErr(err)
	}

	return &categoryOutput{Body: toCategoryBody(category)}, nil
}

func (h *categoryHandler) listCategories(ctx context.Context, in *userInput) (*categoryListOutput, error) {
	userID, err := parseUserID(in.UserID)
	if err != nil {
		return nil, err
	}

	categories, err := h.categories.List(ctx, userID)
	if err != nil {
		return nil, writeErr(err)
	}

	body := make([]categoryBody, 0, len(categories))
	for _, category := range categories {
		body = append(body, toCategoryBody(category))
	}

	return &categoryListOutput{Body: body}, nil
}

func (h *categoryHandler) getCategory(ctx context.Context, in *idInput) (*categoryOutput, error) {
	userID, err := parseUserID(in.UserID)
	if err != nil {
		return nil, err
	}
	categoryID, err := parseCategoryID(in.ID)
	if err != nil {
		return nil, err
	}

	category, err := h.categories.Get(ctx, userID, categoryID)
	if err != nil {
		return nil, writeErr(err)
	}

	return &categoryOutput{Body: toCategoryBody(category)}, nil
}

func (h *categoryHandler) updateCategory(ctx context.Context, in *categoryUpdateInput) (*categoryOutput, error) {
	userID, err := parseUserID(in.UserID)
	if err != nil {
		return nil, err
	}
	categoryID, err := parseCategoryID(in.ID)
	if err != nil {
		return nil, err
	}

	category, err := h.categories.Update(ctx, userID, categoryID, in.Body.Name)
	if err != nil {
		return nil, writeErr(err)
	}

	return &categoryOutput{Body: toCategoryBody(category)}, nil
}

func (h *categoryHandler) deleteCategory(ctx context.Context, in *idInput) (*struct{}, error) {
	userID, err := parseUserID(in.UserID)
	if err != nil {
		return nil, err
	}
	categoryID, err := parseCategoryID(in.ID)
	if err != nil {
		return nil, err
	}

	if err := h.categories.Delete(ctx, userID, categoryID); err != nil {
		return nil, writeErr(err)
	}

	return &struct{}{}, nil
}

func toCategoryBody(category usecase.Category) categoryBody {
	return categoryBody{
		ID:        category.ID,
		UserID:    category.UserID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
	}
}
