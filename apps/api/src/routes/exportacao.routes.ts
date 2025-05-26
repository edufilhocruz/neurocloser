import { Router } from 'express';
import { 
  criarExportacao,
  buscarExportacoes,
  buscarExportacaoPorId,
  downloadArquivo,
  verificarStatusExportacao
} from '../controllers/exportacao.controller';
import { authMiddleware } from '../middleware/auth.middleware';

const router = Router();

// Rotas protegidas por autenticação
router.use(authMiddleware);

router.post('/', criarExportacao);
router.get('/', buscarExportacoes);
router.get('/:id', buscarExportacaoPorId);
router.get('/:id/download', downloadArquivo);
router.get('/:id/status', verificarStatusExportacao);

export default router;