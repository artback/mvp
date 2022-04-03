//go:build integration
// +build integration

package postgres_test

import (
	"context"
	"github.com/artback/mvp/pkg/api/middleware/security"
	"github.com/artback/mvp/pkg/repository/postgres"
	"github.com/artback/mvp/pkg/users"
	"reflect"
	"testing"
)

func userReposity(fn ...func(r users.Repository)) users.Repository {
	repo := postgres.UserRepository{DB: db}
	for _, f := range fn {
		f(repo)
	}
	return repo
}

func TestUserRepository_Delete(t *testing.T) {
	type args struct {
		username string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func(r users.Repository)
	}{
		{
			name: "delete user that exist",
			args: args{username: "alex"},
			setup: func(r users.Repository) {
				r.Insert(context.Background(), users.User{Username: "alex", Password: "pass", Role: security.Seller})
			},
		},
		{
			name: "delete user that don't exist",
			args: args{username: "sven"},
			setup: func(r users.Repository) {
				r.Insert(context.Background(), users.User{Username: "alex", Password: "pass", Role: security.Seller})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := userReposity(tt.setup).Delete(context.Background(), tt.args.username); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_Get(t *testing.T) {
	type args struct {
		username string
	}

	tests := []struct {
		name    string
		args    args
		want    *users.User
		setup   func(r users.Repository)
		wantErr bool
	}{
		{
			name: "Get user that exist",
			args: args{username: "alex"},
			want: &users.User{Username: "alex", Password: "pass", Role: security.Seller},
			setup: func(r users.Repository) {
				r.Insert(context.Background(), users.User{Username: "alex", Password: "pass", Role: security.Seller})
			},
		},
		{
			name: "Get user that don't exist",
			args: args{username: "sven"},
			setup: func(r users.Repository) {
				r.Insert(context.Background(), users.User{Username: "alex", Password: "pass", Role: security.Seller})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userReposity(tt.setup).Get(context.Background(), tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_Insert(t *testing.T) {
	type args struct {
		user users.User
	}

	tests := []struct {
		name    string
		args    args
		setup   func(r users.Repository)
		wantErr bool
	}{
		{
			name:  "insert user without collision",
			args:  args{user: users.User{Username: "NonExisting", Password: "pass", Role: security.Seller}},
			setup: func(r users.Repository) {},
		},
		{
			name: "insert user with collision",
			args: args{user: users.User{Username: "Existing", Password: "pass", Role: security.Seller}},
			setup: func(r users.Repository) {
				r.Insert(context.Background(), users.User{Username: "Existing", Password: "pass", Role: security.Seller})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := userReposity(tt.setup).Insert(context.Background(), tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	type args struct {
		user users.User
	}

	tests := []struct {
		name    string
		args    args
		setup   func(r users.Repository)
		wantErr bool
	}{
		{
			name: "update existing user",
			args: args{user: users.User{Username: "updateExisting", Password: "pass", Role: security.Buyer}},
			setup: func(r users.Repository) {
				r.Insert(context.Background(), users.User{Username: "updateExisting", Password: "pass", Role: security.Seller})
			},
		},
		{
			name: "update non existing user",
			args: args{user: users.User{Username: "updateNotExisting", Password: "pass", Role: security.Buyer}},
			setup: func(r users.Repository) {
			},
			wantErr: true,
		},
		{
			name: "update with empty fields",
			args: args{user: users.User{Username: "updateExistingEmpty", Role: security.Buyer}},
			setup: func(r users.Repository) {
				r.Insert(context.Background(), users.User{Username: "updateExistingEmpty", Password: "pass", Role: security.Seller})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := userReposity(tt.setup).Update(context.Background(), tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
