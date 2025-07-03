package models

// Socio representa a tabela 'socios' no banco de dados.
type Socio struct {
	CNPJ                           string `json:"cnpj" db:"cnpj"`                                                         // "cnpj","text" [cite: 1]
	CNPJBasico                     string `json:"cnpj_basico" db:"cnpj_basico"`                                           // "cnpj_basico","text" [cite: 1]
	IdentificadorDeSocio           string `json:"identificador_de_socio" db:"identificador_de_socio"`                     // "identificador_de_socio","text" [cite: 1]
	NomeSocio                      string `json:"nome_socio" db:"nome_socio"`                                             // "nome_socio","text" [cite: 1]
	CNPJCPFSocio                   string `json:"cnpj_cpf_socio" db:"cnpj_cpf_socio"`                                     // "cnpj_cpf_socio","text" [cite: 1]
	QualificacaoSocio              string `json:"qualificacao_socio" db:"qualificacao_socio"`                             // "qualificacao_socio","text" [cite: 1]
	DataEntradaSociedade           string `json:"data_entrada_sociedade" db:"data_entrada_sociedade"`                     // "data_entrada_sociedade","text" [cite: 1]
	Pais                           string `json:"pais" db:"pais"`                                                         // "pais","text" [cite: 1]
	RepresentanteLegal             string `json:"representante_legal" db:"representante_legal"`                           // "representante_legal","text" [cite: 1]
	NomeRepresentante              string `json:"nome_representante" db:"nome_representante"`                             // "nome_representante","text" [cite: 1]
	QualificacaoRepresentanteLegal string `json:"qualificacao_representante_legal" db:"qualificacao_representante_legal"` // "qualificacao_representante_legal","text" [cite: 1]
	FaixaEtaria                    string `json:"faixa_etaria" db:"faixa_etaria"`                                         // "faixa_etaria","text" [cite: 1]
}
