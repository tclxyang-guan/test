package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"reflect"
	"time"
)

//数据结构一模一样才能转换成功
func DataToAnyData(data interface{}, newData interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, newData)
}

//数据结构可以不一样
func MapToStruct(data map[string]interface{}, obj interface{}) error {
	for k, v := range data {
		err := setField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

//用map的值替换结构的值
func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()        //结构体属性值
	structFieldValue := structValue.FieldByName(name) //结构体单个属性值

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type() //结构体的类型
	val := reflect.ValueOf(value)              //map值的反射值

	var err error
	if structFieldType != val.Type() {
		val, err = typeConversion(fmt.Sprintf("%v", value), structFieldValue.Type().Name()) //类型转换
		if err != nil {
			return err
		}
	}
	structFieldValue.Set(val)
	return nil
}

//类型转换
func typeConversion(value string, ntype string) (reflect.Value, error) {
	switch ntype {
	case "string":
		return reflect.ValueOf(value), nil
	case "time.Time":
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	case "Time":
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	case "uint":
		v, err := cast.ToUintE(value)
		return reflect.ValueOf(v), err
	case "uint8":
		v, err := cast.ToUint8E(value)
		return reflect.ValueOf(v), err
	case "uint16":
		v, err := cast.ToUint16E(value)
		return reflect.ValueOf(v), err
	case "uint32":
		v, err := cast.ToUint32E(value)
		return reflect.ValueOf(v), err
	case "uint64":
		v, err := cast.ToUint64E(value)
		return reflect.ValueOf(v), err
	case "int":
		v, err := cast.ToIntE(value)
		return reflect.ValueOf(v), err
	case "int8":
		v, err := cast.ToInt8E(value)
		return reflect.ValueOf(v), err
	case "int16":
		v, err := cast.ToInt16E(value)
		return reflect.ValueOf(v), err
	case "int32":
		v, err := cast.ToInt32E(value)
		return reflect.ValueOf(v), err
	case "int64":
		v, err := cast.ToInt64E(value)
		return reflect.ValueOf(v), err
	case "float32":
		v, err := cast.ToFloat32E(value)
		return reflect.ValueOf(v), err
	case "float64":
		v, err := cast.ToFloat64E(value)
		return reflect.ValueOf(v), err
	case "Decimal":
		v, err := cast.ToInt64E(value)
		if err != nil {
			return reflect.ValueOf(v), err
		}
		return reflect.ValueOf(decimal.New(v, 0)), nil
	default:
		return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
	}
}
