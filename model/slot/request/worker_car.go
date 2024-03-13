package request

import "oplian/model/common/request"

type QueryWorkerCarReq struct {
	request.PageInfo
	KeyWord    string `json:"keyWord" form:"keyWord"`
	TaskStatus string `json:"taskStatus" form:"taskStatus"`
	MinerId    string `json:"minerId" form:"minerId"`
}

type ModifyWorkerCarReq struct {
	Id            string `json:"id"`
	TaskStatus    string `json:"taskStatus"`
	WorkerTaskNum string `json:"workerTaskNum"`
}

type QueryWorkerCarDetailReq struct {
	request.PageInfo
	Id string `json:"id" form:"id"`
}
