using System.IO;
using MallClient.Interfaces;

namespace MallClient.Services;

public class FileLogger : IAppLogger {
    private static readonly object LockObj = new();
    private readonly string _logFilePath;

    public FileLogger() {
        var appData = Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData);
        var appDir = Path.Combine(appData, "MallClient");

        if (!Directory.Exists(appDir)) {
            Directory.CreateDirectory(appDir);
        }

        _logFilePath = Path.Combine(appDir, "logs.txt");
    }

    public void LogInfo(string message) => WriteLog("INFO", message, null);

    public void LogWarning(string message) => WriteLog("WARN", message, null);

    public void LogError(string message, Exception? ex = null) => WriteLog("ERROR", message, ex);

    private void WriteLog(string level, string message, Exception? ex) {
        lock (LockObj) {
            try {
                using var writer = new StreamWriter(_logFilePath, append: true);
                var timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss");
                writer.WriteLine($"[{timestamp}] [{level}] {message}");

                if (ex != null) {
                    writer.WriteLine($"    Exception: {ex.GetType().Name}");
                    writer.WriteLine($"    Message: {ex.Message}");
                    writer.WriteLine($"    StackTrace: {ex.StackTrace}");
                }
            } catch { }
        }
    }
}