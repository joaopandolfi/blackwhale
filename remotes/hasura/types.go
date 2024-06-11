package hasura

import (
	"encoding/json"
	"fmt"

	"github.com/joaopandolfi/blackwhale/utils/snake_case"
)

const (
	Eq     = "_eq"
	In     = "_in"
	Lte    = "_lte"
	Gte    = "_gte"
	IsNull = "_is_null"
)

type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

type QueryResult struct {
	Data map[string]interface{}
}

// Gets a field from QueryResult
func (qr *QueryResult) Get(field string) interface{} {
	return qr.Data[field]
}

// Cast a field of QueryResult to given model
func (qr *QueryResult) GetTo(field string, model interface{}) error {
	data := qr.Get(field)
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling data: %w", err)
	}
	err = json.Unmarshal(bytes, model)
	if err != nil {
		return fmt.Errorf("unmarshaling to model: %w", err)
	}
	return nil
}

type Variables map[string]interface{}

// Return a Variables object from given `data`. `data` must be a pointer.
func VariablesFrom(data interface{}) (*Variables, error) {
	var vars Variables
	_json, err := snake_case.JsonMarshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling data: %w", err)
	}

	err = json.Unmarshal(_json, &vars)
	if err != nil {
		return nil, fmt.Errorf("error umarshaling json to Variables: %w", err)
	}
	return &vars, nil
}

type Where map[string]interface{}

func (m *Where) AddEquals(key string, value interface{}) *Where {
	(*m)[key] = map[string]interface{}{
		Eq: value,
	}

	return m
}

func (m *Where) AddExp(key, exp string, value interface{}) *Where {
	(*m)[key] = map[string]interface{}{
		exp: value,
	}

	return m
}

func (m *Where) AddIsNull(key string) *Where {
	(*m)[key] = map[string]interface{}{
		IsNull: true,
	}
	return m
}

func NewWhere() *Where {
	return &Where{}
}

type OrderBy map[string]Order

type QueryArgs struct {
	Where   *Where
	OrderBy *OrderBy
	Limit   *int
	Offset  *int
}

func (m *QueryArgs) Variables() Variables {
	vars := Variables{}
	if m.Where != nil {
		vars["where"] = m.Where
	}
	if m.OrderBy != nil {
		vars["order_by"] = m.OrderBy
	}
	if m.Limit != nil {
		vars["limit"] = m.Limit
	}
	if m.Offset != nil {
		vars["offset"] = m.Offset
	}
	return vars
}
