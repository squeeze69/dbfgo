package dbfgo

import (
	"os"
	"testing"
)

func TestDbfgo1(t *testing.T) {
	var tdbf Dbfgo
	nf, err := os.Open("test.dbf")
	if err != nil {
		t.Fatal("Test File not found")
	}
	defer nf.Close()
	tdbf.GetDBFInfo(nf)
	if tdbf.Header.Records != 4 {
		t.Fatal("Numero sbagliato di record ",tdbf.Header.Records," Expected: 4")
	}
}
