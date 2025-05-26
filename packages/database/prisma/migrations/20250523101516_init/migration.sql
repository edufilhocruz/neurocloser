-- CreateTable
CREATE TABLE "_referencia" (
    "referencia" TEXT NOT NULL,
    "valor" TEXT NOT NULL,

    CONSTRAINT "_referencia_pkey" PRIMARY KEY ("referencia")
);

-- CreateTable
CREATE TABLE "cnae" (
    "codigo" TEXT NOT NULL,
    "descricao" TEXT NOT NULL,

    CONSTRAINT "cnae_pkey" PRIMARY KEY ("codigo")
);

-- CreateTable
CREATE TABLE "empresas" (
    "cnpj_basico" TEXT NOT NULL,
    "razao_social" TEXT NOT NULL,
    "natureza_juridica" TEXT,
    "qualificacao" TEXT,
    "capital_social" TEXT,
    "porte" TEXT,
    "ente_federativo" TEXT,
    "telefone_fixo" TEXT,
    "celular" TEXT,
    "email" TEXT,
    "site" TEXT,
    "ativa" BOOLEAN NOT NULL DEFAULT true,
    "enviado_crm" BOOLEAN NOT NULL DEFAULT false,
    "enviado_email" BOOLEAN NOT NULL DEFAULT false,
    "ultima_atualizacao" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "empresas_pkey" PRIMARY KEY ("cnpj_basico")
);

-- CreateTable
CREATE TABLE "estabelecimento" (
    "cnpj_basico" TEXT NOT NULL,
    "cnpj_ordem" TEXT NOT NULL,
    "cnpj_dv" TEXT NOT NULL,
    "matriz_filial" TEXT,
    "nome_fantasia" TEXT,
    "situacao_cadastral" TEXT,
    "data_situacao" TEXT,
    "motivo_situacao" TEXT,
    "cidade_exterior" TEXT,
    "pais" TEXT,
    "data_inicio" TEXT,
    "cnae_principal" TEXT,
    "cnae_secundarios" TEXT,
    "tipo_logradouro" TEXT,
    "logradouro" TEXT,
    "numero" TEXT,
    "complemento" TEXT,
    "bairro" TEXT,
    "cep" TEXT,
    "uf" TEXT,
    "municipio" TEXT,
    "ddd1" TEXT,
    "telefone1" TEXT,
    "ddd2" TEXT,
    "telefone2" TEXT,
    "ddd_fax" TEXT,
    "fax" TEXT,
    "email" TEXT,
    "situacao_especial" TEXT,
    "data_situacao_especial" TEXT,

    CONSTRAINT "estabelecimento_pkey" PRIMARY KEY ("cnpj_basico","cnpj_ordem","cnpj_dv")
);

-- CreateTable
CREATE TABLE "motivo" (
    "codigo" TEXT NOT NULL,
    "descricao" TEXT NOT NULL,

    CONSTRAINT "motivo_pkey" PRIMARY KEY ("codigo")
);

-- CreateTable
CREATE TABLE "municipio" (
    "codigo" TEXT NOT NULL,
    "descricao" TEXT NOT NULL,

    CONSTRAINT "municipio_pkey" PRIMARY KEY ("codigo")
);

-- CreateTable
CREATE TABLE "natureza_juridica" (
    "codigo" TEXT NOT NULL,
    "descricao" TEXT NOT NULL,

    CONSTRAINT "natureza_juridica_pkey" PRIMARY KEY ("codigo")
);

-- CreateTable
CREATE TABLE "pais" (
    "codigo" TEXT NOT NULL,
    "descricao" TEXT NOT NULL,

    CONSTRAINT "pais_pkey" PRIMARY KEY ("codigo")
);

-- CreateTable
CREATE TABLE "qualificacao_socio" (
    "codigo" TEXT NOT NULL,
    "descricao" TEXT NOT NULL,

    CONSTRAINT "qualificacao_socio_pkey" PRIMARY KEY ("codigo")
);

