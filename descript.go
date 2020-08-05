package mnd

type m_addSend struct {
	Sender      string   `json:"sender"`
	Remark      string   `json:"remark"`
	Isextension bool     `json:"isextension"`
	Pwd         string   `json:"pwd"`
	Expire      int      `json:"expire"`
	Recvs       []string `json:"recvs"`
	FileSize    int64    `json:"file_size"`
	FileCount   int      `json:"file_count"`
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
}
