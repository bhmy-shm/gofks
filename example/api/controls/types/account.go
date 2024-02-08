package types

type UserAddParams struct {
	ID         uint64 `json:"id,omitempty"`
	OrgID      int64  `json:"orgID,omitempty"`
	RoleID     uint64 `json:"roleID,omitempty"`
	UserType   uint32 `json:"userType,omitempty"`
	OrgName    string `json:"orgName,omitempty"`
	RoleName   string `json:"roleName,omitempty"`
	Account    string `json:"account,omitempty"`
	Name       string `json:"name,omitempty"`
	Pass       string `json:"pass,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Mobile     string `json:"mobile,omitempty"`
	Phone      string `json:"phone,omitempty"`
	LoginCount int64  `json:"loginCount,omitempty"`
}

type UserListParams struct {
	ID       uint64 `json:"id,omitempty"`
	PageNum  int32  `json:"pageNum"`
	PageSize int32  `json:"pageSize"`
}

type UserListResponse struct {
	Total int64
	List  interface{}
}
