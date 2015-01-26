package table

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func testFormat(t *testing.T) {

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

type T struct{ A, B int }

func (this T) Convert(field interface{}, typeStr string) (str string) {
	switch typeStr {
	case "":
		str = fmt.Sprintf("<%d, %d>", this.A, this.B)
	}
	return str
}

func TestValue(t *testing.T) {
	F(1)
	F("1212")
	F(0.32)
	F(1 + 2i)
	o := T{}
	F(o)
	F([5]T{o, o, o, o})
	F("asasasasa sas")
	F(RawString{"asasasasa sas"})
	F(fmt.Println)
	ColumnSeparator = "\v"
	defer Reset()
	F(map[int]string{1: "123", 2: "234", 3: "345"})
	F(map[T]string{o: "hello"})
	F(map[T]T{o: o})
	F(&map[int]string{1: "123", 2: "234", 3: "345"})

	//fmt.Println(ty, " = ", va, runtime.FuncForPC(reflect.ValueOf(val).Pointer()).Name())
}

func F(obj interface{}) {
	fmt.Print(Format(encode(obj)))
}
func encode(obj interface{}) (str string) {
	//ignore all the panic
	defer func() {
		if r := recover(); r != nil {
			str = createEmptyHeader(1) + createRow(fmt.Sprint(r))
		}
	}()

	v := reflect.ValueOf(obj)

	return encodeAny(v)
}

type RawString struct {
	value string
}

func (this RawString) String() string {
	return this.value
}

func encodeAny(v reflect.Value) (str string) {
	t := v.Type()
	obj := v.Interface()
	if _, ok := obj.(RawString); ok {
		return encodeRawString(v)
	}

	switch t.Kind() {
	case reflect.Struct:
		str = encodeStruct(v)
	case reflect.Array, reflect.Slice:
		str = encodeList(v)
	case reflect.Map:
		str = encodeMap(v)
	case reflect.String:
		str = encodeString(v)
	case reflect.Ptr, reflect.Interface:
		str = encodeAny(v.Elem())
	case reflect.Func:
		str = encodeFunc(v)
	default:
		str = encodePlain(v)
	}
	return str
}

func encodeFunc(v reflect.Value) (str string) {
	var buf bytes.Buffer

	t := v.Type()
	if t.Kind() != reflect.Func {
		return buf.String()
	}

	buf.WriteString(createEmptyHeader(1))
	buf.WriteString(createRow(runtime.FuncForPC(v.Pointer()).Name()))

	return buf.String()
}

func encodeRawString(v reflect.Value) (str string) {
	var buf bytes.Buffer

	obj := v.Interface()

	if o, ok := obj.(RawString); ok {
		buf.WriteString(createEmptyHeader(1))
		buf.WriteString(createRow(o.String()))
	}

	return buf.String()
}

func encodeString(v reflect.Value) (str string) {
	var buf bytes.Buffer

	t := v.Type()
	obj := v.Interface()

	if t.Kind() != reflect.String {
		return buf.String()
	}

	if o, ok := obj.(string); ok {
		buf.WriteString(createRow(o))
	}

	return buf.String()
}

func createEmptyHeader(colNum int) string {
	fields := make([]string, colNum)
	for i, _ := range fields {
		fields[i] = Placeholder
	}
	return createRow(fields...)
}

func createRow(fields ...string) string {
	sep := "\v"
	if ColumnSeparator != "" {
		sep = ColumnSeparator
	}

	var buf bytes.Buffer
	for _, field := range fields {
		buf.WriteString(field + sep)
	}
	buf.WriteString("\n")

	return buf.String()
}

func encodePlain(v reflect.Value) (str string) {
	//reflect
	obj := v.Interface()

	if o, ok := obj.(Convertable); ok {
		str = o.Convert(nil, "")
	} else {
		str = fmt.Sprintf("%v", obj)
	}

	return str
}

func encodeMap(v reflect.Value) (str string) {
	var buf bytes.Buffer

	//reflect
	t := v.Type()

	if t.Kind() != reflect.Map {
		return buf.String()
	}

	buf.WriteString(createEmptyHeader(2))

	keys := v.MapKeys()
	for _, key := range keys {
		value := v.MapIndex(key)
		buf.WriteString(createRow(encodePlain(key), encodePlain(value)))
	}
	return buf.String()
}

func encodeList(v reflect.Value) (str string) {
	var buf bytes.Buffer

	//reflect
	t := v.Type()

	if t.Kind() != reflect.Array && t.Kind() != reflect.Slice {
		return buf.String()
	}

	buf.WriteString(createEmptyHeader(2))

	//format list
	size := v.Len()
	for i := 0; i < size; i++ {
		buf.WriteString(createRow(strconv.Itoa(i+1), encodePlain(v.Index(i))))
	}
	return buf.String()
}

func encodeStruct(v reflect.Value) (str string) {
	var buf bytes.Buffer

	//reflect
	t := v.Type()
	obj := v.Interface()

	if t.Kind() != reflect.Struct {
		return buf.String()
	}

	buf.WriteString(createEmptyHeader(2))

	//format struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		//get field name and value
		name := field.Name
		value := v.FieldByName(field.Name)
		val := value.Interface()

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
					val = o.Convert(val, cmds[1])
				}
			}
		} else {
			val = encodePlain(value)
		}

		buf.WriteString(createRow(name, fmt.Sprintf("%v", val)))
	}
	return buf.String()
}
