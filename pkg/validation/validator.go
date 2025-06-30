package validation

import (
	"fmt"
	"reflect"
)

// 定义验证函数类型
type ValidatorFunc func(value any) error

// 定义验证规则类型
type Rules map[string]ValidatorFunc

func ValidateAllFields(obj any, rules Rules) error {
	return ValidateSelectedFields(obj, rules, GetExportedFieldNames(obj)...)
}

// 通用校验函数
func ValidateSelectedFields(obj any, rules Rules, fields ...string) error {
	// 通过反射获取结构体的值和类型
	objValue := reflect.ValueOf(obj)
	objType := reflect.TypeOf(obj)
	// 确保传入是一个结构体
	if objType.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
		objType = objType.Elem()
	}

	if objType.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct,got %s", objType.Kind())
	}

	for _, field := range fields {
		// 检查字段是否存在
		structField, exists := objType.FieldByName(field)
		if !exists || !structField.IsExported() {
			continue // 跳过不存在或未导出的字段
		}
		// 提取字段值
		fieldValue := objValue.FieldByName(field)
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				// 如果是指针类型但为nil，返回错误
				// return fmt.Errorf("field '%s' cannot be nil",field)
				// 如果是指针类型但为nil，跳过验证
				continue
			}
			fieldValue = fieldValue.Elem()
		}
		validator, ok := rules[field]
		if !ok {
			continue
		}
		if err := validator(fieldValue.Interface()); err != nil {
			return err
		}
	}
	return nil
}

// 返回传入结构体中所有可导出的字段名字
func GetExportedFieldNames(obj any) []string {
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)
	// 校验是否为结构体或结构体指针
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return []string{}
	}

	var fieldNames []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.IsExported() {
			fieldNames = append(fieldNames, field.Name)
		}
	}
	return fieldNames
}
