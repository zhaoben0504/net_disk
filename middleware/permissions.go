package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type (
	// Err return object
	Err struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	// PermissionListFunc define permission list func
	// return:
	// 	nil: token invalid
	// 	[]string: token is valid, and return the permissions
	PermissionListFunc func(k string) []string

	// ContextFunc define context func
	ContextFunc func(k string) map[string]interface{}

	// GetResponseErrFunc get response error func
	GetResponseErrFunc func(lang string) interface{}

	// PermissionItem permission & url relation
	PermissionItem struct {
		Method      string           // http method ,golang Common HTTP methods
		Handler     echo.HandlerFunc // echo handler function
		URL         string           `json,yaml:"url"` // url support exp
		MasterKey   string           `json,yaml:"master_key"`
		Permissions []string         `json,yaml:"permissions"`
		Operation   string           `json,yaml:"operation"` // defines how to handle permission,only support "or"/"and"
	}

	// PermissionMiddlewareConfig permission middleware config
	PermissionMiddlewareConfig struct {
		// Key the key store user info
		Key string

		// Skipper defines a function to skip middleware.
		//Skipper middleware.Skipper

		// IgnoreURLs skip url list
		IgnoreURLs []string

		// PermissionList defines a function get permissions
		GetPermissionList PermissionListFunc

		// SetContext defines a function set context
		GetContext ContextFunc

		// InternalErrFunc sys internal err
		InternalErrFunc GetResponseErrFunc
		// TokenNotExistErrFunc token not exist from request
		TokenNotExistErrFunc GetResponseErrFunc
		// TokenInvalidErrFunc token not exist in redis
		TokenInvalidErrFunc GetResponseErrFunc
		// PermissionErrFunc permission invalid
		PermissionErrFunc GetResponseErrFunc
	}
)

var (
	DefaultPermissionConfig = PermissionMiddlewareConfig{
		Key:                  "token",
		GetPermissionList:    DefaultPermissionList,
		InternalErrFunc:      InternalErr,
		TokenNotExistErrFunc: TokenNotExistErr,
		TokenInvalidErrFunc:  TokenInvalidErr,
		PermissionErrFunc:    PermissionErr,
	}
	InternalErr      = func(lang string) interface{} { return Err{Code: 1, Msg: "Server Internal Error"} }
	TokenNotExistErr = func(lang string) interface{} { return Err{Code: 2, Msg: "Token Not Exist Error"} }
	TokenInvalidErr  = func(lang string) interface{} { return Err{Code: 3, Msg: "Token Invalid Error"} }
	PermissionErr    = func(lang string) interface{} { return Err{Code: 4, Msg: "Permission Error"} }
)

// DefaultPermissionList default PermissionList
func DefaultPermissionList(k string) []string {
	return []string{}
}

// Permission returns a middleware that logs HTTP requests.
func Permission() echo.MiddlewareFunc {
	return PermissionWithConfig(DefaultPermissionConfig)
}

// PermissionWithPermissionList returns a middleware that logs HTTP requests.
func PermissionWithPermissionList(permissionList PermissionListFunc) echo.MiddlewareFunc {
	DefaultPermissionConfig.GetPermissionList = permissionList
	return PermissionWithConfig(DefaultPermissionConfig)
}

// Skipper returns false which processes the middleware.
func (p PermissionMiddlewareConfig) Skipper(c echo.Context) bool {
	if len(p.IgnoreURLs) == 0 {
		return false
	}
	url := c.Request().RequestURI
	for _, k := range p.IgnoreURLs {
		if ok, _ := regexp.MatchString("^"+k+"$", url); ok {
			return true
		}
	}
	return false
}

