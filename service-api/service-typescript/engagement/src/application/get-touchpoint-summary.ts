import { TouchpointFilters } from "../domain/touchpoint.js";
import { TouchpointRepository } from "../domain/touchpoint-repository.js";

export class GetTouchpointSummary {
  constructor(private readonly repository: TouchpointRepository) {}

  async execute(filters?: TouchpointFilters) {
    return this.repository.getSummary(filters);
  }
}
