package models

// CNAE representa a tabela 'cnae' no banco de dados.
type CNAE struct {
	Codigo    string `json:"codigo" db:"codigo"`
	Descricao string `json:"descricao" db:"descricao"`
}
