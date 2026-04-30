import { TemplateRepository } from "../domain/template-repository.js";
import {
  CreateTemplateInput,
  ensureTemplateBody,
  ensureTemplateChannel,
  ensureTemplateKey,
  ensureTemplateName,
  ensureTemplateProvider,
  ensureTemplateStatus,
  ensureTemplateSubject
} from "../domain/template.js";
import { ensureTouchpointText } from "../domain/touchpoint.js";

export class CreateTemplate {
  constructor(private readonly repository: TemplateRepository) {}

  async execute(input: CreateTemplateInput) {
    return this.repository.create({
      tenantSlug: ensureTouchpointText(input.tenantSlug, "tenant_slug_required"),
      key: ensureTemplateKey(input.key),
      name: ensureTemplateName(input.name),
      channel: ensureTemplateChannel(input.channel),
      status: input.status ? ensureTemplateStatus(input.status) : "draft",
      provider: ensureTemplateProvider(input.provider),
      subject: ensureTemplateSubject(input.subject),
      body: ensureTemplateBody(input.body)
    });
  }
}
