package util

import (
	"strings"
	"testing"
)

func TestRedactText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"openai key", "my key is sk-abcdefghijklmnopqrstuvwxyz", "my key is openai_key= [REDACTED]"},
		{"github token", "ghp_abcdefghijklmnopqrstuvwxyztoken", "github_token= [REDACTED]"},
		{"api_key assignment", `api_key = "abcdefghijklmnopqrst"`, "api_key= [REDACTED]"},
		{"password assignment", `password = "supersecret123"`, "password= [REDACTED]"},
		{"bearer token", "Authorization: Bearer eyJhbGciOiJIUzI1.dGhpcyBpcyBh.dGVzdA", "bearer_token= [REDACTED]"},
		{"mysql dsn", "mysql://user:password@localhost:3306/db", "dsn= [REDACTED]localhost:3306/db"},
		{"aws key", "AKIA1234567890ABCDEF", "aws_key= [REDACTED]"},
		{"private key block", "private_key = \"-----BEGIN", "private_key= [REDACTED]"},
		{"plain code unchanged", "func main() { fmt.Println(\"hello\") }", "func main() { fmt.Println(\"hello\") }"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RedactText(tt.input)
			if got != tt.want {
				t.Fatalf("RedactText(%q)\n got = %q\nwant = %q", tt.input, got, tt.want)
			}
			// The redacted result must never contain the raw secret material.
			if tt.want != tt.input && strings.Contains(got, "sk-abcdefghijklmnopqrstuvwxyz") {
				t.Fatalf("raw openai key survived redaction: %q", got)
			}
		})
	}
}

func TestRedactText_DoesNotTouchNormalCode(t *testing.T) {
	code := `package main
import "fmt"
func main() {
	fmt.Println("hello world")
}`
	if got := RedactText(code); got != code {
		t.Fatalf("normal code was modified by redaction:\n%q", got)
	}
}

func TestRedactMap(t *testing.T) {
	data := map[string]interface{}{
		"title":    "report",
		"password": "supersecret",
		"nested": map[string]interface{}{
			"api_key": "sk-abc",
			"safe":    "value",
		},
	}
	got := RedactMap(data)
	if got["title"] != "report" {
		t.Fatalf("non-sensitive field altered: %v", got["title"])
	}
	if got["password"] != "[REDACTED]" {
		t.Fatalf("password not redacted: %v", got["password"])
	}
	nested, ok := got["nested"].(map[string]interface{})
	if !ok {
		t.Fatalf("nested map not preserved: %v", got["nested"])
	}
	if nested["api_key"] != "[REDACTED]" {
		t.Fatalf("nested api_key not redacted: %v", nested["api_key"])
	}
	if nested["safe"] != "value" {
		t.Fatalf("nested safe field altered: %v", nested["safe"])
	}
}
