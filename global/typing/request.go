package typing

type ShData struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type WSData struct {
	Timestamp int64       `json:"timestamp"` //消息时间戳
	Sender    string      `json:"sender"`    //消息发送者
	Type      string      `json:"type"`      //消息种类
	Data      interface{} `json:"data"`      //ws结构化数据
	Detail    string      `json:"detail"`    //ws通知消息
}

type AddListenerRequestStruct struct {
	Name  string `json:"name" binding:"required"`
	LHOST string `json:"lhost" binding:"required"`
	LPORT int    `json:"lport" binding:"required,gt=0,lt=65536"`
}

type LoginRequestStruct struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseStruct struct {
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
	ExpiresAt   int64  `json:"expired_at"`
}
