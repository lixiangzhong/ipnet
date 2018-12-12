package ipnet

import "testing"

func Test_isNumber(t *testing.T) {

	tests := []struct {
		name string
		args string
		want bool
	}{
		{args: "123", want: true},
		{args: "1234567890", want: true},
		{args: "123.", want: false},
		{args: "a", want: false},
		{args: "1a", want: false},
		{args: "b123", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNumber(tt.args); got != tt.want {
				t.Errorf("isNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
