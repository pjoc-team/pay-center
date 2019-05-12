package center

import (
	"flag"
	"github.com/pjoc-team/base-service/pkg/grpc"
	gc "github.com/pjoc-team/base-service/pkg/grpc"
	"github.com/pjoc-team/base-service/pkg/service"
)

const ETCD_DIR_ROOT = "/pub/pjoc/pay"


type PayCenterConfig struct {
	PayGatewayUrl string `yaml:"PayGatewayUrl"`
}

type PayService struct {
	*service.Service
	*grpc.GrpcClientFactory
	*service.GatewayConfig
	*PayCenterConfig
}

func Init(svc *service.Service) *PayService {
	payservice := &PayService{}
	payservice.Service = svc
	flag.Parse()

	gatewayConfig := service.InitGatewayConfig(svc.EtcdPeers, ETCD_DIR_ROOT)
	payservice.GatewayConfig = gatewayConfig

	grpcClientFactory := gc.InitGrpFactory(*svc, payservice.GatewayConfig)
	payservice.GrpcClientFactory = grpcClientFactory
	//gatewayConfig := service.InitGatewayConfig(svc.EtcdPeers, ETCD_DIR_ROOT)
	//payservice.GatewayConfig = gatewayConfig
	return payservice
}
