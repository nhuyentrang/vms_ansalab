package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	dtoMapper "github.com/dranikpg/dto-mapper"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Todo: use gorm smart select, then no need for mapping

var Connected bool = false

var defaultDB *gorm.DB

type SQLQuery[M any, E any] struct {
	expressStr string
	args       []interface{}
	db         *gorm.DB
}

func Connect(host, port, dbName, sslMode, user, password, schema string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s user=%s password=%s search_path=%s",
		host, port, dbName, sslMode, user, password, schema)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(models ...interface{}) error {
	if !Connected {
		return errors.New("database not connected")
	}

	err := defaultDB.AutoMigrate(models...)
	return err
}

func Ping() error {
	if !Connected {
		return errors.New("not connected")
	}
	sqlDB, err := defaultDB.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Ping()
	if err != nil {
		Connected = false
		return err
	}
	return nil
}

func Close() error {
	if !Connected {
		return errors.New("not connected")
	}
	sqlDB, err := defaultDB.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		return err
	}
	Connected = false
	return nil
}

func Stats() (stats sql.DBStats, err error) {
	if !Connected {
		return stats, errors.New("not connected")
	}
	// Returns database statistics
	sqlDB, err := defaultDB.DB()
	if err != nil {
		return stats, err
	}
	return sqlDB.Stats(), nil
}

