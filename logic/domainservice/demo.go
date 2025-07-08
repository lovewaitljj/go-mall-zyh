package domainservice

import (
	"context"
	"github.com/go-study-lab/go-mall/common/errcode"
	util "github.com/go-study-lab/go-mall/common/utils"
	"github.com/go-study-lab/go-mall/dal/dao"
	"github.com/go-study-lab/go-mall/logic/do"
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
func (dds *DemoDomainSvc) GetDemos() ([]*do.Demo, error) {
	demos, err := dds.DemoDao.GetAllDemos()
	if err != nil {
		err = errcode.Wrap("query entity error", err)
		return nil, err
	}

	demoOrders := make([]*do.Demo, 0, len(demos))
	// 后面会介绍工具, Model到Domain Object 可以一键转换
	for _, demo := range demos {
		demoOrder := new(do.Demo)
		util.CopyProperties(demoOrder, demo)
		demoOrders = append(demoOrders, demoOrder)
	}

	return demoOrders, nil
}

func (dds *DemoDomainSvc) CreateDemoOrders(demoOrder *do.Demo) (*do.Demo, error) {
	if operator, ok := dds.ctx.Value("operator").(string); ok {
		demoOrder.Operator = operator
	}
	demoOrderModel, err := dds.DemoDao.CreateDemo(demoOrder)
	if err != nil {
		err = errcode.Wrap("create demo order error", err)
		return nil, err
	}
	err = util.CopyProperties(demoOrder, demoOrderModel)
	return demoOrder, err
}
