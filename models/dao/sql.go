package dao

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// SQLSearchByNameAndCode -> search by name and code: name, code
	SQLSearchByNameAndCode = "lower(name) LIKE lower(?) OR lower(code) LIKE lower(?)"

	listWithArgsMessage = "list with args: %w"

	// PreloadAll data
	PrealoadAll = clause.Associations

	// MaxLimit of query
	MaxLimit = 1000

	// LikeCondition
	LikeCondition = "ilike"
)

// DAO generic public interface
type SQLDAO interface {
	New(v any) error

	// Expose gorm instance
	DB() (*gorm.DB, error)

	Get(id int, v any) error
	GetNested(id int, v any, nesteds []string) error
	GetWithArguments(id int, v any, query string, args ...any) error
	GetGeneric(v any, nested []string, query string, args ...any) error
	GetDeletedWithArguments(id int, v any, query string, args ...any) error

	Last(v any, args ...any) error

	List(v any, limit int) error
	ListConditional(v any, params ListParams, query string, args ...any) error
	ListAll(v any, params ListParams) error
	ListWithArguments(v any, limit int, query string, args ...any) error
	ListWithArgumentsFull(v any, limit int, query string, args ...any) error

	Delete(v any) error
	DeleteWithArguments(v any, args ...any) error

	Upsert(v any) error

	Update(v any) error
	UpdateFull(v any) error
	UpdateWithMap(v any, val map[string]any) error
	UpdateWhere(v any, val map[string]any, query string, args ...any) error

	Raw(v any, query string, args ...any) error

	AppendOnAssociation(v any, association string, data any) error
	DeleteOnAssociation(v any, association string, data any) error
}

type ListParams struct {
	Limit  int
	Offset int
	Nested []string
	Order  string
	Full   bool
}

type dao struct {
	db *gorm.DB
}

// Generate query based on params
// Use wisely
func GenerateQuery(params map[string]any) (string, []any) {
	query := ""
	likeQuery := ""
	data := []any{}
	and := ""
	or := ""
	for k, v := range params {
		if v == "" {
			continue
		}

		if strings.Contains(k, LikeCondition) {
			ks := strings.Split(k, ":")
			column := ks[0]
			likeQuery = fmt.Sprintf("%s %s %s ILIKE ? ", likeQuery, or, column)
			data = append(data, v)
			or = " OR "
			continue
		}

		ks := strings.Split(k, ":")
		key := ks[0]
		condition := ks[1]
		query = fmt.Sprintf("%s %s %s %s ?", query, and, key, condition)
		data = append(data, v)
		and = " AND "
	}

	if likeQuery != "" {
		if query != "" {
			query = fmt.Sprintf("%s AND (%s) ", query, likeQuery)
		} else {
			query = likeQuery
		}
	}

	return query, data
}

// NewDao generic
func Sql(db *gorm.DB) SQLDAO {
	if db == nil {
		panic("Database can't be null")
	}
	return &dao{
		db: db,
	}
}

func withDeleteds(db *gorm.DB) *gorm.DB {
	return db.Unscoped()
}

func (d *dao) processNested(tx *gorm.DB, nested []string) *gorm.DB {
	for _, n := range nested {
		if strings.Contains(n, ":") {
			splited := strings.Split(n, ":")

			args := splited[1:]

			argsInterfaces := make([]any, len(args))

			for index, value := range args {
				if value == "withDeleteds" {
					argsInterfaces[index] = withDeleteds
					continue
				}
				argsInterfaces[index] = value
			}

			tx = tx.Preload(splited[0], argsInterfaces...)
		} else {
			tx = tx.Preload(n)
		}
	}
	return tx
}

func (d *dao) Last(v any, args ...any) error {
	tx := d.db.Preload(clause.Associations).Last(v, args)
	if tx.Error != nil {
		return fmt.Errorf("last: %w", tx.Error)
	}
	return nil
}

func (d *dao) Get(id int, v any) error {
	tx := d.db.Preload(clause.Associations).Find(v, id)
	if tx.Error != nil {
		return fmt.Errorf("get: %w", tx.Error)
	}
	return nil
}

func (d *dao) GetNested(id int, v any, nested []string) error {
	tx := d.processNested(d.db, nested).Preload(clause.Associations).Find(v, id)
	if tx.Error != nil {
		return fmt.Errorf("get nested: %w", tx.Error)
	}
	return nil
}

func (d *dao) GetWithArguments(id int, v any, query string, args ...any) error {
	tx := d.db.Preload(clause.Associations).Where(query, args...).Find(v, id)
	if tx.Error != nil {
		return fmt.Errorf("get with args: %w", tx.Error)
	}

	return nil
}

func (d *dao) GetDeletedWithArguments(id int, v any, query string, args ...any) error {
	tx := d.db.Unscoped().Preload(clause.Associations).Where(query, args...).Find(v, id)
	if tx.Error != nil {
		return fmt.Errorf("get deleted with args: %w", tx.Error)
	}

	return nil
}

func (d *dao) GetGeneric(v any, nested []string, query string, args ...any) error {
	tx := d.processNested(d.db, nested).Preload(clause.Associations).Where(query, args...).Find(v)
	if tx.Error != nil {
		return fmt.Errorf("get generic: %w", tx.Error)
	}

	return nil
}

func (d *dao) List(v any, limit int) error {
	tx := d.db.Limit(limit).Find(v)
	if tx.Error != nil {
		return fmt.Errorf("list: %w", tx.Error)
	}

	return nil
}

