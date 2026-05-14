# STACK STARTERS

Estes templates definem o minimo esperado para qualquer novo servico do ERP.

Cada stack deve partir com:

- middleware de autenticacao JWT/service account;
- tenant resolver explicito;
- propagacao de `X-Correlation-Id` e `traceparent`;
- modelo de erro publico;
- health live/ready/details;
- contrato OpenAPI;
- teste de contrato;
- validacao de configuracao;
- Dockerfile com imagem versionada;
- hooks para logs, metricas e traces.

Os arquivos desta pasta nao sao aplicacoes completas. Eles sao bases copiaveis para manter equivalencia operacional entre stacks.
