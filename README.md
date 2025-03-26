# ATM-Database
=======
# Distributed System Project using Golang

## Overview
This project is a simple distributed banking system implemented in Golang. The system consists of a master server and multiple slave clients. The master server handles database operations, while slave clients connect to the master to perform various banking operations such as login, account creation, balance inquiry, deposits, and withdrawals.

## Features
- **Account Management**: Create a new account with an account number, customer ID, password, and account type.
- **Login**: Secure login using account number and password.
- **Balance Inquiry**: Check the balance of the account.
- **Deposits and Withdrawals**: Perform deposit and withdrawal operations.

## Project Structure
- `master.go`: Master server code that handles database connections and client requests.
- `slave.go`: Slave client code that interacts with the master server to perform banking operations.
- `database.sql`: SQL script to set up the required database and tables.

## Getting Started

### Prerequisites
- Go (Golang) installed on your machine.
- MySQL database installed and running.

# Technical Requirements
1. **Programming Language**: Go (Golang)
2. **Database**: MySQL
3. **Libraries**:
   - `database/sql`: For database operations.
   - `net`: For network operations.
   - `github.com/go-sql-driver/mysql`: MySQL driver for Go.

## Environment Setup

1. **Go Installation**: Ensure Go is installed on your machine. Follow the installation guide [here](https://golang.org/doc/install).
2. **MySQL Installation**: Install MySQL and ensure it is running. Follow the installation guide [here](https://dev.mysql.com/doc/mysql-installation-excerpt/5.7/en/).
3. **Database Setup**: Run the `database.sql` script to set up the required database and tables.

## Running the Project

1. **Start the Master Server**:
    ```sh
    go run master.go
    ```

2. **Start the Slave Client**:
    ```sh
    go run slave.go
    ```

3. **Interaction**: Follow the prompts on the slave client to interact with the system.

### Installation

1. **Clone the repository:**
    ```sh
    git clone https://github.com/yourusername/distributed-system-golang.git
    cd distributed-system-golang
    ```

2. **Set up the MySQL Database:**
    - Start your MySQL server.
    - Create the `BankATM` database and tables by running the SQL commands provided in `database.sql`.

3. **Configure the Database Connection:**
    - Edit the `master.go` file to update the MySQL connection string with your database credentials.
      ```go
      db, err = sql.Open("mysql", "root:yourpassword@tcp(localhost:3306)/BankATM")
      ```

### Usage

1. **Run the Master Server:**
    ```sh
    go run master.go
    ```

2. **Run the Slave Client:**
    ```sh
    go run slave.go
    ```

3. **Interact with the Slave Client:**
    - Follow the on-screen prompts to login, create an account, check balance, deposit, or withdraw funds.

## Database Schema

### Customer Table
```sql
CREATE TABLE Customer (
    CustomerID INT AUTO_INCREMENT PRIMARY KEY,
    CustomerName VARCHAR(100),
    Address VARCHAR(255),
    PhoneNumber VARCHAR(20),
    DateOfBirth DATE
);
```
### Account Table
```sql
CREATE TABLE Account (
    AccountNumber CHAR(16) PRIMARY KEY CHECK (AccountNumber REGEXP '^[0-9]{16}$'),   
    CustomerID INT,
    AccountPassword VARCHAR(4),
    Balance DECIMAL(10, 3),
    AccountType VARCHAR(50),
    DateOpened DATE,
    FOREIGN KEY (CustomerID) REFERENCES Customer(CustomerID),
    CONSTRAINT Check_AccountNumber_Length CHECK (LENGTH(AccountNumber) = 16)
);
```

### Transaction Table
```sql
CREATE TABLE Transaction (
    TransactionID INT AUTO_INCREMENT PRIMARY KEY,
    AccountNumber CHAR(16) CHECK (AccountNumber REGEXP '^[0-9]{16}$'),   
    TransactionType VARCHAR(50),
    Amount DECIMAL(10, 3),
    NewBalance DECIMAL(10, 3),
    TransactionDate DATETIME,
    FOREIGN KEY (AccountNumber) REFERENCES Account(AccountNumber)
);
```

>>>>>>> 7d6eb36 (Initial commit)
