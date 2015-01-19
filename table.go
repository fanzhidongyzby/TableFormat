package table

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

//option config parameters
var (
	//separate rows
	RowSeparator string = "\n"

	//separate columns, empty string means all the space characters
	ColumnSeparator string = ""

	//what means a empty table field
	Placeholder string = "_"

	//what to be filled in blank table field when row's too short
	BlankFilling string = ""

	//what to be filled in blank header field when row's too long
	BlankFillingForHeader string = ""

	//discard more columns or not when row's too long
	ColOverflow bool = true

	//use utf8 character to print board
	UseBoard bool = true

	//what to replace \n \b \t ...
	SpaceAlt byte = ' '

	//what to separator overflow columns
	OverFlowSeparator string = " "

	//what to fill into field in order to centralize
	CenterFilling byte = ' '
)

//split str and filt empty line
func getLines(str string) []string {
	var lines []string
	if RowSeparator == "" {
		lines = strings.Fields(str)
	} else {
		lines = strings.Split(str, RowSeparator)
	}

	//filt empty string
	ret := []string{}
	for _, f := range lines {
		if len(f) > 0 {
			ret = append(ret, f)
		}
	}
	return ret
}

//split line and filt empty elements
func getFields(line string) []string {
	var fields []string
	if ColumnSeparator == "" {
		fields = strings.Fields(line)
	} else {
		fields = strings.Split(line, ColumnSeparator)
	}

	//filt empty string
	ret := []string{}
	for _, f := range fields {
		if len(f) > 0 {
			ret = append(ret, f)
		}
	}
	return ret
}

//change all the space character (\t \n _ \b) to space
func handleSpace(str string) string {
	arr := make([]rune, utf8.RuneCountInString(str))
	index := 0
	for _, c := range str {
		if unicode.IsSpace(c) && c != ' ' {
			c = rune(SpaceAlt)
		}
		arr[index] = c
		index++
	}
	return string(arr)
}

//how long is string in screen, Chinese chararter is 2 length
func width(str string) int {
	sum := 0
	for _, c := range str {
		if utf8.RuneLen(c) > 1 {
			sum += 2
		} else {
			sum++
		}
	}
	return sum
}

//convert string to 2-D slice
func preProcess(data string) [][]string {
	//get non-blank lines
	lines := []string{}
	//for _, line := range strings.Split(data, RowSeparator) {
	for _, line := range getLines(data) {
		if len(getFields(line)) != 0 {
			lines = append(lines, line)
		}
	}

	//handle empty table
	rowNum := len(lines)
	if rowNum == 0 {
		//use place holder to represent a empty table
		return [][]string{{string(CenterFilling) + BlankFillingForHeader + string(CenterFilling)}}
	}
	tb := make([][]string, len(lines))

	//get columns
	colNum := len(getFields(lines[0]))
	//max width of each column
	colWidth := make([]int, colNum)

	for row, line := range lines {
		tb[row] = make([]string, colNum)

		//fillings
		filling := BlankFilling
		if row == 0 {
			filling = BlankFillingForHeader
		}

		//init row as blank filling
		for index, _ := range tb[row] {
			tb[row][index] = filling
		}

		//get fields
		fields := getFields(line)
		for col, val := range fields {
			//handle placeholder
			if val == Placeholder {
				val = filling
			}

			//handle column overflow
			if col >= colNum {
				if ColOverflow {
					col = colNum - 1
					val = tb[row][col] + OverFlowSeparator + val
				} else {
					//discard more cols
					break
				}
			}
			tb[row][col] = handleSpace(val)
		}
	}

	//calcu max width, extend colwidth + 2 to store blank
	for col := 0; col < colNum; col++ {
		for row := 0; row < rowNum; row++ {
			val := tb[row][col]
			size := width(val)
			if size > colWidth[col] {
				colWidth[col] = size
			}
		}
		colWidth[col] += 2
	}

	//middle value with blank
	cfill := string(CenterFilling)
	for row, line := range tb {
		for col, val := range line {
			size := width(val)
			left := (colWidth[col] - size) / 2
			right := colWidth[col] - size - left
			tb[row][col] = strings.Repeat(cfill, left) + val + strings.Repeat(cfill, right)
		}
	}

	return tb

}

