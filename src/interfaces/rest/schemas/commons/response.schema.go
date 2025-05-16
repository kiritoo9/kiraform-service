package commonschema

type ResponseHTTP struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   any    `json:"error"`
}

type ResponseList struct {
	Parameters QueryParams `json:"parameters"`
	TotalPage  int         `json:"total_page"`
	Rows       any         `json:"rows"`
}
