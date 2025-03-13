package nilcheck

import "testing"

func TestAny(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []interface{}
		want bool
	}{
		{
			name: "zero args",
			args: []interface{}{},
			want: false,
		},
		{
			name: "all nils",
			args: []interface{}{nil, nil, nil},
			want: true,
		},
		{
			name: "one of the args is nil",
			args: []interface{}{
				nil, struct{}{},
			},
			want: true,
		},
		{
			name: "all args not nil",
			args: []interface{}{
				struct{}{}, struct{}{}, struct{}{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := Any(tt.args...); got != tt.want {
				t.Errorf("Any() = %v, want %v", got, tt.want)
			}
		})
	}
}
