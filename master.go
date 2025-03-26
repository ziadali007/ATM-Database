package main

import (
	"database/sql"
	"fmt"
	"net"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver package
)

var db *sql.DB

func main() {
	DatabaseConnection() // Open a database connection
	defer db.Close()     // Defer closing the database connection

	listener, err := net.Listen("tcp", ":8088")
	if err != nil {
		fmt.Println("Error starting master server:", err)
		return
	}
	fmt.Println("Master server listening for connections on port 8088...")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection : ", err)
			continue
		}
		fmt.Println("New connection accepted from : ", conn.RemoteAddr())
		go handleSlaveRequest(conn)
	}
}
func DatabaseConnection() {
	var err error
	db, err = sql.Open("mysql", "root:12345678@tcp(localhost:3306)/BankATM")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	db.SetMaxOpenConns(10) // Maximum number of open connections
	db.SetMaxIdleConns(10) // Maximum number of idle connections
}

func handleSlaveRequest(conn net.Conn) {
	defer conn.Close()
	for {
		//Read Request from the slave
		req, err := getRequest(conn)
		if err != nil {
			fmt.Println(req, err)
			return
		}

		if req == "quit" {
			// Check for the special termination message
			fmt.Println("Closing connection with slave")
			return
		} else if req == "login" {
			accNumber, err := getRequest(conn)
			accPass, err2 := getRequest(conn)
			if err != nil || err2 != nil {
				fmt.Println("Error in Get login request")
				sendResponse("Error", conn)
			}
			var right_num, right_pas bool
			right_num, _ = findAccount(accNumber)
			right_pas, _ = correctPassword(accNumber, accPass)
			if !right_num || !right_pas {
				sendResponse("Ther is something wrong in account number or password", conn)
				continue
			}
			sendResponse("1", conn)
			fmt.Println("---> successful login.")
		} else if req == "createAcc" {
			var accNumber, accPass, accType, custID string

			sendResponse("Enter Account new Number", conn)
			accNumber, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in entering account")
				continue
			}

			sendResponse("Enter customer ID", conn)
			custID, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in entering customer ID")
				continue
			}

			sendResponse("Enter Password", conn)
			accPass, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in entering password")
				continue
			}

			sendResponse("Enter Account type", conn)
			accType, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in entering account type")
				continue
			}

			ID, err2 := strconv.Atoi(custID)
			if err2 != nil {
				fmt.Println("Wrong in Customer ID")
				sendResponse("Error in customer ID", conn)
				continue
			}
			err = createAccount(accNumber, ID, accPass, accType)
			if err != nil {
				fmt.Println("Error in creating account")
				sendResponse("Error in creating account", conn)
				continue
			}
			////////////////////DONE ////////////////////////
			sendResponse("OK", conn)
		} else if req == "balance" {
			var accNUmber, accPass string
			var curBalance float64

			// Get acc
			accNUmber, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Account Number")
				continue
			}
			//Get pass
			accPass, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Account Password")
				continue
			}

			curBalance, err = getBalance(accNUmber, accPass)

			if err != nil {
				fmt.Println("Error in Get customer balance")
				sendResponse("something went wrong. please try again", conn)
				continue
			} else {
				sendResponse("1", conn)
				sendResponse(strconv.FormatFloat(curBalance, 'f', -1, 64), conn)
			}
			// DONE
			fmt.Println("Sucssesful Opreation")
		} else if req == "withdraw" {
			var accNUmber, accPass, sAmmount string
			var ammount float64

			accNUmber, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Account Number")
				continue
			}

			accPass, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Account Password")
				continue
			}

			sAmmount, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in get Amount")
				continue
			}
			ammount, err = strconv.ParseFloat(sAmmount, 64)
			if err != nil {
				fmt.Println("Error in Ammount")
				continue
			}
			ok := insertTransaction(accNUmber, req, ammount, accPass)
			if ok {
				sendResponse("OK", conn)
				continue
			}
			sendResponse("Unsuccessful withdraw.", conn)
		} else if req == "deposit" {
			var accNUmber, accPass, sAmmount string
			var ammount float64
			accNUmber, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Account Number")
				continue
			}
			accPass, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Account Password")
				continue
			}
			sAmmount, err = getRequest(conn)
			if err != nil {
				fmt.Println("Error in Account Password")
				sendResponse("ERROR", conn)
				continue
			}
			ammount, err = strconv.ParseFloat(sAmmount, 64)
			if err != nil {
				fmt.Println("Error in Ammount")
				sendResponse("ERROR", conn)
				continue
			}

			ok := insertTransaction(accNUmber, req, ammount, accPass)
			if ok {
				sendResponse("OK", conn)
				continue
			}
			sendResponse("ERROR", conn)
		} else {
			response := "invalid choose"
			_, err = conn.Write([]byte(response))
			if err != nil {
				fmt.Println("Slave connection error : ", err)
			}
		}
	}
}

