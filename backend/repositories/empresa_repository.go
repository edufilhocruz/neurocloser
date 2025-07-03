package repositories

import (
	"database/sql"
	"fmt"

	"github.com/edufilhocruz/neurocloser/backend/models"

	"github.com/jmoiron/sqlx"
)

// EmpresaRepository define a interface para operações de dados da Empresa.
type EmpresaRepository interface {
	GetAllEmpresas(limit *int) ([]*models.Empresa, error)
	GetEmpresaByCNPJBasico(cnpjBasico string) (*models.Empresa, error)
	// NOVO MÉTODO PARA DATALOADER:
	GetEmpresasByCNPJBasicos(cnpjBasicos []string) ([]*models.Empresa, error)
}

// empresaRepository implementa EmpresaRepository para PostgreSQL.
type empresaRepository struct {
	db *sqlx.DB
}

// NewEmpresaRepository cria uma nova instância de EmpresaRepository.
func NewEmpresaRepository(db *sqlx.DB) EmpresaRepository {
	return &empresaRepository{db: db}
}

// GetAllEmpresas busca todas as empresas do banco de dados, com um limite opcional.
func (r *empresaRepository) GetAllEmpresas(limit *int) ([]*models.Empresa, error) {
	empresas := []*models.Empresa{}
	query := `SELECT cnpj_basico, razao_social, natureza_juridica, qualificacao_responsavel, porte_empresa, ente_federativo_responsavel, capital_social FROM empresas`
	if limit != nil && *limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", *limit)
	}

	err := r.db.Select(&empresas, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar empresas: %w", err)
	}
	return empresas, nil
}

// GetEmpresaByCNPJBasico busca uma empresa pelo seu CNPJ Básico.
func (r *empresaRepository) GetEmpresaByCNPJBasico(cnpjBasico string) (*models.Empresa, error) {
	var empresa models.Empresa
	query := `
		SELECT
			cnpj_basico, razao_social, natureza_juridica, qualificacao_responsavel,
			porte_empresa, ente_federativo_responsavel, capital_social
		FROM empresas WHERE cnpj_basico = $1`

	err := r.db.Get(&empresa, query, cnpjBasico)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Retorna nil, nil se nenhuma empresa for encontrada
		}
		return nil, fmt.Errorf("erro ao buscar empresa por CNPJ Básico: %w", err)
	}
	return &empresa, nil
}

// GetEmpresasByCNPJBasicos busca múltiplas empresas por seus CNPJs básicos em uma única consulta.
func (r *empresaRepository) GetEmpresasByCNPJBasicos(cnpjBasicos []string) ([]*models.Empresa, error) {
	if len(cnpjBasicos) == 0 {
		return []*models.Empresa{}, nil
	}

	var empresas []*models.Empresa
	query := `SELECT cnpj_basico, razao_social, natureza_juridica, qualificacao_responsavel, porte_empresa, ente_federativo_responsavel, capital_social FROM empresas WHERE cnpj_basico IN (?)`
	query, args, err := sqlx.In(query, cnpjBasicos)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar query IN para empresas: %w", err)
	}
	query = r.db.Rebind(query) // Rebind para o formato de placeholder do PostgreSQL ($1, $2, etc.)

	err = r.db.Select(&empresas, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar empresas por CNPJs básicos: %w", err)
	}
	return empresas, nil
}
