package dao

import (
	"fmt"
)

type QueryParams struct {
	data       map[string]any
	GenericVal string
}

func (d *QueryParams) init() {
	if d.data == nil {
		d.data = map[string]any{}
	}
}

func (d *QueryParams) Parse() map[string]any {
	return d.data
}

func (d *QueryParams) AddParam(column string, condition string, val any) {
	d.init()
	d.data[fmt.Sprintf("%s:%s", column, condition)] = val
}

func (d *QueryParams) AddLike(column, val string) {
	d.init()

	d.data[fmt.Sprintf("%s:%s", LikeCondition, column)] = val
}

func (d *QueryParams) Data() *map[string]any {
	return &d.data
}
