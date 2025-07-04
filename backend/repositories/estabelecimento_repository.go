// neurocloser/backend/repositories/estabelecimento_repository.go
package repositories

import (
	"backend/models"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

// EstabelecimentoComEmpresa é uma struct auxiliar que combina campos de Estabelecimento e Empresa.
// Usada para escanear resultados de JOINs diretamente com sqlx.Select.
type EstabelecimentoComEmpresa struct {
	models.Estabelecimento // Embeda o struct Estabelecimento. Seus campos são mapeados por nome.
	// Campos da Empresa que vêm do JOIN.
	EmpresaRazaoSocial               sql.NullString `db:"emp_razao_social"`
	EmpresaNaturezaJuridrica         sql.NullString `db:"emp_natureza_juridica"`
	EmpresaQualificacaoResponsavel   sql.NullString `db:"emp_qualificacao_responsavel"`
	EmpresaPorteEmpresa              sql.NullString `db:"emp_porte_empresa"`
	EmpresaEnteFederativoResponsavel sql.NullString `db:"emp_ente_federativo_responsavel"`
	EmpresaCapitalSocialStr          sql.NullString `db:"emp_capital_social"`
}

// EstabelecimentoRepository define a interface para operações de dados do Estabelecimento.
type EstabelecimentoRepository interface {
	GetEstabelecimentoByID(id int) (*models.Estabelecimento, error)
	GetEstabelecimentoByCNPJBasico(cnpjBasico string) (*models.Estabelecimento, error)
	// Retorna uma slice do novo tipo combinado EstabelecimentoComEmpresa
	FindEstabelecimentosByFilters(filters map[string]interface{}, limit *int, offset *int) ([]*EstabelecimentoComEmpresa, error)
}

// estabelecimentoRepository implementa EstabelecimentoRepository para PostgreSQL.
type estabelecimentoRepository struct {
	db *sqlx.DB
}

// NewEstabelecimentoRepository cria uma nova instância de EstabelecimentoRepository.
func NewEstabelecimentoRepository(db *sqlx.DB) EstabelecimentoRepository {
	return &estabelecimentoRepository{db: db}
}

// GetEstabelecimentoByID busca um estabelecimento pelo seu ID.
func (r *estabelecimentoRepository) GetEstabelecimentoByID(id int) (*models.Estabelecimento, error) {
	var (
		estabelecimento         models.Estabelecimento
		idStr                   sql.NullString
		cnpj                    sql.NullString
		cnpjBasico              sql.NullString
		cnpjOrdem               sql.NullString
		cnpjDv                  sql.NullString
		matrizFilial            sql.NullString
		nomeFantasia            sql.NullString
		situacaoCadastral       sql.NullString
		dataSituacaoCadastral   sql.NullString
		motivoSituacaoCadastral sql.NullString
		nomeCidadeExterior      sql.NullString
		pais                    sql.NullString
		dataInicioAtividades    sql.NullString
		cnaeFiscal              sql.NullString
		cnaeFiscalSecundaria    sql.NullString
		tipoLogradouro          sql.NullString
		logradouro              sql.NullString
		numero                  sql.NullString
		complemento             sql.NullString
		bairro                  sql.NullString
		cep                     sql.NullString
		uf                      sql.NullString
		municipio               sql.NullString
		ddd1                    sql.NullString
		telefone1               sql.NullString
		ddd2                    sql.NullString
		telefone2               sql.NullString
		dddFax                  sql.NullString
		fax                     sql.NullString
		correioEletronico       sql.NullString
		situacaoEspecial        sql.NullString
		dataSituacaoEspecial    sql.NullString
	)

	query := `
        SELECT
            id, cnpj, cnpj_basico, cnpj_ordem, cnpj_dv, matriz_filial, nome_fantasia,
            situacao_cadastral, data_situacao_cadastral, motivo_situacao_cadastral,
            nome_cidade_exterior, pais, data_inicio_atividades, cnae_fiscal,
            cnae_fiscal_secundaria, tipo_logradouro, logradouro, numero, complemento,
            bairro, cep, uf, municipio, ddd1, telefone1, ddd2, telefone2,
            ddd_fax, fax, correio_eletronico, situacao_especial, data_situacao_especial
        FROM estabelecimento
        WHERE id = $1
    `
	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&idStr, &cnpj, &cnpjBasico, &cnpjOrdem, &cnpjDv, &matrizFilial, &nomeFantasia,
		&situacaoCadastral, &dataSituacaoCadastral, &motivoSituacaoCadastral,
		&nomeCidadeExterior, &pais, &dataInicioAtividades, &cnaeFiscal,
		&cnaeFiscalSecundaria, &tipoLogradouro, &logradouro, &numero, &complemento,
		&bairro, &cep, &uf, &municipio, &ddd1, &telefone1, &ddd2, &telefone2,
		&dddFax, &fax, &correioEletronico, &situacaoEspecial, &dataSituacaoEspecial,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar estabelecimento por ID: %w", err)
	}

	if idStr.Valid {
		parsedID, _ := strconv.Atoi(idStr.String)
		estabelecimento.ID = parsedID
	} else {
		estabelecimento.ID = 0
	}

	estabelecimento.CNPJ = cnpj.String
	estabelecimento.CNPJBasico = cnpjBasico.String
	estabelecimento.CNPJOrdem = cnpjOrdem.String
	estabelecimento.CNPJDV = cnpjDv.String
	estabelecimento.MatrizFilial = matrizFilial.String
	estabelecimento.NomeFantasia = nomeFantasia.String
	estabelecimento.SituacaoCadastral = situacaoCadastral.String
	estabelecimento.DataSituacaoCadastral = dataSituacaoCadastral.String
	estabelecimento.MotivoSituacaoCadastral = motivoSituacaoCadastral.String
	estabelecimento.NomeCidadeExterior = nomeCidadeExterior.String
	estabelecimento.Pais = pais.String
	estabelecimento.DataInicioAtividades = dataInicioAtividades.String
	estabelecimento.CNAEFiscal = cnaeFiscal.String
	estabelecimento.CNAEFiscalSecundaria = cnaeFiscalSecundaria.String
	estabelecimento.TipoLogradouro = tipoLogradouro.String
	estabelecimento.Logradouro = logradouro.String
	estabelecimento.Numero = numero.String
	estabelecimento.Complemento = complemento.String
	estabelecimento.Bairro = bairro.String
	estabelecimento.CEP = cep.String
	estabelecimento.UF = uf.String
	estabelecimento.Municipio = municipio.String
	estabelecimento.DDD1 = ddd1.String
	estabelecimento.Telefone1 = telefone1.String
	estabelecimento.DDD2 = ddd2.String
	estabelecimento.Telefone2 = telefone2.String
	estabelecimento.DDDFax = dddFax.String
	estabelecimento.Fax = fax.String
	estabelecimento.CorreioEletronico = correioEletronico.String
	estabelecimento.SituacaoEspecial = situacaoEspecial.String
	estabelecimento.DataSituacaoEspecial = dataSituacaoEspecial.String

	estabelecimento.FormatCNPJ() // Chama a função de formatação

	return &estabelecimento, nil
}

