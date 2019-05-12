package center

import (
	"encoding/json"
	"fmt"
	"github.com/pjoc-team/base-service/pkg/date"
	"github.com/pjoc-team/base-service/pkg/logger"
	url2 "github.com/pjoc-team/base-service/pkg/url"
	pb "github.com/pjoc-team/pay-proto/go"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type WechatPay struct {
	OrderAmt uint32 `form:"order_amt" binding:"required"`
}

func (w WechatPay) Request(domain string) (*pb.PayResponse) {
	request := &pb.PayRequest{}
	request.ChannelId = "personal"
	request.Method = pb.Method_name[int32(pb.Method_QR_CODE)]
	request.ProductName = "apple"
	request.ProductDescribe = "apple"
	request.UserIp = "127.0.0.1"
	request.ExtJson = "{\"channel\":\"WECHAT\"}"
	request.PayAmount = w.OrderAmt
	if request.PayAmount <= 0 {
		request.PayAmount = 100
	}
	request.OutTradeNo = g.GenerateOrderId()
	request.OrderTime = date.NowTime()
	request.AppId = "1"
	request.SignType = "RSA"

	bytes, e := json.Marshal(request)
	if e != nil {
		logger.Log.Errorf("marshal error! err: %v", e.Error())
	}

	privateKey := `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA2KaaJp7JeW91WlQCfZeS14US/ot9hIJViutv3JHojdgTx+8A
8psStKaPl2Ac/MTJ/3mHeopCObmgjw/Au/Ne0PS1rveY0Pcazwnp+R1TDP2H9jag
c3GJWS6cvHLB/B4uP3LOnPXN8ctwDVsF19b/howVKUKX6RAX7R2VAEyTIZJIEIQE
0fNvRCWqbVv1RB3LU4cbQmW6nX8dP793fP8s/Lhzcj6vS6UKxLVl5CrCCGIJIBYc
1mI8RbUYvGqwiONEnEwYvOioAoAlkMIXdFndIjngHe7JYfGW1NtPzHLG5yw8anYT
D/3du7hJ/kSN0WM6NLa0P/vbR5+mxVdoRzY+kQIDAQABAoIBABaE2qkBADgbGbuV
19xuENlN/7dtkFJhqbqS1kG6+M0llIjHkvWkoMEePvahCuJLIiPn4ekezdtqLAIy
xPnERiq6BNh26+9sf+DdSvCV17gV8jfpXawiNQCME8aStw8Zo/z8VfWCpzFmz/LT
bzwMIOs/TEPJpDiZb6M52+74BqMKfHTY14YOF8Xr4fiaUFpNTViHeOQXKzoG5PF4
GLlhg7YNgEnjyc578izCoFp/xTjBBHQ7dtu+EnzmXD9QTlz7xUYt4P2TjUEBKy1o
xSxDpgFL+BKYgRazilkrJ2hesbGCvbxDzcd4ivzpfmvqkN74Lq0vF9voL1JSd6D2
3l/R9bECgYEA7nbQgsK3ResReUMumJvE4y1sl2D+rt24QHlu+jOJqVkpAo0L4HLj
vCX0Y8tBfG/hDc5iC12YILCn+EEb9bD2giURg7V+cA+K4IJrbLTbnna/UlvA1PFK
3kHFosdCk5cRlpAppEBLQEUjlf7mjp6k2Xxy71ozg4KlB3wf3QCrs2MCgYEA6JUk
iXd/lntdjb7V/QdwhVFdp/lzst0ClE4q04RNL8ZjwmSrYOAOGO2ktKOBG8lGT3P6
54/BASn9TMOXks8gPE3r/pN+21RGvOq2xtHNOrnV5g1RvlqHtwtv2RUxoEoTKPjB
m6KDeLrPNCuGZ3bzYUUNAys66v3iWM5PK2s1GnsCgYEAjtKyx95/jmzgJlTKj7Sc
E8SdCX2ajHlXZaZVhZ1gkgFIwrJfrqqhI4tH+I1AR5tqm65EorIH72xe7h1w9ZJr
0j8JYm1NsShd8WGrnYwlDZ/prxYtRFzQjpWuHXRit6r/acImbq3jZDcEvU3SIRF7
gpc674iC2f1hgj4hh2hjbikCgYBWfG8ztv34xTMKrHYCOyv6R0FeXwJI9qoo39BJ
Cx9wroMWHD0mLurPFj9y9IHkBTph/SzFwszwU97fFrRcYS0Jf6hL6Cj6AiKzyUvi
Ls30EnqZq0ZEVIG27UfQH3NuuVzalXXZG9trn3vBWJYID1F9UCIAlai5DWOHxl/m
M11x1QKBgCvi33kFLll6SVIKfkwt1Hja/DGlyq/M/xN4qn/wQwGKYzzIU+73SbQR
g44GAiYVMJQrjISg/RVd4ClDxZ+A0cpumfpuSJdcT210L4u5FkuTQAmLZ2HOhTzK
mW9/iR9koFHtTzTKhhYIgSWy9EWkQmcyrOKnEPYqMJjMobDJ1AuG
-----END RSA PRIVATE KEY-----`
	sign := SignString(string(bytes), privateKey)
	request.Sign = sign

	client := http.Client{Timeout: 6 * time.Second}
	path := "/v1/pay/" + request.Method
	url := url2.CompactUrl(domain, path, "")

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		logger.Log.Errorf("err: %v", err.Error())
	}

	reader := strings.NewReader(string(jsonBytes))
	resp, err2 := client.Post(url, "application/json", reader)
	if err2 != nil {
		logger.Log.Errorf("err: %v", err2.Error())
		return nil
	}
	closer := resp.Body
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	response := &pb.PayResponse{}
	all, err2 := ioutil.ReadAll(closer)
	logger.Log.Infof("Get response: %v", string(all))
	err2 = json.Unmarshal(all, response)
	i, err := strconv.Atoi(response.Data["amt"])
	if err != nil {
		return nil
	}
	a := float32(i) / 100.0
	response.Data["amt"] = fmt.Sprintf("%.2f", a)
	response.Data["outTradeNo"] = request.OutTradeNo
	if err2 != nil {
		logger.Log.Errorf("err: %v", err2.Error())
		return nil
	}
	return response
}
