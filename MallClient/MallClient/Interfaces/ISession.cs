using System;
using MallClient.Models;

namespace MallClient.Interfaces;

public interface ISession {
    string Token { get; }
    User CurrentUser { get; }
    bool IsAuthenticated { get; }
    bool IsAdmin { get; }
    bool IsManager { get; }
    bool IsAccountant { get; }
    bool IsTechnician { get; }
    bool IsSupport { get; }
    void Set(string token, User user);
    void Clear();
    void UpdateFullName(string fullName);
    event Action? SessionExpired;
    void RaiseSessionExpired();
}