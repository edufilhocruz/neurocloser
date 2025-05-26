# 🧠 Neurocloser

O **Neurocloser** é uma plataforma de prospecção automatizada que utiliza **inteligência artificial, processamento de voz e automação de CRM** para entrar em contato com decisores de empresas, agendar reuniões e auxiliar no fechamento de vendas por chamadas telefônicas inteligentes.

---

## 📌 Objetivo

Automatizar o processo de **identificação**, **qualificação**, **abordagem** e **agendamento de reuniões comerciais** com empresas, utilizando:

- Banco de dados com **65 milhões de CNPJs**
- Filtros baseados no **ICP (Ideal Customer Profile)**
- **IA conversacional** (OpenAI ou LLaMA)
- Integração com CRM (Twenty)
- Ligações automáticas com gravação (Fonoster + Nari-labs)
- Assistente de **call monitor** com IA para análises e feedbacks

---

## 🧱 Arquitetura

O projeto segue o padrão **monorepo** com múltiplos microserviços, separados por contexto e responsabilidade.

Tecnologias principais:

- **Node.js + TypeScript** (backend)
- **React + Next.js + TailwindCSS** (frontend)
- **Prisma + PostgreSQL** (banco de dados)
- **Python** (ingestão de dados da Receita Federal)
- **OpenAI / LLaMA** (IA para prospecção)
- **Fonoster + Nari-labs** (voz e chamadas)
- **Docker + K3s + Cloud Code** (ambiente e deploy)

---

<pre lang="markdown">
📂 neurocloser/
├── apps/                  # Aplicações executáveis
│   ├── api/               # Backend com filtros de CNPJ e envio para CRM
│   ├── web/               # Frontend em Next.js
│   ├── worker-ai/         # Lógica de IA conversacional e agendamento
│   ├── worker-voice/      # Conversão de texto em voz com Nari-labs
│   └── worker-call-monitor/ # Análise de reuniões com IA
│
├── packages/              # Pacotes reutilizáveis (core, integrações e UI)
│   ├── database/          # Prisma, modelos, repositórios
│   ├── crm/               # Cliente da API do CRM Twenty
│   ├── fonoster/          # Integração com sistema de chamadas
│   ├── listmonk/          # Disparo de campanhas de cold email
│   ├── ai/                # Integrações com OpenAI e LLaMA
│   ├── voice/             # Integração com TTS/Nari-labs
│   └── ui/                # Design system com componentes React
│
├── scripts/               # Scripts auxiliares e ingestão de dados
├── docs/                  # Documentação técnica
├── docker-compose.yml     # Ambiente local com containers
├── tsconfig.base.json     # Configuração base TypeScript
├── .eslintrc.js           # Lint global
└── .prettierrc            # Formatação de código
</pre>


## 🚀 Funcionalidades 

✅ Ingestão de dados da Receita Federal (CSV → SQLite → PostgreSQL)
✅ API para filtragem de empresas por ICP
✅ Dashboard com estatísticas de empresas (ativas, com e-mail, com telefone, etc.)
🔄 Integração com CRM Twenty
🔄 Ligações automáticas com gravação e TTS
🔄 IA para prospecção ativa e resposta de leads
🔄 IA para análise de reuniões com closers humanos

## 🧪 Como Rodar o Projeto Localmente

✅ Docker + Docker Compose
✅ Node.js (v22+)
✅ pnpm (npm install -g pnpm)
✅ Python 3.10+
✅ PostgreSQL 16+

## 🧪 Passos

# Clone o projeto
git clone git@github.com:SEU_USUARIO/neurocloser.git
cd neurocloser

# Instale dependências
pnpm install

# Suba os serviços (PostgreSQL, etc)
docker-compose up -d

# Rode a API
pnpm dev --filter api

# Rode o frontend
pnpm dev --filter web

## 📡 Microserviços

api - Backend REST com filtros e exportações
web - Dashboard com UI para prospecção
worker-ai - IA para dialogar com clientes
worker-voice - Transforma resposta de IA em voz realista
worker-call-monitor - Analisa reuniões com closers e sugere ações

## 🛠️ Em desenvolvimento

🔧 ICP dinâmico com regras personalizadas
🔧 Integração com LLaMA local
🔧 Deploy automatizado via K3s e Portainer
🔧 Dashboard de performance de campanhas
