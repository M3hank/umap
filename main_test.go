package main

import (
	"testing"
)

func TestParametersToNameSet(t *testing.T) {
	tests := []struct {
		params   string
		expected int
	}{
		{"q=apple&page=1", 2},
		{"q=apple", 1},
		{"debug&verbose", 2},
		{"token=&user=admin", 2},
		{"=", 0},
		{"", 0},
	}

	for _, tt := range tests {
		res := parametersToNameSet(tt.params)
		if len(res) != tt.expected {
			t.Errorf("parametersToNameSet(%q) returned %d items; want %d", tt.params, len(res), tt.expected)
		}
	}
}

func TestParamNamesToString(t *testing.T) {
	params := map[string]struct{}{
		"b": {},
		"a": {},
		"c": {},
	}
	expected := "a&b&c"
	got := paramNamesToString(params)
	if got != expected {
		t.Errorf("paramNamesToString returned %q; want %q", got, expected)
	}
}

func TestHasBadExtension(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/static/style.css", true},
		{"/assets/logo.PNG", true},
		{"/docs/paper.pdf?v=1", true},
		{"/api/v1/css", false},
		{"/app.js", false},
		{"/index.html", false},
	}

	for _, tt := range tests {
		got := hasBadExtension(tt.path)
		if got != tt.expected {
			t.Errorf("hasBadExtension(%q) = %v; want %v", tt.path, got, tt.expected)
		}
	}
}

func TestNormalizePathForDedup(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/users/550e8400-e29b-41d4-a716-446655440000", "/users/{slug}"},
		{"/blog/2026-07-27-how-to-fix-bugs-in-go", "/blog/{slug}"},
		{"/wp-admin/admin-ajax.php", "/wp-admin/admin-ajax.php"},
		{"/v1/get-user-profile-info", "/v1/get-user-profile-info"},
		{"/api/v1/users", "/api/v1/users"},
	}

	for _, tt := range tests {
		got := normalizePathForDedup(tt.path)
		if got != tt.expected {
			t.Errorf("normalizePathForDedup(%q) = %q; want %q", tt.path, got, tt.expected)
		}
	}
}
