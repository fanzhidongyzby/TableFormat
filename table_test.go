package table

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestFormat(t *testing.T) {

	//ColOverflow = false
	//UseBoard = false
	//RowSeparator = " "
	//Placeholder = "&"
	//ColumnSeparator = "%"
	//BlankFilling = "nil"
	//BlankFillingForHeader = "<NULL>"
	//SpaceAlt = '^'
	//CenterFilling = '*'
	//OverFlowSeparator = "->"

	str := ` ID _ Num Digit
	1 2 3你好
	4 _ 5 
	7 8 9 10 11`
	fmtStr := Format(str)
	fmt.Print(fmtStr)
	if fmtStr == "" {
		t.Log("nothing returned by format")
	}

	m := map[string]string{"key1xxxxxx": "value1", "key2": "value2"}
	fmtStr = FormatMap(m)
	fmt.Print(fmtStr)
	if fmtStr == "" {
		t.Log("nothing returned by format")
	}

	o := Obj{Key: "NOW", Value: time.Now().UnixNano() / 1e6, Default: []int{1, 2, 3}}
	fmtStr = FormatObj(o)
	fmt.Print(fmtStr)
	if fmtStr == "" {
		t.Log("nothing returned by format")
	}
}

//object format definition
type Obj struct {
	Key     string `table:"Name"`
	Value   int64  `table:"Time,time"`
	Options string `table:"-"`
	Default []int
}

//user-define type convertion
func (this Obj) Convert(field interface{}, typeStr string) (str string) {
	switch typeStr {
	case "time":
		if val, ok := field.(int64); ok {
			str = time.Unix(val/1e3, val%1e3*1e6).Format("2006-01-02 15:04:05")
		} else {
			str = strconv.FormatInt(val, 10)
		}
	default:
		fmt.Println("type " + typeStr + " is not supported")
	}
	return str
}
