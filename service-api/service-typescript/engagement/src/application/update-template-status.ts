import { TemplateRepository } from "../domain/template-repository.js";
import { ensureTemplateStatus } from "../domain/template.js";
import { ensurePublicId } from "../domain/touchpoint.js";

export class UpdateTemplateStatus {
  constructor(private readonly repository: TemplateRepository) {}

  async execute(publicId: string, status: string) {
    return this.repository.updateStatus(
      ensurePublicId(publicId, "template_public_id_invalid"),
      ensureTemplateStatus(status)
    );
  }
}
