package order

import (
	"context"
	"sync"
	"time"

	"github.com/badfan/inno-taxi-order-service/app"
	"github.com/badfan/inno-taxi-order-service/app/models"
	"github.com/badfan/inno-taxi-order-service/app/resources"
	"github.com/badfan/inno-taxi-order-service/app/rpc"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type RPCService struct {
	OrderServiceServer  rpc.UnimplementedOrderServiceServer
	userServiceClient   rpc.UserServiceClient
	driverServiceClient rpc.DriverServiceClient
	resource            resources.IResource
	rpcConfig           *app.RPCConfig
	mtx                 *sync.Mutex
	logger              *zap.SugaredLogger
}

func NewRPCService(
	userClientConn *grpc.ClientConn,
	driverClientConn *grpc.ClientConn,
	resource resources.IResource,
	rpcConfig *app.RPCConfig,
	logger *zap.SugaredLogger) *RPCService {
	return &RPCService{
		OrderServiceServer:  rpc.UnimplementedOrderServiceServer{},
		userServiceClient:   rpc.NewUserServiceClient(userClientConn),
		driverServiceClient: rpc.NewDriverServiceClient(driverClientConn),
		resource:            resource,
		rpcConfig:           rpcConfig,
		mtx:                 &sync.Mutex{},
		logger:              logger,
	}
}

type userInfo struct {
	UUID        string
	origin      string
	destination string
	rating      float32
	taxiType    string
}

type driverInfo struct {
	UUID     string
	rating   float32
	taxiType string
}

type orderInfo struct {
	userUUID     string
	driverUUID   string
	origin       string
	destination  string
	userRating   float32
	driverRating float32
	error
}

var busyUser = make(chan orderInfo)
var busyDriver = make(chan orderInfo)
var freeUsers = make(chan userInfo, 100)
var freeDrivers = make(chan driverInfo, 100)

func (s *RPCService) SetDriverRating(ctx context.Context, req *rpc.SetDriverRatingRequest) (*rpc.EmptyResponse, error) {
	s.logger.Info("setting driver rating by user")

	_, err := s.driverServiceClient.SetDriverRating(ctx, &rpc.SetDriverRatingRequest{Rating: req.GetRating()})
	if err != nil {
		s.logger.Errorf("error occurred while setting driver rating by user: %s", err.Error())
		return nil, err
	}

	return &rpc.EmptyResponse{}, nil
}

func (s *RPCService) SetUserRating(ctx context.Context, req *rpc.SetUserRatingRequest) (*rpc.EmptyResponse, error) {
	s.logger.Info("setting user rating by driver")

	_, err := s.userServiceClient.SetUserRating(ctx, &rpc.SetUserRatingRequest{Rating: req.GetRating()})
	if err != nil {
		s.logger.Errorf("error occurred while setting user rating by driver: %s", err.Error())
		return nil, err
	}

	return &rpc.EmptyResponse{}, nil
}

func (s *RPCService) GetTaxiForUser(ctx context.Context, req *rpc.GetTaxiForUserRequest) (*rpc.GetTaxiForUserResponse, error) {
	s.logger.Infof("driver search is in progress for user uuid: %s", req.GetUserUuid())

	ctxWithTimeOut, cancel := context.WithTimeout(ctx, time.Duration(s.rpcConfig.WaitingTime)*time.Minute)
	defer cancel()

	go takeFreeDriver(ctxWithTimeOut, s)
	freeUsers <- userInfo{
		UUID:        req.GetUserUuid(),
		origin:      req.GetOrigin(),
		destination: req.GetDestination(),
		rating:      req.GetUserRating(),
		taxiType:    req.GetTaxiType(),
	}

	select {
	case info := <-busyUser:
		if info.error != nil {
			s.logger.Errorf("error occurred while getting driver info for user: %s", info.error.Error())
			return nil, info.error
		}
		return &rpc.GetTaxiForUserResponse{DriverUuid: info.driverUUID, DriverRating: info.driverRating}, nil
	}
}

func (s *RPCService) GetOrderForDriver(ctx context.Context, req *rpc.GetOrderForDriverRequest) (*rpc.GetOrderForDriverResponse, error) {
	s.logger.Infof("received order request from driver with uuid: %s", req.GetDriverUuid())

	freeDrivers <- driverInfo{
		UUID:     req.GetDriverUuid(),
		rating:   req.GetDriverRating(),
		taxiType: req.GetTaxiType(),
	}
	select {
	case info := <-busyDriver:
		if info.error != nil {
			s.logger.Errorf("error occurred while getting user info for driver: %s", info.error.Error())
			return nil, info.error
		}
		return &rpc.GetOrderForDriverResponse{
			UserUuid:    info.userUUID,
			Origin:      info.origin,
			Destination: info.destination,
			UserRating:  info.userRating,
		}, nil
	}
}

func (s *RPCService) GetOrderHistory(ctx context.Context, req *rpc.GetOrderHistoryRequest) (*rpc.GetOrderHistoryResponse, error) {
	s.logger.Infof("getting order history for uuid: %s", req.GetUuid())

	reqUUID := uuid.MustParse(req.GetUuid())

	dbOrders, err := s.resource.GetOrderHistory(ctx, reqUUID)
	if err != nil {
		s.logger.Errorf("error occurred while getting order history: %s", err.Error())
		return nil, err
	}

	ordersHistory := grpcOrdersConvert(dbOrders)

	return &rpc.GetOrderHistoryResponse{
		Orders: ordersHistory,
	}, nil
}

func takeFreeDriver(ctx context.Context, s *RPCService) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	user := <-freeUsers
	select {
	case driver := <-freeDrivers:
		if driver.taxiType == user.taxiType {
			_, err := s.resource.CreateOrder(ctx, &models.Order{
				UserUuid:    uuid.MustParse(user.UUID),
				DriverUuid:  uuid.MustParse(driver.UUID),
				Origin:      user.origin,
				Destination: user.destination,
				TaxiType:    models.TaxiType(user.taxiType),
			})
			order := orderInfo{
				userUUID:     user.UUID,
				driverUUID:   driver.UUID,
				origin:       user.origin,
				destination:  user.destination,
				userRating:   user.rating,
				driverRating: driver.rating,
				error:        err,
			}
			busyDriver <- order
			busyUser <- order
			return
		} else {
			freeDrivers <- driver
		}
	case <-ctx.Done():
		s.logger.Infof("driver waiting time exceeded for user %s", user.UUID)
	default:
		s.logger.Infof("user %s is waiting for free drivers...", user.UUID)
	}
}

func grpcOrdersConvert(source []*models.Order) []*rpc.Order {
	var res []*rpc.Order
	for _, order := range source {
		res = append(res, &rpc.Order{
			UserUuid:    order.UserUuid.String(),
			DriverUuid:  order.DriverUuid.String(),
			Origin:      order.Origin,
			Destination: order.Destination,
			TaxiType:    string(order.TaxiType),
			Date:        order.CreatedAt.String(),
			Duration:    order.UpdatedAt.Sub(order.CreatedAt).String()})
	}

	return res
}
