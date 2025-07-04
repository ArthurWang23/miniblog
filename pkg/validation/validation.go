package validation

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"k8s.io/klog/v2"
)

type Validator struct {
	registry map[string]reflect.Value
}

func NewValidator(customValidator any) *Validator {
	return &Validator{registry: extractValidationMethods(customValidator)}
}

func (v *Validator) Validate(ctx context.Context, request any) error {
	validationFunc, ok := v.registry[reflect.TypeOf(request).Elem().Name()]
	if !ok {
		return nil // No validation function found for the request type
	}

	result := validationFunc.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(request)})
	if !result[0].IsNil() {
		return result[0].Interface().(error)
	}

	return nil
}

// extract and return a map of validation functions from the provided custom validation
func extractValidationMethods(customValidator any) map[string]reflect.Value {
	funcs := make(map[string]reflect.Value)
	validatorType := reflect.TypeOf(customValidator)
	validatorValue := reflect.ValueOf(customValidator)

	for i := 0; i < validatorType.NumMethod(); i++ {
		method := validatorType.Method(i)
		methodValue := validatorValue.MethodByName(method.Name)

		if !methodValue.IsValid() || !strings.HasPrefix(method.Name, "Validate") {
			continue
		}
		methodType := methodValue.Type()

		// ensure the method takes a context.Context and a pointer
		if methodType.NumIn() != 2 || methodType.NumOut() != 1 ||
			methodType.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() ||
			methodType.In(1).Kind() != reflect.Pointer {
			continue
		}

		// ensure the method name matches the expected naming convention
		requestTypeName := methodType.In(1).Elem().Name()
		if method.Name != ("Validate" + requestTypeName) {
			continue
		}

		if methodType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		klog.V(4).InfoS("Registering validator", "validator", requestTypeName)
		funcs[requestTypeName] = methodValue
	}
	return funcs
}

// ValidRequired 验证结构体中的必须字段是否存在且不为空
func ValidRequired(obj any, requiredFields ...string) error {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Struct && val.Kind() != reflect.Ptr {
		return fmt.Errorf("input must be a struct or a pointer to struct")
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("input must be a struct or a pointer to struct")
	}

	for _, field := range requiredFields {
		fieldVal := val.FieldByName(field)
		// 确保字段存在
		if !fieldVal.IsValid() {
			return fmt.Errorf("%s is not a valid field", field)
		}

		if fieldVal.IsNil() {
			return fmt.Errorf("%s is nil", field)
		}
	}
	return nil
}
