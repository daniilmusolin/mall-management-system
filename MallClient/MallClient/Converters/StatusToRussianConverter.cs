using System;
using System.Globalization;
using System.Windows.Data;

namespace MallClient.Converters;

public class StatusToRussianConverter : IValueConverter {
    public object Convert(object value, Type targetType, object parameter, CultureInfo culture) {
        var status = value as string;
        return status switch {
            "free" => "Свободно",
            "rented" => "Арендовано",
            "maintenance" => "На ремонте",
            "active" => "Активен",
            "terminated" => "Расторгнут",
            "expired" => "Истек",
            "pending" => "Ожидает",
            "completed" => "Завершен",
            "cancelled" => "Отменен",
            "open" => "Открыта",
            "in_progress" => "В работе",
            "closed" => "Закрыта",
            _ => status ?? string.Empty
        };
    }

    public object ConvertBack(object value, Type targetType, object parameter, CultureInfo culture) {
        throw new NotSupportedException();
    }
}