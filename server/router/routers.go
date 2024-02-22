package router

import (
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

func RegisterHandlers(server *rest.Server, serverCtx *cmd.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/user/login",
				Handler: UserLoginHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/user/details",
				Handler: UserDetailsHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/mail/code/send/register",
				Handler: MailCodeSendRegisterHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/user/register",
				Handler: UserRegisterHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/share/details",
				Handler: ShareDetailsHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Auth},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/file/upload",
					Handler: FileUploadHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/user/repository/save",
					Handler: UserRepositorySaveHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/user/file/list",
					Handler: UserFileListHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/user/file/update",
					Handler: UserFileNameUpdateHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/user/folder/create",
					Handler: UserFolderCreateHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/user/file/delete",
					Handler: UserFileDeleteHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/user/file/move",
					Handler: UserFileMoveHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/user/share/create",
					Handler: ShareBasicCreateHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/file/shared/save",
					Handler: SaveSharedFileHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/refresh/auth",
					Handler: RefrshAuthHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/file/upload/prepare",
					Handler: FileUploadPrepareHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/file/chunk/upload",
					Handler: FileUploadChunkHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/file/chunk/upload/finish",
					Handler: FileUploadChunkFinishHandler(serverCtx),
				},
			}...,
		),
	)
}
