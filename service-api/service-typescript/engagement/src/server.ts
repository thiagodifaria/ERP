import { createServer } from "node:http";
import { route } from "./api/router.js";

const port = Number(process.env.ENGAGEMENT_HTTP_PORT ?? "8088");
const host = process.env.ENGAGEMENT_HTTP_HOST ?? "0.0.0.0";

const server = createServer((request, response) => {
  void route(request, response);
});

server.listen(port, host, () => {
  process.stdout.write(`engagement listening on ${host}:${port}\n`);
});
