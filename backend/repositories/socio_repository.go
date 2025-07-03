// neurocloser/backend/repositories/socio_repository.go
package repositories

import (
	"fmt"

	"github.com/edufilhocruz/neurocloser/backend/models"

	"github.com/jmoiron/sqlx"
)

// SocioRepository define a interface para operações de dados do Sócio.
type SocioRepository interface {
	GetSociosByCNPJBasico(cnpjBasico string) ([]*models.Socio, error)
	// NOVO MÉTODO PARA DATALOADER: Retorna um mapa para facilitar o mapeamento no Dataloader
	GetMultiplesSociosByCNPJBasicos(cnpjBasicos []string) (map[string][]*models.Socio, error)
}

// socioRepository implementa SocioRepository para PostgreSQL.
type socioRepository struct {
	db *sqlx.DB
}

// NewSocioRepository cria uma nova instância de SocioRepository.
func NewSocioRepository(db *sqlx.DB) SocioRepository {
	return &socioRepository{db: db}
}

// GetSociosByCNPJBasico busca os sócios de uma empresa pelo CNPJ Básico.
func (r *socioRepository) GetSociosByCNPJBasico(cnpjBasico string) ([]*models.Socio, error) {
	var socios []*models.Socio
	query := ` 
		SELECT
			cnpj, cnpj_basico, identificador_de_socio, nome_socio, cnpj_cpf_socio,
			qualificacao_socio, data_entrada_sociedade, pais, representante_legal,
			nome_representante, qualificacao_representante_legal, faixa_etaria
		FROM socios
		WHERE cnpj_basico = $1
	`
	err := r.db.Select(&socios, query, cnpjBasico)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar sócios por CNPJ Básico %s: %w", cnpjBasico, err)
	}
	return socios, nil
}

// GetMultiplesSociosByCNPJBasicos busca múltiplos sócios para múltiplos CNPJs básicos em uma única consulta.
func (r *socioRepository) GetMultiplesSociosByCNPJBasicos(cnpjBasicos []string) (map[string][]*models.Socio, error) {
	if len(cnpjBasicos) == 0 {
		return map[string][]*models.Socio{}, nil
	}

	query := `
		SELECT
			cnpj, cnpj_basico, identificador_de_socio, nome_socio, cnpj_cpf_socio,
			qualificacao_socio, data_entrada_sociedade, pais, representante_legal,
			nome_representante, qualificacao_representante_legal, faixa_etaria
		FROM socios
		WHERE cnpj_basico IN (?)
	`
	query, args, err := sqlx.In(query, cnpjBasicos)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar query IN para sócios: %w", err)
	}
	query = r.db.Rebind(query) // Rebind para o formato de placeholder do PostgreSQL ($1, $2, etc.)

	var allSocios []*models.Socio
	err = r.db.Select(&allSocios, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar múltiplos sócios: %w", err)
	}

	sociosMap := make(map[string][]*models.Socio)
	for _, socio := range allSocios {
		// The loop variable 'socio' is already a pointer of type *models.Socio.
		// We append it directly to the map's slice.
		sociosMap[socio.CNPJBasico] = append(sociosMap[socio.CNPJBasico], socio)
	}
	return sociosMap, nil
}
