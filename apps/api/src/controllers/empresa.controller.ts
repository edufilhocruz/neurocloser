import { Request, Response } from 'express';
import { PrismaClient } from '@prisma/client';
import { generateExcel } from '../utils/excel.util';
import { exportToCRM } from '../services/crm.service';
import { exportToListmonk } from '../services/listmonk.service';

const prisma = new PrismaClient();

/**
 * Busca empresas com filtros diversos
 */
export const buscarEmpresas = async (req: Request, res: Response) => {
  try {
    const {
      razaoSocial,
      cnpj,
      naturezaJuridica,
      uf,
      municipio,
      situacaoCadastral,
      porte,
      cnae,
      temEmail,
      temTelefone,
      temCelular,
      ativa,
      page = 1,
      limit = 25
    } = req.query;

    // Construir condições do filtro
    const where: any = {};
    
    if (razaoSocial) {
      where.razao_social = {
        contains: String(razaoSocial),
        mode: 'insensitive'
      };
    }
    
    if (cnpj) {
      where.cnpj_basico = String(cnpj);
    }
    
    if (naturezaJuridica) {
      where.natureza_juridica = String(naturezaJuridica);
    }
    
    if (ativa !== undefined) {
      where.ativa = ativa === 'true';
    }
    
    // Filtros para relacionamentos
    if (uf || municipio || situacaoCadastral || cnae) {
      where.estabelecimentos = {
        some: {
          ...(uf && { uf: String(uf) }),
          ...(municipio && { municipio: String(municipio) }),
          ...(situacaoCadastral && { situacao_cadastral: String(situacaoCadastral) }),
          ...(cnae && { cnae_principal: String(cnae) })
        }
      };
    }
    
    // Filtros para contatos
    if (temEmail === 'true') {
      where.email = { not: null };
    }
    
    if (temTelefone === 'true') {
      where.telefone_fixo = { not: null };
    }
    
    if (temCelular === 'true') {
      where.celular = { not: null };
    }
    
    // Calcular paginação
    const skip = (Number(page) - 1) * Number(limit);
    
    // Buscar empresas com contagem total
    const [empresas, total] = await Promise.all([
      prisma.empresa.findMany({
        where,
        include: {
          estabelecimentos: true,
          simples: true
        },
        skip,
        take: Number(limit)
      }),
      prisma.empresa.count({ where })
    ]);
    
    // Estatísticas adicionais para o dashboard
    const estatisticas = await prisma.$transaction([
      prisma.empresa.count({ where: { ativa: true } }),
      prisma.empresa.count({ where: { ativa: true, celular: { not: null } } }),
      prisma.empresa.count({ where: { ativa: true, telefone_fixo: { not: null } } }),
      prisma.empresa.count({ where: { ativa: true, email: { not: null } } })
    ]);
    
    return res.json({
      empresas,
      total,
      totalPages: Math.ceil(total / Number(limit)),
      currentPage: Number(page),
      estatisticas: {
        totalEmpresas: total,
        empresasAtivas: estatisticas[0],
        empresasAtivasComCelular: estatisticas[1],
        empresasAtivasComTelefoneFixo: estatisticas[2],
        empresasAtivasComEmail: estatisticas[3]
      }
    });
  } catch (error) {
    console.error('Erro ao buscar empresas:', error);
    return res.status(500).json({ error: 'Erro ao buscar empresas' });
  }
};

/**
 * Busca empresa por CNPJ básico
 */
export const buscarEmpresaPorId = async (req: Request, res: Response) => {
  try {
    const { cnpj } = req.params;
    
    const empresa = await prisma.empresa.findUnique({
      where: { cnpj_basico: cnpj },
      include: {
        estabelecimentos: true,
        socios: true,
        simples: true
      }
    });
    
    if (!empresa) {
      return res.status(404).json({ error: 'Empresa não encontrada' });
    }
    
    return res.json(empresa);
  } catch (error) {
    console.error('Erro ao buscar empresa por CNPJ:', error);
    return res.status(500).json({ error: 'Erro ao buscar empresa por CNPJ' });
  }
};

/**
 * Exporta dados de empresas para Excel
 */