//utf8 table characters
const (
	hrLine = "─"
	vtLine = "│"

	topLeft   = "┌"
	topCenter = "┬"
	topRight  = "┐"

	middleLeft   = "├"
	middleCenter = "┼"
	middleRight  = "┤"

	bottomLeft   = "└"
	bottomCenter = "┴"
	bottomRight  = "┘"
)

//form table line
func initLine(left, center, right string, fill []string) []string {
	colNum := len(fill)*2 + 1
	line := make([]string, colNum)
	for i, _ := range line {
		tmp := ""
		switch {
		case i == 0:
			tmp = left
		case i == colNum-1:
			tmp = right
		case i%2 == 0:
			tmp = center
		default:
			tmp = fill[i/2]
		}
		line[i] = tmp
	}
	return line
}

//format with board
func boardFormat(tb [][]string) string {
	if len(tb) == 0 {
		tb = [][]string{{string(CenterFilling) + BlankFillingForHeader + string(CenterFilling)}}
	}
	//table attributes
	rowNum := len(tb)*2 + 1
	colNum := len(tb[0])*2 + 1
	colWidth := make([]int, colNum)
	for i, _ := range tb[0] {
		colWidth[i] = width(tb[0][i])
	}

	//init fill as --- ...
	fill := make([]string, colNum/2)
	for i, _ := range fill {
		fill[i] = strings.Repeat(hrLine, colWidth[i])
	}

	//init top ┌───┬───┐
	topLine := initLine(topLeft, topCenter, topRight, fill)

	//init middle ├───┼───┤
	middleLine := initLine(middleLeft, middleCenter, middleRight, fill)

	//init bottom └───┴───┘
	bottomLine := initLine(bottomLeft, bottomCenter, bottomRight, fill)

	//create board table
	table := make([][]string, rowNum)
	for i, _ := range table {
		switch {
		case i == 0:
			table[i] = topLine
		case i == rowNum-1:
			table[i] = bottomLine
		case i%2 == 0:
			table[i] = middleLine
		default:
			table[i] = initLine(vtLine, vtLine, vtLine, tb[i/2])
		}
	}

	//output table
	var buf bytes.Buffer
	for _, line := range table {
		for _, val := range line {
			buf.WriteString(val)
		}
		buf.WriteString("\n")
	}

	return buf.String()

}

//format without board
func simpleFormat(tb [][]string) string {
	if len(tb) == 0 {
		tb = [][]string{{string(CenterFilling) + BlankFillingForHeader + string(CenterFilling)}}
	}
	//out put table
	var buf bytes.Buffer
	for _, line := range tb {
		for _, val := range line {
			buf.WriteString(val)
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

//table format API
func Format(data string) string {
	//convert string to table
	tb := preProcess(data)

	//print table
	if UseBoard {
		return boardFormat(tb)
	} else {
		return simpleFormat(tb)
	}
}

//table format map
func FormatMap(m map[string]string) string {
	//convert map to string
	var buf bytes.Buffer
	buf.WriteString("Key\vValue\n")
	for key, value := range m {
		buf.WriteString(key + "\v" + value + "\n")
	}

	//format string
	return Format(buf.String())
}

//convert interface for user
type Convertable interface {
	Convert(field interface{}, typeStr string) string
}

//table format object
func FormatObj(obj interface{}) (str string) {
	return Format(encodeObj(obj))
}

//encode object to text
func encodeObj(obj interface{}) (str string) {
	//ignore all the panic
	defer func() {
		if r := recover(); r != nil {

			fmt.Println(r)
			str = ""
		}
	}()

	//reflect type and value
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var buf bytes.Buffer
	buf.WriteString("Name Value\n")

	//format struct
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			//get field name and value
			name := field.Name
			value := v.FieldByName(field.Name).Interface()

			//process tag: `table:"-|<newName>[,<newType>]"`
			tag := field.Tag.Get("table")
			if tag != "" {
				//ignore field
				if tag == "-" {
					continue
				}

				//tokenize
				cmds := strings.Split(tag, ",")
				num := len(cmds)
				if num > 0 && cmds[0] != "" {
					name = cmds[0]
				}
				if num > 1 && cmds[1] != "" {
					if o, ok := obj.(Convertable); ok {
						value = o.Convert(value, cmds[1])
					} 
				}
			}

			buf.WriteString(fmt.Sprintf("%s\v%v\n", name, value))
		}
	}

	return buf.String()
}
