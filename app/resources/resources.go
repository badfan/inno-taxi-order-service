package resources

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/badfan/inno-taxi-order-service/app"
	"github.com/badfan/inno-taxi-order-service/app/models"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3" //nolint:typecheck
	"go.uber.org/zap"
)

type IResource interface {
	CreateOrder(ctx context.Context, order *models.Order) (int, error)
	GetOrderHistory(ctx context.Context, uuid uuid.UUID) ([]*models.Order, error)
}

type Resource struct {
	Db     *sql.DB
	logger *zap.SugaredLogger
}

func NewResource(dbConfig *app.DBConfig, logger *zap.SugaredLogger) (*Resource, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.SSLMode)

	db, err := goose.OpenDBWithDriver("pgx", connStr) //nolint:typecheck
	if err != nil {
		return nil, err
	}

	logger.Info("Migration start")

	err = goose.Up(db, "./migrations/") //nolint:typecheck
	if err != nil {
		return nil, err
	}

	logger.Info("Migration ended")

	return &Resource{Db: db, logger: logger}, nil
}
