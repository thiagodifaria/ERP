import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "node:path";

const serviceTargets: Record<string, string> = {
  edge: "http://localhost:8080",
  identity: "http://localhost:8081",
  "webhook-hub": "http://localhost:8082",
  crm: "http://localhost:8083",
  "workflow-control": "http://localhost:8084",
  "workflow-runtime": "http://localhost:8085",
  analytics: "http://localhost:8086",
  sales: "http://localhost:8087",
  engagement: "http://localhost:8088",
  finance: "http://localhost:8092",
  documents: "http://localhost:8093",
  simulation: "http://localhost:8094",
  billing: "http://localhost:8095",
  rentals: "http://localhost:8096",
  catalog: "http://localhost:8097",
  "platform-control": "http://localhost:8098",
  support: "http://localhost:8099",
  supplier: "http://localhost:8100",
  notification: "http://localhost:8101",
  fiscal: "http://localhost:8102",
  accounting: "http://localhost:8103",
  inventory: "http://localhost:8104",
  procurement: "http://localhost:8105",
  banking: "http://localhost:8106"
};

export default defineConfig({
  plugins: [react()],
  server: {
    fs: {
      allow: [path.resolve(__dirname)]
    },
    proxy: Object.fromEntries(
      Object.entries(serviceTargets).map(([service, target]) => [
        `/__erp/${service}`,
        {
          target,
          changeOrigin: true,
          rewrite: (requestPath: string) => requestPath.replace(`/__erp/${service}`, "")
        }
      ])
    )
  }
});
