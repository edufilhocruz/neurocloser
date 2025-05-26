import { Router } from 'express';
import { login, registro, verificarToken } from '../controllers/auth.controller';
import { authMiddleware } from '../middleware/auth.middleware';

const router = Router();

// Rotas públicas
router.post('/login', login);
router.post('/registro', registro);

// Rota protegida
router.get('/verificar', authMiddleware, verificarToken);

export default router;