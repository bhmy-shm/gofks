package cascade

import "gorm.io/gorm"

type (
	Cascade struct {
		gorm.Model
		CascadeData
	}

	CascadeData struct {
		Gid            string `json:"gid"`
		ApiAddr        string `json:"apiAddr"`
		Name           string `json:"name"`
		SipAddr        string `json:"sipAddr"`
		AuthPwd        string `json:"authPwd"`
		AuthUser       string `json:"authUser"`
		Webhook        string `json:"webhook"`
		CascadeId      string `json:"cascadeId"`
		LowerGroupName string `json:"lowerGroupName"`
		LowerGid       string `json:"lowerGid"`
		LowerSn        string `json:"lowerSn"`
		LowerType      int    `json:"lowerType"`
		CascadeConnStatus
	}

	CascadeConnStatus struct {
		Status  int `json:"status"`
		ErrCode int `json:"errCode"`
	}
)
