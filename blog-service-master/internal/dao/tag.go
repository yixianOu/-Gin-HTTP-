package dao

//用于处理标签模块的 dao 操作,是对数据库的第二层封装
//在 dao 层进行了数据访问对象的封装，并针对业务所需的字段进行了处理。
import (
	"github.com/go-programming-tour-book/blog-service/internal/model"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
)

// GetTag 根据id和state拿tag
func (d *Dao) GetTag(id uint32, state uint8) (model.Tag, error) {
	tag := model.Tag{Model: &model.Model{ID: id}, State: state}
	return tag.Get(d.engine)
}

// GetTagList 从page之后根据name和state查询pageSize条tag
func (d *Dao) GetTagList(name string, state uint8, page, pageSize int) ([]*model.Tag, error) {
	tag := model.Tag{Name: name, State: state}
	pageOffset := app.GetPageOffset(page, pageSize)
	return tag.List(d.engine, pageOffset, pageSize)
}

// GetTagListByIDs 根据ids和state拿tags
func (d *Dao) GetTagListByIDs(ids []uint32, state uint8) ([]*model.Tag, error) {
	tag := model.Tag{State: state}
	return tag.ListByIDs(d.engine, ids)
}

// CountTag 根据name和state拿数量
func (d *Dao) CountTag(name string, state uint8) (int, error) {
	tag := model.Tag{Name: name, State: state}
	return tag.Count(d.engine)
}

// CreateTag 根据name和state和创作者向数据库添加tag
func (d *Dao) CreateTag(name string, state uint8, createdBy string) error {
	tag := model.Tag{
		Name:  name,
		State: state,
		Model: &model.Model{
			CreatedBy: createdBy,
		},
	}

	return tag.Create(d.engine)
}

// UpdateTag 根据id和values更新Dao中的tag
func (d *Dao) UpdateTag(id uint32, name string, state uint8, modifiedBy string) error {
	tag := model.Tag{
		Model: &model.Model{
			ID: id,
		},
	}
	values := map[string]interface{}{
		"state":       state,
		"modified_by": modifiedBy,
	}
	if name != "" {
		values["name"] = name
	}

	return tag.Update(d.engine, values)
}

// DeleteTag 根据id删除tag
func (d *Dao) DeleteTag(id uint32) error {
	tag := model.Tag{Model: &model.Model{ID: id}}
	return tag.Delete(d.engine)
}
