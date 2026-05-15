import { useEffect, useMemo, useState } from "react";
import {
  Activity,
  AlertCircle,
  BookOpen,
  Box,
  Check,
  CheckCircle2,
  ChevronRight,
  Copy,
  Database,
  FileJson,
  Hash,
  Key,
  Layout,
  Map,
  Play,
  RefreshCw,
  Search,
  Server,
  Terminal
} from "lucide-react";
import { docs, markdownHeadings, type DocArticle } from "./data/docs";
import {
  defaultBodyFor,
  defaultEnvironment,
  servicesList,
  stringify,
  type RequestResult,
  type RuntimeEnvironment
} from "./data/runtime";
import { curlFor, sendEndpointRequest } from "./lib/httpClient";
import { endpoints, eventSchemas, services, type EndpointContract } from "./generated/apiCatalog";

type Page = "home" | "api" | "docs" | "contracts" | "env" | "journeys" | "ops";
type Method = EndpointContract["method"];

const storageKey = "erp-control-console-env";
const tokenStorageKey = "erp-control-console-token";

function persistedEnvironment(environment: RuntimeEnvironment): RuntimeEnvironment {
  return { ...environment, bearerToken: "" };
}

function loadEnvironment(): RuntimeEnvironment {
  try {
    const stored = localStorage.getItem(storageKey);
    const bearerToken = sessionStorage.getItem(tokenStorageKey) ?? "";
    return stored ? { ...defaultEnvironment, ...JSON.parse(stored), bearerToken } : { ...defaultEnvironment, bearerToken };
  } catch {
    return defaultEnvironment;
  }
}

function saveEnvironment(environment: RuntimeEnvironment): void {
  localStorage.setItem(storageKey, JSON.stringify(persistedEnvironment(environment)));
  if (environment.bearerToken.trim()) {
    sessionStorage.setItem(tokenStorageKey, environment.bearerToken.trim());
  } else {
    sessionStorage.removeItem(tokenStorageKey);
  }
}

function MethodBadge({ method }: { method: string }) {
  const colors: Record<string, string> = {
    GET: "bg-blue-50 text-blue-700 border-blue-200",
    POST: "bg-green-50 text-green-700 border-green-200",
    PUT: "bg-amber-50 text-amber-700 border-amber-200",
    PATCH: "bg-purple-50 text-purple-700 border-purple-200",
    DELETE: "bg-red-50 text-red-700 border-red-200",
    EVENT: "bg-slate-50 text-slate-700 border-slate-200"
  };

  return (
    <span className={`px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider rounded border ${colors[method] ?? colors.EVENT}`}>
      {method}
    </span>
  );
}

function Topbar({ currentPath, setCurrentPath }: { currentPath: Page; setCurrentPath: (path: Page) => void }) {
  const navItems = [
    { id: "home", label: "Overview", icon: Layout },
    { id: "api", label: "API Explorer", icon: Terminal },
    { id: "docs", label: "Documentação", icon: BookOpen },
    { id: "contracts", label: "Contratos", icon: FileJson },
    { id: "env", label: "Ambientes", icon: Server },
    { id: "journeys", label: "Jornadas", icon: Map },
    { id: "ops", label: "Operações", icon: Activity }
  ] as const;

  return (
    <header className="bg-white border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-[1600px] mx-auto px-6 h-16 flex items-center justify-between">
        <button className="flex items-center gap-2 group" onClick={() => setCurrentPath("home")} type="button">
          <div className="w-8 h-8 bg-blue-600 rounded flex items-center justify-center text-white group-hover:bg-blue-700 transition-colors">
            <Database size={18} />
          </div>
          <span className="font-semibold text-gray-900 tracking-tight text-lg">
            ERP <span className="text-gray-400 font-light px-1">/</span> Control Console
          </span>
        </button>

        <nav className="hidden lg:flex items-center gap-1">
          {navItems.map((item) => {
            const Icon = item.icon;
            const isActive = currentPath === item.id;
            return (
              <button
                key={item.id}
                onClick={() => setCurrentPath(item.id)}
                className={`flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-all duration-200 ${
                  isActive ? "bg-blue-50 text-blue-700" : "text-gray-500 hover:text-gray-900 hover:bg-gray-50"
                }`}
                type="button"
              >
                <Icon size={16} className={isActive ? "text-blue-600" : "text-gray-400"} />
                {item.label}
              </button>
            );
          })}
        </nav>
      </div>
    </header>
  );
}

