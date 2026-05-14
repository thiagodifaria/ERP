using Finance.Api;
using Xunit;

namespace Finance.UnitTests;

public sealed class FinancePoliciesTests
{
  [Theory]
  [InlineData("OPEN", "open")]
  [InlineData(" paid ", "paid")]
  [InlineData("cancelled", "cancelled")]
  [InlineData("invalid", "")]
  public void NormalizePayableStatusShouldAcceptOnlyOperationalStates(string input, string expected)
  {
    Assert.Equal(expected, FinancePolicies.NormalizePayableStatus(input));
  }

  [Fact]
  public void PayableTransitionShouldOnlyAllowOpenToTerminalStates()
  {
    Assert.True(FinancePolicies.CanTransitionPayableStatus("open", "paid"));
    Assert.True(FinancePolicies.CanTransitionPayableStatus("open", "cancelled"));
    Assert.False(FinancePolicies.CanTransitionPayableStatus("paid", "open"));
    Assert.False(FinancePolicies.CanTransitionPayableStatus("cancelled", "paid"));
  }

  [Theory]
  [InlineData("active", "active")]
  [InlineData(" inactive ", "inactive")]
  [InlineData("", "")]
  [InlineData("blocked", "")]
  public void NormalizeCashAccountStatusShouldProtectTreasuryFilters(string input, string expected)
  {
    Assert.Equal(expected, FinancePolicies.NormalizeCashAccountStatus(input));
  }

  [Theory]
  [InlineData("inflow", "inflow")]
  [InlineData(" outflow ", "outflow")]
  [InlineData("sideways", "")]
  public void NormalizeCashMovementDirectionShouldAcceptOnlyLedgerDirections(string input, string expected)
  {
    Assert.Equal(expected, FinancePolicies.NormalizeCashMovementDirection(input));
  }

  [Theory]
  [InlineData("receivable_settlement", "receivable_settlement")]
  [InlineData("payable_payment", "payable_payment")]
  [InlineData("cost_entry", "cost_entry")]
  [InlineData("manual_adjustment", "manual_adjustment")]
  [InlineData("refund", "")]
  public void NormalizeCashMovementTypeShouldAcceptOnlyLedgerMovementTypes(string input, string expected)
  {
    Assert.Equal(expected, FinancePolicies.NormalizeCashMovementType(input));
  }

  [Fact]
  public void ParseUtcOrNowShouldRejectInvalidSettlementTimestamp()
  {
    Assert.Null(FinancePolicies.ParseUtcOrNow("not-a-date"));
    Assert.NotNull(FinancePolicies.ParseUtcOrNow("2026-05-12T12:00:00Z"));
    Assert.NotNull(FinancePolicies.ParseUtcOrNow(null));
  }

  [Theory]
  [InlineData("2026-01", true)]
  [InlineData("2026-12", true)]
  [InlineData("2026-00", false)]
  [InlineData("2026-13", false)]
  [InlineData("202605", false)]
  public void PeriodKeyShouldUseYearMonthShapeWithValidMonth(string periodKey, bool expected)
  {
    Assert.Equal(expected, FinancePolicies.TryResolveCurrentPeriodKey(periodKey, out _));
  }

  [Theory]
  [InlineData("receivable_settlement", true)]
  [InlineData("provider_callback", true)]
  [InlineData("cash_movement", true)]
  [InlineData("read_report", false)]
  public void SensitiveFinancialMutationsShouldRequireIdempotency(string operation, bool expected)
  {
    Assert.Equal(expected, FinancePolicies.RequiresIdempotencyKey(operation));
  }

  [Theory]
  [InlineData("receivable_settlement", true)]
  [InlineData("period_closure", true)]
  [InlineData("manual_adjustment_preview", false)]
  public void LedgerOperationsShouldBeImmutable(string operation, bool expected)
  {
    Assert.Equal(expected, FinancePolicies.IsImmutableLedgerOperation(operation));
  }

  [Theory]
  [InlineData("discount", "discount")]
  [InlineData("commission_release", "commission_release")]
  [InlineData("manual_cash_adjustment", "manual_cash_adjustment")]
  [InlineData("free_text", "")]
  public void ApprovalActionsShouldBeWhitelisted(string input, string expected)
  {
    Assert.Equal(expected, FinancePolicies.NormalizeApprovalAction(input));
  }

  [Theory]
  [InlineData("matched", "matched")]
  [InlineData("manual_review", "manual_review")]
  [InlineData("pending", "")]
  public void ReconciliationStatusShouldBeExplicit(string input, string expected)
  {
    Assert.Equal(expected, FinancePolicies.NormalizeReconciliationStatus(input));
  }

  [Theory]
  [InlineData("hosted_checkout", false, false)]
  [InlineData("direct_card", false, true)]
  [InlineData("redirect", true, true)]
  public void PciScopeReviewShouldTriggerWhenCardDataCouldTouchThePlatform(string providerMode, bool storesCardData, bool expected)
  {
    Assert.Equal(expected, FinancePolicies.RequiresPciScopeReview(providerMode, storesCardData));
  }
}
