package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/atulsm/user-service/internal/repository"
	pb "github.com/atulsm/user-service/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	userRepo   repository.UserRepository
	grpcServer *grpc.Server
}

func NewServer(userRepo repository.UserRepository) *Server {
	return &Server{
		userRepo:   userRepo,
		grpcServer: grpc.NewServer(),
	}
}

func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	pb.RegisterUserServiceServer(s.grpcServer, s)

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

func (s *Server) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	users, total, err := s.userRepo.GetUsers(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}

	pbUsers := make([]*pb.User, len(users))
	for i, user := range users {
		pbUsers[i] = &pb.User{
			Id:          user.ID.String(),
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber.String,
			CreatedAt:   user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &pb.GetUsersResponse{
		Users:    pbUsers,
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
