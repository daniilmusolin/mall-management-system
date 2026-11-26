using System.Text.Json.Serialization;
using MallClient.Providers;

namespace MallClient.Models;

public class Premise {
    [JsonPropertyName("id")] public long Id { get; set; }
    [JsonPropertyName("room_number")] public string RoomNumber { get; set; } = string.Empty;
    [JsonPropertyName("floor")] public int Floor { get; set; }
    [JsonPropertyName("area")] public double Area { get; set; }
    [JsonPropertyName("type")] public string Type { get; set; } = string.Empty;
    [JsonPropertyName("zone")] public string Zone { get; set; } = string.Empty;
    [JsonPropertyName("status")] public string Status { get; set; } = string.Empty;
}

public class Tenant {
    [JsonPropertyName("id")] public long Id { get; set; }
    [JsonPropertyName("name")] public string Name { get; set; } = string.Empty;
    [JsonPropertyName("inn")] public string Inn { get; set; } = string.Empty;
    [JsonPropertyName("legal_address")] public string LegalAddress { get; set; } = string.Empty;
    [JsonPropertyName("contact_phone")] public string ContactPhone { get; set; } = string.Empty;
    [JsonPropertyName("email")] public string Email { get; set; } = string.Empty;
}

public class LeaseContract {
    [JsonPropertyName("id")] public long Id { get; set; }

    [ForeignKey(typeof(PremiseOptionsProvider))]
    [JsonPropertyName("premise_id")] public long PremiseId { get; set; }

    [ForeignKey(typeof(TenantOptionsProvider))]
    [JsonPropertyName("tenant_id")] public long TenantId { get; set; }

    [JsonPropertyName("start_date")] public string StartDate { get; set; } = string.Empty;
    [JsonPropertyName("end_date")] public string EndDate { get; set; } = string.Empty;
    [JsonPropertyName("rate_per_sqm")] public decimal RatePerSqm { get; set; }
    [JsonPropertyName("deposit_amount")] public decimal DepositAmount { get; set; }
    [JsonPropertyName("status")] public string Status { get; set; } = string.Empty;
}

public class Payment {
    [JsonPropertyName("id")] public long Id { get; set; }

    [ForeignKey(typeof(ContractOptionsProvider))]
    [JsonPropertyName("contract_id")] public long ContractId { get; set; }

    [JsonPropertyName("amount")] public decimal Amount { get; set; }
    [JsonPropertyName("payment_date")] public string PaymentDate { get; set; } = string.Empty;
    [JsonPropertyName("payment_type")] public string PaymentType { get; set; } = string.Empty;
    [JsonPropertyName("status")] public string Status { get; set; } = string.Empty;
    [JsonPropertyName("reference_number")] public string ReferenceNumber { get; set; } = string.Empty;
}

public class ServiceRequest {
    [JsonPropertyName("id")] public long Id { get; set; }

    [ForeignKey(typeof(PremiseOptionsProvider))]
    [JsonPropertyName("premise_id")] public long PremiseId { get; set; }

    [JsonPropertyName("request_type")] public string RequestType { get; set; } = string.Empty;
    [JsonPropertyName("description")] public string Description { get; set; } = string.Empty;
    [JsonPropertyName("priority")] public string Priority { get; set; } = string.Empty;
    [JsonPropertyName("status")] public string Status { get; set; } = string.Empty;
    [JsonPropertyName("assigned_to")] public long? AssignedTo { get; set; }
    [JsonPropertyName("created_date")] public string CreatedDate { get; set; } = string.Empty;
    [JsonPropertyName("completed_date")] public string? CompletedDate { get; set; }
}