package models

import "fmt" // Importar fmt para a função de formatação de CNPJ
type Estabelecimento struct {
	ID                      int    `json:"id" db:"id"`
	CNPJ                    string `json:"cnpj" db:"cnpj"`       // CNPJ bruto (sem formatação)
	CNPJFormatado           string `json:"cnpjFormatado" db:"-"` // NOVO: CNPJ formatado, não vindo diretamente do DB
	CNPJBasico              string `json:"cnpj_basico" db:"cnpj_basico"`
	CNPJOrdem               string `json:"cnpj_ordem" db:"cnpj_ordem"`
	CNPJDV                  string `json:"cnpj_dv" db:"cnpj_dv"`
	MatrizFilial            string `json:"matriz_filial" db:"matriz_filial"`
	NomeFantasia            string `json:"nome_fantasia" db:"nome_fantasia"`
	SituacaoCadastral       string `json:"situacao_cadastral" db:"situacao_cadastral"`
	DataSituacaoCadastral   string `json:"data_situacao_cadastral" db:"data_situacao_cadastral"`
	MotivoSituacaoCadastral string `json:"motivo_situacao_cadastral" db:"motivo_situacao_cadastral"`
	NomeCidadeExterior      string `json:"nome_cidade_exterior" db:"nome_cidade_exterior"`
	Pais                    string `json:"pais" db:"pais"`
	DataInicioAtividades    string `json:"data_inicio_atividades" db:"data_inicio_atividades"`
	CNAEFiscal              string `json:"cnae_fiscal" db:"cnae_fiscal"`
	CNAEFiscalSecundaria    string `json:"cnae_fiscal_secundaria" db:"cnae_fiscal_secundaria"`
	TipoLogradouro          string `json:"tipo_logradouro" db:"tipo_logradouro"`
	Logradouro              string `json:"logradouro" db:"logradouro"`
	Numero                  string `json:"numero" db:"numero"`
	Complemento             string `json:"complemento" db:"complemento"`
	Bairro                  string `json:"bairro" db:"bairro"`
	CEP                     string `json:"cep" db:"cep"`
	UF                      string `json:"uf" db:"uf"`
	Municipio               string `json:"municipio" db:"municipio"`
	DDD1                    string `json:"ddd1" db:"ddd1"`
	Telefone1               string `json:"telefone1" db:"telefone1"`
	DDD2                    string `json:"ddd2" db:"ddd2"`
	Telefone2               string `json:"telefone2" db:"telefone2"`
	DDDFax                  string `json:"ddd_fax" db:"ddd_fax"`
	Fax                     string `json:"fax" db:"fax"`
	CorreioEletronico       string `json:"correio_eletronico" db:"correio_eletronico"`
	SituacaoEspecial        string `json:"situacao_especial" db:"situacao_especial"`
	DataSituacaoEspecial    string `json:"data_situacao_especial" db:"data_situacao_especial"`
}

// FormatCNPJ formata o CNPJ do estabelecimento no formato "XX.XXX.XXX/XXXX-XX".
func (e *Estabelecimento) FormatCNPJ() {
	cnpj := e.CNPJ
	if len(cnpj) != 14 {
		e.CNPJFormatado = cnpj // Se não tiver 14 dígitos, usa o CNPJ bruto
		return
	}
	e.CNPJFormatado = fmt.Sprintf("%s.%s.%s/%s-%s", cnpj[0:2], cnpj[2:5], cnpj[5:8], cnpj[8:12], cnpj[12:14])
}