func (d *dao) ListWithArguments(v any, limit int, query string, args ...any) error {
	tx := d.db.Where(query, args...).Limit(limit).Find(v)
	if tx.Error != nil {
		return fmt.Errorf(listWithArgsMessage, tx.Error)
	}

	return nil
}

func (d *dao) ListWithArgumentsFull(v any, limit int, query string, args ...any) error {
	tx := d.db.Preload(clause.Associations).Where(query, args...).Limit(limit).Find(v)
	if tx.Error != nil {
		return fmt.Errorf(listWithArgsMessage, tx.Error)
	}

	return nil
}

func (d *dao) ListWithArgumentsNestedOrdered(v any, listParams ListParams, nested []string, order string, query string, args ...any) error {
	tx := d.processNested(d.db, nested)
	tx.Where(query, args...).Order(order).Offset(listParams.Offset).Limit(listParams.Limit).Find(v)
	if tx.Error != nil {
		return fmt.Errorf(listWithArgsMessage, tx.Error)
	}

	return nil
}

func (d *dao) ListAll(v any, params ListParams) error {
	return d.ListConditional(v, params, "")
}

func (d *dao) ListConditional(v any, params ListParams, query string, args ...any) error {
	tx := d.db

	if params.Full {
		tx = tx.Preload(clause.Associations)
	}

	tx = d.processNested(tx, params.Nested)

	tx = tx.Where(query, args...)

	if params.Order != "" {
		tx = tx.Order(params.Order)
	}

	if params.Offset != 0 {
		tx = tx.Offset(params.Offset)
	}

	if params.Limit != 0 {
		tx = tx.Limit(params.Limit)
	}

	tx = tx.Find(v)
	if tx.Error != nil {
		return fmt.Errorf("listing conditional: %w", tx.Error)
	}

	return nil
}

func (d *dao) ListWithArgumentsFullOrdered(v any, listParams ListParams, order, query string, args ...any) error {
	tx := d.db.Preload(clause.Associations).Where(query, args...).Order(order).Offset(listParams.Offset).Limit(listParams.Offset).Find(v)
	if tx.Error != nil {
		return fmt.Errorf(listWithArgsMessage, tx.Error)
	}

	return nil
}

func (d *dao) New(v any) error {
	tx := d.db.Create(v)
	if tx.Error != nil {
		return fmt.Errorf("saving: %w", tx.Error)
	}

	return nil
}

func (d *dao) Delete(v any) error {
	tx := d.db.Delete(v)
	if tx.Error != nil {
		return fmt.Errorf("deleting: %w", tx.Error)
	}

	return nil
}

func (d *dao) DeleteWithArguments(v any, args ...any) error {
	tx := d.db.Delete(v, args...)
	if tx.Error != nil {
		return fmt.Errorf("deleting (with args): %w", tx.Error)
	}

	return nil
}

func (d *dao) Upsert(v any) error {
	tx := d.db.Clauses(clause.OnConflict{UpdateAll: true}).Session(&gorm.Session{FullSaveAssociations: true}).Save(v)
	if tx.Error != nil {
		return fmt.Errorf("upserting: %w", tx.Error)
	}

	return nil
}

func (d *dao) Update(v any) error {
	tx := d.db.Model(v).Updates(v)
	if tx.Error != nil {
		return fmt.Errorf("updating: %w", tx.Error)
	}

	return nil
}

func (d *dao) UpdateWhere(v any, val map[string]any, query string, args ...any) error {
	tx := d.db.Model(v).Where(query, args...).Updates(val)
	if tx.Error != nil {
		return fmt.Errorf("updating where: %w", tx.Error)
	}

	return nil
}

func (d *dao) UpdateFull(v any) error {
	tx := d.db.Session(&gorm.Session{FullSaveAssociations: true}).Model(v).Updates(v)
	if tx.Error != nil {
		return fmt.Errorf("updating (full): %w", tx.Error)
	}

	return nil
}

func (d *dao) UpdateWithMap(v any, val map[string]any) error {
	tx := d.db.Model(v).Updates(val)
	if tx.Error != nil {
		return fmt.Errorf("updating (map): %w", tx.Error)
	}

	return nil
}

func (d *dao) Raw(v any, query string, args ...any) error {
	tx := d.db.Raw(query, args...).Scan(v)
	if tx.Error != nil {
		return fmt.Errorf("Raw data: %w", tx.Error)
	}

	return nil
}

func (d *dao) DB() (*gorm.DB, error) {
	if d.db == nil {
		return nil, fmt.Errorf("database is not itialized")
	}

	return d.db, nil
}

// v any: the model that has the association. Example.: &User{};
// association string: the association name. Example.: "Languages";
// data any: the data to be associated. Example.: &Language{ Name };
func (d *dao) AppendOnAssociation(v any, association string, data any) error {
	err := d.db.Model(v).Association(association).Append(data)
	if err != nil {
		return fmt.Errorf("association: %w", err)
	}

	return nil
}

// v any: the model that has the association. Example: &User{};
// association string: the association name. Example: "Languages";
// data any: the associated data to be deleted. Example: &Language{ ID, Name };
func (d *dao) DeleteOnAssociation(v any, association string, data any) error {
	err := d.db.Model(v).Association(association).Delete(data)
	if err != nil {
		return fmt.Errorf("association: %w", err)
	}

	return nil
}
