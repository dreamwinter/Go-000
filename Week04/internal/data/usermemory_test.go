package data

import (
	"context"
	"reflect"
	"testing"

	"xll.com/go-000/Week04/internal/biz"
)

func TestUserMemoryRepository_CreateUser(t *testing.T) {
	type args struct {
		user biz.UserDO
	}
	tests := []struct {
		name        string
		repo        *UserMemoryRepository
		args        args
		uuidWrapper func() (string, error)
		want        biz.UserDO
		wantErr     bool
	}{
		{
			name: "Happy Path without user.ID",
			repo: &UserMemoryRepository{
				userMap: map[string]biz.UserDO{},
			},
			args: args{
				user: biz.UserDO{
					Name:   "xll",
					Gender: "M",
					Age:    39,
				},
			},
			uuidWrapper: func() (string, error) {
				return "1", nil
			},
			want: biz.UserDO{
				ID:     "1",
				Name:   "xll",
				Gender: "M",
				Age:    39,
			},
			wantErr: false,
		},
		{
			name: "Happy Path with user.ID",
			repo: &UserMemoryRepository{
				userMap: map[string]biz.UserDO{},
			},
			args: args{
				user: biz.UserDO{
					ID:     "2",
					Name:   "xll",
					Gender: "M",
					Age:    39,
				},
			},
			want: biz.UserDO{
				ID:     "2",
				Name:   "xll",
				Gender: "M",
				Age:    39,
			},
			wantErr: false,
		},
		{
			name: "Happy Path without user.ID but got collision for UUID",
			repo: &UserMemoryRepository{
				userMap: map[string]biz.UserDO{
					"1": biz.UserDO{
						ID:     "1",
						Name:   "cdq",
						Gender: "F",
						Age:    39,
					},
				},
			},
			args: args{
				user: biz.UserDO{
					Name:   "xll",
					Gender: "M",
					Age:    39,
				},
			},
			uuidWrapper: func() func() (string, error) {
				uuids := []string{"1", "2"}
				return func() (string, error) {
					uuid := uuids[0]
					uuids = uuids[1:]
					return uuid, nil
				}
			}(),
			want: biz.UserDO{
				ID:     "2",
				Name:   "xll",
				Gender: "M",
				Age:    39,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.uuidWrapper != nil {
				tempUUIDWraper := uuidWraper
				uuidWraper = tt.uuidWrapper
				defer func() {
					uuidWraper = tempUUIDWraper
				}()
			}
			got, err := tt.repo.CreateUser(context.TODO(), tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserMemoryRepository.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserMemoryRepository.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
