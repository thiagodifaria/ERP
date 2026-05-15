# CONTRATOS

Este documento define a governanca de contratos do ERP. Ele cobre OpenAPI, schemas de eventos, registry, compatibilidade e revisao de mudancas publicas.

## Localizacao

```text
docs/contracts/http/                 OpenAPI por servico
docs/contracts/events/               JSON Schema por evento compartilhado
docs/contracts/registry.json         registry geral de contratos
docs/contracts/schema-registry.json  registry de schemas de eventos
docs/contracts/portal/index.html     portal estatico de navegacao
```

`docs/contracts/` e a fonte atual de contratos. Referencias antigas a `contracts/` devem ser tratadas como historicas.

## Tipos de Contrato

### HTTP

Cada servico com superficie publica deve possuir um arquivo OpenAPI em `docs/contracts/http/<servico>.openapi.yaml`.

O OpenAPI deve declarar:

- paths e metodos;
- parametros relevantes;
- request body quando existir;
- respostas principais;
- exemplos quando ajudam o consumidor;
- tags ou agrupamento funcional quando aplicavel.

### Eventos

Eventos compartilhados devem ter JSON Schema em `docs/contracts/events/`.

Um schema de evento deve declarar:

- nome ou tipo do evento;
- versao;
- payload publico;
- identificadores de tenant/recurso quando aplicavel;
- campos obrigatorios;
- compatibilidade esperada.

## Compatibilidade

### Mudancas Compativeis

- adicionar endpoint novo;
- adicionar campo opcional;
- adicionar enum com fallback seguro;
- adicionar filtro sem alterar comportamento default;
- ampliar descricao, exemplo ou metadata;
- adicionar evento novo sem alterar consumidor existente.

### Mudancas Potencialmente Breaking

- remover endpoint;
- renomear campo publico;
- alterar tipo de campo;
- tornar campo opcional obrigatorio;
- alterar semantica de status HTTP;
- mudar formato de erro publico;
- remover enum sem fallback;
- alterar payload de evento consumido por outro servico.

Mudanca breaking exige decisao explicita: versionar, manter compatibilidade temporaria ou registrar migracao de consumidor.

## Registry

O registry existe para descoberta e auditoria. Ao adicionar ou remover contrato relevante, atualize:

- `docs/contracts/registry.json` para contratos HTTP/docs;
- `docs/contracts/schema-registry.json` para eventos.

O registry deve apontar para arquivos existentes e nao deve conter entrada decorativa sem artefato correspondente.

## Controle Interno De Contratos

Antes de aceitar uma mudanca contratual no codigo mantido internamente, confira:

- O arquivo OpenAPI/event schema mudou junto da implementacao?
- O servico dono e o correto?
- A alteracao e compativel?
- Existe teste de contrato, unidade ou smoke proporcional ao risco?
- O `client-web/client-api` consegue regenerar o catalogo quando a mudanca envolve HTTP?
- A documentacao tem escopo correto e nao replica conteudo de outro arquivo?

## Validacao

```bash
./scripts/test.sh contract
```

Para o console tecnico:

```bash
cd client-web/client-api
npm run generate
```

## Eventos Versionados

Eventos atualmente vivem em `docs/contracts/events/` e cobrem dominios como catalog, CRM enrichment, documents signing, engagement provider events, fiscal consent/document events, platform-control lifecycle/go-live/quota, support cases e webhook-hub inbound/outbound.

Ao criar evento novo:

- use nome estavel;
- inclua versao;
- declare tenant/recurso quando fizer sentido;
- mantenha payload pequeno e objetivo;
- nao use evento como atalho para vazar tabela interna.

## Ciclo de Vida de um Contrato

1. Necessidade publica aparece em um fluxo real.
2. O servico dono e definido.
3. O OpenAPI ou JSON Schema e alterado.
4. A implementacao acompanha o contrato.
5. Teste unitario/contrato/smoke cobre o comportamento proporcionalmente ao risco.
6. Registry e documentacao sao atualizados.
7. O console tecnico pode regenerar catalogo quando a mudanca for HTTP.

Esse ciclo evita que contrato vire documentacao atrasada.

## Politica de Versao

Enquanto o projeto esta em consolidacao interna, os contratos usam versoes de artefato e changelog do repositorio. Quando houver consumidor externo real, mudancas breaking devem ser tratadas com uma destas estrategias:

- criar nova versao de endpoint;
- aceitar os dois formatos por janela de migracao;
- adicionar campo novo e manter campo antigo como deprecated;
- publicar guia de migracao;
- bloquear a mudanca ate consumidor critico ser atualizado.

## Exemplos de Compatibilidade

### Campo opcional novo

Compativel:

```json
{
  "publicId": "sub_123",
  "status": "active",
  "currentPeriodEndsAt": "2026-05-31T23:59:59Z"
}
```

Adicionar `cancellationReason` opcional nao quebra consumidores que ignoram campos desconhecidos.

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

Pode ser compativel se consumidores tiverem fallback para valor desconhecido. Sem fallback, pode quebrar UI, automacao ou report.

## Regras Para Eventos

Eventos devem representar fatos:

- `fiscal.document.cancelled`;
- `platform-control.lifecycle-job.completed`;
- `webhook-hub.outbound-delivery.failed`.

Evite eventos que sejam comandos disfarçados:

- `pleaseCreateInvoice`;
- `updateCustomerSomewhere`;
- `syncEverythingNow`.

Se a intencao e comandar um servico, use HTTP, workflow ou fila de comando explicitamente modelada.

## Contrato e Teste

| Mudanca | Teste minimo esperado |
|---------|-----------------------|
| descricao ou exemplo | contract validation |
| campo opcional novo | contract validation e teste local se usado |
| endpoint novo | teste de handler/regra e contract validation |
| mutacao sensivel | unidade, contract e smoke se afetar fluxo cross-service |
| evento novo | schema validation e teste do produtor |
| breaking change | plano de migracao e validacao de consumidores |

## Contrato e Documentacao

O contrato deve ser detalhado o bastante para maquina e consumidor tecnico. A documentacao humana deve explicar:

- quando usar o endpoint;
- qual o dono funcional;
- quais riscos existem;
- qual suite valida a mudanca.

Nao copie todo o OpenAPI para markdown. Isso gera duplicidade e envelhece rapido.
