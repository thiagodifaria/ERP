// This router starts small and grows with public workflow-control routes.
import { IncomingMessage, ServerResponse } from "node:http";

export function route(_request: IncomingMessage, response: ServerResponse): void {
  response.writeHead(404, { "content-type": "application/json" });
  response.end(JSON.stringify({
    code: "route_not_found",
    message: "Route was not found."
  }));
}
