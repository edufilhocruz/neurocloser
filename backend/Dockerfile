FROM golang:1.24-alpine AS builder

# Define o diretório de trabalho dentro do contêiner.
WORKDIR /app

# Ativa o Go Modules e define o proxy para baixar dependências.
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct

# Copia os arquivos de gerenciamento de dependências.
COPY go.mod go.sum ./

# Baixa as dependências. O Docker irá cachear esta camada.
RUN go mod download

# Copia todo o código-fonte da aplicação.
COPY . .

# Compila a aplicação, criando um executável estático.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o neurocloser-backend ./cmd/main.go

# Estágio 2: Produção
# Usamos uma imagem mínima para a execução, por segurança e tamanho.
FROM alpine:latest
WORKDIR /root/

# Copia o executável compilado do estágio de build.
COPY --from=builder /app/neurocloser-backend .

# Expõe a porta que a aplicação usará.
EXPOSE 8080

# Comando para iniciar a aplicação.
CMD ["./neurocloser-backend"]