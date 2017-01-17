package request

import "reflect"

// convertTo converts struct to any data format(result).
func convertTo(tagName string, param interface{}, result interface{}, fn func(result interface{}, key string, value interface{})) {
	t := toType(param)
	values := toValue(param)
	resultSeed := typeCopy(result)
	for i, max := 0, t.NumField(); i < max; i++ {
		f := t.Field(i)
		if f.PkgPath != "" && !f.Anonymous {
			continue // skip private field
		}
		tag, opts := parseTag(f, tagName)
		if tag == "-" {
			continue // skip `-` tag
		}

		v := values.Field(i)
		if opts.has("omitempty") && isZero(v) {
			continue // skip zero-value when omitempty option exists in tag
		}
		if opts.has("squash") {
			convertTo(tagName, v.Interface(), result, fn)
			continue
		}

		name := getNameFromTag(f, tagName)
		if opts.has("recursive") {
			copy := typeCopy(resultSeed)
			convertTo(tagName, v.Interface(), copy, fn)
			fn(result, name, copy)
			continue
		}

		fn(result, name, v.Interface())
	}
}

// toValue converts any value to reflect.Value.
func toValue(p interface{}) reflect.Value {
	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// toType converts any value to reflect.Type.
func toType(p interface{}) reflect.Type {
	t := reflect.ValueOf(p).Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// isZero checks the value is zero-value or not.
func isZero(v reflect.Value) bool {
	zero := reflect.Zero(v.Type()).Interface()
	value := v.Interface()
	return reflect.DeepEqual(value, zero)
}

// typeCopy returns the copy of zero value from given type.
func typeCopy(v interface{}) interface{} {
	rt := reflect.ValueOf(v).Type()
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	var rv reflect.Value
	switch rt.Kind() {
	case reflect.Map:
		rv = reflect.MakeMap(rt)
	case reflect.Slice:
		rv = reflect.MakeSlice(rt, 0, 0)
	default:
		rv = reflect.New(rt)
	}

	return rv.Interface()
}
