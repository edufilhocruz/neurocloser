/*
  Warnings:

  - You are about to drop the column `ativa` on the `empresas` table. All the data in the column will be lost.
  - You are about to drop the column `celular` on the `empresas` table. All the data in the column will be lost.
  - You are about to drop the column `email` on the `empresas` table. All the data in the column will be lost.
  - You are about to drop the column `ente_federativo` on the `empresas` table. All the data in the column will be lost.
  - You are about to drop the column `enviado_crm` on the `empresas` table. All the data in the column will be lost.
  - You are about to drop the column `enviado_email` on the `empresas` table. All the data in the column will be lost.
  - You are about to drop the column `porte` on the `empresas` table. All the data in the column will be lost.
  - You are about to drop the column `qualificacao` on the `empresas` table. All the data in the column will be lost.
  - You are about to drop the column `site` on the `empresas` table. All the data in the column will be lost.
  - You are about to drop the column `telefone_fixo` on the `empresas` table. All the data in the column will be lost.
  - You are about to drop the column `ultima_atualizacao` on the `empresas` table. All the data in the column will be lost.
  - The `capital_social` column on the `empresas` table would be dropped and recreated. This will lead to data loss if there is data in the column.
  - You are about to drop the `Filtro` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `ListaExportacao` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `LogIntegracao` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `Usuario` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `estabelecimento` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `socios` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "estabelecimento" DROP CONSTRAINT "estabelecimento_cnpj_basico_fkey";

-- DropForeignKey
ALTER TABLE "simples" DROP CONSTRAINT "simples_cnpj_basico_fkey";

-- DropForeignKey
ALTER TABLE "socios" DROP CONSTRAINT "socios_cnpj_basico_fkey";

-- AlterTable
ALTER TABLE "_referencia" ALTER COLUMN "valor" DROP NOT NULL;

-- AlterTable
ALTER TABLE "cnae" ALTER COLUMN "descricao" DROP NOT NULL;

-- AlterTable
ALTER TABLE "empresas" DROP COLUMN "ativa",
DROP COLUMN "celular",
DROP COLUMN "email",
DROP COLUMN "ente_federativo",
DROP COLUMN "enviado_crm",
DROP COLUMN "enviado_email",
DROP COLUMN "porte",
DROP COLUMN "qualificacao",
DROP COLUMN "site",
DROP COLUMN "telefone_fixo",
DROP COLUMN "ultima_atualizacao",
ADD COLUMN     "ente_federativo_responsavel" TEXT,
ALTER COLUMN "razao_social" DROP NOT NULL,
DROP COLUMN "capital_social",
ADD COLUMN     "capital_social" DOUBLE PRECISION;

-- AlterTable
ALTER TABLE "motivo" ALTER COLUMN "descricao" DROP NOT NULL;

-- AlterTable
ALTER TABLE "municipio" ALTER COLUMN "descricao" DROP NOT NULL;

-- AlterTable
ALTER TABLE "natureza_juridica" ALTER COLUMN "descricao" DROP NOT NULL;

-- AlterTable
ALTER TABLE "pais" ALTER COLUMN "descricao" DROP NOT NULL;

-- AlterTable
ALTER TABLE "qualificacao_socio" ALTER COLUMN "descricao" DROP NOT NULL;

-- DropTable
DROP TABLE "Filtro";

-- DropTable
DROP TABLE "ListaExportacao";

-- DropTable
DROP TABLE "LogIntegracao";

-- DropTable
DROP TABLE "Usuario";

-- DropTable
DROP TABLE "estabelecimento";

-- DropTable
DROP TABLE "socios";

-- CreateTable
CREATE TABLE "Estabelecimento" (
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

    CONSTRAINT "Estabelecimento_pkey" PRIMARY KEY ("cnpj_basico","cnpj_ordem","cnpj_dv")
);

-- CreateTable
CREATE TABLE "Socio" (
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

    CONSTRAINT "Socio_pkey" PRIMARY KEY ("cnpj","cnpj_basico","identificador_de_socio")
);

-- CreateIndex
CREATE INDEX "idx_cnae" ON "cnae"("codigo");

-- CreateIndex
CREATE INDEX "idx_empresas_cnpj_basico" ON "empresas"("cnpj_basico");

-- CreateIndex
CREATE INDEX "idx_motivo" ON "motivo"("codigo");

-- CreateIndex
CREATE INDEX "idx_municipio" ON "municipio"("codigo");

-- CreateIndex
CREATE INDEX "idx_natureza_juridica" ON "natureza_juridica"("codigo");

-- CreateIndex
CREATE INDEX "idx_pais" ON "pais"("codigo");

-- CreateIndex
CREATE INDEX "idx_qualificacao_socio" ON "qualificacao_socio"("codigo");

-- CreateIndex
CREATE INDEX "idx_simples_cnpj_basico" ON "simples"("cnpj_basico");

-- RenameIndex
ALTER INDEX "empresas_razao_social_idx" RENAME TO "idx_empresas_razao_social";
