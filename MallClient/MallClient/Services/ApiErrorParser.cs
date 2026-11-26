using System.Collections.Generic;
using System.Text.Json;
using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Services;

public class ApiErrorParser : IApiErrorParser {
    private const string DefaultMessage = "Не удалось выполнить операцию.";

    private readonly IFieldTranslator _translator;

    public ApiErrorParser(IFieldTranslator translator) {
        _translator = translator;
    }

    public string Parse(string json) {
        if (string.IsNullOrWhiteSpace(json))
            return DefaultMessage;

        try {
            var payload = JsonSerializer.Deserialize<ApiErrorResponse>(json);
            if (payload == null) return DefaultMessage;

            if (payload.Fields != null && payload.Fields.Count > 0) {
                var parts = new List<string>();
                foreach (var kv in payload.Fields)
                    parts.Add($"{_translator.Translate(kv.Key)}: {kv.Value}");
                return string.Join("\n", parts);
            }

            if (!string.IsNullOrEmpty(payload.Error))
                return payload.Error;
        } catch (JsonException) {
        }

        return DefaultMessage;
    }
}