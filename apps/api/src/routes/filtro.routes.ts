import { Router } from 'express';
import { 
  criarFiltro, 
  atualizarFiltro, 
  buscarFiltros, 
  buscarFiltroPorId, 
  deletarFiltro 
} from '../controllers/filtro.controller';
import { authMiddleware } from '../middleware/auth.middleware';

const router = Router();

// Rotas protegidas por autenticação
router.use(authMiddleware);

router.post('/', criarFiltro);
router.put('/:id', atualizarFiltro);
router.get('/', buscarFiltros);
router.get('/:id', buscarFiltroPorId);
router.delete('/:id', deletarFiltro);

export default router;