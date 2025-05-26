import { Router } from 'express';
import empresaRoutes from './empresa.routes';
import authRoutes from './auth.routes';
import filtroRoutes from './filtro.routes';
import exportacaoRoutes from './exportacao.routes';

const router = Router();

router.use('/auth', authRoutes);
router.use('/empresas', empresaRoutes);
router.use('/filtros', filtroRoutes);
router.use('/exportacoes', exportacaoRoutes);

export default router;