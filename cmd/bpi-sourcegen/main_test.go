package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestContractTreeFilesNormalizesJSONLineEndings(t *testing.T) {
	root := t.TempDir()
	jsonPath := filepath.Join(root, "sample.json")
	binaryPath := filepath.Join(root, "sample.bin")
	jsonCRLF := []byte("{\r\n  \"code\": 0\r\n}\r\n")
	binary := []byte{0x00, '\r', '\n', 0xff}
	if err := os.WriteFile(jsonPath, jsonCRLF, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binaryPath, binary, 0o600); err != nil {
		t.Fatal(err)
	}

	files, err := contractTreeFiles(root)
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte("{\n  \"code\": 0\n}\n"); !bytes.Equal(files["sample.json"], want) {
		t.Fatalf("JSON bytes = %q, want %q", files["sample.json"], want)
	}
	if !bytes.Equal(files["sample.bin"], binary) {
		t.Fatalf("binary bytes = %v, want %v", files["sample.bin"], binary)
	}
}
