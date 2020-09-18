package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	logs "github.com/sirupsen/logrus"
	"github.com/sulin2018/go-web-base/src/app/config"
)

var db *gorm.DB

func DBInit() {
	logs.Trace("db init")
	var err error
	db, err = gorm.Open(config.AppConfig.DBType,
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.AppConfig.DBUser,
			config.AppConfig.DBPassword,
			config.AppConfig.DBHost,
			config.AppConfig.DBDatabase),
	)
	if err != nil {
		logs.Panicln("models.Setup err: %v", err)
	}

	// Disable table name's pluralization
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Disable association auto update/create
	db.InstantSet("gorm:association_autoupdate", false)
	//db.InstantSet("gorm:association_autocreate", false)

	if config.AppConfig.AppRunMode == "dev" {
		db.LogMode(true)
	}

	DBAutoMigrate()
	logs.Trace("db init complate")
}

func DBAutoMigrate() {
	db.AutoMigrate(
		&User{},
		&Permission{},
		&Group{},
	)
}

func GetDB() *gorm.DB {
	return db
}

func DBPage(page uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := config.AppConfig.PageSize * (page - 1)
		return db.Offset(offset).Limit(config.AppConfig.PageSize)
	}
}

func DBSearch(searchMap map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var searchQS []string
		var args []interface{}
		for k, v := range searchMap {
			searchQS = append(searchQS, fmt.Sprintf("%s LIKE ?", k))
			args = append(args, "%"+v+"%")
		}
		qs := strings.Join(searchQS, " OR ")
		//logs.Println(qs)
		return db.Where(qs, args...)
	}
}

func DBFilter(filterMap map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var filters []string
		var args []interface{}
		for k, v := range filterMap {
			filters = append(filters, fmt.Sprintf("%s = ?", k))
			args = append(args, v)
		}
		qs := strings.Join(filters, " AND ")
		//logs.Println(qs)
		return db.Where(qs, args...)
	}
}

func DBOrder(orderCols []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		qs := strings.Join(orderCols, ",")
		//logs.Println(qs)
		return db.Order(qs)
	}
}

func Detail(result interface{}) error {
	dbResult := db.Find(result)
	return dbResult.Error
}

func Exist(tempModel interface{}) bool {
	var count uint
	//var id interface{}
	//s := reflect.ValueOf(tempModel).Elem()
	//typeOfT := s.Type()
	//for i := 0; i < s.NumField(); i++ {
	//	f := s.Field(i)
	//	if typeOfT.Field(i).Name == "ID" {
	//		id = f.Interface()
	//		break
	//	}
	//	fmt.Printf("%d: %s %s = %v\n", i,
	//		typeOfT.Field(i).Name, f.Type(), f.Interface())
	//}
	//db.Model(tempModel).Where("id=?", id).Count(&count)
	db.Model(tempModel).Count(&count)
	return count == 1
}

func List(results interface{}) {
	db.Find(results)
}

func CreateOrUpdate(anyModel interface{}) error {
	// model must have Update/Create method
	if Exist(anyModel) {
		if v := reflect.ValueOf(anyModel).MethodByName("Update"); v.String() == "<invalid Value>" {
			return errors.New("model must have Update method")
		} else {
			v.Call(nil)
			return nil
		}
	} else {
		if v := reflect.ValueOf(anyModel).MethodByName("Create"); v.String() == "<invalid Value>" {
			return errors.New("model must have Create method")
		} else {
			v.Call(nil)
			return nil
		}
	}
}

func Page(results interface{}, count interface{}, page uint) {
	tempQuery := db
	tempQuery.Model(results).Count(count)

	if page != 0 {
		tempQuery = tempQuery.Scopes(DBPage(page))
	}
	tempQuery.Find(results)
}

func PageColumn(results interface{}, count interface{}, page uint, col string) {
	tempQuery := db
	tempQuery.Model(results).Count(count)

	if page != 0 {
		tempQuery = tempQuery.Scopes(DBPage(page))
	}
	tempQuery.Select(col).Find(results)
}

func ListPageSearchFilter(results interface{}, count interface{}, page uint, searchMap map[string]string, filterMap map[string]interface{}) {
	tempQuery := db

	if searchMap != nil {
		tempQuery = tempQuery.Scopes(DBSearch(searchMap))
	}
	if filterMap != nil {
		tempQuery = tempQuery.Scopes(DBFilter(filterMap))
	}
	tempQuery.Model(results).Count(count)

	if page != 0 {
		tempQuery = tempQuery.Scopes(DBPage(page))
	}
	tempQuery.Find(results)
}

func ListPageSearchFilterOrder(results interface{}, count interface{}, page uint, searchMap map[string]string, filterMap map[string]interface{}, orderCols []string) {
	tempQuery := db

	// filter search
	if searchMap != nil {
		tempQuery = tempQuery.Scopes(DBSearch(searchMap))
	}
	if filterMap != nil {
		tempQuery = tempQuery.Scopes(DBFilter(filterMap))
	}
	tempQuery.Model(results).Count(count)

	// order page
	if page != 0 {
		tempQuery = tempQuery.Scopes(DBPage(page))
	}
	if orderCols != nil {
		tempQuery = tempQuery.Scopes(DBOrder(orderCols))
	}
	tempQuery.Find(results)
}

func Search(results interface{}, searchMap map[string]string) {
	db.Scopes(DBSearch(searchMap)).Find(results)
}
