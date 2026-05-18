# CONTRATOS

Este documento define a governança de contratos do projeto. Ele cobre OpenAPI, schemas de eventos, registry, compatibilidade e revisão de mudanças públicas.

## localização

```text
docs/contracts/http/                 OpenAPI por serviço
docs/contracts/events/               JSON Schema por evento compartilhado
docs/contracts/registry.json         registry geral de contratos
docs/contracts/schema-registry.json  registry de schemas de eventos
docs/contracts/portal/index.html     portal estático de navegação
```

`docs/contracts/` e a fonte atual de contratos. referências antigas a `contracts/` devem ser tratadas como históricas.

## Tipos de Contrato

### HTTP

Cada serviço com superfície pública deve possuir um arquivo OpenAPI em `docs/contracts/http/<serviço>.openapi.yaml`.

O OpenAPI deve declarar:

- paths e métodos;
- parâmetros relevantes;
- request body quando existir;
- respostas principais;
- exemplos quando ajudam o consumidor;
- tags ou agrupamento funcional quando aplicável.

### Eventos

Eventos compartilhados devem ter JSON Schema em `docs/contracts/events/`.

Um schema de evento deve declarar:

- nome ou tipo do evento;
- versão;
- payload público;
- identificadores de tenant/recurso quando aplicável;
- campos obrigatórios;
- compatibilidade esperada.

## Compatibilidade

### mudanças compatíveis

- adicionar endpoint novo;
- adicionar campo opcional;
- adicionar enum com fallback seguro;
- adicionar filtro sem alterar comportamento default;
- ampliar descricao, exemplo ou metadata;
- adicionar evento novo sem alterar consumidor existente.

### mudanças Potencialmente Breaking

- remover endpoint;
- renomear campo público;
- alterar tipo de campo;
- tornar campo opcional obrigatório;
- alterar semantica de status HTTP;
- mudar formato de erro público;
- remover enum sem fallback;
- alterar payload de evento consumido por outro serviço.

mudança breaking exige decisão explicita: versionar, manter compatibilidade temporaria ou registrar migração de consumidor.

## Registry

O registry existe para descoberta e auditoria. Ao adicionar ou remover contrato relevante, atualize:

- `docs/contracts/registry.json` para contratos HTTP/docs;
- `docs/contracts/schema-registry.json` para eventos.

O registry deve apontar para arquivos existentes e não deve conter entrada decorativa sem artefato correspondente.

## Controle Interno De Contratos

Antes de aceitar uma mudança contratual no código mantido internamente, confira:

- O arquivo OpenAPI/event schema mudou junto da implementação?
- O serviço dono e o correto?
- A alteração e compatível?
- Existe teste de contrato, unidade ou smoke proporcional ao risco?
- O `client-web/client-api` consegue regenerar o catálogo quando a mudança envolve HTTP?
- A documentação tem escopo correto e não replica conteudo de outro arquivo?

## validação

```bash
./scripts/test.sh contract
```

Para o console técnico:

```bash
cd client-web/client-api
npm run generate
```

## Eventos Versionados

Eventos atualmente vivem em `docs/contracts/events/` e cobrem domínios como catalog, CRM enrichment, documents signing, engagement provider events, fiscal consent/document events, platform-control lifecycle/go-live/quota, support cases e webhook-hub inbound/outbound. Contratos HTTP também cobrem busca operacional, e-discovery, BI semântico, governança de IA, incident command, policy decision center, approvals, runbooks, timeline, evidence vault, risk/compliance scoring, event mesh, tenant runtime, contract evolution, reconciliation, financial close, master data quality e lakehouse manifest.

Na v1.5.0, mudanças breaking continuam passando pelo Contract & Schema Evolution Center: snapshot, diff, classificação de breaking change, matriz de compatibilidade e aprovação quando houver remocao de operação ou alteração incompatível de schema. Rotas de provider externo também precisam declarar quando retornam `unavailable` por ausência de credencial BYOK ou quando usam API pública, para impedir que o contrato sugira uma integração produtiva inexistente.

Ao criar evento novo:

- use nome estavel;
- inclua versão;
- declare tenant/recurso quando fizer sentido;
- mantenha payload pequeno e objetivo;
- não use evento como atalho para vazar tabela interna.

## Ciclo de Vida de um Contrato

1. Necessidade pública aparece em um fluxo real.
2. O serviço dono e definido.
3. O OpenAPI ou JSON Schema e alterado.
4. A implementação acompanha o contrato.
5. Teste unitario/contrato/smoke cobre o comportamento proporcionalmente ao risco.
6. Registry e documentação são atualizados.
7. O console técnico pode regenerar catálogo quando a mudança for HTTP.

Esse ciclo evita que contrato vire documentação atrasada.

## Politica de versão

Enquanto o projeto está em consolidação interna, os contratos usam versões de artefato e changelog do repositório. Quando houver consumidor externo real, mudanças breaking devem ser tratadas com uma destas estrategias:

- criar nova versão de endpoint;
- aceitar os dois formatos por janela de migração;
- adicionar campo novo e manter campo antigo como deprecated;
- públicar guia de migração;
- bloquear a mudança ate consumidor crítico ser atualizado.

## Exemplos de Compatibilidade

### Campo opcional novo

compatível:

```json
{
  "publicId": "sub_123",
  "status": "active",
  "currentPeriodEndsAt": "2026-05-31T23:59:59Z"
}
```

Adicionar `cancellationReason` opcional não quebra consumidores que ignoram campos desconhecidos.

### Campo renomeado

Potencialmente breaking:

```json
{
  "publicId": "sub_123",
  "state": "active"
}
```

Se antes o contrato usava `status`, trocar para `state` exige estrategia de compatibilidade.

### Enum novo

Pode ser compatível se consumidores tiverem fallback para valor desconhecido. Sem fallback, pode quebrar UI, automação ou report.

## Regras Para Eventos

Eventos devem representar fatos:

- `fiscal.document.cancelled`;
- `platform-control.lifecycle-job.completed`;
- `webhook-hub.outbound-delivery.failed`.

Evite eventos que sejam comandos disfarçados:

- `pleaseCreateInvoice`;
- `updateCustomerSomewhere`;
- `syncEverythingNow`.

Se a intencao e comandar um serviço, use HTTP, workflow ou fila de comando explicitamente modelada.

## Contrato e Teste

| mudança | Teste minimo esperado |
|---------|-----------------------|
| descricao ou exemplo | contract validation |
| campo opcional novo | contract validation e teste local se usado |
| endpoint novo | teste de handler/regra e contract validation |
| mutação sensível | unidade, contract e smoke se afetar fluxo cross-service |
| evento novo | schema validation e teste do produtor |
| breaking change | plano de migração e validação de consumidores |

## Contrato e documentação

O contrato deve ser detalhado o bastante para maquina e consumidor técnico. A documentação humana deve explicar:

- quando usar o endpoint;
- qual o dono funcional;
- quais riscos existem;
- qual suite valida a mudança.

não copie todo o OpenAPI para markdown. Isso gera duplicidade e envelhece rápido.
