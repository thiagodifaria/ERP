import { CreateTemplateInput, EngagementTemplate, TemplateFilters, TemplateStatus } from "./template.js";

export interface TemplateRepository {
  list(filters?: TemplateFilters): Promise<EngagementTemplate[]>;
  getByPublicId(publicId: string): Promise<EngagementTemplate | null>;
  create(input: CreateTemplateInput): Promise<EngagementTemplate>;
  updateStatus(publicId: string, status: TemplateStatus): Promise<EngagementTemplate | null>;
}
