package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/edufilhocruz/neurocloser/backend/models"
	"github.com/jmoiron/sqlx"
)

// EstabelecimentoRepository define a interface para operações de dados do Estabelecimento.
type EstabelecimentoRepository interface {
	GetEstabelecimentoByID(id int) (*models.Estabelecimento, error)
	GetEstabelecimentoByCNPJBasico(cnpjBasico string) (*models.Estabelecimento, error)
	FindEstabelecimentosByFilters(filters map[string]interface{}, limit *int, offset *int) ([]*models.Estabelecimento, error)
}

// estabelecimentoRepository implementa EstabelecimentoRepository para PostgreSQL.
type estabelecimentoRepository struct {
	db *sqlx.DB
}

// NewEstabelecimentoRepository cria uma nova instância de EstabelecimentoRepository.
func NewEstabelecimentoRepository(db *sqlx.DB) EstabelecimentoRepository {
	return &estabelecimentoRepository{db: db}
}

// GetEstabelecimentoByID busca um estabelecimento pelo seu ID.
func (r *estabelecimentoRepository) GetEstabelecimentoByID(id int) (*models.Estabelecimento, error) {
	var estabelecimento models.Estabelecimento
	query := `SELECT * FROM estabelecimento WHERE id = $1`
	err := r.db.Get(&estabelecimento, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar estabelecimento por ID %d: %w", id, err)
	}
	return &estabelecimento, nil
}

// GetEstabelecimentoByCNPJBasico busca um estabelecimento pelo seu CNPJ Básico.
func (r *estabelecimentoRepository) GetEstabelecimentoByCNPJBasico(cnpjBasico string) (*models.Estabelecimento, error) {
	var estabelecimento models.Estabelecimento
	query := `SELECT * FROM estabelecimento WHERE cnpj_basico = $1 ORDER BY matriz_filial DESC LIMIT 1`
	err := r.db.Get(&estabelecimento, query, cnpjBasico)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar estabelecimento por CNPJ Básico %s: %w", cnpjBasico, err)
	}
	return &estabelecimento, nil
}

// FindEstabelecimentosByFilters busca estabelecimentos com base em múltiplos critérios de filtro.
func (r *estabelecimentoRepository) FindEstabelecimentosByFilters(filters map[string]interface{}, limit *int, offset *int) ([]*models.Estabelecimento, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT est.* FROM estabelecimento est")

	args := []interface{}{}
	whereClauses := []string{}
	needsJoin := false

	// Filtros de Estabelecimento
	if cnpj, ok := filters["cnpj"].(string); ok && cnpj != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("est.cnpj = $%d", len(args)+1))
		args = append(args, cnpj)
	}
	if nomeFantasia, ok := filters["nomeFantasia"].(string); ok && nomeFantasia != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("est.nome_fantasia ILIKE $%d", len(args)+1))
		args = append(args, "%"+nomeFantasia+"%")
	}
	if uf, ok := filters["uf"].(string); ok && uf != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("est.uf = $%d", len(args)+1))
		args = append(args, uf)
	}
	if situacaoCadastral, ok := filters["situacaoCadastral"].(string); ok && situacaoCadastral != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("est.situacao_cadastral = $%d", len(args)+1))
		args = append(args, situacaoCadastral)
	}
	if cnaeFiscal, ok := filters["cnaeFiscal"].(string); ok && cnaeFiscal != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("est.cnae_fiscal = $%d", len(args)+1))
		args = append(args, cnaeFiscal)
	}
	if cnaeFiscalSecundaria, ok := filters["cnaeFiscalSecundaria"].(string); ok && cnaeFiscalSecundaria != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("est.cnae_fiscal_secundaria LIKE $%d", len(args)+1))
		args = append(args, "%"+cnaeFiscalSecundaria+"%")
	}

	// Filtros de Empresa (que exigem JOIN)
	if razaoSocial, ok := filters["razaoSocial"].(string); ok && razaoSocial != "" {
		needsJoin = true
		whereClauses = append(whereClauses, fmt.Sprintf("emp.razao_social ILIKE $%d", len(args)+1))
		args = append(args, "%"+razaoSocial+"%")
	}
	if porteEmpresa, ok := filters["porteEmpresa"].(string); ok && porteEmpresa != "" {
		needsJoin = true
		whereClauses = append(whereClauses, fmt.Sprintf("emp.porte_empresa = $%d", len(args)+1))
		args = append(args, porteEmpresa)
	}
	if minCapitalSocial, ok := filters["minCapitalSocial"].(float64); ok && minCapitalSocial >= 0 {
		needsJoin = true
		whereClauses = append(whereClauses, fmt.Sprintf("emp.capital_social >= $%d", len(args)+1))
		args = append(args, minCapitalSocial)
	}
	if maxCapitalSocial, ok := filters["maxCapitalSocial"].(float64); ok && maxCapitalSocial >= 0 {
		needsJoin = true
		whereClauses = append(whereClauses, fmt.Sprintf("emp.capital_social <= $%d", len(args)+1))
		args = append(args, maxCapitalSocial)
	}

	if needsJoin {
		queryBuilder.WriteString(" JOIN empresas emp ON est.cnpj_basico = emp.cnpj_basico")
	}

	if len(whereClauses) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(whereClauses, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY est.id ASC")

	if limit != nil && *limit > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d", len(args)+1))
		args = append(args, *limit)
	}
	if offset != nil && *offset >= 0 {
		queryBuilder.WriteString(fmt.Sprintf(" OFFSET $%d", len(args)+1))
		args = append(args, *offset)
	}

	var estabelecimentos []*models.Estabelecimento
	err := r.db.Select(&estabelecimentos, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar estabelecimentos com filtros: %w", err)
	}

	return estabelecimentos, nil
}
