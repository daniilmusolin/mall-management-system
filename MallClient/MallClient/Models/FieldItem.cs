using System;
using System.Collections.Generic;

namespace MallClient.Models;

public class FieldItem {
    public string Name { get; set; } = string.Empty;
    public string Label { get; set; } = string.Empty;
    public string Value { get; set; } = string.Empty;

    public bool IsDropdown { get; set; } = false;
    public bool IsNumeric { get; set; } = false;
    public bool IsDate { get; set; } = false;
    public bool IsPhone { get; set; } = false;
    public bool IsEmail { get; set; } = false;

    public int MaxLength { get; set; } = 0;

    public DateTime? DateValue {
        get {
            if (string.IsNullOrWhiteSpace(Value)) return null;
            if (DateTime.TryParse(Value, out var dt)) return dt;
            return null;
        }
        set {
            Value = value?.ToString("yyyy-MM-dd") ?? string.Empty;
        }
    }

    public List<OptionItem> Options { get; set; } = new();
    public bool IsValid { get; set; } = true;
}

public class OptionItem {
    public string Value { get; set; } = string.Empty;
    public string DisplayName { get; set; } = string.Empty;

    public override string ToString() {
        return DisplayName;
    }
}