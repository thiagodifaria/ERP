import { z } from "zod";

const workflowActionSchema = z.object({
  stepId: z.string().trim().min(1),
  actionKey: z.string().trim().min(1),
  label: z.string().trim().min(1),
  delaySeconds: z.number().int().nonnegative().nullable().default(null),
  compensationActionKey: z.string().trim().min(1).nullable().default(null)
}).strict();

export const createWorkflowRunSchema = z.object({
  workflowDefinitionKey: z.string().trim().min(1),
  subjectType: z.string().trim().min(1),
  subjectPublicId: z.string().uuid(),
  initiatedBy: z.string().trim().min(1)
}).strict();

export const createWorkflowRunEventSchema = z.object({
  body: z.string().trim().min(1),
  createdBy: z.string().trim().min(1)
}).strict();

export const createWorkflowDefinitionSchema = z.object({
  key: z.string().trim().min(1),
  name: z.string().trim().min(1),
  description: z.string().trim().nullable().optional(),
  trigger: z.string().trim().min(1),
  actions: z.array(workflowActionSchema).optional()
}).strict();

export const updateWorkflowDefinitionSchema = z.object({
  name: z.string().trim().min(1).optional(),
  description: z.string().trim().nullable().optional(),
  trigger: z.string().trim().min(1).optional(),
  actions: z.array(workflowActionSchema).optional()
}).strict();

export const updateWorkflowDefinitionStatusSchema = z.object({
  status: z.string().trim().min(1)
}).strict();
