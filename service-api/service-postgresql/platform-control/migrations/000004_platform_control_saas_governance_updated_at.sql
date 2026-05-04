DROP TRIGGER IF EXISTS trg_platform_control_quotas_updated_at ON platform_control.quotas;
CREATE TRIGGER trg_platform_control_quotas_updated_at
BEFORE UPDATE ON platform_control.quotas
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_platform_control_tenant_blocks_updated_at ON platform_control.tenant_blocks;
CREATE TRIGGER trg_platform_control_tenant_blocks_updated_at
BEFORE UPDATE ON platform_control.tenant_blocks
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
