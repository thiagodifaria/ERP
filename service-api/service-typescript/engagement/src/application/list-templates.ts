import { TemplateRepository } from "../domain/template-repository.js";
import { TemplateFilters, ensureTemplateChannel, ensureTemplateProvider, ensureTemplateStatus } from "../domain/template.js";
import { ensureTouchpointText } from "../domain/touchpoint.js";

export class ListTemplates {
  constructor(private readonly repository: TemplateRepository) {}

  async execute(filters: TemplateFilters = {}) {
    return this.repository.list({
      tenantSlug: filters.tenantSlug ? ensureTouchpointText(filters.tenantSlug, "tenant_slug_required") : undefined,
      channel: filters.channel ? ensureTemplateChannel(filters.channel) : undefined,
      status: filters.status ? ensureTemplateStatus(filters.status) : undefined,
      provider: filters.provider ? ensureTemplateProvider(filters.provider) : undefined,
      q: filters.q?.trim() || undefined
    });
  }
}