// /////////////////Dealing with Slaves/////////////////////
func sendResponse(msg string, conn net.Conn) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("Slave connection error : ", err)
	}
}
func getRequest(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "Error reading from slave : ", err
	}
	req := string(buf[:n])
	return req, err
}

// ////////////////////////////////////////SQL//////////////////

func findAccount(accountNumber string) (bool, error) {
	var accountIN string
	err := db.QueryRow("SELECT AccountNumber FROM Account WHERE AccountNumber = (?)", accountNumber).Scan(&accountIN)
	if err != nil {
		return false, err
	}
	if accountIN == accountNumber {
		return true, nil
	}
	return false, nil
}

func correctPassword(accountNumber string, accountPassword string) (bool, error) {
	var passWord string
	err := db.QueryRow("SELECT AccountPassword FROM Account WHERE AccountNumber = (?)", accountNumber).Scan(&passWord)
	if err != nil {
		return false, err
	}
	if accountPassword == passWord {
		return true, nil
	}
	return false, nil
}

func createAccount(AccNUmber string, customerId int, accPass string, accType string) error {
	// Function of date missing
	_, err := db.Exec("INSERT INTO account (AccountNumber,CustomerID , AccountPassword , Balance ,AccountType , DateOpened) VALUES (?, ? ,?, 0 ,?, NOW() )", AccNUmber, customerId, accPass, accType)
	if err != nil {
		fmt.Println("Error in insert account:", err)
		return err
	}
	fmt.Println("account inserted successfully!")
	return err
}
func getBalance(accountNumber string, accountPassword string) (float64, error) {
	var balance float64
	err := db.QueryRow("SELECT Balance FROM Account WHERE AccountNumber = (?) AND AccountPassword = (?)", accountNumber, accountPassword).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func insertTransaction(Acc_Number string, Tran_Type string, Amount_ float64, acc_pass string) bool {

	if Tran_Type == "deposit" {
		samePass, errr := correctPassword(Acc_Number, acc_pass)
		if errr != nil {
			fmt.Println("Error in check password : ", errr)
		}
		if !samePass {
			fmt.Println("Wrong Password Or user account")
			return false
		}
		/////////////////////////////////////////////////////////////////////
		_, eerr := db.Exec("update account  set Balance = Balance + (?) where AccountNumber = (?)", Amount_, Acc_Number)
		if eerr != nil {
			fmt.Println("Error update balance :", eerr)
			return false
		}

		_, err := db.Exec("INSERT INTO transaction (AccountNumber, TransactionType, Amount, NewBalance, TransactionDate) SELECT (?) AS AccountNumber, 'Deposit' AS TransactionType, (?) AS Amount,(Account.Balance ) AS NewBalance, NOW() AS TransactionDate FROM  account AS Account WHERE Account.AccountNumber = (?)", Acc_Number, Amount_, Acc_Number)
		if err != nil {
			fmt.Println("Error inserting data:", err)
			return false
		}

		fmt.Println("Sucsseful deposit.")
		balance, err := getBalance(Acc_Number, acc_pass)
		if err != nil {
			fmt.Println("Error : ", err)
			return false
		}
		fmt.Print("Now Your Balance = ", balance)

	} else if Tran_Type == "withdraw" {
		balance, errr := getBalance(Acc_Number, acc_pass)
		if errr != nil {
			fmt.Println("error happen in withdraw : ", errr)
			return false
		}
		if Amount_ > balance {
			fmt.Println("The ammount the custmer need to withdraw the more than his balance")
			return false
		}
		Amount_ = Amount_ * -1
		_, eerr := db.Exec("update account  set Balance = Balance + (?) where AccountNumber = (?)", Amount_, Acc_Number)
		if eerr != nil {
			fmt.Println("Error update balance :", eerr)
			return false
		}
		_, err := db.Exec("INSERT INTO transaction (AccountNumber, TransactionType, Amount, NewBalance, TransactionDate) SELECT (?) AS AccountNumber, 'Withdraw' AS TransactionType, (?) AS Amount,(Account.Balance ) AS NewBalance, NOW() AS TransactionDate FROM  account AS Account WHERE Account.AccountNumber = (?)", Acc_Number, Amount_, Acc_Number)
		if err != nil {
			fmt.Println("Error inserting data:", err)
			return false
		}
		fmt.Println("Sucsseful withdraw.")
	}

	return true
}
