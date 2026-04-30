import { randomUUID } from "node:crypto";
import { TemplateRepository } from "../domain/template-repository.js";
import { CreateTemplateInput, EngagementTemplate, TemplateFilters, TemplateStatus } from "../domain/template.js";

function bootstrapTemplates(tenantSlug: string): EngagementTemplate[] {
  return [
    {
      id: 1,
      publicId: "00000000-0000-0000-0000-00000000f101",
      tenantSlug,
      key: "lead-follow-up-whatsapp",
      name: "Lead Follow-Up WhatsApp",
      channel: "whatsapp",
      status: "active",
      provider: "manual",
      subject: null,
      body: "Ola {{firstName}}, vimos seu interesse e queremos avancar com a proxima conversa.",
      createdAt: "2026-04-20T00:02:00.000Z",
      updatedAt: "2026-04-20T00:02:00.000Z"
    }
  ];
}

export class InMemoryTemplateRepository implements TemplateRepository {
  private readonly templates: EngagementTemplate[];

  constructor(tenantSlug: string) {
    this.templates = bootstrapTemplates(tenantSlug);
  }

  async list(filters: TemplateFilters = {}): Promise<EngagementTemplate[]> {
    return this.templates.filter((template) => {
      if (filters.tenantSlug && template.tenantSlug !== filters.tenantSlug) {
        return false;
      }

      if (filters.channel && template.channel !== filters.channel) {
        return false;
      }

      if (filters.status && template.status !== filters.status) {
        return false;
      }

      if (filters.provider && template.provider !== filters.provider) {
        return false;
      }

      if (filters.q) {
        const query = filters.q.toLowerCase();
        const haystack = `${template.key} ${template.name} ${template.body} ${template.subject ?? ""}`.toLowerCase();
        return haystack.includes(query);
      }

      return true;
    });
  }

  async getByPublicId(publicId: string): Promise<EngagementTemplate | null> {
    return this.templates.find((template) => template.publicId === publicId) ?? null;
  }

  async create(input: CreateTemplateInput): Promise<EngagementTemplate> {
    if (this.templates.some((template) => template.tenantSlug === input.tenantSlug && template.key === input.key)) {
      throw new Error("template_key_conflict");
    }

    const now = new Date().toISOString();
    const template: EngagementTemplate = {
      id: this.templates.length + 1,
      publicId: randomUUID(),
      tenantSlug: input.tenantSlug,
      key: input.key,
      name: input.name,
      channel: input.channel,
      status: input.status ?? "draft",
      provider: input.provider,
      subject: input.subject ?? null,
      body: input.body,
      createdAt: now,
      updatedAt: now
    };

    this.templates.push(template);
    return template;
  }

  async updateStatus(publicId: string, status: TemplateStatus): Promise<EngagementTemplate | null> {
    const template = this.templates.find((item) => item.publicId === publicId);

    if (!template) {
      return null;
    }

    template.status = status;
    template.updatedAt = new Date().toISOString();
    return template;
  }
}