function HomePage({ setCurrentPath }: { setCurrentPath: (path: Page) => void }) {
  const totalEndpoints = endpoints.length;

  return (
    <div className="max-w-[1200px] mx-auto px-6 py-12">
      <section className="mb-16">
        <h1 className="text-4xl font-semibold text-gray-900 tracking-tight mb-4">Console Técnico Unificado</h1>
        <p className="text-xl text-gray-500 font-light max-w-2xl mb-8 leading-relaxed">
          O ambiente para desenvolvedores, QA e operações interagirem com o ecossistema modular do ERP. Teste APIs reais,
          leia documentação versionada, valide contratos e monitore jornadas em um só lugar.
        </p>
        <div className="flex flex-wrap items-center gap-4">
          <button onClick={() => setCurrentPath("api")} className="bg-blue-600 text-white px-6 py-2.5 rounded-lg font-medium hover:bg-blue-700 transition-colors shadow-sm flex items-center gap-2" type="button">
            <Terminal size={18} />
            Testar API
          </button>
          <button onClick={() => setCurrentPath("docs")} className="bg-white text-gray-700 border border-gray-300 px-6 py-2.5 rounded-lg font-medium hover:bg-gray-50 transition-colors shadow-sm flex items-center gap-2" type="button">
            <BookOpen size={18} />
            Ler Documentação
          </button>
          <button onClick={() => setCurrentPath("contracts")} className="bg-white text-gray-700 border border-gray-300 px-6 py-2.5 rounded-lg font-medium hover:bg-gray-50 transition-colors shadow-sm flex items-center gap-2" type="button">
            <FileJson size={18} />
            Validar Contratos
          </button>
        </div>
      </section>

      <section className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-16">
        {[
          { title: "API Explorer", desc: "Interface interativa para rotas geradas dos contratos OpenAPI.", icon: Terminal, path: "api" as Page },
          { title: "Documentação Unificada", desc: "Guias de arquitetura, serviços, integração e operações.", icon: BookOpen, path: "docs" as Page },
          { title: "Contratos e Schemas", desc: "Governança de OpenAPI, eventos e cobertura por serviço.", icon: FileJson, path: "contracts" as Page },
          { title: "Jornadas de Teste", desc: "Fluxos de negócio atravessando múltiplos domínios.", icon: Map, path: "journeys" as Page }
        ].map((card) => {
          const Icon = card.icon;
          return (
            <button key={card.title} onClick={() => setCurrentPath(card.path)} className="text-left group bg-white border border-gray-200 rounded-2xl p-6 hover:shadow-md hover:border-blue-200 transition-all duration-300" type="button">
              <div className="w-10 h-10 bg-gray-50 group-hover:bg-blue-50 rounded-lg flex items-center justify-center mb-4 transition-colors">
                <Icon size={20} className="text-gray-600 group-hover:text-blue-600" />
              </div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">{card.title}</h3>
              <p className="text-sm text-gray-500 leading-relaxed">{card.desc}</p>
            </button>
          );
        })}
      </section>

      <section className="grid grid-cols-1 lg:grid-cols-3 gap-12 mb-16">
        <div>
          <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-6">Visão Geral da Plataforma</h3>
          <div className="flex flex-col gap-4">
            {[
              { label: "Serviços Ativos", value: services.length, icon: Box },
              { label: "Endpoints Testáveis", value: totalEndpoints, icon: Terminal },
              { label: "Schemas de Eventos", value: eventSchemas.length, icon: Hash },
              { label: "Contratos Incompletos", value: 3, icon: AlertCircle, warn: true }
            ].map((metric) => (
              <div key={metric.label} className="flex items-center justify-between p-4 bg-gray-50 rounded-xl border border-gray-100">
                <div className="flex items-center gap-3">
                  <metric.icon size={18} className={metric.warn ? "text-amber-500" : "text-gray-400"} />
                  <span className="text-sm font-medium text-gray-700">{metric.label}</span>
                </div>
                <span className={`text-lg font-semibold ${metric.warn ? "text-amber-600" : "text-gray-900"}`}>{metric.value}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="lg:col-span-2">
          <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-6">Domínios do ERP</h3>
          <div className="flex flex-wrap gap-2">
            {servicesList.map((service) => (
              <span key={service} className="px-3 py-1.5 bg-white border border-gray-200 text-gray-600 text-sm font-medium rounded-full shadow-sm hover:border-gray-300 transition-colors">
                {service}
              </span>
            ))}
          </div>
        </div>
      </section>

      <section>
        <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-6">Como Iniciar</h3>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          {[
            { step: "01", title: "Configure o Ambiente", desc: "Defina URLs, token, tenant e modo proxy/direct." },
            { step: "02", title: "Escolha um Endpoint", desc: "Filtre por serviço, método ou texto no API Explorer." },
            { step: "03", title: "Envie a Request", desc: "Use body editável, path params e query params reais." },
            { step: "04", title: "Inspecione", desc: "Leia status, headers, body e cURL da chamada executada." }
          ].map((step) => (
            <div key={step.step}>
              <div className="text-3xl font-light text-gray-200 mb-2">{step.step}</div>
              <h4 className="text-base font-medium text-gray-900 mb-1">{step.title}</h4>
              <p className="text-sm text-gray-500">{step.desc}</p>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}

function ApiExplorerPage({ environment }: { environment: RuntimeEnvironment }) {
  const [searchTerm, setSearchTerm] = useState("");
  const [methodFilter, setMethodFilter] = useState<"ALL" | Method>("ALL");
  const [selectedEndpoint, setSelectedEndpoint] = useState<EndpointContract>(endpoints[0]);
  const [activeTab, setActiveTab] = useState<"parameters" | "headers" | "body">("body");
  const [bodyText, setBodyText] = useState(defaultBodyFor(endpoints[0], environment.tenantSlug));
  const [queryText, setQueryText] = useState("");
  const [pathValues, setPathValues] = useState<Record<string, string>>({});
  const [response, setResponse] = useState<RequestResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [history, setHistory] = useState<RequestResult[]>([]);

  useEffect(() => {
    setBodyText(defaultBodyFor(selectedEndpoint, environment.tenantSlug));
    setQueryText(selectedEndpoint.path.startsWith("/api/") && selectedEndpoint.method === "GET" ? `tenantSlug=${environment.tenantSlug}` : "");
    setPathValues(Object.fromEntries(selectedEndpoint.pathParams.map((param) => [param, param.toLowerCase().includes("tenant") ? environment.tenantSlug : ""])));
    setResponse(null);
  }, [selectedEndpoint, environment.tenantSlug]);

  const filteredEndpoints = useMemo(() => {
    const query = searchTerm.toLowerCase();
    return endpoints.filter((endpoint) => {
      const matchesSearch = `${endpoint.path} ${endpoint.service} ${endpoint.method} ${endpoint.description}`.toLowerCase().includes(query);
      const matchesMethod = methodFilter === "ALL" || endpoint.method === methodFilter;
      return matchesSearch && matchesMethod;
    });
  }, [methodFilter, searchTerm]);

  const groupedEndpoints = useMemo(() => {
    const groups: Record<string, EndpointContract[]> = {};
    for (const endpoint of filteredEndpoints) {
      groups[endpoint.service] ??= [];
      groups[endpoint.service].push(endpoint);
    }
    return groups;
  }, [filteredEndpoints]);

  async function handleSend() {
    setLoading(true);
    const result = await sendEndpointRequest({ endpoint: selectedEndpoint, environment, bodyText, pathValues, queryText });
    setResponse(result);
    setHistory((current) => [result, ...current].slice(0, 8));
    setLoading(false);
  }

  const requestCurl = curlFor({ endpoint: selectedEndpoint, environment, bodyText, pathValues, queryText });

  return (
    <div className="flex h-[calc(100vh-64px)] bg-white overflow-hidden">
      <aside className="w-80 border-r border-gray-200 flex flex-col bg-gray-50/50">
        <div className="p-4 border-b border-gray-200 bg-white space-y-3">
          <div className="relative">
            <Search className="absolute left-3 top-2.5 text-gray-400" size={16} />
            <input type="text" placeholder="Buscar endpoint..." className="w-full pl-9 pr-4 py-2 bg-gray-100 border-transparent rounded-lg text-sm focus:bg-white focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-all outline-none" value={searchTerm} onChange={(event) => setSearchTerm(event.target.value)} />
          </div>
          <div className="flex flex-wrap gap-1">
            {(["ALL", "GET", "POST", "PUT", "PATCH", "DELETE"] as const).map((method) => (
              <button key={method} onClick={() => setMethodFilter(method)} className={`px-2 py-1 rounded border text-[11px] font-semibold ${methodFilter === method ? "bg-blue-600 text-white border-blue-600" : "bg-white text-gray-500 border-gray-200 hover:text-gray-900"}`} type="button">
                {method}
              </button>
            ))}
          </div>
        </div>
        <div className="flex-1 overflow-y-auto p-4 space-y-6 console-scrollbar">
          {Object.entries(groupedEndpoints).map(([service, serviceEndpoints]) => (
            <div key={service}>
              <h4 className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3 ml-2 flex items-center gap-2">
                <Box size={12} /> {service}
              </h4>
              <div className="space-y-1">
                {serviceEndpoints.map((endpoint) => (
                  <button key={endpoint.id} onClick={() => setSelectedEndpoint(endpoint)} className={`w-full text-left flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors ${selectedEndpoint.id === endpoint.id ? "bg-blue-50 text-blue-900 shadow-sm border border-blue-100" : "text-gray-600 hover:bg-gray-100 hover:text-gray-900 border border-transparent"}`} type="button">
                    <MethodBadge method={endpoint.method} />
                    <span className="truncate font-mono text-[13px]">{endpoint.path}</span>
                  </button>
                ))}
              </div>
            </div>
          ))}
        </div>
      </aside>

      <main className="flex-1 flex flex-col overflow-hidden bg-white">
        <div className="p-8 border-b border-gray-100">
          <div className="flex items-center gap-3 mb-4">
            <span className="px-2.5 py-1 bg-gray-100 text-gray-600 text-xs font-medium rounded-full border border-gray-200">{selectedEndpoint.service}</span>
            <span className="text-gray-300">•</span>
            <span className="text-sm text-gray-500">{selectedEndpoint.tag}</span>
            <span className="text-gray-300">•</span>
            <span className="text-xs text-gray-400">{selectedEndpoint.source}</span>
          </div>
          <h2 className="text-2xl font-semibold text-gray-900 mb-2">{selectedEndpoint.description}</h2>
          <div className="flex items-center gap-4 bg-gray-50 p-3 rounded-lg border border-gray-200 font-mono text-sm mt-6">
            <MethodBadge method={selectedEndpoint.method} />
            <span className="text-gray-700 break-all">{environment.mode === "proxy" ? `/__erp/${selectedEndpoint.service}` : environment.baseUrls[selectedEndpoint.service]}{selectedEndpoint.path}</span>
            <button className="ml-auto text-gray-400 hover:text-gray-600" title="Copiar cURL" onClick={() => navigator.clipboard?.writeText(requestCurl)} type="button">
              <Copy size={16} />
            </button>
          </div>
        </div>

        <div className="flex-1 overflow-y-auto flex flex-col xl:flex-row">
          <section className="flex-1 p-8 xl:border-r border-gray-100">
            <div className="flex items-center gap-6 border-b border-gray-200 mb-6">
              {(["parameters", "headers", "body"] as const).map((tab) => (
                <button key={tab} onClick={() => setActiveTab(tab)} className={`pb-3 text-sm font-medium border-b-2 transition-colors ${activeTab === tab ? "border-blue-600 text-blue-600" : "border-transparent text-gray-500 hover:text-gray-900"}`} type="button">
                  {tab === "parameters" ? "Parameters" : tab === "headers" ? "Headers" : "Body"}
                </button>
              ))}
            </div>

            {activeTab === "parameters" && (
              <div className="space-y-5">
                <label className="block">
                  <span className="text-sm font-medium text-gray-700">Query string</span>
                  <input className="mt-2 w-full border border-gray-200 rounded-lg p-3 font-mono text-sm outline-none focus:border-blue-500" value={queryText} onChange={(event) => setQueryText(event.target.value)} placeholder={`tenantSlug=${environment.tenantSlug}`} />
                </label>
                {selectedEndpoint.pathParams.length ? selectedEndpoint.pathParams.map((param) => (
                  <label className="block" key={param}>
                    <span className="text-sm font-medium text-gray-700">{param}</span>
                    <input className="mt-2 w-full border border-gray-200 rounded-lg p-3 font-mono text-sm outline-none focus:border-blue-500" value={pathValues[param] ?? ""} onChange={(event) => setPathValues((current) => ({ ...current, [param]: event.target.value }))} placeholder={`valor para {${param}}`} />
                  </label>
                )) : <div className="text-sm text-gray-500 p-4 bg-gray-50 rounded-lg border border-dashed border-gray-200">Este path não declara parâmetros.</div>}
              </div>
            )}

            {activeTab === "headers" && (
              <pre className="bg-gray-50 border border-gray-200 rounded-lg p-4 font-mono text-sm text-gray-700 overflow-auto">{stringify({
                "content-type": "application/json",
                authorization: environment.bearerToken ? "Bearer ***" : undefined,
                "x-correlation-id": environment.correlationId,
                "idempotency-key": selectedEndpoint.method === "GET" ? undefined : environment.idempotencyKey
              })}</pre>
            )}

            {activeTab === "body" && (
              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-sm font-medium text-gray-700">JSON Body</span>
                  <button className="text-xs text-blue-600 hover:underline" onClick={() => setBodyText(defaultBodyFor(selectedEndpoint, environment.tenantSlug))} type="button">Restaurar exemplo</button>
                </div>
                <textarea className="w-full h-64 p-4 bg-gray-50 border border-gray-200 rounded-lg font-mono text-sm text-gray-800 focus:ring-2 focus:ring-blue-200 focus:border-blue-500 outline-none resize-none" value={bodyText} onChange={(event) => setBodyText(event.target.value)} disabled={!selectedEndpoint.hasBody} placeholder={!selectedEndpoint.hasBody ? "Não aplicável para este método" : "Insira o JSON aqui..."} />
              </div>
            )}

            <div className="mt-8 pt-6 border-t border-gray-100">
              <button onClick={handleSend} disabled={loading} className="w-full bg-blue-600 text-white py-2.5 rounded-lg font-medium hover:bg-blue-700 transition-colors flex items-center justify-center gap-2 disabled:opacity-70" type="button">
                {loading ? <RefreshCw size={18} className="animate-spin" /> : <Play size={18} />}
                {loading ? "Enviando..." : "Enviar Requisição"}
              </button>
            </div>
          </section>

          <section className="flex-1 bg-gray-50 p-8 flex flex-col min-h-[520px]">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-semibold text-gray-700">Response Console</h3>
              {response && (
                <div className="flex items-center gap-4 text-xs font-mono">
                  <span className={`font-bold px-2 py-1 rounded ${response.status >= 200 && response.status < 300 ? "text-green-600 bg-green-100" : "text-red-600 bg-red-100"}`}>{response.status || "ERR"} {response.statusText}</span>
                  <span className="text-gray-500">{response.durationMs}ms</span>
                </div>
              )}
            </div>

            <div className="flex-1 bg-white border border-gray-200 rounded-lg relative overflow-hidden flex flex-col">
              {!response && !loading ? (
                <div className="flex-1 flex flex-col items-center justify-center text-gray-400">
                  <Terminal size={32} className="mb-3 opacity-50" />
                  <p className="text-sm">Clique em Enviar para ver a resposta real</p>
                </div>
              ) : loading ? (
                <div className="flex-1 flex items-center justify-center">
                  <RefreshCw size={24} className="text-blue-500 animate-spin" />
                </div>
              ) : (
                <>
                  <div className="bg-gray-100 px-4 py-2 border-b border-gray-200 flex justify-between gap-3">
                    <span className="text-xs text-gray-500 font-mono truncate">{response?.url}</span>
                    <button className="text-xs text-gray-500 hover:text-gray-900 flex items-center gap-1" onClick={() => navigator.clipboard?.writeText(stringify(response?.body))} type="button">
                      <Copy size={14} /> Copiar
                    </button>
                  </div>
                  <pre className="p-4 text-sm font-mono text-gray-800 overflow-auto flex-1 bg-gray-50 console-scrollbar">{stringify(response?.body)}</pre>
                </>
              )}
            </div>

            <div className="mt-5">
              <h4 className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2">Histórico</h4>
              <div className="space-y-2 max-h-32 overflow-auto console-scrollbar">
                {history.map((item, index) => (
                  <div key={`${item.url}-${index}`} className="text-xs bg-white border border-gray-200 rounded p-2 flex justify-between gap-3">
                    <span className="truncate font-mono">{item.url}</span>
                    <span className={item.status >= 200 && item.status < 300 ? "text-green-600" : "text-red-600"}>{item.status || "ERR"}</span>
                  </div>
                ))}
              </div>
            </div>
          </section>
        </div>
      </main>
    </div>
  );
}

function MarkdownBlock({ article }: { article: DocArticle }) {
  const lines = article.content.split(/\r?\n/).slice(0, 220);
  return (
    <div className="max-w-3xl">
      {lines.map((line, index) => {
        if (line.startsWith("# ")) return <h1 key={index} className="text-3xl font-bold text-gray-900 mb-6">{line.replace(/^#\s+/, "")}</h1>;
        if (line.startsWith("## ")) return <h2 key={index} className="text-2xl font-semibold text-gray-900 mt-10 mb-4 pb-2 border-b border-gray-100">{line.replace(/^##\s+/, "").replace(/`/g, "")}</h2>;
        if (line.startsWith("### ")) return <h3 key={index} className="text-xl font-medium text-gray-900 mt-8 mb-3">{line.replace(/^###\s+/, "").replace(/`/g, "")}</h3>;
        if (line.startsWith("- ")) return <p key={index} className="text-gray-600 mb-2 pl-4">• {line.replace(/^-\s+/, "")}</p>;
        if (line.trim().startsWith("|")) return <pre key={index} className="text-xs bg-gray-50 border border-gray-100 px-3 py-1 overflow-x-auto font-mono text-gray-600">{line}</pre>;
        if (!line.trim()) return <div key={index} className="h-3" />;
        return <p key={index} className="text-gray-600 mb-4 leading-relaxed">{line.replace(/`/g, "")}</p>;
      })}
    </div>
  );
}

function DocumentationPage() {
  const [activeId, setActiveId] = useState(docs[0].id);
  const active = docs.find((doc) => doc.id === activeId) ?? docs[0];
  const headings = markdownHeadings(active.content);

  return (
    <div className="flex h-[calc(100vh-64px)] bg-white overflow-hidden">
      <aside className="w-64 border-r border-gray-200 bg-gray-50/50 flex flex-col p-4">
        <h3 className="text-xs font-bold text-gray-400 uppercase tracking-wider mb-4">Conteúdo</h3>
        <div className="space-y-1">
          {docs.map((doc) => (
            <button key={doc.id} onClick={() => setActiveId(doc.id)} className={`block w-full text-left px-3 py-2 rounded-md text-sm ${doc.id === active.id ? "bg-blue-50 text-blue-700 font-medium" : "text-gray-600 hover:bg-gray-100"}`} type="button">
              {doc.title}
            </button>
          ))}
        </div>
      </aside>
      <main className="flex-1 overflow-y-auto p-12 lg:px-24 console-scrollbar">
        <MarkdownBlock article={active} />
      </main>
      <aside className="hidden xl:block w-72 p-8 border-l border-gray-100 overflow-auto console-scrollbar">
        <h4 className="text-xs font-semibold text-gray-900 uppercase tracking-wider mb-4">Nesta página</h4>
        <p className="text-xs text-gray-400 mb-5 font-mono">{active.source}</p>
        <ul className="space-y-3 text-sm text-gray-500">
          {headings.map((heading) => <li key={heading} className="hover:text-gray-900 cursor-pointer">{heading}</li>)}
        </ul>
      </aside>
    </div>
  );
}

function ContractsPage() {
  const [tab, setTab] = useState<"http" | "events" | "coverage">("http");
  const partial = new Set(["crm", "sales", "billing", "finance", "identity"]);

  return (
    <div className="max-w-[1200px] mx-auto px-6 py-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">Governança de Contratos</h1>
        <p className="text-gray-500">Validação entre implementação, OpenAPI, schemas de eventos e registry.</p>
      </div>
      <div className="flex items-center gap-6 border-b border-gray-200 mb-8">
        {[
          ["http", "HTTP OpenAPI"],
          ["events", "Event Schemas"],
          ["coverage", "Coverage"]
        ].map(([id, label]) => (
          <button key={id} onClick={() => setTab(id as typeof tab)} className={`pb-3 text-sm font-medium border-b-2 ${tab === id ? "border-blue-600 text-blue-600" : "border-transparent text-gray-500"}`} type="button">
            {label}
          </button>
        ))}
      </div>

      {tab === "http" && (
        <div className="bg-white border border-gray-200 rounded-xl shadow-sm overflow-hidden mb-8">
          <table className="w-full text-left text-sm">
            <thead className="bg-gray-50 border-b border-gray-200 text-gray-600">
              <tr>
                <th className="px-6 py-4 font-medium">Serviço</th>
                <th className="px-6 py-4 font-medium">Versão</th>
                <th className="px-6 py-4 font-medium">Endpoints</th>
                <th className="px-6 py-4 font-medium">Status</th>
                <th className="px-6 py-4 font-medium">Contrato</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {services.map((service) => (
                <tr key={service.slug} className="hover:bg-gray-50 transition-colors">
                  <td className="px-6 py-4 font-medium text-gray-900 flex items-center gap-2"><Box size={16} className="text-gray-400" /> {service.slug}</td>
                  <td className="px-6 py-4 text-gray-600 font-mono text-xs">{service.version}</td>
                  <td className="px-6 py-4 text-gray-600">{service.endpointCount} declarados</td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium ${partial.has(service.slug) ? "bg-amber-50 text-amber-700" : "bg-green-50 text-green-700"}`}>
                      {partial.has(service.slug) ? <AlertCircle size={14} /> : <CheckCircle2 size={14} />}
                      {partial.has(service.slug) ? "Atenção" : "OK"}
                    </span>
                    {partial.has(service.slug) && <p className="text-xs text-gray-400 mt-1">Auditoria indicou rotas implementadas além do OpenAPI atual.</p>}
                  </td>
                  <td className="px-6 py-4 text-xs text-blue-600 font-mono">{service.contractFile}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {tab === "events" && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {eventSchemas.map((schema) => (
            <div key={schema.fileName} className="bg-white border border-gray-200 rounded-xl p-5 shadow-sm">
              <FileJson className="text-blue-600 mb-4" size={22} />
              <h3 className="font-semibold text-gray-900">{schema.name}</h3>
              <p className="text-sm text-gray-500 font-mono mt-2">{schema.source}</p>
            </div>
          ))}
        </div>
      )}

      {tab === "coverage" && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-white border border-gray-200 rounded-xl p-6 shadow-sm"><div className="text-3xl font-light">{services.length}</div><p className="text-sm text-gray-500">Contratos HTTP</p></div>
          <div className="bg-white border border-gray-200 rounded-xl p-6 shadow-sm"><div className="text-3xl font-light">{endpoints.length}</div><p className="text-sm text-gray-500">Endpoints disponíveis no console</p></div>
          <div className="bg-white border border-amber-200 rounded-xl p-6 shadow-sm bg-amber-50/40"><div className="text-3xl font-light text-amber-600">{partial.size}</div><p className="text-sm text-amber-700">Serviços com cobertura parcial conhecida</p></div>
        </div>
      )}
    </div>
  );
}

function EnvironmentsPage({ environment, setEnvironment }: { environment: RuntimeEnvironment; setEnvironment: (env: RuntimeEnvironment) => void }) {
  function updateBaseUrl(service: string, value: string) {
    setEnvironment({ ...environment, baseUrls: { ...environment.baseUrls, [service]: value } });
  }

  return (
    <div className="max-w-[1000px] mx-auto px-6 py-8">
      <div className="flex justify-between items-end mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Ambientes e Contexto</h1>
          <p className="text-gray-500">Configure variáveis globais e base URLs para testes no console.</p>
        </div>
        <div className="flex gap-3">
          <button className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50" onClick={() => setEnvironment(defaultEnvironment)} type="button">Resetar</button>
          <button className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 flex items-center gap-2" onClick={() => saveEnvironment(environment)} type="button"><Check size={16} /> Salvar</button>
        </div>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        <section className="space-y-6">
          <div className="bg-white border border-gray-200 p-5 rounded-xl shadow-sm">
            <h3 className="font-semibold text-gray-900 mb-4 flex items-center gap-2"><Server size={18} className="text-gray-400" /> Perfil Ativo</h3>
            <select className="w-full border border-gray-300 rounded-lg p-2 text-sm bg-gray-50 outline-none focus:border-blue-500" value={environment.mode} onChange={(event) => setEnvironment({ ...environment, mode: event.target.value as RuntimeEnvironment["mode"] })}>
              <option value="proxy">Local Docker via Vite proxy</option>
              <option value="direct">Direct localhost</option>
            </select>
            <p className="text-xs text-gray-400 mt-3">Proxy evita CORS durante desenvolvimento.</p>
          </div>
          <div className="bg-white border border-gray-200 p-5 rounded-xl shadow-sm">
            <h3 className="font-semibold text-gray-900 mb-4 flex items-center gap-2"><Key size={18} className="text-gray-400" /> Contexto Global</h3>
            <div className="space-y-4">
              {[
                ["tenantSlug", "Tenant Slug"],
                ["bearerToken", "Bearer Token"],
                ["correlationId", "Correlation ID"],
                ["idempotencyKey", "Idempotency Key"]
              ].map(([key, label]) => (
                <label key={key} className="block">
                  <span className="block text-xs font-medium text-gray-500 mb-1">{label}</span>
                  <input type={key === "bearerToken" ? "password" : "text"} value={String(environment[key as keyof RuntimeEnvironment])} onChange={(event) => setEnvironment({ ...environment, [key]: event.target.value })} className="w-full border border-gray-300 rounded-lg p-2 text-sm outline-none focus:border-blue-500" />
                </label>
              ))}
            </div>
          </div>
        </section>

        <section className="md:col-span-2 bg-white border border-gray-200 rounded-xl shadow-sm overflow-hidden">
          <div className="p-5 border-b border-gray-200 bg-gray-50">
            <h3 className="font-semibold text-gray-900">Base URLs por Serviço</h3>
            <p className="text-xs text-gray-500 mt-1">Usadas no modo direct. No modo proxy, Vite redireciona por serviço.</p>
          </div>
          <div className="p-5 space-y-3 max-h-[560px] overflow-y-auto console-scrollbar">
            {servicesList.map((service) => (
              <div key={service} className="flex items-center gap-4">
                <span className="w-36 text-sm font-medium text-gray-700">{service}</span>
                <input type="text" value={environment.baseUrls[service] ?? ""} onChange={(event) => updateBaseUrl(service, event.target.value)} className="flex-1 border border-gray-200 rounded-md p-1.5 text-sm font-mono text-gray-600 bg-gray-50 focus:bg-white focus:border-blue-500 outline-none" />
              </div>
            ))}
          </div>
        </section>
      </div>
    </div>
  );
}

function JourneysPage() {
  const [step, setStep] = useState(3);
  const steps = [
    { step: 1, name: "Criar Cliente (CRM)", status: step > 1 ? "success" : "ready", service: "crm", method: "POST" },
    { step: 2, name: "Criar Plano Base (Billing)", status: step > 2 ? "success" : "ready", service: "billing", method: "POST" },
    { step: 3, name: "Gerar Assinatura (Billing)", status: step > 3 ? "success" : "ready", service: "billing", method: "POST", path: "/api/billing/subscriptions" },
    { step: 4, name: "Aguardar Fatura (Webhook)", status: step > 4 ? "success" : "pending", service: "webhook-hub", method: "EVENT" },
    { step: 5, name: "Liquidar Recebível (Finance)", status: step > 5 ? "success" : "pending", service: "finance", method: "POST" }
  ];

  return (
    <div className="max-w-[1000px] mx-auto px-6 py-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">Jornadas de Teste E2E</h1>
        <p className="text-gray-500">Execute fluxos de negócio completos que atravessam múltiplos domínios.</p>
      </div>
      <div className="bg-white border border-gray-200 rounded-xl shadow-sm overflow-hidden">
        <div className="p-6 border-b border-gray-100 bg-gray-50/50 flex justify-between items-center">
          <div><h3 className="text-lg font-semibold text-gray-900 mb-1">Order to Cash (Assinatura)</h3><p className="text-sm text-gray-500">CRM {"->"} Sales {"->"} Billing {"->"} Finance</p></div>
          <button className="bg-white border border-gray-300 text-gray-700 px-4 py-2 rounded-lg text-sm font-medium hover:bg-gray-50 shadow-sm" onClick={() => setStep(3)} type="button">Resetar Jornada</button>
        </div>
        <div className="p-6 relative">
          <div className="absolute left-[39px] top-10 bottom-10 w-px bg-gray-200 z-0" />
          <div className="space-y-8 relative z-10">
            {steps.map((item) => {
              const active = item.step === step;
              return (
                <div key={item.step} className="flex gap-6">
                  <div className={`w-8 h-8 rounded-full flex items-center justify-center border-2 shrink-0 bg-white ${item.status === "success" ? "border-green-500 text-green-500" : active ? "border-blue-600 text-blue-600" : "border-gray-300 text-gray-400"}`}>
                    {item.status === "success" ? <Check size={16} /> : <span className="text-sm font-medium">{item.step}</span>}
                  </div>
                  <div className={`flex-1 pt-1 ${item.status === "pending" && !active ? "opacity-50" : ""}`}>
                    <div className="flex items-center justify-between mb-2">
                      <h4 className="font-medium text-gray-900">{item.name}</h4>
                      {active && <button className="text-xs bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700 flex items-center gap-1" onClick={() => setStep((current) => Math.min(current + 1, 6))} type="button"><Play size={12} /> Executar Passo</button>}
                    </div>
                    {active && <div className="bg-gray-50 border border-gray-200 rounded-lg p-3 text-xs font-mono text-gray-600"><div className="flex items-center gap-2 mb-2 pb-2 border-b border-gray-200"><MethodBadge method={item.method} /> <span>{item.path ?? "/api/console/journey"}</span></div>{`{ "tenantSlug": "bootstrap-ops", "source": "api-console" }`}</div>}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}

function OperationsPage({ environment }: { environment: RuntimeEnvironment }) {
  const [checks, setChecks] = useState<Record<string, RequestResult>>({});
  const [loading, setLoading] = useState(false);

  async function runHealthChecks() {
    setLoading(true);
    const healthEndpoints = services.map((service) => endpoints.find((endpoint) => endpoint.service === service.slug && endpoint.path === "/health/ready")).filter(Boolean) as EndpointContract[];
    const results = await Promise.all(
      healthEndpoints.map(async (endpoint) => [endpoint.service, await sendEndpointRequest({ endpoint, environment, bodyText: "", pathValues: {}, queryText: "" })] as const)
    );
    setChecks(Object.fromEntries(results));
    setLoading(false);
  }

  return (
    <div className="max-w-[1200px] mx-auto px-6 py-8">
      <div className="flex items-end justify-between mb-8">
        <div><h1 className="text-2xl font-bold text-gray-900 mb-2">Painel Operacional</h1><p className="text-gray-500">Health, readiness, providers e fila de testes.</p></div>
        <button onClick={runHealthChecks} className="bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-blue-700 flex items-center gap-2" type="button">{loading ? <RefreshCw size={16} className="animate-spin" /> : <Activity size={16} />} Checar serviços</button>
      </div>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
        <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm"><h3 className="font-semibold text-gray-700 mb-4">Serviços prontos</h3><div className="text-3xl font-light text-gray-900">{Object.values(checks).filter((check) => check.status >= 200 && check.status < 300).length}/{services.length}</div></div>
        <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm"><h3 className="font-semibold text-gray-700 mb-4">Latência média</h3><div className="text-3xl font-light text-gray-900">{Object.values(checks).length ? Math.round(Object.values(checks).reduce((sum, item) => sum + item.durationMs, 0) / Object.values(checks).length) : "-"}ms</div></div>
        <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm"><h3 className="font-semibold text-gray-700 mb-4">Modo</h3><div className="text-3xl font-light text-gray-900">{environment.mode}</div></div>
      </div>
      <div className="bg-white border border-gray-200 rounded-xl shadow-sm overflow-hidden">
        <div className="p-5 border-b border-gray-200 bg-gray-50"><h3 className="font-semibold text-gray-900">Readiness por serviço</h3></div>
        <div className="divide-y divide-gray-100">
          {services.map((service) => {
            const check = checks[service.slug];
            const ok = check && check.status >= 200 && check.status < 300;
            return <div key={service.slug} className="p-4 flex items-center justify-between"><div><h4 className="font-medium text-gray-900 text-sm">{service.slug}</h4><p className="text-xs text-gray-500">{environment.mode === "proxy" ? `/__erp/${service.slug}/health/ready` : `${environment.baseUrls[service.slug]}/health/ready`}</p></div><span className={`inline-flex items-center gap-1 text-xs font-medium px-2 py-1 rounded ${!check ? "text-gray-500 bg-gray-100" : ok ? "text-green-700 bg-green-50" : "text-red-700 bg-red-50"}`}>{!check ? "não testado" : `${check.status || "ERR"} ${check.durationMs}ms`}</span></div>;
          })}
        </div>
      </div>
    </div>
  );
}

export default function App() {
  const [currentPath, setCurrentPath] = useState<Page>("home");
  const [environment, setEnvironment] = useState<RuntimeEnvironment>(loadEnvironment);

  useEffect(() => {
    saveEnvironment(environment);
  }, [environment]);

  return (
    <div className="min-h-screen bg-white text-gray-900 font-sans selection:bg-blue-100 selection:text-blue-900">
      <Topbar currentPath={currentPath} setCurrentPath={setCurrentPath} />
      <main>
        {currentPath === "home" && <HomePage setCurrentPath={setCurrentPath} />}
        {currentPath === "api" && <ApiExplorerPage environment={environment} />}
        {currentPath === "docs" && <DocumentationPage />}
        {currentPath === "contracts" && <ContractsPage />}
        {currentPath === "env" && <EnvironmentsPage environment={environment} setEnvironment={setEnvironment} />}
        {currentPath === "journeys" && <JourneysPage />}
        {currentPath === "ops" && <OperationsPage environment={environment} />}
      </main>
    </div>
  );
}
