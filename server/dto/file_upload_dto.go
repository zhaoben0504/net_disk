package dto

type FileUploadRequest struct {
	Hash string `json:"hash, optional"`
	Name string `json:"name, optional"`
	Ext  string `json:"ext, optional"`
	Size int64  `json:"size, optional"`
	Path string `json:"path, optional"`
}

type FileUploadResponse struct {
	Identity string `json:"identity"`
	Msg      string `json:"msg"`
	Ext      string `json:"ext"`
	Name     string `json:"name"`
}
