package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMainFunction(t *testing.T) {
	t.Run("Should print error message when no zip code is provided", func(t *testing.T) {
		// Simulate no arguments
		os.Args = []string{"main"}
		output := captureOutput(main)
		expected := "Add the zip code to the command: go run main.go <zip code>\n"
		if output != expected {
			t.Errorf("Expected %q, got %q", expected, output)
		}
	})

	t.Run("Should print address when address fetching succeeds", func(t *testing.T) {
		// simulate valid arguments
		os.Args = []string{"main", "01001000"}
		output := captureOutput(main)
		expected := "Received from brasilapi: source:BrasilAPI - Praça da Sé, Sé - São Paulo, SP, 01001000" // Replace with actual address format
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q, got %q", expected, output)
		}
	})

	t.Run("Should print timeout message when API call times out", func(t *testing.T) {
		// simulate valid arguments
		os.Args = []string{"main", "01001000"}
		// override apiTimeout to a very short duration to force timeout
		apiTimeout = 1 * time.Nanosecond
		output := captureOutput(main)
		expected := "Timeout\n"
		if output != expected {
			t.Errorf("Expected %q, got %q", expected, output)
		}
	})
}

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = stdout

	return buf.String()
}
