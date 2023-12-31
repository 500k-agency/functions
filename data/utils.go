package data

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func IsUnixZero(t time.Time) bool {
	return t.Unix() == 0
}

// GetTimeUTCPointer gets utc time pointer
func GetTimeUTCPointer() *time.Time {
	t := time.Now().UTC()
	return &t
}

func TimeToPointer(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func UnixToTimePointer(v int64) *time.Time {
	if v == 0 {
		return nil
	}
	return TimeToPointer(time.Unix(v, 0))
}

// RandString generates random sting from characters in letters
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

// Generic min
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Int64 returns a pointer to the int64 value passed in.
func Int64(v int64) *int64 {
	return &v
}

func Int64Any(v any) *int64 {
	vv, ok := v.(int64)
	if !ok {
		vy, ok := v.(*int64)
		if !ok {
			return nil
		}
		return vy
	}
	return Int64(vv)
}

// Bool returns a pointer to the bool value passed in.
func Bool(v bool) *bool {
	return &v
}

func PointerTo[T any](v T) *T {
	return &v
}

func ValueFromPtr[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

// Time returns a pointer to the time value passed in.
func Time(v time.Time) *time.Time {
	return &v
}

func Coalesce[T comparable](first T, other ...T) T {
	var zero T
	if first != zero {
		return first
	}
	for _, v := range other {
		if v != zero {
			return v
		}
	}
	return zero
}

func Filter[T any](items []T, fn func(item T) bool) []T {
	filteredItems := []T{}
	for _, value := range items {
		if fn(value) {
			filteredItems = append(filteredItems, value)
		}
	}
	return filteredItems
}

func Has[T any](items []T, fn func(item T) bool) bool {
	if len(items) == 0 {
		return false
	}
	return len(Filter(items, fn)) != 0
}

func StructToMap(v interface{}) map[string]interface{} {
	vo := reflect.ValueOf(v)

	if vo.Kind() == reflect.Map {
		// already a map
		if t, ok := v.(map[string]interface{}); ok {
			return t
		}
		return nil
	}

	if vo.Kind() == reflect.Ptr {
		vo = vo.Elem()
	}

	if vo.Kind() != reflect.Struct {
		return nil
	}

	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}

	return m
}

// Helper functions
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

var (
	titleCaser = cases.Title(language.English)
)

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToSentenceCase(str string) string {
	return titleCaser.String(strings.ReplaceAll(str, "_", " "))
}

type Name struct {
	FirstName string
	LastName  string
	FullName  string
}

func SplitName(name string) (ret Name) {
	ret.FullName = name
	if name == "" {
		return
	}
	nameParts := strings.Split(name, " ")
	if len(nameParts) == 1 {
		ret.FirstName = name
		return
	}

	ret.FirstName = nameParts[0]
	ret.LastName = strings.Join(nameParts[1:len(nameParts)-1], " ")
	return
}
