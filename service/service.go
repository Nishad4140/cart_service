package service

import (
	"context"
	"fmt"

	"github.com/Nishad4140/cart_service/adapter"
	"github.com/Nishad4140/cart_service/entities"
	"github.com/Nishad4140/proto_files/pb"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var (
	Tracer        opentracing.Tracer
	ProductClient pb.ProductServiceClient
)

func RetrieveTracer(tr opentracing.Tracer) {
	Tracer = tr
}

type CartService struct {
	Adapter adapter.AdapterInterface
	pb.UnimplementedCartServiceServer
}

func NewCartService(adapter adapter.AdapterInterface) *CartService {
	return &CartService{
		Adapter: adapter,
	}
}

func (cart *CartService) CreateCart(ctx context.Context, req *pb.CartCreate) (*pb.CartResponse, error) {
	span := Tracer.StartSpan("create cart grpc")
	defer span.Finish()

	err := cart.Adapter.CreateCart(uint(req.UserId))
	if err != nil {
		return &pb.CartResponse{}, err
	}
	res := &pb.CartResponse{
		UserId: req.UserId,
	}
	return res, nil
}

func (cart *CartService) AddToCart(ctx context.Context, req *pb.AddToCartRequest) (*pb.CartResponse, error) {
	productData, err := ProductClient.GetProduct(context.TODO(), &pb.GetProductByID{Id: uint32(req.ProdId)})
	if err != nil {
		return &pb.CartResponse{}, fmt.Errorf("error fetching product data")
	}
	if productData.Name == "" {
		return &pb.CartResponse{}, fmt.Errorf("product not found")
	}
	if productData.Quantity < req.Quantity {
		return &pb.CartResponse{}, fmt.Errorf("not enough quantity is available to add in cart please decrease the quantity")
	}
	reqEntity := entities.CartItems{
		ProductId: uint(req.ProdId),
		Total:     float64(productData.Price),
		Quantity:  int(req.Quantity),
	}
	fmt.Println(reqEntity)
	err = cart.Adapter.AddToCart(reqEntity, uint(req.UserId))
	if err != nil {
		return &pb.CartResponse{}, err
	}
	res := &pb.CartResponse{
		UserId:  req.UserId,
		IsEmpty: false,
	}

	return res, nil
}

type HealthChecker struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *HealthChecker) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	fmt.Println("check called")
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthChecker) Watch(in *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watching is not supported")
}
