package util

import (
	"errors"
	"fmt"
	"gokit/log"
	"reflect"
	"strconv"
	"strings"

	"github.com/ziutek/mymysql/mysql"
)

/*
ctx: 必须为指向结构体的指针
*/
func Parse(row mysql.Row, res mysql.Result, ctx interface{}) error {
	if nil == row || nil == res || nil == ctx {
		return errors.New("args cannot be nil")
	}

	tp := reflect.TypeOf(ctx)
	if tp.Kind() != reflect.Ptr {
		return errors.New("non-ptr")
	}

	val := reflect.ValueOf(ctx)
	size := tp.Elem().NumField()
	for i := 0; i < size; i++ {
		addr := val.Elem().Field(i)
		tag := tp.Elem().Field(i).Tag.Get("db")
		if "" == tag || "-" == tag {
			continue
		}
		index := res.Map(tag)
		switch tp.Elem().Field(i).Type.Kind() {
		case reflect.String:
			baz := row.Str(index)
			addr.SetString(row.Str(index))
			addr.SetString(baz)
		case reflect.Int64:
			addr.SetInt(row.Int64(index))
		case reflect.Int:
			addr.SetInt(int64(row.Int(index)))
		case reflect.Bool:
			addr.SetBool(row.Bool(index))
		case reflect.Float64, reflect.Float32:
			addr.SetFloat(row.Float(index))
		case reflect.Uint64, reflect.Uint32:
			addr.SetUint(row.Uint64(index))
		case reflect.Uint8: //for byte
			addr.SetUint(row.Uint64(index))
			//            case reflect.Struct: TODO  next version
		}
	}

	return nil
}

/*
i: 必须为结构体或者结构体指针
*/
func Columns(i interface{}) (string, error) {
	if nil == i {
		return "", errors.New("args cannot be nil")
	}

	sqlStr := ""
	tp := reflect.TypeOf(i)
	if tp.Kind() == reflect.Ptr {
		length := tp.Elem().NumField()
		for i := 0; i < length; i++ {
			baz := tp.Elem().Field(i).Tag.Get("db")
			if "" != baz && "-" != baz {
				sqlStr += baz + ","
			}
		}
		return strings.TrimRight(sqlStr, ","), nil
	}

	if tp.Kind() == reflect.Struct {
		length := tp.NumField()
		for i := 0; i < length; i++ {
			baz := tp.Field(i).Tag.Get("db")
			if "" != baz && "-" != baz {
				sqlStr += baz + ","
			}
		}
		return strings.TrimRight(sqlStr, ","), nil
	}

	log.Errorln("non-ptr or not-struct")
	return "", errors.New("non-ptr or not-struct")
}

/*
i:必须为结构体指针，或者是结构体指针的Slice
*/
func InsertSql(table string, i interface{}) (string, error) {
	typ := reflect.TypeOf(i)
	val := reflect.ValueOf(i)

	if typ.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		baz := "INSERT INTO " + table + " SET "
		length := typ.Elem().NumField()
		for i := 0; i < length; i++ {
			boo := typ.Elem().Field(i).Tag.Get("db")
			if "" != boo && "-" != boo {
				switch typ.Elem().Field(i).Type.Kind() {
				case reflect.String:
					baz += boo + "= '" + val.Elem().Field(i).String() + "',"
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					baz += boo + "= " + fmt.Sprintf("%d", val.Elem().Field(i).Int()) + ","
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					baz += boo + "= " + fmt.Sprintf("%d", val.Elem().Field(i).Uint()) + ","
				case reflect.Float32, reflect.Float64:
					baz += boo + "= " + strconv.FormatFloat(val.Elem().Field(i).Float(), 'f', 3, 64) + ","
				case reflect.Bool:
					baz += boo + "= " + fmt.Sprintf("%t", val.Elem().Field(i).Bool()) + ","
				}

			}
			continue
		}

		return strings.TrimRight(baz, ","), nil
	}

	if typ.Kind() == reflect.Slice {
		size := val.Len()
		if size <= 0 {
			return "", errors.New("empty slice")
		}

		baz := "INSERT INTO " + table + "("
		baz2 := ""
		length := val.Index(0).Type().Elem().NumField()
		for i := 0; i < length; i++ {
			baz3 := val.Index(0).Type().Elem().Field(i).Tag.Get("db")
			if "" != baz3 && "-" != baz3 {
				baz2 += baz3 + ","
			}
		}
		baz += strings.TrimRight(baz2, ",") + ") VALUES"

		baz6 := ""
		for i := 0; i < size; i++ {
			length := val.Index(i).Type().Elem().NumField()
			baz4 := "("
			for j := 0; j < length; j++ {
				baz5 := val.Index(i).Type().Elem().Field(j).Tag.Get("db")
				if "" != baz5 && "-" != baz5 {
					switch val.Index(i).Type().Elem().Field(j).Type.Kind() {
					case reflect.String:
						baz4 += "'" + val.Index(i).Elem().Field(j).String() + "',"
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						baz4 += fmt.Sprintf("%d,", val.Index(i).Elem().Field(j).Int())
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						baz4 += fmt.Sprintf("%d,", val.Index(i).Elem().Field(j).Uint())
					case reflect.Float32, reflect.Float64:
						baz4 += fmt.Sprintf("%0.3f,", val.Index(i).Elem().Field(j).Float())
					case reflect.Bool:
						baz4 += fmt.Sprintf("%t,", val.Index(i).Elem().Field(j).Bool())
					}

				}
				continue
			}

			baz6 += strings.TrimRight(baz4, ",") + "),"
		}

		return baz + strings.TrimRight(baz6, ","), nil
	}

	log.Errorln("non-ptr and not-struct-slice,accepted type:", typ.Kind().String())
	return "", errors.New("non-ptr or not-struct-slice")
}
