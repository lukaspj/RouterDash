package main

import (
	"os/exec"
	"testing"
)

func TestJavaScriptSyntax(t *testing.T) {
	// This calls "bun test" on the JS validator, which extracts and
	// syntax-checks the embedded <script> block in static/index.html.
	// If the JS has a parse error (stray brace, missing catch, etc.)
	// the test fails with a clear message pointing at the broken line.
	cmd := exec.Command("bun", "test", "static/js/validate.test.js")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("JavaScript validation failed:\n%s", string(out))
	}
	t.Log(string(out))
}
