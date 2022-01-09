package httptest

import "net/http"

var (
	RequestParams = map[string]interface{}{"page": 1, "pageSize": 10}
	LoginParams   = map[string]interface{}{"username": "admin", "password": "123456"}
	LoginResponse = Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: "操作成功"},
		{Key: "data",
			Value: Responses{
				{Key: "accessToken", Value: "", Type: "notempty"},
			},
		},
	}
	LogoutResponse = Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: "操作成功"},
	}
	ResponseDatas = Responses{
		{Key: "pageSize", Value: 10},
		{Key: "page", Value: 1},
	}
)
