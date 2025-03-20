package util

import (
	"database/sql"
	"reflect"
	"time"
)

const DateTime = "2006-01-02 15:04:05"

func CopyFields(src interface{}, dst interface{}) {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst).Elem()

	// 检查 src 是否为指针类型，如果是则解引用
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		dstField := dstVal.FieldByName(srcVal.Type().Field(i).Name)
		if dstField.IsValid() && dstField.CanSet() {
			switch srcField.Type() {
			case dstField.Type():
				dstField.Set(srcField)
			case reflect.TypeOf(time.Time{}):
				// 处理 time.Time 到 string 的转换
				if !srcField.IsZero() {
					dstField.SetString(srcField.Interface().(time.Time).Format(DateTime))
				}
			case reflect.TypeOf(sql.NullString{}):
				// 处理 sql.NullString 到 string 的转换
				dstField.SetString(srcField.Field(0).String())
			case reflect.TypeOf(sql.NullTime{}):
				if !srcField.IsZero() && srcField.Field(1).Bool() {
					dstField.SetString(srcField.Field(0).Interface().(time.Time).Format(DateTime))
				}
			}

		}
	}
}
