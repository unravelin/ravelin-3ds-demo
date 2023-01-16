package handler

import "testing"

func Test_getLastFour(t *testing.T) {
	tests := []struct {
		pan  string
		want string
	}{
		{pan: "0123456789", want: "6789"},
		{pan: "0000000000001234", want: "1234"},
		{pan: "", want: ""},
		{pan: "1234", want: "1234"},
		{pan: "12", want: "12"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := getLastFour(tt.pan); got != tt.want {
				t.Errorf("getLastFour() = %v, want %v", got, tt.want)
			}
		})
	}
}
