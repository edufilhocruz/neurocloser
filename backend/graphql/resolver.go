package graphql

import (
	"github.com/edufilhocruz/neurocloser/backend/repositories"

	"github.com/jmoiron/sqlx"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB                  *sqlx.DB
	EmpresaRepo         repositories.EmpresaRepository
	EstabelecimentoRepo repositories.EstabelecimentoRepository
	SocioRepo           repositories.SocioRepository
	CNAERepo            repositories.CNAERepository
}