-- CreateTable
CREATE TABLE "simples" (
    "cnpj_basico" TEXT NOT NULL,
    "opcao_simples" TEXT,
    "data_opcao_simples" TEXT,
    "data_exclusao_simples" TEXT,
    "opcao_mei" TEXT,
    "data_opcao_mei" TEXT,
    "data_exclusao_mei" TEXT,

    CONSTRAINT "simples_pkey" PRIMARY KEY ("cnpj_basico")
);

-- CreateTable
CREATE TABLE "socios" (
    "cnpj" TEXT NOT NULL,
    "cnpj_basico" TEXT NOT NULL,
    "identificador_socio" TEXT NOT NULL,
    "nome_socio" TEXT,
    "cnpj_cpf_socio" TEXT,
    "qualificacao_socio" TEXT,
    "data_entrada" TEXT,
    "pais" TEXT,
    "representante_legal" TEXT,
    "nome_representante" TEXT,
    "qualificacao_representante" TEXT,
    "faixa_etaria" TEXT,

    CONSTRAINT "socios_pkey" PRIMARY KEY ("cnpj","cnpj_basico","identificador_socio")
);

-- CreateTable
CREATE TABLE "Usuario" (
    "id" SERIAL NOT NULL,
    "nome" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "senha" TEXT NOT NULL,
    "role" TEXT NOT NULL DEFAULT 'user',
    "criado_em" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "atualizado_em" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Usuario_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Filtro" (
    "id" SERIAL NOT NULL,
    "nome" TEXT NOT NULL,
    "usuario_id" INTEGER NOT NULL,
    "condicoes" JSONB NOT NULL,
    "criado_em" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "atualizado_em" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Filtro_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "ListaExportacao" (
    "id" SERIAL NOT NULL,
    "nome" TEXT NOT NULL,
    "descricao" TEXT,
    "usuario_id" INTEGER NOT NULL,
    "filtro_id" INTEGER,
    "total_registros" INTEGER NOT NULL DEFAULT 0,
    "tipo" TEXT NOT NULL,
    "status" TEXT NOT NULL DEFAULT 'pendente',
    "arquivo_url" TEXT,
    "criado_em" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "atualizado_em" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "ListaExportacao_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "LogIntegracao" (
    "id" SERIAL NOT NULL,
    "tipo" TEXT NOT NULL,
    "status" TEXT NOT NULL,
    "mensagem" TEXT,
    "detalhes" JSONB,
    "empresa_id" TEXT,
    "criado_em" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "LogIntegracao_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "empresas_razao_social_idx" ON "empresas"("razao_social");

-- CreateIndex
CREATE INDEX "estabelecimento_cnpj_basico_idx" ON "estabelecimento"("cnpj_basico");

-- CreateIndex
CREATE INDEX "socios_cnpj_basico_idx" ON "socios"("cnpj_basico");

-- CreateIndex
CREATE INDEX "socios_cnpj_cpf_socio_idx" ON "socios"("cnpj_cpf_socio");

-- CreateIndex
CREATE INDEX "socios_nome_socio_idx" ON "socios"("nome_socio");

-- CreateIndex
CREATE UNIQUE INDEX "Usuario_email_key" ON "Usuario"("email");

-- AddForeignKey
ALTER TABLE "estabelecimento" ADD CONSTRAINT "estabelecimento_cnpj_basico_fkey" FOREIGN KEY ("cnpj_basico") REFERENCES "empresas"("cnpj_basico") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "simples" ADD CONSTRAINT "simples_cnpj_basico_fkey" FOREIGN KEY ("cnpj_basico") REFERENCES "empresas"("cnpj_basico") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "socios" ADD CONSTRAINT "socios_cnpj_basico_fkey" FOREIGN KEY ("cnpj_basico") REFERENCES "empresas"("cnpj_basico") ON DELETE RESTRICT ON UPDATE CASCADE;
