# Business Operating System

![Business Operating System](https://img.shields.io/badge/Business%20Operating%20System-Platform-111827?style=for-the-badge&logo=github&logoColor=white)

**Um sistema completo para organizar as operações de uma empresa em um só ecossistema, conectando vendas, clientes, cobranças, financeiro, fiscal, documentos, análises, automações e integrações com serviços externos.**

[![Versão](https://img.shields.io/badge/Versão-1.4.6-2563EB?style=flat)](docs/CHANGELOG.md)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-contratos-6BA539?style=flat&logo=openapiinitiative&logoColor=white)](docs/contracts/http)
[![Serviços](https://img.shields.io/badge/Serviços-26%20HTTP%20APIs-111827?style=flat)](docs/SERVICOS.md)
[![Console](https://img.shields.io/badge/API%20Console-client--api-2563EB?style=flat)](client-web/client-api)
[![Runtime](https://img.shields.io/badge/Runtime-Docker%20Compose-2496ED?style=flat&logo=docker&logoColor=white)](infra)

## Visão Geral

O projeto é um sistema operacional de negócios. Em vez de tratar cada área da empresa como uma ferramenta separada, ele organiza tudo como partes de uma mesma operação. A venda conversa com a cobrança, a cobrança conversa com o financeiro, o financeiro conversa com documentos e relatórios, e as integrações externas entram nesse fluxo de forma controlada.

A ideia é cobrir a jornada completa de uma empresa. Um cliente pode entrar como oportunidade comercial, passar por negociação, virar contrato, gerar cobrança, movimentar o financeiro, exigir documentos, acionar tarefas automáticas, alimentar relatórios e deixar evidências para auditoria. O valor do projeto está justamente nessa visão ponta a ponta.

Embora o nome histórico do repositório ainda lembre ERP, o escopo é mais amplo. O projeto reúne funções parecidas com CRM para relacionamento com clientes, ERP para operação interna, plataformas de cobrança para assinaturas e faturas, ferramentas de automação para tarefas repetitivas, painéis de análise para decisões e uma camada de integração para conversar com serviços de terceiros.

## Documentação

| Arquivo | Finalidade |
|---------|------------|
| [README.md](README.md) | visão geral objetiva do repositório |
| [README_EN.md](README_EN.md) | visão detalhada em inglês |
| [docs/ARQUITETURA.md](docs/ARQUITETURA.md) | arquitetura, fronteiras e decisões técnicas |
| [docs/API.md](docs/API.md) | regras da API e índice de contratos |
| [docs/SERVICOS.md](docs/SERVICOS.md) | mapa dos serviços e responsabilidades |
| [docs/CONTRATOS.md](docs/CONTRATOS.md) | contratos OpenAPI, eventos e compatibilidade |
| [docs/INTEGRACOES.md](docs/INTEGRACOES.md) | integrações internas e externas |
| [docs/OPERACOES.md](docs/OPERACOES.md) | execução local, validação e diagnóstico |
| [docs/PADROES.md](docs/PADROES.md) | padrões de engenharia |
| [docs/CHANGELOG.md](docs/CHANGELOG.md) | histórico de evolução por versão |

## O Que O Projeto Entrega

O projeto entrega uma base para operar uma empresa com rastreabilidade. Isso significa que ações importantes não ficam perdidas em planilhas, telas isoladas ou processos manuais sem histórico. Cada etapa pode ser registrada, consultada, validada e conectada a outras áreas.

Na área comercial, ele ajuda a organizar leads, clientes, oportunidades, propostas e contratos. Na área financeira, cobre cobranças, faturas, recebíveis, pagamentos, comissões e conciliações. Na área fiscal e documental, mantém documentos, certificados, anexos, evidências e rotinas de auditoria. Na operação, permite fluxos de trabalho, filas, suporte, notificações, buscas e relatórios.

O projeto também foi preparado para conversar com serviços externos. Isso inclui gateways de pagamento, serviços de inteligência artificial, leitura automática de documentos, consulta cadastral, notícias, câmbio, bancos, comunicação e assinatura digital. Quando uma chave de acesso não existe, a funcionalidade deve aparecer como indisponível ou manual, sem fingir que está integrada de verdade.

## Conceitos Em Linguagem Simples

| Termo | Significado neste projeto |
|-------|---------------------------|
| CRM | Parte responsável por relacionamento comercial, clientes, leads e oportunidades. |
| ERP | Parte responsável por processos internos, cadastros, fiscal, estoque, compras e financeiro. |
| Billing | Cobranças, assinaturas, faturas e tentativas de pagamento. |
| Automação | Tarefas que o sistema executa seguindo regras, como criar uma etapa, avisar outro serviço ou continuar um fluxo. |
| Analytics | Relatórios e indicadores que ajudam a entender operação, risco, qualidade e progresso. |
| Aviso automático, ou webhook | Mensagem enviada ou recebida quando algo importante acontece em outro sistema. |
| Serviço externo | Ferramenta de terceiro usada pelo sistema, como Stripe, OpenAI, WhatsApp, DocuSign, BrasilAPI ou um provedor fiscal. |
| Chave própria, ou BYOK | Modelo em que o usuário informa a própria chave de acesso para ativar um serviço externo. |
| Fluxo de trabalho | Sequência organizada de etapas para executar um processo de negócio. |
| Governança | Controles que mostram quem pode fazer o quê, em qual empresa ou ambiente, com quais limites, riscos e evidências. |

## Principais Áreas

| Área | O que cobre |
|------|-------------|
| Comercial | Leads, clientes, oportunidades, propostas, contratos recorrentes e catálogo comercial. |
| Cobrança e financeiro | Assinaturas, faturas, tentativas de pagamento, recebíveis, contas a pagar, tesouraria, comissões e conciliação. |
| Fiscal e banking | Documentos fiscais, certificados, SPED, Pix, boletos, Open Finance, extratos e reconciliação bancária. |
| Documentos | Anexos, versões, armazenamento, assinatura, leitura automática de documentos e trilhas de auditoria. |
| Fluxos de trabalho | Processos automatizados entre áreas, com registro de execução, falha, nova tentativa e correção. |
| Análises | Relatórios operacionais, indicadores de qualidade, risco, fechamento financeiro e visão de plataforma. |
| Integrações externas | Pagamentos, IA, leitura automática de documentos, consulta cadastral, mercado, notícias, comunicação, assinatura e avisos automáticos entre sistemas. |
| Governança SaaS | Empresas, usuários, permissões, autenticação reforçada, limites de uso, catálogo de ativações e ciclo de vida operacional. |
| Busca e evidência | Busca operacional, descoberta de informações, retenção legal, exportações controladas e provas de auditoria. |

## Fluxos De Negócio

1. Fluxo comercial. O sistema acompanha a jornada desde o primeiro contato com um possível cliente até a proposta e o contrato.

2. Fluxo de cobrança. Um contrato pode gerar faturas, tentativas de pagamento, eventos de cobrança e reflexos financeiros.

3. Fluxo financeiro. Recebíveis, contas a pagar, tesouraria, comissões, extratos e conciliações aparecem como partes do mesmo acompanhamento.

4. Fluxo documental. Anexos, contratos, documentos fiscais, assinaturas e evidências ficam associados à operação que os gerou.

5. Fluxo de automação. Processos repetitivos podem seguir etapas previsíveis, com registro de execução, falha, nova tentativa e compensação.

6. Fluxo de integração. Serviços externos podem ser conectados quando o usuário informa as chaves necessárias, mantendo transparência sobre o que está ativo.

7. Fluxo de governança. O projeto mantém visibilidade sobre empresas, ambientes, permissões, limites, riscos, aprovações, incidentes e evidências.

## Console Técnico Da API

`client-web/client-api` é uma central técnica para testar e entender a API. Ele funciona como uma versão mais completa e específica de um Swagger UI, com catálogo de rotas, documentação, ambientes, contratos, jornadas e telas de operação.

Esse console não é o produto visual final para usuários de negócio. Ele existe para desenvolvedores, QA, operação e validação técnica enquanto um frontend corporativo próprio pode ser construído separadamente.

```bash
cd client-web/client-api
npm install
npm run generate
npm run dev
```

## Runtime Local

```bash
./scripts/build.sh
./scripts/build.sh up
./scripts/build.sh ps
./scripts/build.sh logs edge
./scripts/build.sh down
```

Banco:

```bash
./scripts/build.sh migrate all
./scripts/build.sh seed all
./scripts/build.sh psql
./scripts/build.sh backup /tmp/erp-local-backup.sql
./scripts/build.sh restore /tmp/erp-local-backup.sql
```

## Validação

| Escopo | Comando |
|--------|---------|
| testes unitários | `./scripts/test.sh unit` |
| testes de integração | `./scripts/test.sh integration` |
| contratos HTTP e eventos | `./scripts/test.sh contract` |
| checks de plataforma | `./scripts/test.sh platform` |
| smoke tests | `./scripts/test.sh smoke` |
| performance | `./scripts/test.sh performance` |
| backup e restore | `./scripts/test.sh backup-restore` |
| hardening | `./scripts/test.sh hardening` |
| aceite produtivo | `./scripts/test.sh production-readiness` |

## Estrutura Do Repositório

```text
client-web/client-api/     console técnico da API
docs/                      documentação do projeto
docs/contracts/            OpenAPI, eventos, registry e portal
infra/                     Docker Compose e Kubernetes
scripts/                   entradas de runtime, build e validação
service-api/               serviços backend e contextos PostgreSQL
```

## Ownership Privado

Este repositório é mantido de forma privada. Mudanças de código são controladas diretamente pelo mantenedor.

## Contato

**Thiago Di Faria**  
thiagodifaria@gmail.com

[GitHub](https://github.com/thiagodifaria)  
[LinkedIn](https://linkedin.com/in/thiagodifaria)
