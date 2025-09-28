package repo

import (
	"context"
	"encoding/json"

	"orderservice/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Order interface {
	CreateOrder(ctx context.Context, order *model.Order) (string, error)
	GetOrderByID(ctx context.Context, id string) (*model.Order, error)
	GetAllOrders(ctx context.Context) ([]*model.Order, error)
}
type order struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Order {
	return &order{
		db: db,
	}
}

func (o order) CreateOrder(ctx context.Context, ord *model.Order) (string, error) {
	tx, err := o.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	// Insert order
	_, err = tx.Exec(ctx, `
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES (
			@order_uid, @track_number, @entry, @locale, @internal_signature,
			@customer_id, @delivery_service, @shardkey, @sm_id, @date_created, @oof_shard
		)`,
		pgx.NamedArgs{
			"order_uid":          ord.OrderUID,
			"track_number":       ord.TrackNumber,
			"entry":              ord.Entry,
			"locale":             ord.Locale,
			"internal_signature": ord.InternalSignature,
			"customer_id":        ord.CustomerID,
			"delivery_service":   ord.DeliveryService,
			"shardkey":           ord.ShardKey,
			"sm_id":              ord.SMID,
			"date_created":       ord.DateCreated,
			"oof_shard":          ord.OOFShard,
		},
	)
	if err != nil {
		return "", err
	}

	// Insert delivery
	_, err = tx.Exec(ctx, `
		INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
		VALUES (@order_uid, @name, @phone, @zip, @city, @address, @region, @email)`,
		pgx.NamedArgs{
			"order_uid": ord.OrderUID,
			"name":      ord.Delivery.Name,
			"phone":     ord.Delivery.Phone,
			"zip":       ord.Delivery.Zip,
			"city":      ord.Delivery.City,
			"address":   ord.Delivery.Address,
			"region":    ord.Delivery.Region,
			"email":     ord.Delivery.Email,
		},
	)
	if err != nil {
		return "", err
	}

	// Insert payment
	_, err = tx.Exec(ctx, `
		INSERT INTO payments (
			order_uid, transaction, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES (
			@order_uid, @transaction, @request_id, @currency, @provider,
			@amount, @payment_dt, @bank, @delivery_cost, @goods_total, @custom_fee
		)`,
		pgx.NamedArgs{
			"order_uid":     ord.OrderUID,
			"transaction":   ord.Payment.Transaction,
			"request_id":    ord.Payment.RequestID,
			"currency":      ord.Payment.Currency,
			"provider":      ord.Payment.Provider,
			"amount":        ord.Payment.Amount,
			"payment_dt":    ord.Payment.PaymentDt,
			"bank":          ord.Payment.Bank,
			"delivery_cost": ord.Payment.DeliveryCost,
			"goods_total":   ord.Payment.GoodsTotal,
			"custom_fee":    ord.Payment.CustomFee,
		},
	)
	if err != nil {
		return "", err
	}

	// Insert items
	for _, item := range ord.Items {
		_, err = tx.Exec(ctx, `
			INSERT INTO items (
				order_uid, chrt_id, track_number, price, rid,
				name, sale, size, total_price, nm_id, brand, status
			) VALUES (
				@order_uid, @chrt_id, @track_number, @price, @rid,
				@name, @sale, @size, @total_price, @nm_id, @brand, @status
			)`,
			pgx.NamedArgs{
				"order_uid":    ord.OrderUID,
				"chrt_id":      item.ChrtID,
				"track_number": item.TrackNumber,
				"price":        item.Price,
				"rid":          item.RID,
				"name":         item.Name,
				"sale":         item.Sale,
				"size":         item.Size,
				"total_price":  item.TotalPrice,
				"nm_id":        item.NMID,
				"brand":        item.Brand,
				"status":       item.Status,
			},
		)
		if err != nil {
			return "", err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return "", err
	}

	return ord.OrderUID, nil
}

func (o order) GetOrderByID(ctx context.Context, id string) (*model.Order, error) {
	ord := &model.Order{}

	row := o.db.QueryRow(ctx, `
		SELECT order_uid, track_number, entry, locale, internal_signature,
		       customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders
		WHERE order_uid = $1
	`, id)

	err := row.Scan(
		&ord.OrderUID,
		&ord.TrackNumber,
		&ord.Entry,
		&ord.Locale,
		&ord.InternalSignature,
		&ord.CustomerID,
		&ord.DeliveryService,
		&ord.ShardKey,
		&ord.SMID,
		&ord.DateCreated,
		&ord.OOFShard,
	)
	if err != nil {
		return nil, err
	}

	row = o.db.QueryRow(ctx, `
		SELECT name, phone, zip, city, address, region, email
		FROM deliveries
		WHERE order_uid = $1
	`, id)

	err = row.Scan(
		&ord.Delivery.Name,
		&ord.Delivery.Phone,
		&ord.Delivery.Zip,
		&ord.Delivery.City,
		&ord.Delivery.Address,
		&ord.Delivery.Region,
		&ord.Delivery.Email,
	)
	if err != nil {
		return nil, err
	}

	row = o.db.QueryRow(ctx, `
		SELECT transaction, request_id, currency, provider, amount,
		       payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payments
		WHERE order_uid = $1
	`, id)

	err = row.Scan(
		&ord.Payment.Transaction,
		&ord.Payment.RequestID,
		&ord.Payment.Currency,
		&ord.Payment.Provider,
		&ord.Payment.Amount,
		&ord.Payment.PaymentDt,
		&ord.Payment.Bank,
		&ord.Payment.DeliveryCost,
		&ord.Payment.GoodsTotal,
		&ord.Payment.CustomFee,
	)
	if err != nil {
		return nil, err
	}

	rows, err := o.db.Query(ctx, `
		SELECT chrt_id, track_number, price, rid, name,
		       sale, size, total_price, nm_id, brand, status
		FROM items
		WHERE order_uid = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err = rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NMID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	ord.Items = items

	return ord, nil
}

func (o order) GetAllOrders(ctx context.Context) ([]*model.Order, error) {
	rows, err := o.db.Query(ctx, `
		SELECT 
			o.order_uid,
			o.track_number,
			o.entry,
			row_to_json(d.*) as delivery,
			row_to_json(p.*) as payment,
			COALESCE(json_agg(i.*) FILTER (WHERE i.chrt_id IS NOT NULL), '[]') as items,
			o.locale,
			o.internal_signature,
			o.customer_id,
			o.delivery_service,
			o.shardkey,
			o.sm_id,
			o.date_created,
			o.oof_shard
		FROM orders o
		LEFT JOIN deliveries d ON o.order_uid = d.order_uid
		LEFT JOIN payments   p ON o.order_uid = p.order_uid
		LEFT JOIN items      i ON o.order_uid = i.order_uid
		GROUP BY o.order_uid, d.*, p.*
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order

	for rows.Next() {
		var (
			order        model.Order
			deliveryJSON []byte
			paymentJSON  []byte
			itemsJSON    []byte
		)

		err = rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&deliveryJSON,
			&paymentJSON,
			&itemsJSON,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SMID,
			&order.DateCreated,
			&order.OOFShard,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(deliveryJSON, &order.Delivery); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(paymentJSON, &order.Payment); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	return orders, nil
}
