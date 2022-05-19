package gocsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type option struct {
	SuppressError bool
}

func WithSuppressError(suppressError bool) func(*option) {
	return func(o *option) {
		o.SuppressError = suppressError
	}
}

func Read[T any](r io.Reader, options ...func(*option)) ([]*T, error) {
	cr := csv.NewReader(r)
	records, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}

	var o option
	for _, f := range options {
		f(&o)
	}

	var out T

	// only allow struct type, TODO: how to specify type restriction with generic
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("invalid generic type %v", rv.Kind()))
	}

	// get pointer of type T instance
	// which is addressable and could set value into it
	p := &out

	ret := make([]*T, 0)

	// empty file
	if len(records) <= 1 {
		return nil, nil
	}

	headers := records[0]

	e := reflect.ValueOf(p).Elem()

	mapping := make(map[int]int, 0)

	// parse the struct, get the name of fields,
	// build the cache for headers index mapping to field name
	for i := 0; i < e.NumField(); i++ {
		field := e.Type().Field(i)
		name := field.Name
		override := field.Tag.Get("csv")

		if override == "" {
			idx := findByName(headers, name)
			if idx != -1 {
				// store record index to field index cache
				mapping[idx] = i
			}
			continue
		}

		for _, name := range strings.Split(override, ",") {
			idx := findByName(headers, name)
			if idx != -1 {
				// store record index to field index cache
				mapping[idx] = i
			}
		}
	}

	for i := 1; i < len(records); i++ {
		for idx, val := range records[i] {
			if i, ok := mapping[idx]; ok {
				if !e.Field(i).CanSet() {
					// this should not happen
					panic("reflection can not set")
				}

				// handle empty value
				if val == "" {
					e.Field(i).Set(e.Field(i))
					continue
				}

				switch e.Type().Field(i).Type.Kind() {
				case reflect.String:
					e.Field(i).SetString(val)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					x, err := strconv.ParseUint(val, 10, 64)
					if err != nil {
						if o.SuppressError {
							log.Printf("[gocsv] error: %+v\n", err)
							// set default value
							e.Field(i).Set(e.Field(i))
							continue
						}
						return nil, err
					}
					e.Field(i).SetUint(x)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					x, err := strconv.ParseInt(val, 10, 64)
					if err != nil {
						if o.SuppressError {
							log.Printf("[gocsv] error: %+v\n", err)
							// set default value
							e.Field(i).Set(e.Field(i))
							continue
						}
						return nil, err
					}
					e.Field(i).SetInt(x)
				case reflect.Float32, reflect.Float64:
					x, err := strconv.ParseFloat(val, 64)
					if err != nil {
						if o.SuppressError {
							log.Printf("[gocsv] error: %+v\n", err)
							// set default value
							e.Field(i).Set(e.Field(i))
							continue
						}
						return nil, err
					}
					e.Field(i).SetFloat(x)
				case reflect.Bool:
					x, err := strconv.ParseBool(val)
					if err != nil {
						if o.SuppressError {
							log.Printf("[gocsv] error: %+v\n", err)
							// set default value
							e.Field(i).Set(e.Field(i))
							continue
						}
						return nil, err
					}
					e.Field(i).SetBool(x)
				case reflect.Complex64, reflect.Complex128:
				case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
					// TODO: handle Ptr
				case reflect.Ptr:
					return nil, errors.New("not implement")
				}

			}
		}

		ret = append(ret, p)
	}

	return ret, nil
}

func findByName(src []string, name string) int {
	for idx, key := range src {
		if strings.EqualFold(name, key) {
			return idx
		}
	}
	return -1
}
