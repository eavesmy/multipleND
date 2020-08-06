package mnd

/*
{"sender":"","remark":"","isextension":false,"pwd":"","expire":2,"recvs":["social","public"],"file_size":4496,"file_count":1,"notSaveTo":false,"trafficStatus":0,"downPreCountLimit":0}
*/
type m_addSend struct {
	Sender            string   `json:"sender"`
	Remark            string   `json:"remark"`
	Isextension       bool     `json:"isextension"`
	Pwd               string   `json:"pwd"`
	Expire            int      `json:"expire"`
	Recvs             []string `json:"recvs"`
	FileSize          int64    `json:"file_size"`
	FileCount         int      `json:"file_count"`
	NotSaveTo         bool     `json:"notSaveTo"`
	TrafficStatus     int      `json:"trafficStatus"`
	DownPreCountLimit int      `json:"downPreCountLimit"`
}

type m_getUpId struct {
	Preid      string `json:"preid"`
	Boxid      string `json:"boxid"`
	Linkid     string `json:"linkid"`
	Utype      string `json:"utype"`
	OriginUpid string `json:"originUpid"`
	Length     int64  `json:"length"`
	Count      int    `json:"count"`
}

type m_upload struct {
	IsPart bool   `json:"ispart"`
	Fname  string `json:"fname"`
	Fsize  int64  `json:"fsize"`
	Upid   string `json:"upId"`
	PartNu int    `json:"partnu,omitempty"`
}

type m_fast struct {
	Hash map[string]string `json:"hash"`
	Uf   map[string]string `json:"uf"`
	UpId string            `json:"upId"`
}

type m_r_fast struct {
	Status int               `json:"status"`
	Ufile  map[string]string `json:"ufile"`
}

type m_send struct {
	Bid     string `json:"bid"`
	Tid     string `json:"tid"`
	UfileId string `json:"ufileid"`
}

type m_complete struct {
	IsPart   bool              `json:"ispart"`
	Fname    string            `json:"fname"`
	UpId     string            `json:"upId"`
	Location map[string]string `json:"location"`
}
type m_token struct {
	DevInfo string `json:"dev_info"`
}
