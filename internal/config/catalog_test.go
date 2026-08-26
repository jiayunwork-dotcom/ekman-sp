package config

import (
	"os"
	"testing"
)

func TestCatalogFind(t *testing.T) {
	c, ok := FindExample("midlat-wind")
	if !ok || c.Tau != 0.2 {
		t.Fatalf("%v %v", ok, c)
	}
	_, ok = FindExample("missing")
	if ok {
		t.Fatal("expected missing")
	}
}

func TestLoadExample(t *testing.T) {
	if err := os.Chdir("../../"); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir("internal/config")
	res, err := LoadExample("north-trade")
	if err != nil {
		t.Fatal(err)
	}
	if res.Params.Tau <= 0 {
		t.Fatal()
	}
}
