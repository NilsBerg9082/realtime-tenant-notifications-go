package main

import "testing"

func TestOnboardingChannel(t *testing.T) {
	tests := []struct{ tenant, want string }{{"acme", "tenant-acme"}, {"northwind", "tenant-northwind"}}
	for _, tt := range tests {
		if got := onboardingChannel(tt.tenant); got != tt.want {
			t.Fatalf("%q: got %q, want %q", tt.tenant, got, tt.want)
		}
	}
}
