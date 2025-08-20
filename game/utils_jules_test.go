package sgc7game

import "testing"

func Test_IsRngString_jules(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid string with numbers and commas",
			args: args{str: "1,2,3,4,5"},
			want: true,
		},
		{
			name: "valid string with numbers, commas, and spaces",
			args: args{str: "1, 2, 3, 4, 5"},
			want: true,
		},
		{
			name: "valid string with only numbers",
			args: args{str: "12345"},
			want: true,
		},
		{
			name: "valid string with only commas",
			args: args{str: ",,,"},
			want: true,
		},
		{
			name: "valid string with only spaces",
			args: args{str: "   "},
			want: true,
		},
		{
			name: "empty string",
			args: args{str: ""},
			want: true,
		},
		{
			name: "invalid string with letters",
			args: args{str: "1,2,a,4,5"},
			want: false,
		},
		{
			name: "invalid string with special characters",
			args: args{str: "1,2,$,4,5"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRngString(tt.args.str); got != tt.want {
				t.Errorf("IsRngString() = %v, want %v", got, tt.want)
			}
		})
	}
}