// GetEstabelecimentoByCNPJBasico busca um estabelecimento pelo seu CNPJ Básico.
func (r *estabelecimentoRepository) GetEstabelecimentoByCNPJBasico(cnpjBasico string) (*models.Estabelecimento, error) {
	var (
		estabelecimento         models.Estabelecimento
		idStr                   sql.NullString
		cnpj                    sql.NullString
		cnpjBasicoCol           sql.NullString
		cnpjOrdem               sql.NullString
		cnpjDv                  sql.NullString
		matrizFilial            sql.NullString
		nomeFantasia            sql.NullString
		situacaoCadastral       sql.NullString
		dataSituacaoCadastral   sql.NullString
		motivoSituacaoCadastral sql.NullString
		nomeCidadeExterior      sql.NullString
		pais                    sql.NullString
		dataInicioAtividades    sql.NullString
		cnaeFiscal              sql.NullString
		cnaeFiscalSecundaria    sql.NullString
		tipoLogradouro          sql.NullString
		logradouro              sql.NullString
		numero                  sql.NullString
		complemento             sql.NullString
		bairro                  sql.NullString
		cep                     sql.NullString
		uf                      sql.NullString
		municipio               sql.NullString
		ddd1                    sql.NullString
		telefone1               sql.NullString
		ddd2                    sql.NullString
		telefone2               sql.NullString
		dddFax                  sql.NullString
		fax                     sql.NullString
		correioEletronico       sql.NullString
		situacaoEspecial        sql.NullString
		dataSituacaoEspecial    sql.NullString
	)

	query := `
        SELECT
            id, cnpj, cnpj_basico, cnpj_ordem, cnpj_dv, matriz_filial, nome_fantasia,
            situacao_cadastral, data_situacao_cadastral, motivo_situacao_cadastral,
            nome_cidade_exterior, pais, data_inicio_atividades, cnae_fiscal,
            cnae_fiscal_secundaria, tipo_logradouro, logradouro, numero, complemento,
            bairro, cep, uf, municipio, ddd1, telefone1, ddd2, telefone2,
            ddd_fax, fax, correio_eletronico, situacao_especial, data_situacao_especial
        FROM estabelecimento
        WHERE cnpj_basico = $1
        ORDER BY matriz_filial DESC
        LIMIT 1
    `
	row := r.db.QueryRow(query, cnpjBasico)

	err := row.Scan(
		&idStr, &cnpj, &cnpjBasicoCol, &cnpjOrdem, &cnpjDv, &matrizFilial, &nomeFantasia,
		&situacaoCadastral, &dataSituacaoCadastral, &motivoSituacaoCadastral,
		&nomeCidadeExterior, &pais, &dataInicioAtividades, &cnaeFiscal,
		&cnaeFiscalSecundaria, &tipoLogradouro, &logradouro, &numero, &complemento,
		&bairro, &cep, &uf, &municipio, &ddd1, &telefone1, &ddd2, &telefone2,
		&dddFax, &fax, &correioEletronico, &situacaoEspecial, &dataSituacaoEspecial,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar estabelecimento por CNPJ Básico: %w", err)
	}

	if idStr.Valid {
		parsedID, _ := strconv.Atoi(idStr.String)
		estabelecimento.ID = parsedID
	} else {
		estabelecimento.ID = 0
	}

	estabelecimento.CNPJ = cnpj.String
	estabelecimento.CNPJBasico = cnpjBasicoCol.String
	estabelecimento.CNPJOrdem = cnpjOrdem.String
	estabelecimento.CNPJDV = cnpjDv.String
	estabelecimento.MatrizFilial = matrizFilial.String
	estabelecimento.NomeFantasia = nomeFantasia.String
	estabelecimento.SituacaoCadastral = situacaoCadastral.String
	estabelecimento.DataSituacaoCadastral = dataSituacaoCadastral.String
	estabelecimento.MotivoSituacaoCadastral = motivoSituacaoCadastral.String
	estabelecimento.NomeCidadeExterior = nomeCidadeExterior.String
	estabelecimento.Pais = pais.String
	estabelecimento.DataInicioAtividades = dataInicioAtividades.String
	estabelecimento.CNAEFiscal = cnaeFiscal.String
	estabelecimento.CNAEFiscalSecundaria = cnaeFiscalSecundaria.String
	estabelecimento.TipoLogradouro = tipoLogradouro.String
	estabelecimento.Logradouro = logradouro.String
	estabelecimento.Numero = numero.String
	estabelecimento.Complemento = complemento.String
	estabelecimento.Bairro = bairro.String
	estabelecimento.CEP = cep.String
	estabelecimento.UF = uf.String
	estabelecimento.Municipio = municipio.String
	estabelecimento.DDD1 = ddd1.String
	estabelecimento.Telefone1 = telefone1.String
	estabelecimento.DDD2 = ddd2.String
	estabelecimento.Telefone2 = telefone2.String
	estabelecimento.DDDFax = dddFax.String
	estabelecimento.Fax = fax.String
	estabelecimento.CorreioEletronico = correioEletronico.String
	estabelecimento.SituacaoEspecial = situacaoEspecial.String
	estabelecimento.DataSituacaoEspecial = dataSituacaoEspecial.String

	estabelecimento.FormatCNPJ() // Chama a função de formatação

	return &estabelecimento, nil
}

