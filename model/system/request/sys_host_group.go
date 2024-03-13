package request

type CreateSysHostGroupReq struct {
	GroupNames []string `json:"groupNames"`
}

type UpdateSysHostGroupReq struct {
	ID        int    `json:"id" form:"id"`
	GroupName string `json:"groupName" form:"groupName"`
}
