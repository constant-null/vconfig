package vconfig

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func Unmarshal(out interface{}) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("Value should be a pointer")
	}

	err := unmarshal(v, tagInfo{})

	return err
}

func unmarshal(v reflect.Value, t tagInfo) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if t.Default != "" {
		viper.SetDefault(t.Name, t.Default)
	}

	if !viper.IsSet(t.Name) && t.Required {
		return fmt.Errorf("Variable %s is missing", t.Name)
	}

	switch v.Kind() {
	case reflect.Bool:
		v.Set(reflect.ValueOf(viper.GetBool(t.Name)))
	case reflect.String:
		v.Set(reflect.ValueOf(viper.GetString(t.Name)))
	case reflect.Float64:
		v.Set(reflect.ValueOf(viper.GetFloat64(t.Name)))
	case reflect.Int:
		v.Set(reflect.ValueOf(viper.GetInt(t.Name)))
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.String {
			v.Set(reflect.ValueOf(viper.GetStringSlice(t.Name)))
		}
	case reflect.Struct:
		st := v.Type()
		for i := 0; i < st.NumField(); i++ {
			if !v.Field(i).CanSet() {
				continue
			}

			ft := parseTag(st.Field(i), t.Name)
			err := unmarshal(v.Field(i), ft)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type tagInfo struct {
	Name     string
	Default  string
	Required bool
}

func parseTag(f reflect.StructField, prefix string) tagInfo {
	name := f.Tag.Get("vconfig")
	if name == "" {
		name = f.Name
	}

	tag := tagInfo{}
	tag.Name = strings.Trim(strings.Join([]string{prefix, strings.ToLower(name)}, "."), ".")
	tag.Required, _ = strconv.ParseBool(f.Tag.Get("required"))
	tag.Default = f.Tag.Get("default")

	return tag
}
