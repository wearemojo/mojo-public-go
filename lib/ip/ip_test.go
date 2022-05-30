package ip

import (
	"testing"
)

func TestGetIP(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "ipv4",
			in:   "192.168.254.5",
			out:  "192.168.254.5",
		},
		{
			name: "ipv4 with port",
			in:   "127.0.0.1:2564",
			out:  "127.0.0.1",
		},
		{
			name: "ipv6",
			in:   "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			out:  "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		},
		{
			name: "ipv6 with brackets",
			in:   "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
			out:  "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		},
		{
			name: "ipv6 with port",
			in:   "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:4553",
			out:  "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			ip := GetIP(test.in)

			if test.out != ip {
				t.Errorf("expected %s, got %s", test.out, ip)
			}
		})
	}
}
