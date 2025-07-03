// neurocloser/backend/graphql/schema.resolvers.go
package graphql

import (
	"context"
	"fmt"
	"log" // Para a conversão do CapitalSocial
	"strings"

	"github.com/edufilhocruz/neurocloser/backend/dataloaders"
	"github.com/edufilhocruz/neurocloser/backend/models"

	"github.com/edufilhocruz/neurocloser/backend/graphql/generated"
	"github.com/edufilhocruz/neurocloser/backend/graphql/model"

	"github.com/graph-gophers/dataloader"
)

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// Empresas é o resolver para a query 'empresas'.
func (r *queryResolver) Empresas(ctx context.Context, limit *int, offset *int) ([]*models.Empresa, error) {
	empresas, err := r.EmpresaRepo.GetAllEmpresas(limit)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar todas as empresas: %w", err)
	}
	return empresas, nil
}

// Empresa é o resolver para a query 'empresa' (por CNPJ Básico).
func (r *queryResolver) Empresa(ctx context.Context, cnpjBasico string) (*models.Empresa, error) {
	empresa, err := r.EmpresaRepo.GetEmpresaByCNPJBasico(cnpjBasico)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar empresa por CNPJ Básico %s: %w", cnpjBasico, err)
	}
	return empresa, nil
}

// Estabelecimento é o resolver para a query 'estabelecimento' (por ID).
func (r *queryResolver) Estabelecimento(ctx context.Context, id int) (*models.Estabelecimento, error) {
	estabelecimento, err := r.EstabelecimentoRepo.GetEstabelecimentoByID(id)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar estabelecimento por ID %d: %w", id, err)
	}
	return estabelecimento, nil
}

// SociosByCnpjBasico para buscar sócios diretamente.
func (r *queryResolver) SociosByCnpjBasico(ctx context.Context, cnpjBasico string) ([]*models.Socio, error) {
	// AQUI TAMBÉM PODEMOS USAR O DATALOADER se esta query for chamada múltiplas vezes para diferentes CNPJs Básicos
	loaders := dataloaders.ForContext(ctx)
	if loaders == nil {
		return nil, fmt.Errorf("dataloaders não disponíveis no contexto")
	}

	// FIX: Chame o thunk retornado por Load() para obter o resultado e o erro.
	sociosResult, err := loaders.SociosByCNPJBasico.Load(ctx, dataloader.StringKey(cnpjBasico))()
	if err != nil {
		return nil, fmt.Errorf("falha no dataloader ao buscar sócios por CNPJ Básico %s: %w", cnpjBasico, err)
	}

	// O resultado do dataloader é interface{}, precisamos fazer um type assertion
	if s, ok := sociosResult.([]*models.Socio); ok {
		return s, nil
	}
	return nil, fmt.Errorf("tipo inesperado de resultado do dataloader para sócios: %T", sociosResult)
}

// CnaeByCodigo para buscar um CNAE diretamente.
func (r *queryResolver) CnaeByCodigo(ctx context.Context, codigo string) (*models.CNAE, error) {
	loaders := dataloaders.ForContext(ctx)
	if loaders == nil {
		return nil, fmt.Errorf("dataloaders não disponíveis no contexto")
	}

	// FIX: Chame o thunk retornado por Load() para obter o resultado e o erro.
	cnaeResult, err := loaders.CNAEByCodigo.Load(ctx, dataloader.StringKey(codigo))()
	if err != nil {
		return nil, fmt.Errorf("falha no dataloader ao buscar CNAE por código %s: %w", codigo, err)
	}

	if c, ok := cnaeResult.(*models.CNAE); ok {
		return c, nil
	}
	return nil, fmt.Errorf("tipo inesperado de resultado do dataloader para CNAE: %T", cnaeResult)
}

