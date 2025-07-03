package models

// Simples representa a tabela 'simples' no banco de dados.
type Simples struct {
	CNPJBasico          string `json:"cnpj_basico" db:"cnpj_basico"`                     // "cnpj_basico","text" [cite: 1]
	OpcaoSimples        string `json:"opcao_simples" db:"opcao_simples"`                 // "opcao_simples","text" [cite: 1]
	DataOpcaoSimples    string `json:"data_opcao_simples" db:"data_opcao_simples"`       // "data_opcao_simples","text" [cite: 1]
	DataExclusaoSimples string `json:"data_exclusao_simples" db:"data_exclusao_simples"` // "data_exclusao_simples","text" [cite: 1]
	OpcaoMEI            string `json:"opcao_mei" db:"opcao_mei"`                         // "opcao_mei","text" [cite: 1]
	DataOpcaoMEI        string `json:"data_opcao_mei" db:"data_opcao_mei"`               // "data_opcao_mei","text" [cite: 1]
	DataExclusaoMEI     string `json:"data_exclusao_mei" db:"data_exclusao_mei"`         // "data_exclusao_mei","text" [cite: 1]
}
