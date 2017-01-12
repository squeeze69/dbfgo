package dbfgo

import (
	"fmt"
	"os"
	"strings"
)

//DbfHead informations about dbf file
type DbfHead struct {
	Version    []byte
	Updatedate string
	Records    int64
	Headerlen  int64
	Recordlen  int64
}

//Field field informations
type Field struct {
	Name             string
	Fieldtype        string
	FieldDataaddress []byte
	FieldLen         int64
	DecimalCount     []byte
	Workareaid       []byte
}

//Record structure holding single record
type Record struct {
	NotDeleted bool
	//Data   string
	Data map[string]string
}

//GetDbfHead retrieve the Dbf Header informations
func GetDbfHead(reader *os.File) (dbfhead DbfHead) {
	buf := make([]byte, 16)
	reader.Seek(0, 0)
	_, err := reader.Read(buf)
	if err != nil {
		panic(err)
	}
	dbfhead.Version = buf[0:1]
	dbfhead.Updatedate = fmt.Sprintf("%d", buf[1:4])
	dbfhead.Headerlen = Changebytetoint(buf[8:10])
	dbfhead.Recordlen = Changebytetoint(buf[10:12])
	dbfhead.Records = Changebytetoint(buf[4:8])
	return dbfhead
}

//CopyUpToNull - copy till the first 0 byte (excluded) - like C's strcpy
func CopyUpToNull(b []byte) (s string) {
	for _, val := range b {
		if val == 0 {
			return
		}
		s = s + string(val)
	}
	return
}

//GetFields retrieve fields informations
func GetFields(reader *os.File) []Field {
	var curbuf []byte
	dbfhead := GetDbfHead(reader)

	off := dbfhead.Headerlen - 32 - 264
	fieldlist := make([]Field, off/32)
	buf := make([]byte, off)
	_, err := reader.ReadAt(buf, 32)
	if err != nil {
		panic(err)
	}
	for i, val := range fieldlist {
		a := i * 32
		curbuf = buf[a:]
		val.Name = CopyUpToNull(curbuf[0:11])
		val.Fieldtype = fmt.Sprintf("%s", curbuf[11:12])
		val.FieldDataaddress = curbuf[12:16]
		val.FieldLen = Changebytetoint(curbuf[16:17])
		val.DecimalCount = curbuf[17:18]
		val.Workareaid = curbuf[20:21]
		fieldlist[i] = val

	}
	return fieldlist
}

//Changebytetoint - transform a byte array in an int64 value
func Changebytetoint(b []byte) (x int64) {
	for i, val := range b {
		if i == 0 {
			x = x + int64(val)
		} else {
			x = x + int64(2<<7*int64(i)*int64(val))
		}
	}

	return
}

//GetRecords retrieves ALL dbf records
func GetRecords(fp *os.File) (records map[int]Record) {
	dbfhead := GetDbfHead(fp)
	fp.Seek(0, 0)
	fields := GetFields(fp)
	recordlen := dbfhead.Recordlen
	start := dbfhead.Headerlen
	buf := make([]byte, recordlen)
	i := 1
	temp := map[int]Record{}
	for {
		_, err := fp.ReadAt(buf, start)
		if err != nil {
			return temp
		}
		record := Record{}
		if string(buf[0:1]) == " " {
			record.NotDeleted = true
		} else if string(buf[0:1]) == "*" {
			record.NotDeleted = false
		}
		tempdata := map[string]string{}
		a := int64(1)
		for _, val := range fields {
			fieldlen := val.FieldLen
			tempdata[val.Name] = strings.Trim(fmt.Sprintf("%s", buf[a:a+fieldlen]), " ")
			a = a + fieldlen
		}
		record.Data = tempdata
		temp[i] = record
		start = start + recordlen
		i++
	}
}

//GetRecordbyField - retrieve records by a field value (simple sequential search)
func GetRecordbyField(fieldname string, fieldval string, fp *os.File) (record map[int]Record) {
	fields := GetFields(fp)
	records := GetRecords(fp)
	temp := map[int]Record{}
	i := 1
	for _, val := range records {
		for _, val1 := range fields {
			if val1.Name == fieldname && val.NotDeleted {
				if val.Data[val1.Name] == fieldval || val.Data[val1.Name] == " " {
					temp[i] = val
				}
			}
		}
		i++
	}
	return temp
}
