package dao

import (
	"context"
	util "github.com/go-study-lab/go-mall/common/utils"
	"github.com/go-study-lab/go-mall/dal/model"
	"github.com/go-study-lab/go-mall/logic/do"
)

type DemoDao struct {
	ctx context.Context
}

func NewDemoDao(ctx context.Context) *DemoDao {
	return &DemoDao{ctx: ctx}
}

func (demo *DemoDao) GetAllDemos() (demos []*model.Demo, err error) {

	err = DB().WithContext(demo.ctx).Find(&demos).Error
	if err != nil {
		return nil, err
	}

	return demos, err
}

func (demo *DemoDao) CreateDemo(demoObj *do.Demo) (*model.Demo, error) {
	model := new(model.Demo)
	err := util.CopyProperties(model, demoObj)
	if err != nil {
		return nil, err
	}
	err = DBMaster().WithContext(demo.ctx).Create(model).Error
	return model, err
}
