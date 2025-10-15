package utils_test

import (
	tu "ssb/internal/commands/testUtils"
	"ssb/internal/commands/utils"
	"testing"
)

func TestBuildEndpoint(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		parts []string
	}{
		{
			name:  "empty",
			want:  "base-url",
			parts: []string{},
		},
		{
			name:  "auth",
			want:  "base-url/auth",
			parts: []string{"auth"},
		},
		{
			name:  "auth-with-id",
			want:  "base-url/auth/uuid",
			parts: []string{"auth", "uuid"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := utils.CLIConfig{
				URL:      "base-url",
				Username: "username",
			}
			tu.SetConfig(t, cfg)
			got, err := utils.BuildEndpoint(tt.parts...)
			if err != nil {
				t.Fatalf("couldn't build the url %v", err)
			}

			if tt.want != got {
				t.Fatalf("want %s, got %s", tt.want, got)
			}
		})
	}
}

func TestMakeRequest(t *testing.T) {
}
