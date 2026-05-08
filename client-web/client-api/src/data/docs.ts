import api from "../../../../docs/API.md?raw";
import arquitetura from "../../../../docs/ARQUITETURA.md?raw";
import contratos from "../../../../docs/CONTRATOS.md?raw";
import integracoes from "../../../../docs/INTEGRACOES.md?raw";
import operacoes from "../../../../docs/OPERACOES.md?raw";
import padroes from "../../../../docs/PADROES.md?raw";
import servicos from "../../../../docs/SERVICOS.md?raw";

export type DocArticle = {
  id: string;
  title: string;
  source: string;
  content: string;
};

export const docs: DocArticle[] = [
  { id: "api", title: "API", source: "docs/API.md", content: api },
  { id: "servicos", title: "Serviços", source: "docs/SERVICOS.md", content: servicos },
  { id: "arquitetura", title: "Arquitetura", source: "docs/ARQUITETURA.md", content: arquitetura },
  { id: "contratos", title: "Contratos", source: "docs/CONTRATOS.md", content: contratos },
  { id: "integracoes", title: "Integrações", source: "docs/INTEGRACOES.md", content: integracoes },
  { id: "operacoes", title: "Operações", source: "docs/OPERACOES.md", content: operacoes },
  { id: "padroes", title: "Padrões", source: "docs/PADROES.md", content: padroes }
];

export function markdownHeadings(markdown: string): string[] {
  return markdown
    .split(/\r?\n/)
    .filter((line) => line.startsWith("## "))
    .slice(0, 12)
    .map((line) => line.replace(/^##\s+/, "").replace(/`/g, ""));
}
