package generator

import (
	"time"

	"orderservice/internal/model"

	"github.com/brianvoe/gofakeit/v7"
)

func init() {
	gofakeit.Seed(time.Now().UnixNano())
}

func RandomOrder() *model.Order {
	orderUID := gofakeit.UUID()
	trackNumber := gofakeit.LetterN(2) + gofakeit.DigitN(10)
	customerID := gofakeit.Username()
	chrtID := gofakeit.Number(1000000, 9999999)
	itemRID := gofakeit.UUID()

	return &model.Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       gofakeit.LetterN(4),
		Delivery: model.Delivery{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Street(),
			Region:  gofakeit.State(),
			Email:   gofakeit.Email(),
		},
		Payment: model.Payment{
			Transaction:  gofakeit.UUID(),
			RequestID:    gofakeit.UUID(),
			Currency:     gofakeit.CurrencyShort(),
			Provider:     gofakeit.Company(),
			Amount:       gofakeit.Number(100, 5000),
			PaymentDt:    int(time.Now().Unix()),
			Bank:         gofakeit.BankName(),
			DeliveryCost: gofakeit.Number(500, 2000),
			GoodsTotal:   gofakeit.Number(100, 3000),
			CustomFee:    gofakeit.Number(0, 100),
		},
		Items: []model.Item{
			{
				ChrtID:      chrtID,
				TrackNumber: trackNumber,
				Price:       gofakeit.Number(50, 500),
				RID:         itemRID,
				Name:        gofakeit.ProductName(),
				Sale:        gofakeit.Number(0, 50),
				Size:        gofakeit.DigitN(1),
				TotalPrice:  gofakeit.Number(50, 500),
				NMID:        gofakeit.Number(1000000, 9999999),
				Brand:       gofakeit.Company(),
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        customerID,
		DeliveryService:   gofakeit.Company(),
		ShardKey:          gofakeit.DigitN(1),
		SMID:              gofakeit.Number(1, 100),
		DateCreated:       time.Now().UTC(),
		OOFShard:          gofakeit.DigitN(1),
	}
}
