package repo

import (
	"context"
	"testing"

	"orderservice/internal/model"
	"orderservice/mocks"

	gomock "go.uber.org/mock/gomock"
)

func TestMockRepo_GetOrderByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mr := mocks.NewMockRepo(ctrl)
	ctx := context.Background()

	expected := &model.Order{OrderUID: "order1"}
	mr.EXPECT().GetOrderByID(ctx, "order1").Return(expected, nil)

	got, err := mr.GetOrderByID(ctx, "order1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.OrderUID != expected.OrderUID {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestMockRepo_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mr := mocks.NewMockRepo(ctrl)
	ctx := context.Background()

	ord := &model.Order{OrderUID: "c1"}
	mr.EXPECT().CreateOrder(ctx, ord).Return("c1", nil)

	id, err := mr.CreateOrder(ctx, ord)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "c1" {
		t.Fatalf("unexpected id: %s", id)
	}
}
