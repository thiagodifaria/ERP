// The workflow-control bootstrap keeps runtime concerns thin and explicit.
import { createServer } from "node:http";
import { route } from "./api/router.js";

const port = Number(process.env.WORKFLOW_CONTROL_HTTP_PORT ?? "8084");
const host = process.env.WORKFLOW_CONTROL_HTTP_HOST ?? "0.0.0.0";

const server = createServer((request, response) => {
  route(request, response);
});

server.listen(port, host, () => {
  process.stdout.write(`workflow-control listening on ${host}:${port}\n`);
});
