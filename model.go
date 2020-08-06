package mnd

type Ret struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}

type RetFast struct {
	Data *m_r_fast `json:"data"`
}