export const exportarParaExcel = async (req: Request, res: Response) => {
  try {
    const { filtro, nome } = req.body;
    const userId = (req as any).user.id;
    
    // Criar registro de exportação
    const exportacao = await prisma.listaExportacao.create({
      data: {
        nome: nome || `Exportação Excel ${new Date().toISOString()}`,
        usuario_id: userId,
        tipo: 'excel',
        status: 'processando',
        filtro_id: filtro?.id
      }
    });
    
    // Iniciar processo de exportação em background
    setTimeout(async () => {
      try {
        // Construir where baseado no filtro
        const where = filtro ? JSON.parse(JSON.stringify(filtro.condicoes)) : {};
        
        // Buscar empresas
        const empresas = await prisma.empresa.findMany({
          where,
          include: {
            estabelecimentos: true
          },
          take: 10000 // Limitar para evitar problemas de memória
        });
        
        // Gerar arquivo Excel
        const filePath = await generateExcel(empresas, nome);
        
        // Atualizar exportação
        await prisma.listaExportacao.update({
          where: { id: exportacao.id },
          data: {
            status: 'concluido',
            arquivo_url: filePath,
            total_registros: empresas.length
          }
        });
        
        // Log de sucesso
        await prisma.logIntegracao.create({
          data: {
            tipo: 'excel',
            status: 'sucesso',
            mensagem: `Exportação Excel concluída: ${empresas.length} registros`,
            detalhes: { exportacaoId: exportacao.id }
          }
        });
      } catch (error) {
        console.error('Erro na exportação Excel em background:', error);
        
        // Atualizar como erro
        await prisma.listaExportacao.update({
          where: { id: exportacao.id },
          data: {
            status: 'erro'
          }
        });
        
        // Log de erro
        await prisma.logIntegracao.create({
          data: {
            tipo: 'excel',
            status: 'erro',
            mensagem: `Erro na exportação Excel: ${(error as Error).message}`,
            detalhes: { exportacaoId: exportacao.id, erro: (error as Error).stack }
          }
        });
      }
    }, 100);
    
    return res.status(202).json({ 
      message: 'Exportação iniciada com sucesso', 
      exportacaoId: exportacao.id 
    });
  } catch (error) {
    console.error('Erro ao iniciar exportação para Excel:', error);
    return res.status(500).json({ error: 'Erro ao iniciar exportação para Excel' });
  }
};

/**
 * Exporta dados de empresas para o CRM
 */
export const exportarParaCRM = async (req: Request, res: Response) => {
  try {
    const { filtro, nome } = req.body;
    const userId = (req as any).user.id;
    
    // Criar registro de exportação
    const exportacao = await prisma.listaExportacao.create({
      data: {
        nome: nome || `Exportação CRM ${new Date().toISOString()}`,
        usuario_id: userId,
        tipo: 'crm',
        status: 'processando',
        filtro_id: filtro?.id
      }
    });
    
    // Iniciar processo de exportação em background
    setTimeout(async () => {
      try {
        // Construir where baseado no filtro
        const where = filtro ? JSON.parse(JSON.stringify(filtro.condicoes)) : {};
        
        // Buscar empresas
        const empresas = await prisma.empresa.findMany({
          where,
          include: {
            estabelecimentos: true,
            socios: true
          },
          take: 5000 // Limitar para evitar problemas
        });
        
        // Exportar para o CRM
        const resultado = await exportToCRM(empresas);
        
        // Atualizar exportação
        await prisma.listaExportacao.update({
          where: { id: exportacao.id },
          data: {
            status: 'concluido',
            total_registros: empresas.length
          }
        });
        
        // Atualizar empresas como enviadas ao CRM
        await prisma.$transaction(
          empresas.map(empresa => 
            prisma.empresa.update({
              where: { cnpj_basico: empresa.cnpj_basico },
              data: { enviado_crm: true }
            })
          )
        );
        
        // Log de sucesso
        await prisma.logIntegracao.create({
          data: {
            tipo: 'crm',
            status: 'sucesso',
            mensagem: `Exportação CRM concluída: ${empresas.length} registros`,
            detalhes: { exportacaoId: exportacao.id, resultado }
          }
        });
      } catch (error) {
        console.error('Erro na exportação CRM em background:', error);
        
        // Atualizar como erro
        await prisma.listaExportacao.update({
          where: { id: exportacao.id },
          data: {
            status: 'erro'
          }
        });
        
        // Log de erro
        await prisma.logIntegracao.create({
          data: {
            tipo: 'crm',
            status: 'erro',
            mensagem: `Erro na exportação CRM: ${(error as Error).message}`,
            detalhes: { exportacaoId: exportacao.id, erro: (error as Error).stack }
          }
        });
      }
    }, 100);
    
    return res.status(202).json({ 
      message: 'Exportação para CRM iniciada com sucesso', 
      exportacaoId: exportacao.id 
    });
  } catch (error) {
    console.error('Erro ao iniciar exportação para CRM:', error);
    return res.status(500).json({ error: 'Erro ao iniciar exportação para CRM' });
  }
};

