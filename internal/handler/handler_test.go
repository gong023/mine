package handler

import "testing"

func TestRound(t *testing.T) {
	cases := []struct {
		src  float64
		want float64
	}{
		{
			src:  10.994,
			want: 10.99,
		},
		{
			src:  10.9949999999999,
			want: 10.99,
		},
		{
			src:  10.995,
			want: 11.00,
		},
		{
			src:  -10.995,
			want: -11.00,
		},
		{
			src:  0.0,
			want: 0.0,
		},
	}
	for _, c := range cases {
		t.Run(t.Name(), func(t *testing.T) {
			got := round(c.src)
			if c.want != got {
				t.Fatalf("want: %f, got: %f", c.want, got)
			}
		})
	}
}