// struct auxiliar interna para escanear resultados de JOIN (Estabelecimento + Empresa)
// IMPORTANTE: Esta struct deve estar no nível de pacote, não dentro de uma função.
type estabelecimentoWithEmpresa struct {
	models.Estabelecimento `db:""` // Embeda Estabelecimento. `db:""` permite sqlx mapear por nome.

	// Campos da Empresa que vêm do JOIN (com alias únicos para evitar conflito de nome)
	EmpresaRazaoSocial               sql.NullString `db:"emp_razao_social"`
	EmpresaNaturezaJuridrica         sql.NullString `db:"emp_natureza_juridica"`
	EmpresaQualificacaoResponsavel   sql.NullString `db:"emp_qualificacao_responsavel"`
	EmpresaPorteEmpresa              sql.NullString `db:"emp_porte_empresa"`
	EmpresaEnteFederativoResponsavel sql.NullString `db:"emp_ente_federativo_responsavel"`
	EmpresaCapitalSocialStr          sql.NullString `db:"emp_capital_social"`
}

// FindEstabelecimentosByFilters busca estabelecimentos com base em múltiplos critérios de filtro.
// Retorna uma lista de EstabelecimentoComEmpresa, que inclui os dados de Empresa já carregados via JOIN.
func (r *estabelecimentoRepository) FindEstabelecimentosByFilters(filters map[string]interface{}, limit *int, offset *int) ([]*EstabelecimentoComEmpresa, error) {
	var results []estabelecimentoWithEmpresa // Vamos escanear para esta slice de structs combinadas

	// Aprimorando o SELECT para usar aliases e selecionar todas as colunas de 'e' e 'emp'.
	// Os aliases para colunas da empresa são CRUCIAIS para o sqlx mapear corretamente.
	baseQuery := `
		SELECT
            e.id, e.cnpj, e.cnpj_basico, e.cnpj_ordem, e.cnpj_dv, e.matriz_filial, e.nome_fantasia,
            e.situacao_cadastral, e.data_situacao_cadastral, e.motivo_situacao_cadastral,
            e.nome_cidade_exterior, e.pais, e.data_inicio_atividades, e.cnae_fiscal,
            e.cnae_fiscal_secundaria, e.tipo_logradouro, e.logradouro, e.numero, e.complemento,
            e.bairro, e.cep, e.uf, e.municipio, e.ddd1, e.telefone1, e.ddd2, e.telefone2,
            e.ddd_fax, e.fax, e.correio_eletronico, e.situacao_especial, e.data_situacao_especial,
            emp.razao_social AS emp_razao_social,
            emp.natureza_juridica AS emp_natureza_juridica,
            emp.qualificacao_responsavel AS emp_qualificacao_responsavel,
            emp.porte_empresa AS emp_porte_empresa,
            emp.ente_federativo_responsavel AS emp_ente_federativo_responsavel,
            emp.capital_social AS emp_capital_social
		FROM estabelecimento e
		JOIN empresas emp ON e.cnpj_basico = emp.cnpj_basico
		WHERE 1=1
	`
	queryParts := []string{baseQuery}
	args := []interface{}{}
	argCounter := 1

	// Adicione os filtros
	if cnpj, ok := filters["cnpj"].(string); ok && cnpj != "" {
		queryParts = append(queryParts, fmt.Sprintf(" AND e.cnpj = $%d", argCounter))
		args = append(args, cnpj)
		argCounter++
	}
	if nomeFantasia, ok := filters["nomeFantasia"].(string); ok && nomeFantasia != "" {
		queryParts = append(queryParts, fmt.Sprintf(" AND e.nome_fantasia ILIKE $%d", argCounter))
		args = append(args, "%"+nomeFantasia+"%")
		argCounter++
	}
	if uf, ok := filters["uf"].(string); ok && uf != "" {
		queryParts = append(queryParts, fmt.Sprintf(" AND e.uf = $%d", argCounter))
		args = append(args, uf)
		argCounter++
	}
	if situacaoCadastral, ok := filters["situacaoCadastral"].(string); ok && situacaoCadastral != "" {
		queryParts = append(queryParts, fmt.Sprintf(" AND e.situacao_cadastral = $%d", argCounter))
		args = append(args, situacaoCadastral)
		argCounter++
	}
	if cnaeFiscal, ok := filters["cnaeFiscal"].(string); ok && cnaeFiscal != "" {
		queryParts = append(queryParts, fmt.Sprintf(" AND e.cnae_fiscal = $%d", argCounter))
		args = append(args, cnaeFiscal)
		argCounter++
	}
	if cnaeFiscalSecundaria, ok := filters["cnaeFiscalSecundaria"].(string); ok && cnaeFiscalSecundaria != "" {
		queryParts = append(queryParts, fmt.Sprintf(" AND e.cnae_fiscal_secundaria LIKE $%d", argCounter)) // Busca por substring
		args = append(args, "%"+cnaeFiscalSecundaria+"%")
		argCounter++
	}
	if razaoSocial, ok := filters["razaoSocial"].(string); ok && razaoSocial != "" {
		queryParts = append(queryParts, fmt.Sprintf(" AND emp.razao_social ILIKE $%d", argCounter))
		args = append(args, "%"+razaoSocial+"%")
		argCounter++
	}
	if porteEmpresa, ok := filters["porteEmpresa"].(string); ok && porteEmpresa != "" {
		queryParts = append(queryParts, fmt.Sprintf(" AND emp.porte_empresa = $%d", argCounter))
		args = append(args, porteEmpresa)
		argCounter++
	}
	if minCapitalSocial, ok := filters["minCapitalSocial"].(float64); ok && minCapitalSocial >= 0 {
		queryParts = append(queryParts, fmt.Sprintf(" AND emp.capital_social >= $%d", argCounter))
		args = append(args, minCapitalSocial)
		argCounter++
	}
	if maxCapitalSocial, ok := filters["maxCapitalSocial"].(float64); ok && maxCapitalSocial >= 0 {
		queryParts = append(queryParts, fmt.Sprintf(" AND emp.capital_social <= $%d", argCounter))
		args = append(args, maxCapitalSocial)
		argCounter++
	}
	if municipio, ok := filters["municipio"].(string); ok && municipio != "" {
		queryParts = append(queryParts, fmt.Sprintf(" AND e.municipio ILIKE $%d", argCounter))
		args = append(args, "%"+municipio+"%")
		argCounter++
	}
	if naturezaJuridica, ok := filters["naturezaJuridica"].(string); ok && naturezaJuridica != "" {
		queryParts = append(queryParts, fmt.Sprintf(" AND emp.natureza_juridica = $%d", argCounter))
		args = append(args, naturezaJuridica)
		argCounter++
	}
	if dataSituacaoCadastralMin, ok := filters["dataSituacaoCadastralMin"].(string); ok && dataSituacaoCadastralMin != "" {
		queryParts = append(queryParts, fmt.Sprintf(" AND e.data_situacao_cadastral >= $%d", argCounter))
		args = append(args, dataSituacaoCadastralMin)
		argCounter++
	}
	if dataSituacaoCadastralMax, ok := filters["dataSituacaoCadastralMax"].(string); ok && dataSituacaoCadastralMax != "" {
		queryParts = append(queryParts, fmt.Sprintf(" AND e.data_situacao_cadastral <= $%d", argCounter))
		args = append(args, dataSituacaoCadastralMax)
		argCounter++
	}
	if dataInicioAtividadesMin, ok := filters["dataInicioAtividadesMin"].(string); ok && dataInicioAtividadesMin != "" {
		queryParts = append(queryParts, fmt.Sprintf(" AND e.data_inicio_atividades >= $%d", argCounter))
		args = append(args, dataInicioAtividadesMin)
		argCounter++
	}
	if dataInicioAtividadesMax, ok := filters["dataInicioAtividadesMax"].(string); ok && dataInicioAtividadesMax != "" {
		queryParts = append(queryParts, fmt.Sprintf(" AND e.data_inicio_atividades <= $%d", argCounter))
		args = append(args, dataInicioAtividadesMax)
		argCounter++
	}

	fullQuery := strings.Join(queryParts, " ") + " ORDER BY e.cnpj ASC"

	if limit != nil && *limit > 0 {
		fullQuery += fmt.Sprintf(" LIMIT $%d", argCounter)
		args = append(args, *limit)
		argCounter++
	}
	if offset != nil && *offset >= 0 {
		fullQuery += fmt.Sprintf(" OFFSET $%d", argCounter)
		args = append(args, *offset)
		argCounter++
	}

	// Usamos sqlx.Select para escanear diretamente para a slice da struct combinada.
	err := r.db.Select(&results, fullQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar estabelecimentos com filtros: %w", err)
	}

	var finalResults []*EstabelecimentoComEmpresa
	for i := range results {
		results[i].Estabelecimento.FormatCNPJ() // Chama o método FormatCNPJ() do Estabelecimento
		e := EstabelecimentoComEmpresa(results[i])
		finalResults = append(finalResults, &e)
	}

	return finalResults, nil
}
