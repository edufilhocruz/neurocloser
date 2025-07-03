package models

type ProspeccaoDetalhada struct {
	Empresa         *Empresa         `json:"empresa"`
	Estabelecimento *Estabelecimento `json:"estabelecimento"`
	Socios          []*Socio         `json:"socios"`
	CNAEFiscal      *CNAE            `json:"cnaeFiscal"`     // CNAE Fiscal Principal
	CNAESecundaria  []*CNAE          `json:"cnaeSecundaria"` // CNAEs Secundários
}
