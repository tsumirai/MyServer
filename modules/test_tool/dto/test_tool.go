package dto

type DivideTableReq struct {
	TabelName string `json:"table_name"`
	TableNum  int64  `json:"table_num"`
}
