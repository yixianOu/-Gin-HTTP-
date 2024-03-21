package model

//针对标签模块的模型操作，如对数据库的增删改查函数，进行封装，并且只与实体产生关系
import (
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/jinzhu/gorm"
)

// Tag 创建标签 model
type Tag struct {
	*Model
	Name  string `json:"name"`
	State uint8  `json:"state"`
}

func (t Tag) TableName() string {
	return "blog_tag"
}

type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}

func (t Tag) Count(db *gorm.DB) (int, error) {
	var count int
	//查询条件：该标签的命名
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	//查询条件：该标签的状态
	db = db.Where("state = ?", t.State)
	//查询条件：是否删除，对以t为模型的表进行查询，将结果存在count中
	if err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (t Tag) List(db *gorm.DB, pageOffset, pageSize int) ([]*Tag, error) {
	var tags []*Tag
	var err error
	//
	if pageOffset >= 0 && pageSize > 0 {
		//跳过开头的pageOffset个文章，显示之后的pageSize个文章
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	//设置查询条件：标签名和标签状态
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)
	//在数据库中根据条件，将查询结果存储到tags中
	if err = db.Where("is_del = ?", 0).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (t Tag) ListByIDs(db *gorm.DB, ids []uint32) ([]*Tag, error) {
	var tags []*Tag
	db = db.Where("state = ? AND is_del = ?", t.State, 0)
	err := db.Where("id IN (?)", ids).Find(&tags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}

func (t Tag) Get(db *gorm.DB) (Tag, error) {
	var tag Tag
	err := db.Where("id = ? AND is_del = ? AND state = ?", t.ID, 0, t.State).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return tag, err
	}

	return tag, nil
}

func (t Tag) Create(db *gorm.DB) error {
	return db.Create(&t).Error
}

func (t Tag) Update(db *gorm.DB, values interface{}) error {
	return db.Model(&t).Where("id = ? AND is_del = ?", t.ID, 0).Updates(values).Error
}

func (t Tag) Delete(db *gorm.DB) error {
	return db.Where("id = ? AND is_del = ?", t.Model.ID, 0).Delete(&t).Error
}
