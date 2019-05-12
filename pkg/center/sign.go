package center

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pjoc-team/base-service/pkg/sign"
	"sort"
	"strings"
)

func MapToString(params map[string]interface{}) string {
	key := make([]string, 0)
	for k, _ := range params {
		key = append(key, k)
	}

	sort.Strings(key)

	builder := strings.Builder{}
	delimiter := ""
	//for _, key := range p.SortedKeyFieldNames {
	//	value := params[key]
	//	if p.IgnoreEmptyValue && value == "" {
	//		continue
	//	}
	//	builder.WriteString(delimiter)
	//	builder.WriteString(key)
	//	builder.WriteString(p.KeyValueDelimiter)
	//	builder.WriteString(value)
	//	delimiter = p.PairsDelimiter
	//}
	for _, k := range key {
		value := params[k]
		if value == "" {
			continue
		}
		builder.WriteString(delimiter)
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(fmt.Sprintf("%v", value))
		delimiter = "&"
	}
	return builder.String()
}

func SignString(parametersJson string, privateKey string) string {
	if parametersJson == "" {
		//*parametersJson = `{"pay_amount":"1", "out_trade_no": "201810281512542133254", "order_time":"2018-10-28 12:00:00", "app_id":"1", "sign":"zn72Zc/r5gFjSByynONkZLua3VETMFZkpNd9yxCXW8H+aFxEVe5M7xefUBVHl2zmMjDbwoTGzbpOhI43k157VqyPUyXR0ReSxRlDN1QMMK/YSjMzJILLAVYfwUcjkhOzbvIF4r0QnI1RpbRCtGRpyLcDZnjr7xuEPcAmaje+GONFAw5qsjtxnHJr3Eb3RQUGwj5A0wfwVPVpqEuaKvMy/Obd8KwRpu0Vfo5i6LZRAhn3IqZLoatVon6jv4H95gywkyg1D03wZf8KYGBiBI5wplqpxgvWBwrupZfnffnIecMFqbKu6AKY0a16PXKzZuQGuH9nLp2imHdX6DYEGBAO2w==", "channel_id":"demo", "sign_type":"RSA"}`
		flag.Usage()
		panic("params required")
		//*parametersJson = "{}"
	}
	params := make(map[string]interface{})
	e := json.Unmarshal([]byte(parametersJson), &params)
	if e != nil {
		panic(e)
	}
	message := MapToString(params)
	//compacter := sign.NewParamsCompacter(AType{}, "json", []string{"sign"}, true, "&", "=")
	if privateKey == "" {
		privateKey = `-----BEGIN RSA PRIVATE KEY-----
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
	}
	//	publicKey := `-----BEGIN PUBLIC KEY-----
	//MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2KaaJp7JeW91WlQCfZeS
	//14US/ot9hIJViutv3JHojdgTx+8A8psStKaPl2Ac/MTJ/3mHeopCObmgjw/Au/Ne
	//0PS1rveY0Pcazwnp+R1TDP2H9jagc3GJWS6cvHLB/B4uP3LOnPXN8ctwDVsF19b/
	//howVKUKX6RAX7R2VAEyTIZJIEIQE0fNvRCWqbVv1RB3LU4cbQmW6nX8dP793fP8s
	///Lhzcj6vS6UKxLVl5CrCCGIJIBYc1mI8RbUYvGqwiONEnEwYvOioAoAlkMIXdFnd
	//IjngHe7JYfGW1NtPzHLG5yw8anYTD/3du7hJ/kSN0WM6NLa0P/vbR5+mxVdoRzY+
	//kQIDAQAB
	//-----END PUBLIC KEY-----`

	sign, err := sign.SignPKCS1v15([]byte(message), []byte(privateKey), crypto.SHA256)
	if err != nil {
		fmt.Println("sign error: ", err.Error())
		return ""
	}
	signBase64 := string(base64.StdEncoding.EncodeToString(sign))
	fmt.Println(signBase64)
	return signBase64
}
