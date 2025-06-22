package service

import (
	"context"
	"github.com/goflower-io/example/postgres/api"
	"github.com/goflower-io/example/postgres/crud"
	"github.com/goflower-io/example/postgres/crud/user"
	"math"
	"strings"
	"time"

	"github.com/goflower-io/xsql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserServiceImpl UserServiceImpl
type UserServiceImpl struct {
	api.UnimplementedUserServiceServer
	Client *crud.Client
}

type IValidateUser interface {
	ValidateUser(a *api.User) error
}

// CreateUser CreateUser
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *api.User) (*api.User, error) {
	if checker, ok := interface{}(s).(IValidateUser); ok {
		if err := checker.ValidateUser(req); err != nil {
			return nil, err
		}
	}

	a := &user.User{
		Id:      0,
		Name:    req.GetName(),
		Age:     req.GetAge(),
		Address: req.GetAddress(),
	}
	var err error
	if a.Ctime, err = time.ParseInLocation("2006-01-02 15:04:05", req.GetCtime(), time.Local); err != nil {
		return nil, err
	}
	if a.Mtime, err = time.ParseInLocation("2006-01-02 15:04:05", req.GetMtime(), time.Local); err != nil {
		return nil, err
	}
	_, err = s.Client.User.
		Create().
		SetUser(a).
		Save(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// query after create and return
	a2, err := s.Client.Master.User.
		Find().
		Where(
			user.IdOp.EQ(a.Id),
		).
		One(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertUser(a2), nil
}

// DeleteUser DeleteUser
func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *api.UserId) (*emptypb.Empty, error) {
	_, err := s.Client.User.
		Delete().
		Where(
			user.IdOp.EQ(req.GetId()),
		).
		Exec(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// Updateuser UpdateUser
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *api.UpdateUserReq) (*api.User, error) {
	if checker, ok := interface{}(s).(IValidateUser); ok {
		if err := checker.ValidateUser(req.User); err != nil {
			return nil, err
		}
	}
	if len(req.GetMasks()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty filter condition")
	}
	update := s.Client.User.Update()
	for _, v := range req.GetMasks() {
		switch v {
		case api.UserField_User_name:
			update.SetName(req.GetUser().GetName())
		case api.UserField_User_age:
			update.SetAge(req.GetUser().GetAge())
		case api.UserField_User_address:
			update.SetAddress(req.GetUser().GetAddress())
		case api.UserField_User_ctime:
			t, err := time.ParseInLocation("2006-01-02 15:04:05", req.GetUser().GetCtime(), time.Local)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
			update.SetCtime(t)
		case api.UserField_User_mtime:
			t, err := time.ParseInLocation("2006-01-02 15:04:05", req.GetUser().GetMtime(), time.Local)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
			update.SetMtime(t)
		}
	}
	_, err := update.
		Where(
			user.IdOp.EQ(req.GetUser().GetId()),
		).
		Save(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// query after update and return
	a, err := s.Client.Master.User.
		Find().
		Where(
			user.IdOp.EQ(req.GetUser().GetId()),
		).
		One(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertUser(a), nil
}

// GetUser GetUser
func (s *UserServiceImpl) GetUser(ctx context.Context, req *api.UserId) (*api.User, error) {
	a, err := s.Client.User.
		Find().
		Where(
			user.IdOp.EQ(req.GetId()),
		).
		One(ctx)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return convertUser(a), nil
}

// ListUsers ListUsers
func (s *UserServiceImpl) ListUsers(ctx context.Context, req *api.ListUsersReq) (*api.ListUsersResp, error) {
	page := req.GetPage()
	size := req.GetPageSize()
	if size <= 0 {
		size = 20
	}
	offset := size * (page - 1)
	if offset < 0 {
		offset = 0
	}
	if len(req.GetFields()) == 0 {
		for field := range api.UserField_name {
			if field > 0 {
				req.Fields = append(req.Fields, api.UserField(field))
			}
		}
	}

	selectFields := make([]string, 0, len(req.GetFields()))
	for _, v := range req.GetFields() {
		selectFields = append(selectFields, strings.TrimPrefix(v.String(), "User_"))
	}
	finder := s.Client.User.
		Find().
		Select(selectFields...).
		Offset(offset).
		Limit(size)

	if req.GetOrderby() == api.UserField_User_unknow {
		req.Orderby = api.UserField_User_id
	}
	odb := strings.TrimPrefix(req.GetOrderby().String(), "User_")
	if req.GetDesc() {
		finder.OrderDesc(odb)
	} else {
		finder.OrderAsc(odb)
	}
	counter := s.Client.User.
		Find().
		Count()

	var ps []*xsql.Predicate
	for _, v := range req.GetFilters() {
		p, err := xsql.GenP(strings.TrimPrefix(v.Field.String(), "User_"), v.Op, v.Val)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	if len(ps) > 0 {
		p := xsql.And(ps...)
		finder.WhereP(p)
		counter.WhereP(p)
	}
	list, err := finder.All(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	count, err := counter.Int64(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pageCount := int32(math.Ceil(float64(count) / float64(size)))

	return &api.ListUsersResp{Users: convertUserList(list), TotalCount: int32(count), PageCount: pageCount, PageSize: size, Page: page}, nil
}

func convertUser(a *user.User) *api.User {
	return &api.User{
		Id:      a.Id,
		Name:    a.Name,
		Age:     a.Age,
		Address: a.Address,
		Ctime:   a.Ctime.Format("2006-01-02 15:04:05"),
		Mtime:   a.Mtime.Format("2006-01-02 15:04:05"),
	}
}

func convertUserList(list []*user.User) []*api.User {
	ret := make([]*api.User, 0, len(list))
	for _, v := range list {
		ret = append(ret, convertUser(v))
	}
	return ret
}
