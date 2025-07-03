// neurocloser/backend/database/connection.go
package database

import (
	"fmt"
	"log"
	"os"
	"time" // Importa o pacote time

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Driver PostgreSQL
)

// DB é a instância global do banco de dados.
var DB *sqlx.DB

// InitDB inicializa a conexão com o banco de dados PostgreSQL com retries.
func InitDB() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL não configurada. Por favor, defina a variável de ambiente.")
	}

	maxRetries := 5                  // Número máximo de tentativas de conexão
	retryInterval := 5 * time.Second // Intervalo entre as tentativas

	for i := 0; i < maxRetries; i++ {
		var err error
		DB, err = sqlx.Connect("postgres", connStr)
		if err == nil {
			// Se a conexão for bem-sucedida, tente um ping
			err = DB.Ping()
			if err == nil {
				fmt.Println("Conexão com o banco de dados PostgreSQL estabelecida com sucesso!")
				return // Conexão bem-sucedida, sai da função
			}
		}

		log.Printf("Tentativa %d/%d: Erro ao conectar ou pingar o banco de dados: %v. Tentando novamente em %v...", i+1, maxRetries, err, retryInterval)
		time.Sleep(retryInterval) // Espera antes de tentar novamente
	}

	log.Fatalf("Falha total na conexão com o banco de dados após %d tentativas.", maxRetries)
}

// CloseDB fecha a conexão com o banco de dados.
func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Printf("Erro ao fechar a conexão com o banco de dados: %v", err)
		} else {
			fmt.Println("Conexão com o banco de dados PostgreSQL fechada.")
		}
	}
}
