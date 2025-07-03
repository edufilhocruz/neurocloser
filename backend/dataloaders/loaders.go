package dataloaders

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/edufilhocruz/neurocloser/backend/models"
	"github.com/edufilhocruz/neurocloser/backend/repositories"
	"github.com/graph-gophers/dataloader" // Importa a biblioteca Dataloader
)

// Loaders é uma struct que contém todos os Dataloaders da sua aplicação.
type Loaders struct {
	EmpresaByCNPJBasico *dataloader.Loader
	SociosByCNPJBasico  *dataloader.Loader
	CNAEByCodigo        *dataloader.Loader
}

// NewLoaders cria e inicializa todos os Dataloaders.
func NewLoaders(empresaRepo repositories.EmpresaRepository,
	socioRepo repositories.SocioRepository,
	cnaeRepo repositories.CNAERepository) *Loaders {

	// Configurações comuns para os Dataloaders
	loaderOptions := []dataloader.Option{
		dataloader.WithCache(dataloader.NewCache()),
		dataloader.WithBatchCapacity(100),
		dataloader.WithWait(1 * time.Millisecond), // Pequeno delay para permitir batching
	}

	// Dataloader para Empresas por CNPJ Básico
	empresaLoader := dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		cnpjBasicos := make([]string, len(keys))
		for i, key := range keys {
			cnpjBasicos[i] = key.String() // .String() para converter dataloader.Key para string
		}

		empresas, err := empresaRepo.GetEmpresasByCNPJBasicos(cnpjBasicos)
		if err != nil {
			return errorResults(err, len(keys))
		}

		empresaMap := make(map[string]*models.Empresa)
		for _, emp := range empresas {
			empresaMap[emp.CNPJBasico] = emp
		}

		results := make([]*dataloader.Result, len(keys))
		for i, key := range keys {
			if emp, ok := empresaMap[key.String()]; ok {
				results[i] = &dataloader.Result{Data: emp}
			} else {
				results[i] = &dataloader.Result{Error: fmt.Errorf("empresa com CNPJ Básico %s não encontrada", key.String())}
			}
		}
		return results
	}, loaderOptions...) // Aplica as opções

	// Dataloader para Sócios por CNPJ Básico
	// Retorna map[string][]*models.Socio para o Dataloader
	socioLoader := dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		cnpjBasicos := make([]string, len(keys))
		for i, key := range keys {
			cnpjBasicos[i] = key.String()
		}

		sociosMap, err := socioRepo.GetMultiplesSociosByCNPJBasicos(cnpjBasicos)
		if err != nil {
			return errorResults(err, len(keys))
		}

		results := make([]*dataloader.Result, len(keys))
		for i, key := range keys {
			// Garante que sempre retorna uma slice, mesmo que vazia, se o CNPJ não tiver sócios
			if s, ok := sociosMap[key.String()]; ok {
				results[i] = &dataloader.Result{Data: s}
			} else {
				results[i] = &dataloader.Result{Data: []*models.Socio{}}
			}
		}
		return results
	}, loaderOptions...)

	// Dataloader para CNAEs por Código
	cnaeLoader := dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		codigos := make([]string, len(keys))
		for i, key := range keys {
			codigos[i] = key.String()
		}

		cnaes, err := cnaeRepo.GetCNAEsByCodigos(codigos)
		if err != nil {
			return errorResults(err, len(keys))
		}

		cnaeMap := make(map[string]*models.CNAE)
		for _, c := range cnaes {
			cnaeMap[c.Codigo] = c
		}

		results := make([]*dataloader.Result, len(keys))
		for i, key := range keys {
			if c, ok := cnaeMap[key.String()]; ok {
				results[i] = &dataloader.Result{Data: c}
			} else {
				results[i] = &dataloader.Result{Error: fmt.Errorf("CNAE com código %s não encontrado", key.String())}
			}
		}
		return results
	}, loaderOptions...)

	return &Loaders{
		EmpresaByCNPJBasico: empresaLoader,
		SociosByCNPJBasico:  socioLoader,
		CNAEByCodigo:        cnaeLoader,
	}
}

// errorResults é uma função auxiliar para retornar erros em um formato compatível com dataloader.
func errorResults(err error, count int) []*dataloader.Result {
	results := make([]*dataloader.Result, count)
	for i := 0; i < count; i++ {
		results[i] = &dataloader.Result{Error: err}
	}
	return results
}

// DataloaderMiddleware é um middleware HTTP que injeta os dataloaders no contexto da requisição.
func DataloaderMiddleware(empresaRepo repositories.EmpresaRepository,
	socioRepo repositories.SocioRepository,
	cnaeRepo repositories.CNAERepository) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loaders := NewLoaders(empresaRepo, socioRepo, cnaeRepo)
			ctx := context.WithValue(r.Context(), loadersKey, loaders)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// contextKey é um tipo para evitar colisões de chave de contexto.
type contextKey string

const loadersKey contextKey = "dataloaders"

// ForContext é uma função auxiliar para obter os Dataloaders do contexto.
func ForContext(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
