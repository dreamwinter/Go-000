package service

import (
	"context"
	"log"

	"xll.com/go-000/Week04/internal/data"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"xll.com/go-000/Week04/api"
	"xll.com/go-000/Week04/internal/biz"
)

// UserService is controller of user related functions
type UserService struct {
	api.UnimplementedUserServer
	UserRepo biz.UserRepository
}

// CreateUser is
func (s *UserService) CreateUser(ctx context.Context, in *api.CreateUserReq) (*api.CreateUserRes, error) {
	user, err := s.UserRepo.CreateUser(ctx, biz.UserDO{
		ID:     in.GetUser().Id,
		Name:   in.GetUser().Name,
		Gender: in.GetUser().Gender,
		Age:    int(in.GetUser().Age),
	})
	if err != nil {
		log.Printf("Error encountered during create user:%+v, error: %v\n", in.GetUser(), err)
		if data.IsErrUserAlreadyExisted(err) {
			return nil, status.Errorf(codes.AlreadyExists, "User already exists")
		}
		return nil, status.Errorf(codes.Internal, "Internal Server error")
	}
	return &api.CreateUserRes{
		User: &api.UserMessage{
			Id:     user.ID,
			Name:   user.Name,
			Gender: user.Gender,
			Age:    int32(user.Age),
		},
	}, nil
}

// GetUser is
func (s *UserService) GetUser(ctx context.Context, in *api.GetUserReq) (*api.GetUserRes, error) {
	user, err := s.UserRepo.GetUser(ctx, in.GetId())
	if err != nil {
		log.Printf("Error encountered during get user with id:%v, error: %v\n", in.GetId(), err)
		if data.IsErrUserNotFound(err) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "Internal Server error")
	}
	return &api.GetUserRes{
		User: &api.UserMessage{
			Id:     user.ID,
			Name:   user.Name,
			Gender: user.Gender,
			Age:    int32(user.Age),
		},
	}, nil
}

// UpdateUser is
func (s *UserService) UpdateUser(ctx context.Context, in *api.UpdateUserReq) (*api.UpdateUserRes, error) {
	user, err := s.UserRepo.UpdateUser(ctx, biz.UserDO{
		ID:     in.GetUser().Id,
		Name:   in.GetUser().Name,
		Gender: in.GetUser().Gender,
		Age:    int(in.GetUser().Age),
	})
	if err != nil {
		log.Printf("Error encountered during update user:%+v, error: %v\n", in.GetUser(), err)
		if data.IsErrUserNotFound(err) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "Internal Server error")
	}
	return &api.UpdateUserRes{
		User: &api.UserMessage{
			Id:     user.ID,
			Name:   user.Name,
			Gender: user.Gender,
			Age:    int32(user.Age),
		},
	}, nil
}

// DeleteUser is
func (s *UserService) DeleteUser(ctx context.Context, in *api.DeleteUserReq) (*api.DeleteUserRes, error) {
	err := s.UserRepo.DeleteUser(ctx, in.GetId())
	if err != nil {
		log.Printf("Error encountered during delete user with id:%v, error: %v\n", in.GetId(), err)
		if data.IsErrUserNotFound(err) {
			return &api.DeleteUserRes{
				Success: false,
			}, status.Errorf(codes.NotFound, err.Error())
		}
		return &api.DeleteUserRes{
			Success: false,
		}, status.Errorf(codes.Internal, "Internal Server error")
	}
	return &api.DeleteUserRes{
		Success: true,
	}, nil
}

// ListUsers is
func (s *UserService) ListUsers(req *api.ListUsersReq, stream api.User_ListUsersServer) error {
	userList, err := s.UserRepo.ListUsers(context.Background())
	if err != nil {
		log.Printf("Error encountered during get user list, error: %v\n", err)
		return status.Errorf(codes.Internal, "Internal Server error")
	}

	for _, u := range userList {
		stream.Send(&api.ListUsersRes{
			User: &api.UserMessage{
				Id:     u.ID,
				Name:   u.Name,
				Gender: u.Gender,
				Age:    int32(u.Age),
			},
		})
	}
	return nil
}

// NewUserService is the constructor function
func NewUserService(repo biz.UserRepository) api.UserServer {
	return &UserService{
		UserRepo: repo,
	}
}
