package request

import (
	"reflect"
	"strings"
)

// getNameFromTag return the value in tag or field name in the struct field.
func getNameFromTag(f reflect.StructField, tagName string) string {
	tag, _ := parseTag(f, tagName)
	if tag != "" {
		return tag
	}
	return f.Name
}

// parseTag returns the first tag value of the struct field.
func parseTag(f reflect.StructField, tag string) (string, tagOptions) {
	return splitTags(getTagValues(f, tag))
}

// getTagValues returns tag value of the struct field.
func getTagValues(f reflect.StructField, tag string) string {
	return f.Tag.Get(tag)
}

// splitTags returns the first tag value and slice of rest.
func splitTags(tags string) (string, tagOptions) {
	res := strings.Split(tags, ",")
	return res[0], res[1:]
}

// tagOptions is wrapper struct for rest tag values.
type tagOptions []string

// has checks the value exists in the rest values or not.
func (t tagOptions) has(tag string) bool {
	for _, opt := range t {
		if opt == tag {
			return true
		}
	}
	return false
}
