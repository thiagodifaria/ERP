# catalog

Contexto de catalogo para produtos e servicos do ERP.

## Escopo atual

- categorias por tenant
- itens de catalogo com SKU, unidade, preco base, tipo e atributos
- ativacao e desativacao de itens
- leitura por `publicId`
- paginação por cursor para categorias e itens
- criação em lote com `partial success`
- persistencia PostgreSQL ou bootstrap em memoria

## Rotas publicas

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/catalog/categories`
- `GET /api/catalog/categories/page`
- `POST /api/catalog/categories`
- `GET /api/catalog/items`
- `GET /api/catalog/items/page`
- `POST /api/catalog/items`
- `POST /api/catalog/items/bulk`
- `GET /api/catalog/items/{publicId}`
- `PATCH /api/catalog/items/{publicId}`
