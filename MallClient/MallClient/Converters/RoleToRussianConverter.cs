using System;
using System.Globalization;
using System.Windows.Data;

namespace MallClient.Converters;

public class RoleToRussianConverter : IValueConverter {
    public object Convert(object value, Type targetType, object parameter, CultureInfo culture) {
        var role = value as string;
        return role switch {
            "admin" => "Администратор",
            "manager" => "Менеджер",
            "accountant" => "Бухгалтер",
            "technician" => "Техник",
            "support" => "Поддержка",
            _ => role ?? string.Empty
        };
    }

    public object ConvertBack(object value, Type targetType, object parameter, CultureInfo culture) {
        throw new NotSupportedException();
    }
}