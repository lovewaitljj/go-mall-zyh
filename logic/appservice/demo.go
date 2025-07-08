package appservice

import (
	"context"
	"fmt"
	"github.com/go-study-lab/go-mall/api/request"
	"github.com/go-study-lab/go-mall/api/response"
	"github.com/go-study-lab/go-mall/common/errcode"
	util "github.com/go-study-lab/go-mall/common/utils"
	"github.com/go-study-lab/go-mall/logic/do"
	"github.com/go-study-lab/go-mall/logic/domainservice"
)

// 演示DEMO, 后期使用时删掉

type DemoAppSvc struct {
	ctx           context.Context
	demoDomainSvc *domainservice.DemoDomainSvc
}

func NewDemoAppSvc(ctx context.Context) *DemoAppSvc {
	return &DemoAppSvc{
		ctx:           ctx,
		demoDomainSvc: domainservice.NewDemoDomainSvc(ctx),
	}
}

//func (das *DemoAppSvc)DoSomething() {
//	demo, err := das.demoDomainSvc.GetDemoEntity(id)
//	if err != nil {
//		logger.New(das.ctx).Error("DemoAppSvc DoSomething err", err)
//		return err
//	}
//	......
//}

// GetDemoIdentities 配置GORM时的演示方法, 显的有点脑残,
// 后面章节再解释怎么用ApplicationService 进行逻辑解耦
func (das *DemoAppSvc) GetDemoIdentities() ([]uint32, error) {
	demos, err := das.demoDomainSvc.GetDemos()
	if err != nil {
		return nil, err
	}
	identities := make([]uint32, 0, len(demos))

	for _, demo := range demos {
		identities = append(identities, demo.ID)
	}
	return identities, nil
}

func (das *DemoAppSvc) CreateDemoOrder(req *request.DemoCreate) (*response.Demo, error) {
	// 1.转化为领域对象
	demoOrder := new(do.Demo)
	err := util.CopyProperties(demoOrder, req)
	if err != nil {
		errcode.Wrap("请求转化成demoOrderDo失败", err)
		return nil, err
	}
	demoOrderDo, err := das.demoDomainSvc.CreateDemoOrders(demoOrder)
	if err != nil {
		return nil, err
	}
	//todo 做一些外围逻辑，例如异步飞书推送等,可能涉及第三方lib
	fmt.Println("异步逻辑等")
	replyDemoOrder := new(response.Demo)
	err = util.CopyProperties(replyDemoOrder, demoOrderDo)
	if err != nil {
		errcode.Wrap("请求转化成demoOrderResp失败", err)
		return nil, err
	}

	return replyDemoOrder, nil
}
