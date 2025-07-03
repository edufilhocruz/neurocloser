package repositories

import (
	"database/sql"
	"fmt"

	"github.com/edufilhocruz/neurocloser/backend/models"

	"github.com/jmoiron/sqlx"
)

// CNAERepository define a interface para operações de dados de CNAE.
type CNAERepository interface {
	GetCNAEByCodigo(codigo string) (*models.CNAE, error)
	GetCNAEsByCodigos(codigos []string) ([]*models.CNAE, error)
}

// cnaeRepository implementa CNAERepository para PostgreSQL.
type cnaeRepository struct {
	db *sqlx.DB
}

// NewCNAERepository cria uma nova instância de CNAERepository.
func NewCNAERepository(db *sqlx.DB) CNAERepository {
	return &cnaeRepository{db: db}
}

// GetCNAEByCodigo busca um CNAE pela sua código.
func (r *cnaeRepository) GetCNAEByCodigo(codigo string) (*models.CNAE, error) {
	var cnae models.CNAE
	query := "SELECT codigo, descricao FROM cnae WHERE codigo = $1"
	err := r.db.Get(&cnae, query, codigo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // CNAE não encontrado, não é um erro para o chamador.
		}
		return nil, fmt.Errorf("erro ao buscar CNAE por código '%s': %w", codigo, err)
	}
	return &cnae, nil
}

// GetCNAEsByCodigos busca múltiplos CNAEs por uma lista de códigos.
// Utiliza sqlx.In para construir a cláusula IN de forma segura.
func (r *cnaeRepository) GetCNAEsByCodigos(codigos []string) ([]*models.CNAE, error) {
	if len(codigos) == 0 {
		return []*models.CNAE{}, nil
	}

	var cnaes []*models.CNAE
	// Constrói a cláusula IN dinamicamente para a query SQL
	query := "SELECT codigo, descricao FROM cnae WHERE codigo IN (?)"
	query, args, err := sqlx.In(query, codigos)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar query IN para CNAEs: %w", err)
	}
	query = r.db.Rebind(query) // Rebind para o formato de placeholder do PostgreSQL ($1, $2, etc.)

	err = r.db.Select(&cnaes, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar CNAEs por códigos: %w", err)
	}
	return cnaes, nil
}