// NewQuery create new query instance
func NewQuery[M any, E any]( /*db *gorm.DB*/ dbInstances ...interface{}) *SQLQuery[M, E] {
	query := &SQLQuery[M, E]{}

	var isDBInitiallized bool = false
	for _, db := range dbInstances {
		if db != nil {
			switch t := db.(type) {
			case *gorm.DB:
				if t != nil {
					query.db = t
					isDBInitiallized = true
				}
			default:
			}
		}
	}

	// Assign default DB
	if !isDBInitiallized && defaultDB != nil {
		query.db = defaultDB
		isDBInitiallized = true
	}
	if !isDBInitiallized {
		panic(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> database is not initialized")
	}

	query.expressStr = ""
	query.args = make([]interface{}, 0)
	return query
}

// AddConditionOfTextField add one filter condition of normal text field into query
func (query *SQLQuery[M, E]) AddConditionOfTextField(cascadingLogic string, fieldName string, comparisonOperator string, value interface{}) {
	if fieldName == "" {
		return
	}

	if query.expressStr == "" {
		if comparisonOperator == "LIKE" {
			query.expressStr = fmt.Sprintf("lower(\"%s\") %s ?", fieldName, comparisonOperator)
		} else {
			query.expressStr = fmt.Sprintf("\"%s\" %s ?", fieldName, comparisonOperator)
		}
	} else {
		if comparisonOperator == "LIKE" {

			query.expressStr = fmt.Sprintf("%s %s lower(\"%s\") %s ?", query.expressStr, cascadingLogic, fieldName, comparisonOperator)
		} else {
			query.expressStr = fmt.Sprintf("%s %s \"%s\" %s ?", query.expressStr, cascadingLogic, fieldName, comparisonOperator)

		}
	}

	if comparisonOperator == "LIKE" {
		s, ok := value.(string)
		if ok {
			s = strings.ToLower(s)
			query.args = append(query.args, "%"+s+"%")
		}
	} else {
		query.args = append(query.args, value)
	}
}

// AddTwoConditionOfTextField add two filter condition of two normal text field into query
func (query *SQLQuery[M, E]) AddTwoConditionOfTextField(cascadingLogic string, fieldName1 string, comparisonOperator1 string, value1 interface{}, combineLogic string, fieldName2 string, comparisonOperator2 string, value2 interface{}) {
	if fieldName1 == "" || fieldName2 == "" {
		return
	}

	if query.expressStr == "" {
		if comparisonOperator1 == "LIKE" && comparisonOperator2 != "LIKE" {
			query.expressStr = fmt.Sprintf("lower(\"%s\") %s ? %s \"%s\" %s ?", fieldName1, comparisonOperator1, combineLogic, fieldName2, comparisonOperator2)
		} else if comparisonOperator1 != "LIKE" && comparisonOperator2 == "LIKE" {
			query.expressStr = fmt.Sprintf("\"%s\" %s ? %s lower(\"%s\") %s ?", fieldName1, comparisonOperator1, combineLogic, fieldName2, comparisonOperator2)
		} else if comparisonOperator1 == "LIKE" && comparisonOperator2 == "LIKE" {
			query.expressStr = fmt.Sprintf("lower(\"%s\") %s ? %s lower(\"%s\") %s ?", fieldName1, comparisonOperator1, combineLogic, fieldName2, comparisonOperator2)
		} else {
			query.expressStr = fmt.Sprintf("\"%s\" %s ? %s \"%s\" %s ?", fieldName1, comparisonOperator1, combineLogic, fieldName2, comparisonOperator2)
		}
	} else {
		if comparisonOperator1 == "LIKE" && comparisonOperator2 != "LIKE" {
			query.expressStr = fmt.Sprintf("%s %s (lower(\"%s\") %s ? %s \"%s\" %s ?)", query.expressStr, cascadingLogic, fieldName1, comparisonOperator1, combineLogic, fieldName2, comparisonOperator2)
		} else if comparisonOperator1 != "LIKE" && comparisonOperator2 == "LIKE" {
			query.expressStr = fmt.Sprintf("%s %s (\"%s\" %s ? %s lower(\"%s\") %s ?)", query.expressStr, cascadingLogic, fieldName1, comparisonOperator1, combineLogic, fieldName2, comparisonOperator2)
		} else if comparisonOperator1 == "LIKE" && comparisonOperator2 == "LIKE" {
			query.expressStr = fmt.Sprintf("%s %s( lower(\"%s\") %s ? %s lower(\"%s\") %s ?)", query.expressStr, cascadingLogic, fieldName1, comparisonOperator1, combineLogic, fieldName2, comparisonOperator2)
		} else {
			query.expressStr = fmt.Sprintf("%s %s (\"%s\" %s ? %s \"%s\" %s ?)", query.expressStr, cascadingLogic, fieldName1, comparisonOperator1, combineLogic, fieldName2, comparisonOperator2)

		}
	}

	if comparisonOperator1 == "LIKE" {
		s, ok := value1.(string)
		if ok {
			s = strings.ToLower(s)
			query.args = append(query.args, "%"+s+"%")
		}
	} else {
		query.args = append(query.args, value1)
	}

	if comparisonOperator2 == "LIKE" {
		s, ok := value2.(string)
		if ok {
			s = strings.ToLower(s)
			query.args = append(query.args, "%"+s+"%")
		}
	} else {
		query.args = append(query.args, value2)
	}
}

// AddConditionOfJsonbField add one filter condition of jsonb field into query
func (query *SQLQuery[M, E]) AddConditionOfJsonbField(cascadingLogic string, fieldName string, key string, comparisonOperator string, value interface{}) {
	if fieldName == "" {
		return
	}

	if query.expressStr == "" {
		if comparisonOperator == "LIKE" {
			query.expressStr = fmt.Sprintf("lower(\"%s\") ->> '%s' %s ?", fieldName, key, comparisonOperator)
		} else {
			query.expressStr = fmt.Sprintf("\"%s\" ->> '%s' %s ?", fieldName, key, comparisonOperator)
		}
	} else {
		if comparisonOperator == "LIKE" {
			query.expressStr = fmt.Sprintf("%s %s lower(\"%s\") ->> '%s' %s ?", query.expressStr, cascadingLogic, fieldName, key, comparisonOperator)
		} else {
			query.expressStr = fmt.Sprintf("%s %s \"%s\" ->> '%s' %s ?", query.expressStr, cascadingLogic, fieldName, key, comparisonOperator)
		}
	}

	if comparisonOperator == "LIKE" {
		s, ok := value.(string)
		if ok {
			s = strings.ToLower(s)
			query.args = append(query.args, "%"+s+"%")
		}
	} else {
		query.args = append(query.args, value)
	}
}

// Exec run the the query to get all items with current filter, no paging
func (query *SQLQuery[M, E]) ExecNoPaging(sort string) (dtos []M, count int64, err error) {
	if !Connected {
		return dtos, 0, errors.New("database not connected")
	}
	count = 0

	if strings.HasPrefix(sort, "-") {
		sort = "\"" + strings.TrimPrefix(sort, "-") + "\"" + " desc"
	} else if strings.HasPrefix(sort, "+") {
		sort = "\"" + strings.TrimPrefix(sort, "+") + "\"" + " asc"
	} else {
		sort = "\"created_at\"" + " desc"
	}

	// Query with filter
	var items []E
	result := defaultDB.Order(sort).Where(query.expressStr, query.args...).Find(&items)
	if result.Error != nil {
		return dtos, count, result.Error
	}
	//count = result.RowsAffected

	// Map entity item to DTO model
	dtos = make([]M, 0)
	for _, item := range items {
		// Mapping from entity model to DTO model
		var dto M
		if err := dtoMapper.Map(&dto, item); err != nil {
			return dtos, count, err
		}
		dtos = append(dtos, dto)
		count++
	}

	// Return
	return dtos, count, result.Error
}
func (query *SQLQuery[M, E]) ExecWithoutPaging(sort string) (dtos []M, count int64, err error) {
	if !Connected {
		return dtos, 0, errors.New("database not connected")
	}

	// Xử lý sort
	if strings.HasPrefix(sort, "-") {
		sort = "\"" + strings.TrimPrefix(sort, "-") + "\"" + " desc"
	} else if strings.HasPrefix(sort, "+") {
		sort = "\"" + strings.TrimPrefix(sort, "+") + "\"" + " asc"
	} else {
		sort = "\"created_at\"" + " desc"
	}

	// Đếm tổng số kết quả
	var entityModel E
	result := query.db.Model(entityModel).Where(query.expressStr, query.args...).Count(&count)
	if result.Error != nil {
		return dtos, 0, result.Error
	}

	// Truy vấn tất cả dữ liệu không phân trang
	var items []E
	result = query.db.Order(sort).Where(query.expressStr, query.args...).Find(&items)
	if result.Error != nil {
		return dtos, 0, result.Error
	}

	// Ánh xạ từ entity sang DTO
	dtos = make([]M, 0)
	for _, item := range items {
		var dto M
		if err := dtoMapper.Map(&dto, item); err != nil {
			return dtos, count, err
		}
		dtos = append(dtos, dto)
	}

	return dtos, count, nil
}

// ExecPaging run the the query to get items with current filter, with paging
func (query *SQLQuery[M, E]) ExecWithPaging(sort string, limit int, page int) (dtos []M, count int64, err error) {
	if !Connected {
		return dtos, 0, errors.New("database not connected")
	}

	// Validate query param
	if limit < 1 {
		limit = 100
	}
	if page < 1 {
		page = 1
	}
	if strings.HasPrefix(sort, "-") {
		sort = "\"" + strings.TrimPrefix(sort, "-") + "\"" + " desc"
	} else if strings.HasPrefix(sort, "+") {
		sort = "\"" + strings.TrimPrefix(sort, "+") + "\"" + " asc"
	} else {
		sort = "\"created_at\"" + " desc"
	}

	// Calculate offset
	offset := limit * (page - 1)
	var result *gorm.DB

	// Count total number
	var entityModel E
	result = query.db.Model(entityModel).Where(query.expressStr, query.args...).Count(&count)
	if result.Error != nil {
		return dtos, 0, result.Error
	}

	// Query with filter
	var items []E
	result = query.db.Limit(limit).Offset(offset).Order(sort).Where(query.expressStr, query.args...).Find(&items)

	// Map entity item to DTO model
	dtos = make([]M, 0)
	for _, item := range items {
		// Mapping from entity model to DTO model
		var dto M
		if err := dtoMapper.Map(&dto, item); err != nil {
			return dtos, count, err
		}
		dtos = append(dtos, dto)
	}

	return dtos, count, result.Error
}

// CreateItemFromDTO map dto (data transfer object) to new database's item struct
// and write that item into database , accepts generic types
//
// It return created item and error
func CreateItemFromDTO[M any, E any](dto M) (M, error) {
	if !Connected {
		return dto, errors.New("database not connected")
	}

	// Validate dto object  input
	validate := validator.New()
	err := validate.Struct(dto)
	if err != nil {
		return dto, err
	}

	// Mapping from DTO to entity model
	var item E
	if err := dtoMapper.Map(&item, dto); err != nil {
		return dto, err
	}

	// Create new entity using smart select

	var entity E
	if result := defaultDB.Model(entity).Create(&item); result.Error != nil {
		return dto, result.Error
	}

	// Mapping from entity model to DTO model
	if err := dtoMapper.Map(&dto, item); err != nil {
		return dto, err
	}
	return dto, nil
}

// ReadItemIntoDTO read an item by ID from database then map resutl into dto (data transfer object),
// accepts generic types
//
// It return read dto and error
func ReadItemByIDIntoDTO[M any, E any](id string) (dto M, err error) {
	if !Connected {
		return dto, errors.New("database not connected")
	}
	var item E
	if err := defaultDB.Where("id = ?", id).First(&item).Error; err != nil {
		return dto, err
	}

	// Mapping from entity model to DTO model
	if err := dtoMapper.Map(&dto, item); err != nil {
		return dto, err
	}
	return dto, nil
}

func ReadItemNonDeletedByIDIntoDTO[M any, E any](id string) (dto M, err error) {
	if !Connected {
		return dto, errors.New("database not connected")
	}

	var item E

	// Thêm điều kiện "deleted = false" trong câu truy vấn
	if err := defaultDB.Where("id = ? AND deleted = ?", id, false).First(&item).Error; err != nil {
		return dto, err
	}

	// Mapping từ entity model sang DTO model
	if err := dtoMapper.Map(&dto, item); err != nil {
		return dto, err
	}

	return dto, nil
}

func ReadItemByNameIntoDTO[M any, E any](id string) (dto M, err error) {
	if !Connected {
		return dto, errors.New("database not connected")
	}
	var item E
	if err := defaultDB.Where("name = ?", id).First(&item).Error; err != nil {
		return dto, err
	}

	// Mapping from entity model to DTO model
	if err := dtoMapper.Map(&dto, item); err != nil {
		return dto, err
	}
	return dto, nil
}
func ReadItemByCodeIntoDTO[M any, E any](code string) (dto M, err error) {
	if !Connected {
		return dto, errors.New("database not connected")
	}
	var item E
	if err := defaultDB.Where("code = ? AND deleted = ?", code, false).First(&item).Error; err != nil {
		return dto, err
	}

	// Mapping from entity model to DTO model
	if err := dtoMapper.Map(&dto, item); err != nil {
		return dto, err
	}
	return dto, nil
}

// ReadItemIntoDTO read an item by ID from database then map resutl into dto (data transfer object),
// accepts generic types
//
// It return read dtos and error
func ReadMultiItemsByIDIntoDTO[M any, E any](ids []string, sort string) (dtos []M, count int64, err error) {
	if !Connected {
		return dtos, 0, errors.New("database not connected")
	}
	count = 0

	if strings.HasPrefix(sort, "-") {
		sort = "\"" + strings.TrimPrefix(sort, "-") + "\"" + " desc"
	} else if strings.HasPrefix(sort, "+") {
		sort = "\"" + strings.TrimPrefix(sort, "+") + "\"" + " asc"
	} else {
		sort = "\"created_at\"" + " desc"
	}

	var items []E
	result := defaultDB.Order(sort).Where("id IN ?", ids).Find(&items)
	if result.Error != nil {
		return dtos, 0, result.Error
	}
	//count = result.RowsAffected

	dtos = make([]M, 0)
	for _, item := range items {
		// Mapping from entity model to DTO model
		var dto M
		if err := dtoMapper.Map(&dto, item); err != nil {
			return dtos, count, err
		}
		dtos = append(dtos, dto)
		count++
	}

	return dtos, count, nil
}

// ReadItemIntoDTO read an item by ID from database then map resutl into dto (data transfer object),
// accepts generic types
//
// It return read dtos and error
func ReadAllItemsIntoDTO[M any, E any](sort string) (dtos []M, count int64, err error) {
	if !Connected {
		return dtos, 0, errors.New("database not connected")
	}
	count = 0

	if strings.HasPrefix(sort, "-") {
		sort = "\"" + strings.TrimPrefix(sort, "-") + "\"" + " desc"
	} else if strings.HasPrefix(sort, "+") {
		sort = "\"" + strings.TrimPrefix(sort, "+") + "\"" + " asc"
	} else {
		sort = "\"created_at\"" + " desc"
	}

	var items []E

	result := defaultDB.Order(sort).Find(&items)
	if result.Error != nil {
		return dtos, 0, result.Error
	}
	//count = result.RowsAffected

	// Mapping from entity model to DTO model
	dtos = make([]M, 0)
	for _, item := range items {
		var dto M
		if err := dtoMapper.Map(&dto, item); err != nil {
			return dtos, count, err
		}
		dtos = append(dtos, dto)
		count++
	}

	return dtos, count, nil
}
func ReadAllNonDeletedItemsIntoDTO[M any, E any](sort string) (dtos []M, count int64, err error) {
	if !Connected {
		return dtos, 0, errors.New("database not connected")
	}
	count = 0

	if strings.HasPrefix(sort, "-") {
		sort = "\"" + strings.TrimPrefix(sort, "-") + "\"" + " desc"
	} else if strings.HasPrefix(sort, "+") {
		sort = "\"" + strings.TrimPrefix(sort, "+") + "\"" + " asc"
	} else {
		sort = "\"created_at\"" + " desc"
	}

	var items []E

	// Thêm điều kiện lọc deleted = false
	result := defaultDB.Where("deleted = ?", false).Order(sort).Find(&items)
	if result.Error != nil {
		return dtos, 0, result.Error
	}

	dtos = make([]M, 0)
	for _, item := range items {
		var dto M
		if err := dtoMapper.Map(&dto, item); err != nil {
			return dtos, count, err
		}
		dtos = append(dtos, dto)
		count++
	}

	return dtos, count, nil
}

// ReadItemWithFilterIntoDTO read an item with Filter from database then map resutl into dto (data transfer object),
// accepts generic types
//
// It return read dto and error
func ReadItemWithFilterIntoDTO[M any, E any](query string, args ...interface{}) (dto M, err error) {
	if !Connected {
		return dto, errors.New("database not connected")
	}

	var item E
	result := defaultDB.Where(query, args...).First(&item)
	if result.Error != nil {
		return dto, result.Error
	}

	// Mapping from entity model to DTO model
	if err := dtoMapper.Map(&dto, item); err != nil {
		return dto, err
	}

	return dto, nil
}

// UpdateItemByIDIntoDTO check if item ID exist in database, then map dto to item struct for updating it (actually patching),
// accepts generic types. Empty (null) field will not be updated
//
// It return updated item (dto) and error
func UpdateItemByIDFromDTO[M any, E any](id string, dto M) (M, error) {
	if !Connected {
		return dto, errors.New("database not connected")
	}

	// Check item exist by ID
	var item E
	if err := defaultDB.Where("id = ?", id).First(&item).Error; err != nil {
		return dto, err
	}

	// Mapping from DTO to entity model
	if err := dtoMapper.Map(&item, dto); err != nil {
		return dto, err
	}

	// Update item
	if err := defaultDB.Model(item).Where("id = ?", id).Updates(&item).Error; err != nil {
		return dto, err
	}

	// Mapping back from updated entity to DTO
	if err := dtoMapper.Map(&dto, item); err != nil {
		return dto, err
	}

	// Todo: uuid of dto is not updated here, please make dto's id updated here
	return dto, nil
}

func UpdateItemByCodeFromDTO[M any, E any](id string, dto M) (M, error) {
	if !Connected {
		return dto, errors.New("database not connected")
	}

	// Check item exist by ID
	var item E
	if err := defaultDB.Where("code = ?", id).First(&item).Error; err != nil {
		return dto, err
	}

	// Mapping from DTO to entity model
	if err := dtoMapper.Map(&item, dto); err != nil {
		return dto, err
	}

	// Update item
	if err := defaultDB.Model(item).Where("code = ?", id).Updates(&item).Error; err != nil {
		return dto, err
	}

	// Mapping back from updated entity to DTO
	if err := dtoMapper.Map(&dto, item); err != nil {
		return dto, err
	}

	// Todo: uuid of dto is not updated here, please make dto's id updated here
	return dto, nil
}

// DeleteItemByID delete item by ID,
// accepts generic types.
//
// It return error if there is any
func DeleteItemByID[E any](id string) (err error) {
	if !Connected {
		return errors.New("database not connected")
	}

	var item E
	if err = defaultDB.Where("id = ?", id).Delete(&item).Error; err != nil {
		return err
	}

	return nil
}

// DeleteAllItem delete all item,
// accepts generic types.
//
// It return error if there is any
func DeleteAllItem[E any](softDelete bool) (err error) {
	if !Connected {
		return errors.New("database not connected")
	}

	var item E
	if softDelete {
		// Softdelete: the record WON’T be removed from the database,
		// but GORM will set the DeletedAt‘s value to the current time,
		// and the data is not findable with normal Query methods anymore.
		// You can find soft deleted records with Unscoped:
		// db.Unscoped().Where("age = 20").Find(&user)
		if err = defaultDB.Where("created_at > ?", "2000-01-01 00:00:00").Delete(&item).Error; err != nil {
			return err
		}
	} else {
		if err = defaultDB.Unscoped().Where("created_at > ?", "2000-01-01 00:00:00").Delete(&item).Error; err != nil {
			return err
		}
	}

	return nil
}

// CheckItemExistedByID check item is existed by ID,
// accepts generic types.
//
// It return true if item is existed
func CheckItemExistedByID[E any](id string) (exists bool, err error) {
	if !Connected {
		return exists, errors.New("database not connected")
	}

	var item E
	if err = defaultDB.Model(item).Select("count(*) > 0").Where("id = ?", id).Find(&exists).Error; err != nil {
		return exists, err
	}

	return exists, nil
}

func UpdateSingleColumn[E any](id string, columnName string, value interface{}) error {
	if !Connected {
		return errors.New("database not connected")
	}

	// Check item exist by ID
	var item E
	if err := defaultDB.Where("id = ?", id).First(&item).Error; err != nil {
		return err
	}

	// Update item
	if err := defaultDB.Model(item).Where("id = ?", id).Update(columnName, value).Error; err != nil {
		return err
	}

	return nil
}
