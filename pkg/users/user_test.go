package users

import "testing"

func TestUser_IsRole(t *testing.T) {
	t.Parallel()

	type args struct {
		roles []Role
	}

	tests := []struct {
		name string
		user User
		args args
		want bool
	}{
		{
			name: "user is of role",
			user: User{Username: "ben", Role: Buyer},
			args: args{
				roles: []Role{Buyer, Seller},
			},
			want: true,
		},
		{
			name: "user is not of role",
			user: User{Username: "ben", Role: Buyer},
			args: args{
				roles: []Role{Seller},
			},
			want: false,
		},
		{
			name: "user is not of role",
			user: User{Username: "ben", Role: Buyer},
			args: args{
				roles: []Role{Seller},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.user.IsRole(tt.args.roles...); got != tt.want {
				t.Errorf("IsRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
