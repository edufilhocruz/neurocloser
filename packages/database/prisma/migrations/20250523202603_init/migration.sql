/*
  Warnings:

  - You are about to drop the `Estabelecimento` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `Socio` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropTable
DROP TABLE "Estabelecimento";

-- DropTable
DROP TABLE "Socio";

-- CreateTable
CREATE TABLE "estabelecimento" (
    "cnpj_basico" TEXT NOT NULL,
    "cnpj_ordem" TEXT NOT NULL,
    "cnpj_dv" TEXT NOT NULL,
    "matriz_filial" TEXT,
    "nome_fantasia" TEXT,
    "situacao_cadastral" TEXT,
    "data_situacao_cadastral" TEXT,
    "motivo_situacao_cadastral" TEXT,
    "nome_cidade_exterior" TEXT,
    "pais" TEXT,
    "data_inicio_atividades" TEXT,
    "cnae_fiscal" TEXT,
    "cnae_fiscal_secundaria" TEXT,
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
    "correio_eletronico" TEXT,
    "situacao_especial" TEXT,
    "data_situacao_especial" TEXT,
    "cnpj" TEXT,

    CONSTRAINT "estabelecimento_pkey" PRIMARY KEY ("cnpj_basico","cnpj_ordem","cnpj_dv")
);

-- CreateTable
CREATE TABLE "socios" (
    "cnpj" TEXT NOT NULL,
    "cnpj_basico" TEXT NOT NULL,
    "identificador_de_socio" TEXT NOT NULL,
    "nome_socio" TEXT,
    "cnpj_cpf_socio" TEXT,
    "qualificacao_socio" TEXT,
    "data_entrada_sociedade" TEXT,
    "pais" TEXT,
    "representante_legal" TEXT,
    "nome_representante" TEXT,
    "qualificacao_representante_legal" TEXT,
    "faixa_etaria" TEXT,

    CONSTRAINT "socios_pkey" PRIMARY KEY ("cnpj","cnpj_basico","identificador_de_socio")
);
