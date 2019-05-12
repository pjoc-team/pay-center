package main

import (
	"flag"
	"github.com/pjoc-team/base-service/pkg/config/yaml"
	"github.com/pjoc-team/base-service/pkg/logger"
	"github.com/pjoc-team/base-service/pkg/service"
	"github.com/pjoc-team/pay-center/pkg/center"
)

var (
	listenAddr            = flag.String("listen-addr", ":8088", "HTTP listen address.")
	configURI             = flag.String("c", "config.yaml", "uri to load config")
	tlsEnable             = flag.Bool("tls", false, "enable tls")
	logLevel              = flag.String("log-level", "debug", "logger level")
	logFormat             = flag.String("log-format", "text", "text or json")
	caCert                = flag.String("ca-cert", service.WithConfigDir("ca.pem"), "Trusted CA certificate.")
	tlsCert               = flag.String("tls-cert", service.WithConfigDir("cert.pem"), "TLS server certificate.")
	tlsKey                = flag.String("tls-key", service.WithConfigDir("key.pem"), "TLS server key.")
	serviceName           = flag.String("s", "", "PayGatewayService name in service discovery.")
	registerServiceToEtcd = flag.Bool("r", true, "Register service to etcd.")
	etcdPeers             = flag.String("etcd-peers", "", "Etcd peers. example: 127.0.0.1:2379,127.0.0.1:12379")
)


func main() {
	flag.Parse()
	serviceDir := center.ETCD_DIR_ROOT + "/services"
	svc := service.InitService(*listenAddr,
		*configURI,
		*tlsEnable,
		*logLevel,
		*logFormat,
		*caCert,
		*tlsCert,
		*tlsKey,
		*serviceName,
		*registerServiceToEtcd,
		*etcdPeers,
		serviceDir)

	payCenterConfig := &center.PayCenterConfig{}
	err := yaml.UnmarshalFromFile("conf/config.yaml", payCenterConfig)
	if err != nil {
		logger.Log.Errorf("failed to unmarshal payCenterConfig, error: %v", err.Error())
		panic(err)
	}

	notifyService := center.Init(svc)
	notifyService.PayCenterConfig = payCenterConfig

	center.StartGin(notifyService, *listenAddr)
}
