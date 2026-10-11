package fa3

import (
	"path/filepath"
	"testing"
)

func TestValidateFixture(t *testing.T) {
	path := filepath.Join("..", "..", "fixtures", "invoice-fa3.xml")
	f, err := ValidateFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(f) > 0 {
		t.Fatalf("expected no findings for fixture, got %+v", f)
	}
}

func TestValidateInvalid(t *testing.T) {
	_, err := ValidateXML([]byte("not xml"))
	if err == nil {
		t.Fatal("expected error")
	}
}
