package transformers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/joaopandolfi/blackwhale/remotes/hasura"
)

// QueryResultTo converts an hasura query result into a given struct
func QueryResultTo[T any](result *hasura.QueryResult) (T, error) {
	var data T
	typeOf := reflect.TypeOf(data)
	key := typeOf.Name()
	if key == "" {
		key = typeOf.Elem().Name()
	}
	key = strings.ToLower(key)

	err := result.GetTo(key, &data)
	if err != nil {
		return data, fmt.Errorf("parsing QueryResult to %s: %w", key, err)
	}

	return data, nil
}
