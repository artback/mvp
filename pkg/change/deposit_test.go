package change

import (
	"github.com/artback/mvp/pkg/coin"
	"reflect"
	"testing"
)

func TestToDeposit(t *testing.T) {
	type args struct {
		coin.Coins
		amount int
	}
	tests := []struct {
		name string
		args args
		want Deposit
	}{
		{
			name: "even",
			args: args{
				Coins:  coin.Coins{5, 50, 100},
				amount: 53,
			},
			want: Deposit{50: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.Coins, tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeposit_New(t *testing.T) {
	tests := []struct {
		name string
		Deposit
		want int
	}{
		{
			name: "want 210 deposit",
			Deposit: Deposit{
				5:   2,
				100: 2,
			},
			want: 210,
		},
		{
			name: "want 190 deposit",
			Deposit: Deposit{
				50: 3,
				20: 2,
			},
			want: 190,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.Deposit.ToAmount(); got != tt.want {
				t.Errorf("ToAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}
