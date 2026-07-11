// Package notify builds human-readable change notifications from domain models.
//
// Diff replaces the per-feature hand-written "if before.X != after.X { ... }"
// slack formatters: instead, each notifiable field carries a `notify:"..."` tag
// and Diff reflects over the struct to emit a "label: old => new" summary. New
// fields become notifiable by adding a tag, not by editing formatter code.
package notify

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Diff compares two values of the same struct type and returns a
// "label: old => new\n" summary of the tagged fields that changed. Only fields
// with a `notify:"..."` tag are considered.
//
// Tag grammar (comma-separated; the first element is the label):
//
//	notify:"ラベル"                     → "ラベル: <old> => <new>"
//	notify:"ラベル,unit=Kbps"           → numeric values get the unit suffix
//	notify:"ラベル,true=許可,false=禁止"  → bool values render as the given words
//	notify:"ラベル,date"                → *time.Time rendered as YYYY-MM-DD
//
// Pointers are dereferenced; if both sides are nil the field is skipped, and a
// lone nil renders as "(なし)". Untagged/unexported fields are ignored. before
// and after must be the same struct type (pointers are accepted and elem'd);
// a type mismatch or non-struct yields an empty string.
func Diff(before, after any) string {
	bv, ok1 := structValue(before)
	av, ok2 := structValue(after)
	if !ok1 || !ok2 || bv.Type() != av.Type() {
		return ""
	}

	var b strings.Builder
	t := bv.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		tag, ok := field.Tag.Lookup("notify")
		if !ok || tag == "-" {
			continue
		}
		opts := parseTag(tag)

		bVal := bv.Field(i)
		aVal := av.Field(i)
		bNil, aNil := isNilPtr(bVal), isNilPtr(aVal)
		if bNil && aNil {
			continue
		}
		if !bNil && !aNil && equalValue(bVal, aVal) {
			continue
		}
		b.WriteString(opts.label)
		b.WriteString(": ")
		b.WriteString(render(bVal, opts))
		b.WriteString(" => ")
		b.WriteString(render(aVal, opts))
		b.WriteString("\n")
	}
	return b.String()
}

type tagOpts struct {
	label     string
	unit      string
	trueText  string
	falseText string
	asDate    bool
}

func parseTag(tag string) tagOpts {
	parts := strings.Split(tag, ",")
	o := tagOpts{label: parts[0]}
	for _, p := range parts[1:] {
		switch {
		case p == "date":
			o.asDate = true
		case strings.HasPrefix(p, "unit="):
			o.unit = strings.TrimPrefix(p, "unit=")
		case strings.HasPrefix(p, "true="):
			o.trueText = strings.TrimPrefix(p, "true=")
		case strings.HasPrefix(p, "false="):
			o.falseText = strings.TrimPrefix(p, "false=")
		}
	}
	return o
}

func structValue(v any) (reflect.Value, bool) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return reflect.Value{}, false
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}
	return rv, true
}

func isNilPtr(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr && v.IsNil()
}

func deref(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}

func equalValue(a, b reflect.Value) bool {
	return render(a, tagOpts{}) == render(b, tagOpts{}) &&
		reflect.DeepEqual(deref(a).Interface(), deref(b).Interface())
}

func render(v reflect.Value, o tagOpts) string {
	if isNilPtr(v) {
		return "(なし)"
	}
	v = deref(v)
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			if o.trueText != "" {
				return o.trueText
			}
			return "true"
		}
		if o.falseText != "" {
			return o.falseText
		}
		return "false"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10) + o.unit
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10) + o.unit
	case reflect.String:
		return v.String()
	default:
		if t, ok := v.Interface().(time.Time); ok {
			if o.asDate {
				return t.Format("2006-01-02")
			}
			return t.Format(time.RFC3339)
		}
		return fmt.Sprintf("%v", v.Interface())
	}
}
