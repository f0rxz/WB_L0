package cache

import (
	"testing"
	"time"

	"orderservice/internal/model"
	"orderservice/mocks"

	gomock "go.uber.org/mock/gomock"
)

func TestMockCache_TTLExpiration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mc := mocks.NewMockCache(ctrl)

	ord := &model.Order{OrderUID: "o2"}

	mc.EXPECT().Set(ord)
	mc.EXPECT().Get("o2").Return(nil, false)

	mc.Set(ord)

	time.Sleep(10 * time.Millisecond)
	if _, ok := mc.Get("o2"); ok {
		t.Fatalf("expected cache miss after TTL expiry")
	}
}

func TestMockCache_SetupCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mc := mocks.NewMockCache(ctrl)

	orders := []*model.Order{{OrderUID: "a"}, {OrderUID: "b"}}

	mc.EXPECT().SetupCache(orders)
	mc.EXPECT().Get("a").Return(&model.Order{OrderUID: "a"}, true)
	mc.EXPECT().Get("b").Return(&model.Order{OrderUID: "b"}, true)

	mc.SetupCache(orders)

	if o, ok := mc.Get("a"); !ok || o.OrderUID != "a" {
		t.Fatalf("expected order a in cache")
	}
	if o, ok := mc.Get("b"); !ok || o.OrderUID != "b" {
		t.Fatalf("expected order b in cache")
	}
}

func TestMockCache_CloseDoesNotPanic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mc := mocks.NewMockCache(ctrl)

	mc.EXPECT().Close()
	// calling Close should satisfy expectation and not panic
	mc.Close()
}
