import { Router } from 'express';
import { 
  buscarEmpresas, 
  buscarEmpresaPorId, 
  exportarParaExcel, 
  exportarParaCRM, 
  exportarParaEmail 
} from '../controllers/empresa.controller';
import { authMiddleware } from '../middleware/auth.middleware';

const router = Router();

// Rotas protegidas por autenticação
router.use(authMiddleware);

// Busca de empresas com filtros
router.get('/', buscarEmpresas);

// Busca empresa por CNPJ básico
router.get('/:cnpj', buscarEmpresaPorId);

// Exportação de dados
router.post('/exportar/excel', exportarParaExcel);
router.post('/exportar/crm', exportarParaCRM);
router.post('/exportar/email', exportarParaEmail);

export default router;