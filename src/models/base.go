package models

import (
	"fmt"
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
		logs.Panicln("models.Setup err: ", err)
	}

	// Disable table name's pluralization
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Disable association auto update/create
	// db.InstantSet("gorm:association_autoupdate", false)
	// db.InstantSet("gorm:association_autocreate", false)

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

func DBPage(page uint, pageSize uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := pageSize * (page - 1)
		return db.Offset(offset).Limit(pageSize)
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

func Detail(tempModel interface{}) error {
	result := db.Find(tempModel)
	if result.Error != nil {
		logs.Error(result.Error)
	}
	return result.Error
}

func Exist(tempModel interface{}) bool {
	result := db.First(tempModel, tempModel)
	return result.RowsAffected == 1
}

// results 接收数据指针
// count 总数指针
// page 页码, 0表示获取所有/不分页
// pageSize 页大小, 0表示使用配置大小
func Page(results interface{}, count interface{}, page uint, pageSize uint) error {
	tempQuery := db
	tempQuery.Model(results).Count(count)

	if page != 0 {
		if pageSize == 0 {
			pageSize = config.AppConfig.PageSize
		}
		tempQuery = tempQuery.Scopes(DBPage(page, pageSize))
	}
	result := tempQuery.Find(results)
	if result.Error != nil {
		logs.Error(result.Error)
		return result.Error
	}
	return nil
}

// 分页获取 限定字段
func PageColumns(results interface{}, count interface{}, page uint, pageSize uint, col string) error {
	tempQuery := db
	tempQuery.Model(results).Count(count)

	if page != 0 {
		if pageSize == 0 {
			pageSize = config.AppConfig.PageSize
		}
		tempQuery = tempQuery.Scopes(DBPage(page, pageSize))
	}
	result := tempQuery.Select(col).Find(results)
	if result.Error != nil {
		logs.Error(result.Error)
		return result.Error
	}
	return nil
}

func ListPageSearchFilterOrder(results interface{}, count interface{}, page uint, pageSize uint, searchMap map[string]string, filterMap map[string]interface{}, orderCols []string) error {
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
		if pageSize == 0 {
			pageSize = config.AppConfig.PageSize
		}
		tempQuery = tempQuery.Scopes(DBPage(page, pageSize))
	}
	if orderCols != nil {
		tempQuery = tempQuery.Scopes(DBOrder(orderCols))
	}
	result := tempQuery.Find(results)
	if result.Error != nil {
		logs.Error(result.Error)
		return result.Error
	}
	return nil
}
