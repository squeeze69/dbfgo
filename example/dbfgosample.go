package main

import (
	"fmt"
	"os"

	"github.com/jeffycf/dbfgo"
)

func main() {
	var infile = "111.dbf"
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
	}
	records := dbfgo.GetRecordbyField("****", "****", fp)
	for _, val := range records {
		fmt.Println(val.Data["****"])
	}
	records1 := dbfgo.GetRecords(fp)
	for _, val := range records1 {
		if val.NotDeleted {
			fmt.Println(val.Data["CONFIGCODE"], val.Data["DIRECTORY"])
		}

	}

}
