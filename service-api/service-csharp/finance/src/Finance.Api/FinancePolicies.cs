using System.Globalization;

namespace Finance.Api;

public static class FinancePolicies
{
  public static string NormalizePayableStatus(string? status)
  {
    var normalized = (status ?? string.Empty).Trim().ToLowerInvariant();
    return normalized switch
    {
      "open" or "paid" or "cancelled" => normalized,
      _ => string.Empty
    };
  }

  public static string NormalizeCashAccountStatus(string? status)
  {
    var normalized = (status ?? string.Empty).Trim().ToLowerInvariant();
    return normalized switch
    {
      "" => string.Empty,
      "active" or "inactive" => normalized,
      _ => string.Empty
    };
  }

  public static string NormalizeCashMovementDirection(string? direction)
  {
    var normalized = (direction ?? string.Empty).Trim().ToLowerInvariant();
    return normalized switch
    {
      "" => string.Empty,
      "inflow" or "outflow" => normalized,
      _ => string.Empty
    };
  }

  public static string NormalizeCashMovementType(string? movementType)
  {
    var normalized = (movementType ?? string.Empty).Trim().ToLowerInvariant();
    return normalized switch
    {
      "" => string.Empty,
      "receivable_settlement" or "payable_payment" or "cost_entry" or "manual_adjustment" => normalized,
      _ => string.Empty
    };
  }

  public static bool CanTransitionPayableStatus(string currentStatus, string targetStatus)
  {
    return currentStatus switch
    {
      "open" => targetStatus is "paid" or "cancelled",
      "paid" => false,
      "cancelled" => false,
      _ => false
    };
  }

  public static DateTime? ParseUtcOrNow(string? value)
  {
    if (string.IsNullOrWhiteSpace(value))
    {
      return DateTime.UtcNow;
    }

    if (!DateTimeOffset.TryParse(value, CultureInfo.InvariantCulture, DateTimeStyles.AssumeUniversal, out var parsed))
    {
      return null;
    }

    return parsed.UtcDateTime;
  }

  public static bool TryResolveCurrentPeriodKey(string periodKey, out string currentPeriodKey)
  {
    currentPeriodKey = DateTime.UtcNow.ToString("yyyy-MM", CultureInfo.InvariantCulture);
    return periodKey.Length == 7
      && periodKey[4] == '-'
      && int.TryParse(periodKey[..4], out _)
      && int.TryParse(periodKey[5..], out var month)
      && month is >= 1 and <= 12;
  }
}
