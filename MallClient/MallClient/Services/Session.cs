using System;
using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Services;

public class Session : ISession {
    public string Token { get; private set; } = string.Empty;
    public User CurrentUser { get; private set; } = new();

    public bool IsAuthenticated => !string.IsNullOrEmpty(Token);
    public bool IsAdmin => CurrentUser.Role == "admin";
    public bool IsManager => CurrentUser.Role == "manager";
    public bool IsAccountant => CurrentUser.Role == "accountant";
    public bool IsTechnician => CurrentUser.Role == "technician";
    public bool IsSupport => CurrentUser.Role == "support";

    public event Action? SessionExpired;

    public void Set(string token, User user) {
        Token = token;
        CurrentUser = user;
    }

    public void Clear() {
        Token = string.Empty;
        CurrentUser = new User();
    }

    public void UpdateFullName(string fullName) {
        CurrentUser.FullName = fullName;
    }

    public void RaiseSessionExpired() {
        SessionExpired?.Invoke();
    }
}