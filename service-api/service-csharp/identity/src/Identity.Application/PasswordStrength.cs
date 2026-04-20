namespace Identity.Application;

internal static class PasswordStrength
{
  public static bool IsStrong(string password)
  {
    return !string.IsNullOrWhiteSpace(password)
      && password.Length >= 10
      && password.Any(char.IsUpper)
      && password.Any(char.IsLower)
      && password.Any(char.IsDigit);
  }
}
