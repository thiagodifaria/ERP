// LeadRepository define a persistencia minima do agregado de lead.
// Regras de negocio devem depender desta abstracao.
package repository

import "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"

type LeadRepository interface {
  List() []entity.Lead
  FindByPublicID(publicID string) *entity.Lead
  FindByEmail(email string) *entity.Lead
  Save(lead entity.Lead) entity.Lead
  Update(lead entity.Lead) entity.Lead
}
