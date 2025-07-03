// neurocloser/backend/cmd/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/edufilhocruz/neurocloser/backend/database"
	"github.com/edufilhocruz/neurocloser/backend/dataloaders"
	"github.com/edufilhocruz/neurocloser/backend/graphql"
	"github.com/edufilhocruz/neurocloser/backend/graphql/generated"
	"github.com/edufilhocruz/neurocloser/backend/repositories"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	// Inicializa a conexão com o banco de dados
	database.InitDB()
	defer database.CloseDB() // Garante que a conexão será fechada ao final do programa

	// Inicializa TODOS os repositórios necessários
	empresaRepo := repositories.NewEmpresaRepository(database.DB)
	estabelecimentoRepo := repositories.NewEstabelecimentoRepository(database.DB)
	socioRepo := repositories.NewSocioRepository(database.DB)
	cnaeRepo := repositories.NewCNAERepository(database.DB)

	// Cria uma nova instância de resolver e injeta os repositórios
	resolver := &graphql.Resolver{
		DB:                  database.DB,
		EmpresaRepo:         empresaRepo,
		EstabelecimentoRepo: estabelecimentoRepo,
		SocioRepo:           socioRepo,
		CNAERepo:            cnaeRepo,
	}

	// Configuração do Servidor GraphQL
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver, Directives: generated.DirectiveRoot{}}))

	// Aplica o middleware do Dataloader ao servidor GraphQL
	// O middleware deve vir ANTES do servidor GraphQL para que os loaders estejam no contexto.
	http.Handle("/query", dataloaders.DataloaderMiddleware(empresaRepo, socioRepo, cnaeRepo)(srv))

	// Rota para o Playground GraphQL (não precisa do dataloader middleware para o playground)
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Porta padrão se não estiver definida
	}

	fmt.Printf("Servidor backend NeuroCloser iniciado na porta :%s\n", port)
	fmt.Println("Você pode acessar o playground GraphQL em http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":"+port, nil)) // Inicia o servidor HTTP
}
