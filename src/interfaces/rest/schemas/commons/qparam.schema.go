package commonschema

type QueryParams struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}