// BuscarProspeccao é o resolver para a query 'buscarProspeccao' com filtros e paginação.
func (r *queryResolver) BuscarProspeccao(ctx context.Context, filter *model.ProspeccaoFilter, limit *int, offset *int) ([]*models.ProspeccaoDetalhada, error) {
	repoFilters := make(map[string]interface{})

	if filter != nil {
		if filter.Cnpj != nil {
			repoFilters["cnpj"] = *filter.Cnpj
		}
		if filter.RazaoSocial != nil {
			repoFilters["razaoSocial"] = *filter.RazaoSocial
		}
		if filter.NomeFantasia != nil {
			repoFilters["nomeFantasia"] = *filter.NomeFantasia
		}
		if filter.Uf != nil {
			repoFilters["uf"] = *filter.Uf
		}
		if filter.SituacaoCadastral != nil {
			repoFilters["situacaoCadastral"] = *filter.SituacaoCadastral
		}
		if filter.PorteEmpresa != nil {
			repoFilters["porteEmpresa"] = *filter.PorteEmpresa
		}
		if filter.CnaeFiscal != nil {
			repoFilters["cnaeFiscal"] = *filter.CnaeFiscal
		}
		if filter.CnaeFiscalSecundaria != nil {
			repoFilters["cnaeFiscalSecundaria"] = *filter.CnaeFiscalSecundaria
		}
		if filter.MinCapitalSocial != nil {
			repoFilters["minCapitalSocial"] = *filter.MinCapitalSocial
		}
		if filter.MaxCapitalSocial != nil {
			repoFilters["maxCapitalSocial"] = *filter.MaxCapitalSocial
		}
	}

	// 1. Buscar a lista primária de estabelecimentos. O repositório agora retorna o tipo simples.
	estabelecimentos, err := r.EstabelecimentoRepo.FindEstabelecimentosByFilters(repoFilters, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estabelecimentos com filtros: %w", err)
	}
	if len(estabelecimentos) == 0 {
		return []*models.ProspeccaoDetalhada{}, nil
	}

	// Obtém os dataloaders do contexto da requisição
	loaders := dataloaders.ForContext(ctx)
	if loaders == nil {
		return nil, fmt.Errorf("dataloaders não disponíveis no contexto")
	}

	var resultados []*models.ProspeccaoDetalhada
	for _, estabelecimento := range estabelecimentos {
		cnpjBasico := estabelecimento.CNPJBasico
		if len(cnpjBasico) < 8 {
			log.Printf("CNPJ Básico inválido para estabelecimento ID %d: %s", estabelecimento.ID, cnpjBasico)
			continue
		}

		// 2. Usar Dataloaders para buscar dados relacionados de forma eficiente.
		// FIX: Chamar o thunk com () e tratar o erro.
		empresaResult, err := loaders.EmpresaByCNPJBasico.Load(ctx, dataloader.StringKey(cnpjBasico))()
		if err != nil {
			log.Printf("Erro no dataloader ao buscar empresa para CNPJ Básico %s: %v", cnpjBasico, err)
			continue // Pula esta prospecção se a empresa não puder ser carregada
		}
		empresa, _ := empresaResult.(*models.Empresa)

		sociosResult, err := loaders.SociosByCNPJBasico.Load(ctx, dataloader.StringKey(cnpjBasico))()
		var socios []*models.Socio
		if err != nil {
			log.Printf("Erro no dataloader ao buscar sócios para CNPJ Básico %s: %v", cnpjBasico, err)
		} else {
			socios, _ = sociosResult.([]*models.Socio)
		}

		var cnaeFiscal *models.CNAE
		if estabelecimento.CNAEFiscal != "" {
			cnaeResult, err := loaders.CNAEByCodigo.Load(ctx, dataloader.StringKey(estabelecimento.CNAEFiscal))()
			if err != nil {
				log.Printf("Erro no dataloader ao buscar CNAE Fiscal principal %s: %v", estabelecimento.CNAEFiscal, err)
			} else {
				cnaeFiscal, _ = cnaeResult.(*models.CNAE)
			}
		}

		var cnaeSecundaria []*models.CNAE
		if estabelecimento.CNAEFiscalSecundaria != "" {
			cnaeSecundariosCodigos := strings.Split(estabelecimento.CNAEFiscalSecundaria, ",")
			var cleanKeys []dataloader.Key
			for _, code := range cnaeSecundariosCodigos {
				trimmedCode := strings.TrimSpace(code)
				if trimmedCode != "" {
					cleanKeys = append(cleanKeys, dataloader.StringKey(trimmedCode))
				}
			}

			if len(cleanKeys) > 0 {
				// FIX: Usar LoadMany em vez de LoadAll
				cnaesResults, errs := loaders.CNAEByCodigo.LoadMany(ctx, cleanKeys)()
				if len(errs) > 0 {
					log.Printf("Erros no dataloader ao buscar CNAEs Secundários: %v", errs)
				}
				for _, res := range cnaesResults {
					if c, ok := res.(*models.CNAE); ok && c != nil {
						cnaeSecundaria = append(cnaeSecundaria, c)
					}
				}
			}
		}

		resultados = append(resultados, &models.ProspeccaoDetalhada{
			Empresa:         empresa,
			Estabelecimento: estabelecimento,
			Socios:          socios,
			CNAEFiscal:      cnaeFiscal,
			CNAESecundaria:  cnaeSecundaria,
		})
	}

	return resultados, nil
}

// O bloco de Mutation (func (r *Resolver) Mutation() ...) foi removido temporariamente.
