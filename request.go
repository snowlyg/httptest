package httptest

import (
	"net/http"
)

var (
	// default page request params
	GetRequestFunc = NewWithQueryObjectParamFunc(map[string]interface{}{"page": 1, "pageSize": 10})

	// default page request params
	PostRequestFunc = NewWithJsonParamFunc(map[string]interface{}{"page": 1, "pageSize": 10})

	// default login request params
	LoginFunc = NewWithJsonParamFunc(map[string]interface{}{"username": "admin", "password": "123456"})

	// default login response params
	LoginResponse = Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: "操作成功"},
		{Key: "data",
			Value: Responses{
				{Key: "accessToken", Value: "", Type: "notempty"},
			},
		},
	}
	// default logout response params
	LogoutResponse = Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: "操作成功"},
	}

	SuccessResponse = Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: "操作成功"},
	}

	// default data response params
	ResponsePage = Responses{
		{Key: "status", Value: http.StatusOK},
		{Key: "message", Value: "操作成功"},
		{Key: "data", Value: Responses{
			{Key: "pageSize", Value: 10},
			{Key: "page", Value: 1},
		}},
	}
)
