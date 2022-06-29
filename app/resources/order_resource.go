package resources

import (
	"context"

	"github.com/badfan/inno-taxi-order-service/app/models"
	"github.com/badfan/inno-taxi-order-service/app/models/sqlc"
	"github.com/google/uuid"
)

func (r *Resource) CreateOrder(ctx context.Context, order *models.Order) (int, error) {
	queries := sqlc.New(r.Db)

	res, err := queries.CreateOrder(ctx, sqlc.CreateOrderParams{
		UserUuid:    order.UserUuid,
		DriverUuid:  order.DriverUuid,
		Origin:      order.Origin,
		Destination: order.Destination,
		TaxiType:    sqlc.TaxiType(order.TaxiType),
	})
	if err != nil {
		return 0, err
	}

	return int(res.ID), nil
}

func (r *Resource) GetOrderHistory(ctx context.Context, uuid uuid.UUID) ([]*models.Order, error) {
	queries := sqlc.New(r.Db)

	res, err := queries.GetOrdersByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return sqlcOrderArrConvert(res), nil
}

func sqlcOrderArrConvert(source []sqlc.Order) []*models.Order {
	var res []*models.Order

	for i, item := range source {
		res[i] = &models.Order{
			ID:          item.ID,
			UserUuid:    item.UserUuid,
			DriverUuid:  item.DriverUuid,
			Origin:      item.Origin,
			Destination: item.Destination,
			TaxiType:    models.TaxiType(item.TaxiType),
			IsFinished:  item.IsFinished,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
	}

	return res
}
