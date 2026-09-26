package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestCreateRejectsBlankName(t *testing.T) {
	svc := NewCategoryUsecase(&fakeRepo{})
	_, err := svc.Create(context.Background(), uuid.New(), "  ")
	if !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("err %v", err)
	}
}

func TestCreateGetUpdateDeleteAndList(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewCategoryUsecase(repo)
	user := uuid.New()
	ctx := context.Background()

	created, err := svc.Create(ctx, user, " food ")
	if err != nil {
		t.Fatal(err)
	}
	if created.Name != "food" || created.UserID != user || created.ID == uuid.Nil {
		t.Fatalf("created %+v", created)
	}

	got, err := svc.Get(ctx, user, created.ID)
	if err != nil || got.Name != "food" {
		t.Fatalf("get %+v err %v", got, err)
	}

	updated, err := svc.Update(ctx, user, created.ID, "rent")
	if err != nil || updated.Name != "rent" {
		t.Fatalf("update %+v err %v", updated, err)
	}

	listed, err := svc.List(ctx, user)
	if err != nil || len(listed) != 1 || listed[0].Name != "rent" {
		t.Fatalf("list %+v err %v", listed, err)
	}

	if err := svc.Delete(ctx, user, created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Get(ctx, user, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("get after delete %v", err)
	}
}

func TestOtherUserCannotSeeOrChangeCategory(t *testing.T) {
	svc := NewCategoryUsecase(&fakeRepo{})
	owner := uuid.New()
	other := uuid.New()
	ctx := context.Background()

	created, err := svc.Create(ctx, owner, "food")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := svc.Get(ctx, other, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("get %v", err)
	}
	listed, err := svc.List(ctx, other)
	if err != nil || len(listed) != 0 {
		t.Fatalf("list %+v err %v", listed, err)
	}
	if _, err := svc.Update(ctx, other, created.ID, "rent"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("update %v", err)
	}
	if err := svc.Delete(ctx, other, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("delete %v", err)
	}

	still, err := svc.Get(ctx, owner, created.ID)
	if err != nil || still.Name != "food" {
		t.Fatalf("owner get %+v err %v", still, err)
	}
}

func TestDuplicateNameForSameUser(t *testing.T) {
	svc := NewCategoryUsecase(&fakeRepo{})
	user := uuid.New()
	ctx := context.Background()
	if _, err := svc.Create(ctx, user, "food"); err != nil {
		t.Fatal(err)
	}
	_, err := svc.Create(ctx, user, "food")
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("err %v", err)
	}
}

type fakeRepo struct {
	rows map[uuid.UUID]Category
}

var _ CategoryRepo = (*fakeRepo)(nil)

func (f *fakeRepo) CreateOne(_ context.Context, category Category) (Category, error) {
	for _, row := range f.rows {
		if row.UserID == category.UserID && row.Name == category.Name {
			return Category{}, ErrConflict
		}
	}
	if f.rows == nil {
		f.rows = map[uuid.UUID]Category{}
	}
	f.rows[category.ID] = category
	return category, nil
}

func (f *fakeRepo) FindOneByUserIDAndCategoryID(_ context.Context, userID, categoryID uuid.UUID) (Category, error) {
	row, ok := f.rows[categoryID]
	if !ok || row.UserID != userID {
		return Category{}, ErrNotFound
	}
	return row, nil
}

func (f *fakeRepo) FindManyByUserID(_ context.Context, userID uuid.UUID) ([]Category, error) {
	out := make([]Category, 0)
	for _, row := range f.rows {
		if row.UserID == userID {
			out = append(out, row)
		}
	}
	return out, nil
}

func (f *fakeRepo) UpdateNameByUserIDAndCategoryID(_ context.Context, userID, categoryID uuid.UUID, name string) (Category, error) {
	row, ok := f.rows[categoryID]
	if !ok || row.UserID != userID {
		return Category{}, ErrNotFound
	}
	for _, other := range f.rows {
		if other.ID != categoryID && other.UserID == userID && other.Name == name {
			return Category{}, ErrConflict
		}
	}
	row.Name = name
	f.rows[categoryID] = row
	return row, nil
}

func (f *fakeRepo) DeleteOneByUserIDAndCategoryID(_ context.Context, userID, categoryID uuid.UUID) error {
	row, ok := f.rows[categoryID]
	if !ok || row.UserID != userID {
		return ErrNotFound
	}
	delete(f.rows, categoryID)
	return nil
}
