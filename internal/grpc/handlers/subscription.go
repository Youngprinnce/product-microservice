package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/youngprinnce/product-microservice/internal/service"
	"github.com/youngprinnce/product-microservice/internal/service/subscription"
	pb "github.com/youngprinnce/product-microservice/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SubscriptionHandler implements the SubscriptionService gRPC interface
type SubscriptionHandler struct {
	pb.UnimplementedSubscriptionServiceServer
	subscriptionService subscription.SubscriptionBC
}

// NewSubscriptionHandler creates a new subscription gRPC handler
func NewSubscriptionHandler(subscriptionService subscription.SubscriptionBC) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// CreateSubscriptionPlan creates a new subscription plan
func (h *SubscriptionHandler) CreateSubscriptionPlan(ctx context.Context, req *pb.CreateSubscriptionPlanRequest) (*pb.CreateSubscriptionPlanResponse, error) {
	createReq := subscription.CreateSubscriptionPlanRequest{
		ProductID: req.ProductId,
		PlanName:  req.PlanName,
		Duration:  int(req.Duration),
		Price:     req.Price,
	}

	plan, err := h.subscriptionService.CreateSubscriptionPlan(ctx, createReq)
	if err != nil {
		return nil, convertSubscriptionToGRPCError(err)
	}

	return &pb.CreateSubscriptionPlanResponse{
		Plan: convertToProtobufSubscriptionPlan(plan),
	}, nil
}

// GetSubscriptionPlan retrieves a subscription plan by ID
func (h *SubscriptionHandler) GetSubscriptionPlan(ctx context.Context, req *pb.GetSubscriptionPlanRequest) (*pb.GetSubscriptionPlanResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid subscription plan ID")
	}

	plan, err := h.subscriptionService.GetSubscriptionPlan(ctx, id)
	if err != nil {
		return nil, convertSubscriptionToGRPCError(err)
	}

	return &pb.GetSubscriptionPlanResponse{
		Plan: convertToProtobufSubscriptionPlan(plan),
	}, nil
}

// UpdateSubscriptionPlan updates a subscription plan
func (h *SubscriptionHandler) UpdateSubscriptionPlan(ctx context.Context, req *pb.UpdateSubscriptionPlanRequest) (*pb.UpdateSubscriptionPlanResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid subscription plan ID")
	}

	updateReq := subscription.UpdateSubscriptionPlanRequest{
		PlanName: req.PlanName,
		Duration: &[]int{int(req.Duration)}[0],
		Price:    &req.Price,
	}

	plan, err := h.subscriptionService.UpdateSubscriptionPlan(ctx, id, updateReq)
	if err != nil {
		return nil, convertSubscriptionToGRPCError(err)
	}

	return &pb.UpdateSubscriptionPlanResponse{
		Plan: convertToProtobufSubscriptionPlan(plan),
	}, nil
}

// DeleteSubscriptionPlan deletes a subscription plan
func (h *SubscriptionHandler) DeleteSubscriptionPlan(ctx context.Context, req *pb.DeleteSubscriptionPlanRequest) (*pb.DeleteSubscriptionPlanResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid subscription plan ID")
	}

	err = h.subscriptionService.DeleteSubscriptionPlan(ctx, id)
	if err != nil {
		return nil, convertSubscriptionToGRPCError(err)
	}

	return &pb.DeleteSubscriptionPlanResponse{
		Success: true,
	}, nil
}

// ListSubscriptionPlans lists subscription plans for a product
func (h *SubscriptionHandler) ListSubscriptionPlans(ctx context.Context, req *pb.ListSubscriptionPlansRequest) (*pb.ListSubscriptionPlansResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	plans, total, err := h.subscriptionService.ListSubscriptionPlans(ctx, productID, page, pageSize)
	if err != nil {
		return nil, convertSubscriptionToGRPCError(err)
	}

	pbPlans := make([]*pb.SubscriptionPlan, len(plans))
	for i, plan := range plans {
		pbPlans[i] = convertToProtobufSubscriptionPlan(plan)
	}

	return &pb.ListSubscriptionPlansResponse{
		Plans:    pbPlans,
		Total:    total,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

// convertToProtobufSubscriptionPlan converts domain subscription plan to protobuf
func convertToProtobufSubscriptionPlan(plan *subscription.SubscriptionPlan) *pb.SubscriptionPlan {
	return &pb.SubscriptionPlan{
		Id:        plan.ID.String(),
		ProductId: plan.ProductID.String(),
		PlanName:  plan.PlanName,
		Duration:  int32(plan.Duration),
		Price:     plan.Price,
		CreatedAt: timestamppb.New(plan.CreatedAt),
		UpdatedAt: timestamppb.New(plan.UpdatedAt),
	}
}

// convertSubscriptionToGRPCError converts service errors to gRPC errors
func convertSubscriptionToGRPCError(err error) error {
	switch err.(type) {
	case service.NotFound:
		return status.Error(codes.NotFound, err.Error())
	case service.BadRequest:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
