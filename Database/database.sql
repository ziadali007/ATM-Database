create database BankATM;
use BankATM;
CREATE TABLE Customer (
    CustomerID INT AUTO_INCREMENT PRIMARY KEY,
    CustomerName VARCHAR(100),
    Address VARCHAR(255),
    PhoneNumber VARCHAR(20),
    DateOfBirth DATE
);

CREATE TABLE Account (
  AccountNumber CHAR(16) PRIMARY KEY CHECK (AccountNumber REGEXP '^[0-9]{16}$'),   
    CustomerID INT,
    AccountPassword varchar(4),
    Balance DECIMAL(10, 3),
    AccountType VARCHAR(50),
    DateOpened DATE,
    FOREIGN KEY (CustomerID) REFERENCES Customer(CustomerID),
    CONSTRAINT Check_AccountNumber_Length CHECK (LENGTH(AccountNumber) = 16)
);

CREATE TABLE Transaction (
    TransactionID INT AUTO_INCREMENT PRIMARY KEY,
  AccountNumber CHAR(16) CHECK (AccountNumber REGEXP '^[0-9]{16}$'),   
    TransactionType VARCHAR(50),
    Amount DECIMAL(10, 3),
    NewBalance DECIMAL(10, 3),
    TransactionDate DATETIME,
    FOREIGN KEY (AccountNumber) REFERENCES Account(AccountNumber)
);

select * from Customer;