/**
 * Exporta dados de empresas para lista de email
 */
export const exportarParaEmail = async (req: Request, res: Response) => {
  try {
    const { filtro, nome, listaNome } = req.body;
    const userId = (req as any).user.id;
    
    // Criar registro de exportação
    const exportacao = await prisma.listaExportacao.create({
      data: {
        nome: nome || `Exportação Email ${new Date().toISOString()}`,
        usuario_id: userId,
        tipo: 'email',
        status: 'processando',
        filtro_id: filtro?.id
      }
    });
    
    // Iniciar processo de exportação em background
    setTimeout(async () => {
      try {
        // Construir where baseado no filtro
        const where = filtro ? JSON.parse(JSON.stringify(filtro.condicoes)) : {};
        
        // Adicionar filtro para garantir que há email
        where.email = { not: null };
        
        // Buscar empresas
        const empresas = await prisma.empresa.findMany({
          where,
          select: {
            cnpj_basico: true,
            razao_social: true,
            email: true,
            estabelecimentos: {
              select: {
                nome_fantasia: true
              }
            }
          },
          take: 10000 // Limitar para evitar problemas
        });
        
        // Filtrar apenas empresas com email válido
        const empresasComEmail = empresas.filter(e => e.email && e.email.includes('@'));
        
        // Preparar dados para o Listmonk
        const dados = empresasComEmail.map(empresa => ({
          email: empresa.email as string,
          name: empresa.razao_social,
          cnpj: empresa.cnpj_basico,
          nome_fantasia: empresa.estabelecimentos[0]?.nome_fantasia || ''
        }));
        
        // Exportar para o Listmonk
        const resultado = await exportToListmonk(dados, listaNome || nome || 'Prospecção CNPJ');
        
        // Atualizar exportação
        await prisma.listaExportacao.update({
          where: { id: exportacao.id },
          data: {
            status: 'concluido',
            total_registros: empresasComEmail.length
          }
        });
        
        // Atualizar empresas como enviadas para email
        await prisma.$transaction(
          empresasComEmail.map(empresa => 
            prisma.empresa.update({
              where: { cnpj_basico: empresa.cnpj_basico },
              data: { enviado_email: true }
            })
          )
        );
        
        // Log de sucesso
        await prisma.logIntegracao.create({
          data: {
            tipo: 'email',
            status: 'sucesso',
            mensagem: `Exportação Email concluída: ${empresasComEmail.length} registros`,
            detalhes: { exportacaoId: exportacao.id, resultado }
          }
        });
      } catch (error) {
        console.error('Erro na exportação Email em background:', error);
        
        // Atualizar como erro
        await prisma.listaExportacao.update({
          where: { id: exportacao.id },
          data: {
            status: 'erro'
          }
        });
        
        // Log de erro
        await prisma.logIntegracao.create({
          data: {
            tipo: 'email',
            status: 'erro',
            mensagem: `Erro na exportação Email: ${(error as Error).message}`,
            detalhes: { exportacaoId: exportacao.id, erro: (error as Error).stack }
          }
        });
      }
    }, 100);
    
    return res.status(202).json({ 
      message: 'Exportação para lista de email iniciada com sucesso', 
      exportacaoId: exportacao.id 
    });
  } catch (error) {
    console.error('Erro ao iniciar exportação para lista de email:', error);
    return res.status(500).json({ error: 'Erro ao iniciar exportação para lista de email' });
  }
};