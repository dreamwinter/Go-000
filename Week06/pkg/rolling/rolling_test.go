package rolling

import (
	"testing"
	"time"
)

func TestNumber_Max(t *testing.T) {
	n := NewNumber(10)
	for _, x := range []float64{10, 11, 9} {
		n.UpdateMax(x)
		time.Sleep(1 * time.Second)
	}
	tests := []struct {
		name string
		r    *Number
		want float64
	}{
		{
			name: "Happy Path",
			r:    n,
			want: float64(11),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Max(time.Now()); got != tt.want {
				t.Errorf("Number.Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Avg(t *testing.T) {
	n := NewNumber(20)
	for _, x := range []float64{0, 0.5, 1.5, 2.5, 3.5, 4.5} {
		n.Increment(x)
		time.Sleep(1 * time.Second)
	}
	tests := []struct {
		name string
		r    *Number
		want float64
	}{
		{
			name: "Happy Path",
			r:    n,
			want: float64(0.625),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Avg(time.Now()); got != tt.want {
				t.Errorf("Number.Max() = %v, want %v", got, tt.want)
			}
		})
	}
}
