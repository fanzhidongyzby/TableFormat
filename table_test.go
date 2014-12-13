package table

import (
	"fmt"
	"testing"
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
}
