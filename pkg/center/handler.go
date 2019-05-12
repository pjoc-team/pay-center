package center

import (
	context2 "context"
	"github.com/gin-gonic/gin"
	"github.com/pjoc-team/base-service/pkg/logger"
	"github.com/pjoc-team/pay-proto/go"
	"net/http"
)

func StartGin(service *PayService, listenAddr string) {
	engine := gin.New()
	engine.LoadHTMLGlob("html/templates/*")
	engine.GET("/", func(c *gin.Context) {
		//c.String(200, "pang pong porobong")
		c.File("./html/index.html")
	})
	engine.GET("/bank", func(context *gin.Context) {
		response := Request(service.PayCenterConfig.PayGatewayUrl)
		logger.Log.Infof("Response: %v", response)
		context.HTML(http.StatusOK, "bank.html", response.Data)
	})

	engine.POST("/wechat", func(context *gin.Context) {
		pay := &WechatPay{}
		if err := context.Bind(pay); err != nil {
			logger.Log.Errorf("Failed to bind: %v", err.Error())
		}
		logger.Log.Infof("Amt: %v", pay.OrderAmt)
		response := pay.Request(service.PayCenterConfig.PayGatewayUrl)
		logger.Log.Infof("Response: %v", response)
		context.HTML(http.StatusOK, "wechat.html", response.Data)
	})

	engine.GET("/query/:order_id", func(context *gin.Context) {
		orderId := context.Param("order_id")
		if orderId == "" {
			logger.Log.Errorf("OrderId is null")
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}
		serviceClient, e := service.GrpcClientFactory.GetDatabaseClient()
		if e != nil || serviceClient == nil{
			logger.Log.Errorf("failed to get grpc client!")
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}
		baseOrder := &pay.BasePayOrder{OutTradeNo: orderId}
		order := &pay.PayOrder{BasePayOrder: baseOrder}
		response, err := serviceClient.FindPayOrder(context2.TODO(), order)
		if err != nil{
			logger.Log.Errorf("Failed to find order! error: %v", err)
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}
		context.JSON(http.StatusOK, response)

	})

	engine.GET("/notify/:gateway_order_id", handleGatewayOrderIdFunc(service)).
		POST("/notify/:gateway_order_id", handleGatewayOrderIdFunc(service))
	//listenAddr := fmt.Sprintf(":%d", port)
	engine.Run(listenAddr)
}

func handleGatewayOrderIdFunc(service *PayService) func(*gin.Context) {
	return func(context *gin.Context) {
	}
}
