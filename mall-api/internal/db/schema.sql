CREATE DATABASE IF NOT EXISTS shopping_center;
CREATE DATABASE IF NOT EXISTS shopping_center_test;

CREATE TABLE IF NOT EXISTS Users (
    UserId INT AUTO_INCREMENT PRIMARY KEY,
    Email VARCHAR(254) NOT NULL UNIQUE,
    PasswordHash VARCHAR(255) NOT NULL,
    FullName VARCHAR(255) NOT NULL,
    Role ENUM('admin', 'manager', 'accountant', 'technician', 'support') NOT NULL,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_users_email (Email)
);

CREATE TABLE IF NOT EXISTS RefreshTokens (
    TokenId INT AUTO_INCREMENT PRIMARY KEY,
    UserId INT NOT NULL,
    Token VARCHAR(255) NOT NULL UNIQUE,
    ExpiresAt TIMESTAMP NOT NULL,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (UserId) REFERENCES Users(UserId) ON DELETE CASCADE,
    INDEX idx_refresh_token (Token)
);

CREATE TABLE IF NOT EXISTS Premises (
    PremiseId INT AUTO_INCREMENT PRIMARY KEY,
    RoomNumber VARCHAR(20) NOT NULL UNIQUE,
    Floor INT NOT NULL,
    Area DECIMAL(10, 2) NOT NULL,
    Type VARCHAR(50) NOT NULL,
    Zone VARCHAR(50),
    Status ENUM('free', 'rented', 'maintenance') DEFAULT 'free',
    INDEX idx_premises_status (Status)
);

CREATE TABLE IF NOT EXISTS Tenants (
    TenantId INT AUTO_INCREMENT PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    INN VARCHAR(12) NOT NULL UNIQUE,
    LegalAddress VARCHAR(500),
    ContactPhone VARCHAR(20),
    Email VARCHAR(254),
    INDEX idx_tenants_inn (INN)
);

CREATE TABLE IF NOT EXISTS LeaseContracts (
    ContractId INT AUTO_INCREMENT PRIMARY KEY,
    PremiseId INT NOT NULL,
    TenantId INT NOT NULL,
    StartDate DATE NOT NULL,
    EndDate DATE NOT NULL,
    RatePerSqm DECIMAL(10, 2) NOT NULL,
    DepositAmount DECIMAL(10, 2) DEFAULT 0,
    Status ENUM('active', 'terminated', 'expired') DEFAULT 'active',
    FOREIGN KEY (PremiseId) REFERENCES Premises(PremiseId),
    FOREIGN KEY (TenantId) REFERENCES Tenants(TenantId),
    INDEX idx_contracts_status (Status)
);

CREATE TABLE IF NOT EXISTS Payments (
    PaymentId INT AUTO_INCREMENT PRIMARY KEY,
    ContractId INT NOT NULL,
    Amount DECIMAL(10, 2) NOT NULL,
    PaymentDate DATE NOT NULL,
    PaymentType ENUM('rent', 'deposit', 'utilities', 'penalty') NOT NULL,
    Status ENUM('pending', 'completed', 'cancelled') DEFAULT 'pending',
    ReferenceNumber VARCHAR(100),
    FOREIGN KEY (ContractId) REFERENCES LeaseContracts(ContractId) ON DELETE CASCADE,
    INDEX idx_payments_date (PaymentDate)
);

CREATE TABLE IF NOT EXISTS ServiceRequests (
    RequestId INT AUTO_INCREMENT PRIMARY KEY,
    PremiseId INT NOT NULL,
    RequestType VARCHAR(100) NOT NULL,
    Description TEXT,
    Priority ENUM('low', 'medium', 'high', 'critical') DEFAULT 'medium',
    Status ENUM('open', 'in_progress', 'closed', 'cancelled') DEFAULT 'open',
    AssignedTo INT NULL,
    CreatedDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CompletedDate TIMESTAMP NULL,
    FOREIGN KEY (PremiseId) REFERENCES Premises(PremiseId),
    FOREIGN KEY (AssignedTo) REFERENCES Users(UserId) ON DELETE SET NULL,
    INDEX idx_requests_status (Status)
);

ALTER TABLE Users MODIFY Email VARCHAR(254) NOT NULL;
ALTER TABLE Users MODIFY FullName VARCHAR(255) NOT NULL;
ALTER TABLE Tenants MODIFY LegalAddress VARCHAR(500);
ALTER TABLE Tenants MODIFY Email VARCHAR(254);
ALTER TABLE LeaseContracts MODIFY Status ENUM('active','terminated','expired') DEFAULT 'active';
ALTER TABLE Payments MODIFY Status ENUM('pending','completed','cancelled') DEFAULT 'pending';
ALTER TABLE Payments MODIFY ReferenceNumber VARCHAR(100);
ALTER TABLE ServiceRequests MODIFY Priority ENUM('low','medium','high','critical') DEFAULT 'medium';
ALTER TABLE ServiceRequests MODIFY Status ENUM('open','in_progress','closed','cancelled') DEFAULT 'open';