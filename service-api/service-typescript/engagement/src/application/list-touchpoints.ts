import { TouchpointFilters } from "../domain/touchpoint.js";
import { TouchpointRepository } from "../domain/touchpoint-repository.js";

export class ListTouchpoints {
  constructor(private readonly repository: TouchpointRepository) {}

  async execute(filters?: TouchpointFilters) {
    return this.repository.list(filters);
  }
}
