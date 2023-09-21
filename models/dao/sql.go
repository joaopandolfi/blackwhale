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
	New(v interface{}) error

	// Expose gorm instance
	DB() (*gorm.DB, error)

	Get(id int, v interface{}) error
	GetNested(id int, v interface{}, nesteds []string) error
	GetWithArguments(id int, v interface{}, query string, args ...interface{}) error
	GetGeneric(v interface{}, nested []string, query string, args ...interface{}) error
	GetDeletedWithArguments(id int, v interface{}, query string, args ...interface{}) error

	Last(v interface{}, args ...interface{}) error

	List(v interface{}, limit int) error
	ListConditional(v interface{}, params ListParams, query string, args ...interface{}) error
	ListAll(v interface{}, params ListParams) error
	ListWithArguments(v interface{}, limit int, query string, args ...interface{}) error
	ListWithArgumentsFull(v interface{}, limit int, query string, args ...interface{}) error

	Delete(v interface{}) error
	DeleteWithArguments(v interface{}, args ...interface{}) error

	Upsert(v interface{}) error

	Update(v interface{}) error
	UpdateFull(v interface{}) error
	UpdateWithMap(v interface{}, val map[string]interface{}) error
	UpdateWhere(v interface{}, val map[string]interface{}, query string, args ...interface{}) error

	Raw(v interface{}, query string, args ...interface{}) error

	AppendOnAssociation(v interface{}, association string, data interface{}) error
	DeleteOnAssociation(v interface{}, association string, data interface{}) error
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
func GenerateQuery(params map[string]interface{}) (string, []interface{}) {
	query := ""
	likeQuery := ""
	data := []interface{}{}
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

			argsInterfaces := make([]interface{}, len(args))

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

func (d *dao) Last(v interface{}, args ...interface{}) error {
	tx := d.db.Preload(clause.Associations).Last(v, args)
	if tx.Error != nil {
		return fmt.Errorf("last: %w", tx.Error)
	}
	return nil
}

func (d *dao) Get(id int, v interface{}) error {
	tx := d.db.Preload(clause.Associations).Find(v, id)
	if tx.Error != nil {
		return fmt.Errorf("get: %w", tx.Error)
	}
	return nil
}

func (d *dao) GetNested(id int, v interface{}, nested []string) error {
	tx := d.processNested(d.db, nested).Preload(clause.Associations).Find(v, id)
	if tx.Error != nil {
		return fmt.Errorf("get nested: %w", tx.Error)
	}
	return nil
}

func (d *dao) GetWithArguments(id int, v interface{}, query string, args ...interface{}) error {
	tx := d.db.Preload(clause.Associations).Where(query, args...).Find(v, id)
	if tx.Error != nil {
		return fmt.Errorf("get with args: %w", tx.Error)
	}

	return nil
}

func (d *dao) GetDeletedWithArguments(id int, v interface{}, query string, args ...interface{}) error {
	tx := d.db.Unscoped().Preload(clause.Associations).Where(query, args...).Find(v, id)
	if tx.Error != nil {
		return fmt.Errorf("get deleted with args: %w", tx.Error)
	}

	return nil
}

func (d *dao) GetGeneric(v interface{}, nested []string, query string, args ...interface{}) error {
	tx := d.processNested(d.db, nested).Preload(clause.Associations).Where(query, args...).Find(v)
	if tx.Error != nil {
		return fmt.Errorf("get generic: %w", tx.Error)
	}

	return nil
}

func (d *dao) List(v interface{}, limit int) error {
	tx := d.db.Limit(limit).Find(v)
	if tx.Error != nil {
		return fmt.Errorf("list: %w", tx.Error)
	}

	return nil
}

func (d *dao) ListWithArguments(v interface{}, limit int, query string, args ...interface{}) error {
	tx := d.db.Where(query, args...).Limit(limit).Find(v)
	if tx.Error != nil {
		return fmt.Errorf(listWithArgsMessage, tx.Error)
	}

	return nil
}

func (d *dao) ListWithArgumentsFull(v interface{}, limit int, query string, args ...interface{}) error {
	tx := d.db.Preload(clause.Associations).Where(query, args...).Limit(limit).Find(v)
	if tx.Error != nil {
		return fmt.Errorf(listWithArgsMessage, tx.Error)
	}

	return nil
}

func (d *dao) ListWithArgumentsNestedOrdered(v interface{}, listParams ListParams, nested []string, order string, query string, args ...interface{}) error {
	tx := d.processNested(d.db, nested)
	tx.Where(query, args...).Order(order).Offset(listParams.Offset).Limit(listParams.Limit).Find(v)
	if tx.Error != nil {
		return fmt.Errorf(listWithArgsMessage, tx.Error)
	}

	return nil
}

func (d *dao) ListAll(v interface{}, params ListParams) error {
	return d.ListConditional(v, params, "")
}

func (d *dao) ListConditional(v interface{}, params ListParams, query string, args ...interface{}) error {
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

func (d *dao) ListWithArgumentsFullOrdered(v interface{}, listParams ListParams, order, query string, args ...interface{}) error {
	tx := d.db.Preload(clause.Associations).Where(query, args...).Order(order).Offset(listParams.Offset).Limit(listParams.Offset).Find(v)
	if tx.Error != nil {
		return fmt.Errorf(listWithArgsMessage, tx.Error)
	}

	return nil
}

func (d *dao) New(v interface{}) error {
	tx := d.db.Create(v)
	if tx.Error != nil {
		return fmt.Errorf("saving: %w", tx.Error)
	}

	return nil
}

func (d *dao) Delete(v interface{}) error {
	tx := d.db.Delete(v)
	if tx.Error != nil {
		return fmt.Errorf("deleting: %w", tx.Error)
	}

	return nil
}

func (d *dao) DeleteWithArguments(v interface{}, args ...interface{}) error {
	tx := d.db.Delete(v, args...)
	if tx.Error != nil {
		return fmt.Errorf("deleting (with args): %w", tx.Error)
	}

	return nil
}

func (d *dao) Upsert(v interface{}) error {
	tx := d.db.Clauses(clause.OnConflict{UpdateAll: true}).Session(&gorm.Session{FullSaveAssociations: true}).Save(v)
	if tx.Error != nil {
		return fmt.Errorf("upserting: %w", tx.Error)
	}

	return nil
}

func (d *dao) Update(v interface{}) error {
	tx := d.db.Model(v).Updates(v)
	if tx.Error != nil {
		return fmt.Errorf("updating: %w", tx.Error)
	}

	return nil
}

func (d *dao) UpdateWhere(v interface{}, val map[string]interface{}, query string, args ...interface{}) error {
	tx := d.db.Model(v).Where(query, args...).Updates(val)
	if tx.Error != nil {
		return fmt.Errorf("updating where: %w", tx.Error)
	}

	return nil
}

func (d *dao) UpdateFull(v interface{}) error {
	tx := d.db.Session(&gorm.Session{FullSaveAssociations: true}).Model(v).Updates(v)
	if tx.Error != nil {
		return fmt.Errorf("updating (full): %w", tx.Error)
	}

	return nil
}

func (d *dao) UpdateWithMap(v interface{}, val map[string]interface{}) error {
	tx := d.db.Model(v).Updates(val)
	if tx.Error != nil {
		return fmt.Errorf("updating (map): %w", tx.Error)
	}

	return nil
}

func (d *dao) Raw(v interface{}, query string, args ...interface{}) error {
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

// v interface{}: the model that has the association. Example.: &User{};
// association string: the association name. Example.: "Languages";
// data interface{}: the data to be associated. Example.: &Language{ Name };
func (d *dao) AppendOnAssociation(v interface{}, association string, data interface{}) error {
	err := d.db.Model(v).Association(association).Append(data)
	if err != nil {
		return fmt.Errorf("association: %w", err)
	}

	return nil
}

// v interface{}: the model that has the association. Example: &User{};
// association string: the association name. Example: "Languages";
// data interface{}: the associated data to be deleted. Example: &Language{ ID, Name };
func (d *dao) DeleteOnAssociation(v interface{}, association string, data interface{}) error {
	err := d.db.Model(v).Association(association).Delete(data)
	if err != nil {
		return fmt.Errorf("association: %w", err)
	}

	return nil
}
