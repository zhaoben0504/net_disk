package router

import (
	"net/http"

	"suhc-gitlab-01.inovance.local/mnk/server/lcdp.git/middleware"
)

func initApplicationRouter() {
	list := []middleware.PermissionItem{
		{
			Method:      http.MethodPost,
			Handler:     applicationHandler.Add,
			URL:         "/lcdp/app",
			Permissions: []string{"A"},
		},
		{
			Method:      http.MethodGet,
			Handler:     applicationHandler.Test,
			URL:         "/lcdp/app/resources/:appid/:filename",
			Permissions: []string{"A"},
		},
	}

	middleware.GenerateHandler(Echo, list)
}
