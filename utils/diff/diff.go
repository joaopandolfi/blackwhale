package diff

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/joaopandolfi/blackwhale/v2/utils/snake_case"
)

// Diff structs tool
// returns the b value
func Diff(a, b any) (map[string]any, error) {
	return diff(a, b, false)
}

func DiffIgnoreVoid(a, b any) (map[string]any, error) {
	return diff(a, b, true)
}

func diff(a, b any, igonoreVoid bool) (map[string]any, error) {
	var A, B, mdiff map[string]any

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return nil, fmt.Errorf("a and b need to be the same type: a (%v) b(%v)", reflect.TypeOf(a), reflect.TypeOf(b))
	}

	sa, err := snake_case.JsonMarshal(a)
	//sa, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("first argument is not marshable: %w", err)
	}

	json.Unmarshal(sa, &A)

	sb, err := snake_case.JsonMarshal(b)
	//sb, err := json.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("second argument is not marshable: %w", err)
	}

	json.Unmarshal(sb, &B)

	mdiff = map[string]any{}

	for k := range A {
		if reflect.ValueOf(B[k]).Kind() == reflect.Invalid || reflect.ValueOf(A[k]).Kind() == reflect.Invalid {
			continue
		}

		if reflect.ValueOf(A[k]).Kind() == reflect.Map && reflect.ValueOf(B[k]).Kind() == reflect.Map {
			val, _ := diff(A[k], B[k], igonoreVoid)
			if len(val) > 0 {
				mdiff[k] = val
			}
			continue
		}

		if B[k] != nil && A[k] != B[k] {
			if igonoreVoid && isVoid(A[k]) {
				continue
			}

			if !isVoid(B[k]) {
				mdiff[k] = B[k]
			}
		}
	}

	return mdiff, nil
}

func isVoid(v any) bool {
	switch reflect.TypeOf(v).Kind() {
	case reflect.String:
		vs := v.(string)
		return vs == "" || vs == "0001-01-01T00:00:00Z"
	case reflect.Int:
		return v.(int) == 0
	case reflect.Float64:
		return v.(float64) == 0
	case reflect.Invalid:
		return true
	}
	return false
}

func ExtractFromKeyMap(m map[string]any, v any) map[string]any {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return ExtractMap(keys, v)
}

func ExtractMap(keys []string, v any) map[string]any {
	temp := map[string]any{}
	result := map[string]any{}

	b, _ := snake_case.JsonMarshal(&v)
	json.Unmarshal(b, &temp)

	for _, k := range keys {
		if temp[k] != nil && !isVoid(temp[k]) {
			result[k] = temp[k]
		}
	}

	return result
}
