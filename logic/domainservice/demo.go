package domainservice

import (
	"context"
	"github.com/go-study-lab/go-mall/common/errcode"
	"github.com/go-study-lab/go-mall/dal/dao"
	"github.com/go-study-lab/go-mall/dal/model"
)

// 演示DEMO, 后期使用时删掉

type DemoDomainSvc struct {
	ctx     context.Context
	DemoDao *dao.DemoDao
}

func NewDemoDomainSvc(ctx context.Context) *DemoDomainSvc {
	return &DemoDomainSvc{
		ctx:     ctx,
		DemoDao: dao.NewDemoDao(ctx),
	}
}

// GetDemos 配置GORM时的演示方法
func (dds *DemoDomainSvc) GetDemos() ([]*model.Demo, error) {
	demos, err := dds.DemoDao.GetAllDemos()
	if err != nil {
		err = errcode.Wrap("query entity error", err)
		return nil, err
	}

	demoOrders := make([]*model.Demo, 0, len(demos))
	// 后面会介绍工具, Model到Domain Object 可以一键转换
	for _, demo := range demos {
		demoOrders = append(demoOrders, &model.Demo{
			ID:         demo.ID,
			Account:    demo.Account,
			UpperLimit: demo.UpperLimit,
			Game:       demo.Game,
			Lang:       demo.Lang,
			GameCat:    demo.GameCat,
			Operator:   demo.Operator,
			CreatedAt:  demo.CreatedAt,
			UpdatedAt:  demo.UpdatedAt,
		})
	}

	return demoOrders, nil
}