// PermissionWithConfig returns a middleware that logs HTTP requests.
func PermissionWithConfig(config PermissionMiddlewareConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			token := ""
			auth := make(map[string]interface{})
			switch req.Method {
			case http.MethodGet, http.MethodDelete:
				token = req.URL.Query().Get(config.Key)
				tool.Logger.Infof("url: %s, method: %s", req.RequestURI, req.Method)
			case http.MethodPost, http.MethodPut:
				if strings.Contains(req.Header.Get("Content-Type"), "multipart/form-data") {
					token = c.FormValue(config.Key)
					query := ""
					q, err := c.FormParams()
					if err != nil {
						tool.Logger.Error(err.Error())
						return err
					}
					for a, b := range q {
						if len(query) > 0 {
							query += "&"
						}
						query += a + "=" + strings.Join(b, ",")
					}
					tool.Logger.Infof("url: %s, method: %s, content: %s", req.URL.Path, req.Method, query)
				} else {
					body, err := ioutil.ReadAll(req.Body)
					_ = req.Body.Close()
					if nil != err {
						tool.Logger.Error(err.Error())
						_ = c.JSON(http.StatusOK, config.InternalErrFunc(tool.GetHeaderLanguage(c.Request().Header)))
						return err
					}
					tool.Logger.Infof("url: %s, method: %s, content: %s", req.URL.Path, req.Method, string(body))
					err = json.NewDecoder(bytes.NewReader(body)).Decode(&auth)
					if nil != err {
						tool.Logger.Error(err.Error())
						_ = c.JSON(http.StatusOK, config.InternalErrFunc(tool.GetHeaderLanguage(c.Request().Header)))
						return err
					}
					req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
					if t, ok := auth[config.Key]; ok {
						token = fmt.Sprintf("%s", t)
					}
				}
			default:
				body, err := ioutil.ReadAll(req.Body)
				_ = req.Body.Close()
				if nil != err {
					tool.Logger.Error(err.Error())
					_ = c.JSON(http.StatusOK, config.InternalErrFunc(tool.GetHeaderLanguage(c.Request().Header)))
					return err
				}
				tool.Logger.Infof("url: %s, method: %s, content: %s", req.URL.Path, req.Method, string(body))
				err = json.NewDecoder(bytes.NewReader(body)).Decode(&auth)
				if nil != err {
					tool.Logger.Error(err.Error())
					_ = c.JSON(http.StatusOK, config.PermissionErrFunc(tool.GetHeaderLanguage(c.Request().Header)))
					return err
				}
				req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				if t, ok := auth[config.Key]; ok {
					token = fmt.Sprintf("%s", t)
				}
			}

			if config.Skipper(c) {
				return next(c)
			}

			if token == "" {
				_ = c.JSON(http.StatusOK, config.TokenNotExistErrFunc(tool.GetHeaderLanguage(c.Request().Header)))
				return errors.New("token nil")
			}

			if len(permissionCache) == 0 {
				return next(c)
			}
			str := ""
			if config.GetPermissionList != nil {
				permissions := config.GetPermissionList(token)
				if permissions == nil {
					return c.JSON(http.StatusOK, config.TokenInvalidErrFunc(tool.GetHeaderLanguage(c.Request().Header)))
				}
				str = strings.Join(permissions, ",")
			}

			for _, p := range permissionCache {
				if ok, _ := regexp.MatchString("^"+p.URL+"$", req.URL.Path); ok {
					if len(p.MasterKey) != 0 && strings.Contains(str, p.MasterKey) {
						break
					}
					if len(p.Permissions) == 0 {
						break
					}
					if p.Operation == "and" {
						count := 0
						for _, i := range p.Permissions {
							if strings.Contains(str, i) {
								count++
							}
						}
						if len(p.Permissions) != count {
							return c.JSON(http.StatusOK, config.PermissionErrFunc(tool.GetHeaderLanguage(c.Request().Header)))
						}
					} else {
						// or
						flag := false
						for _, i := range p.Permissions {
							if strings.Contains(str, i) {
								flag = true
								break
							}
						}
						if !flag {
							return c.JSON(http.StatusOK, config.PermissionErrFunc(tool.GetHeaderLanguage(c.Request().Header)))
						}
					}
					break
				}
			}

			if config.GetContext != nil {
				context := config.GetContext(token)
				for k, v := range context {
					c.Set(k, v)
				}
			}
			return next(c)
		}
	}
}

var permissionCache []PermissionItem

// GenerateHandler set handler to echo
func GenerateHandler(e *echo.Echo, list []PermissionItem) {
	reg := "[0-9a-zA-Z_\\.\\-]+"

	if len(list) == 0 || e == nil {
		return
	}
	for _, i := range list {
		switch i.Method {
		case http.MethodGet:
			e.GET(i.URL, i.Handler)
		case http.MethodHead:
			e.HEAD(i.URL, i.Handler)
		case http.MethodPost:
			e.POST(i.URL, i.Handler)
		case http.MethodPut:
			e.PUT(i.URL, i.Handler)
		case http.MethodPatch:
			e.PATCH(i.URL, i.Handler)
		case http.MethodDelete:
			e.DELETE(i.URL, i.Handler)
		case http.MethodConnect:
			e.CONNECT(i.URL, i.Handler)
		case http.MethodOptions:
			e.OPTIONS(i.URL, i.Handler)
		case http.MethodTrace:
			e.TRACE(i.URL, i.Handler)
		}

		if strings.Contains(i.URL, ":") {
			strs := strings.Split(i.URL, "/")
			for j := range strs {
				if strings.Contains(strs[j], ":") {
					strs[j] = reg
				}
			}
			i.URL = strings.Join(strs, "/")
		}
		permissionCache = append(permissionCache, i)
	}
}
