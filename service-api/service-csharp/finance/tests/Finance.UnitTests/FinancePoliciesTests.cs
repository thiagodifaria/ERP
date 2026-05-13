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
}
