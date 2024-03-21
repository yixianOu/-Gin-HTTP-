package model

import (
	"fmt"
	"time"

	otgorm "github.com/eddycjy/opentracing-gorm"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/pkg/setting"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	STATE_OPEN  = 1
	STATE_CLOSE = 0
)

// Model 创建公共 model
type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

// NewDBEngine 创建 DB 实例的 NewDBEngine 方法，同时增加了 gorm 开源库的引入和 MySQL 驱动库的初始化
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	s := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	db, err := gorm.Open(databaseSetting.DBType, fmt.Sprintf(s,
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))
	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunMode == "debug" {
		db.LogMode(true)
	}
	db.SingularTable(true)
	//对 Callback 方法进行回调注册，替换gorm中默认的创删改函数
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)
	otgorm.AddGormCallbacks(db)
	return db, nil
}

func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		//在scope中创建字段并将时间信息填入其中
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				_ = createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			//在scope中创建字段并将名字信息填入其中
			if modifyTimeField.IsBlank {
				_ = modifyTimeField.Set(nowTime)
			}
		}
	}
}

func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	//Get("gorm:update_column") 获取当前设置的标识 gorm:update_column 的字段。
	if _, ok := scope.Get("gorm:update_column"); !ok {
		//若不存在，也就是没有自定义设置 update_column，
		//那么将会在更新回调内设置默认字段 ModifiedOn 的值为当前的时间戳。
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		//调用get获取gorm:delete_option的字段属性
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}
		//判断是否存在 DeletedOn 和 IsDel 字段，若存在则执行 UPDATE 操作进行软删除
		//（修改 DeletedOn 和 IsDel 的值），否则执行 DELETE 进行硬删除。
		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")
		isDelField, hasIsDelField := scope.FieldByName("IsDel")
		if !scope.Search.Unscoped && hasDeletedOnField && hasIsDelField {
			now := time.Now().Unix()
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v,%v=%v%v%v",
				scope.QuotedTableName(),                            //拿表名
				scope.Quote(deletedOnField.DBName),                 //拿字段名
				scope.AddToVars(now),                               //赋值
				scope.Quote(isDelField.DBName),                     //拿字段名
				scope.AddToVars(1),                                 //赋值
				addExtraSpaceIfExist(scope.CombinedConditionSql()), //删除的对象
				addExtraSpaceIfExist(extraOption),                  //额外的字段信息
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),                            //表名
				addExtraSpaceIfExist(scope.CombinedConditionSql()), //删除的对象
				addExtraSpaceIfExist(extraOption),                  //额外的字段信息
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
