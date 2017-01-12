package dbfgo

import (
	"os"
	"testing"
)

func TESTdbfgo1(t *testing.T) {
	var tdbf Dbfgo
	nf, err := os.Open("test.dbf")
	if err != nil {
		t.Fatal("Test File not found")
	}
	defer nf.Close()
	tdbf.GetDBFInfo(nf)
	if tdbf.Header.Records != 3 {
		t.Fatal("Numero sbagliato di record")
	}
}
