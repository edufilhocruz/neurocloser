package models

type Empresa struct {
	CNPJBasico                string  `json:"cnpj_basico" db:"cnpj_basico"`                                 // "cnpj_basico","text" [cite: 1]
	RazaoSocial               string  `json:"razao_social" db:"razao_social"`                               // "razao_social","text" [cite: 1]
	NaturezaJuridica          string  `json:"natureza_juridica" db:"natureza_juridica"`                     // "natureza_juridica","text" [cite: 1]
	QualificacaoResponsavel   string  `json:"qualificacao_responsavel" db:"qualificacao_responsavel"`       // "qualificacao_responsavel","text" [cite: 1]
	PorteEmpresa              string  `json:"porte_empresa" db:"porte_empresa"`                             // "porte_empresa","text" [cite: 1]
	EnteFederativoResponsavel string  `json:"ente_federativo_responsavel" db:"ente_federativo_responsavel"` // "ente_federativo_responsavel","text" [cite: 1]
	CapitalSocial             float64 `json:"capital_social" db:"capital_social"`                           // "capital_social","real" [cite: 1]
}
