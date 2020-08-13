package model

// 前端上传返回json
type ImgInfo struct {
	Imgnames string `json:"imgnames"`
	Imgurls  string `json:"imgurls"`
}

type UploadResp struct {
	Code       string    `json:"code"`
	Info       string    `json:"info"`
	Data       []ImgInfo `json:"data"`
	Exceptions string    `json:"exceptions"`
}

// 通用api响应
type CommonResp struct {
	Code int         `json:"code"` //响应编码 0 成功 500 错误 403 无权限  -1  失败
	Msg  string      `json:"msg"`  //消息
	Data interface{} `json:"data"` //数据内容
}

// layui table 响应
type LayuiResp struct {
	Code  int                      `json:"code"` //响应编码 0 成功 500 错误 403 无权限  -1  失败
	Msg   string                   `json:"msg"`  //消息
	Count int                      `json:"count"`
	Data  []map[string]interface{} `json:"data"`
}

// 通用分页表格响应
type TableDataInfo struct {
	Total int         `json:"total"` //总数
	Rows  interface{} `json:"rows"`  //数据
	Code  int         `json:"code"`  //响应编码 0 成功 500 错误 403 无权限
	Msg   string      `json:"msg"`   //消息
}

// 通用的删除请求
type RemoveReq struct {
	Ids string `form:"ids"  binding:"required"`
}

// 通用详情请求
type DetailReq struct {
	Id int64 `json:"id"` //主键ID
}

// 通用修改请求
type EditReq struct {
	Id int64 `json:"id"` //主键ID
}

// 登陆Ajax判断
type IsLogin struct {
	Lgoinret int `json:"lgoinret"`
}

// 通用返回
type AjaxResp struct {
	ResultCode int    `json:"resultCode"`
	Url        string `json:"url,omitempty"`
	ErrorMsg   string `json:"errorMsg,omitempty"`
	Msg        string `json:"msg,omitempty"`
}

// 注册返回
type RegisterResp struct {
	Ret    int    `json:"ret"`
	Zctype int    `json:"zctype"`
	Msg    string `json:"msg,omitempty"`
}
