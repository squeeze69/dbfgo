//sample usage of dbfgo - forlk
//reads 111.dbf or a file passed by command line
package main

import (
	"fmt"
	"os"

	"github.com/squeeze69/dbfgo"
)

func main() {
	var infile = "111.dbf"
	namefields := make([]string, 0)
	if len(os.Args) > 1 {
		infile = os.Args[1]
	}
	fp, err := os.OpenFile(infile, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	fields := dbfgo.GetFields(fp)
	for _, val := range fields {
		fmt.Println(val.Name, val.Fieldtype, val.FieldLen)
		namefields = append(namefields, val.Name)
	}
	records1 := dbfgo.GetRecords(fp)
	for _, val := range records1 {
		if val.NotDeleted {
			for _, nm := range namefields {
				fmt.Print(val.Data[nm], ";")
			}
			fmt.Println("")
		}
	}
}
