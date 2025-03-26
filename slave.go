// slave
package main

import (
	"fmt"
	"net"
)

func main() {
	// Connect to the master server
	conn, err := net.Dial("tcp", "192.168.120.32:8088")
	if err != nil {
		fmt.Println("Error connecting to master server:", err)
		return
	}
	defer conn.Close()
	for {
		fmt.Println("--------------------------Welcome to ABC Bank-------------------------------")
		fmt.Println("--->> Press 1 to login to your account. ")
		fmt.Println("--->> Press 2 to Create an account. ")
		fmt.Println("--->> Press 3 to Shutdown The System. ")
		fmt.Println("----------------------------------------------------------------------------")
		var userInput, masterResponse string
		fmt.Scan(&userInput)

		if userInput == "1" {
			sendRequest("login", conn)
			var accNumb, accPass string

			fmt.Println("--->> Please Enter your account number : ")
			fmt.Scan(&accNumb)
			sendRequest(accNumb, conn)
			fmt.Println("----------------------------------------------------------------------------")

			fmt.Println("--->> Please Enter your account password : ")
			fmt.Scan(&accPass)
			sendRequest(accPass, conn)
			fmt.Println("----------------------------------------------------------------------------")

			masterResponse = receiveResponse(conn)

			if masterResponse != "1" {
				fmt.Println(masterResponse)
				fmt.Println(masterResponse)
				continue
			}

			fmt.Println("--------------------Successful login, welcome back.-------------------------")

			for {
				fmt.Println("--->> Please Select what you want to do.")
				fmt.Println("--->> Press 1 to see you balance.")
				fmt.Println("--->> Press 2 to deposit.")
				fmt.Println("--->> Press 3 to withdraw.")
				fmt.Println("--->> Press on any other key to quit.")
				fmt.Println("----------------------------------------------------------------------------")
				fmt.Scan(&userInput)
				fmt.Println("----------------------------------------------------------------------------")

				if userInput == "1" {
					sendRequest("balance", conn)
					sendRequest(accNumb, conn)
					sendRequest(accPass, conn)
					masterResponse = receiveResponse(conn)

					if masterResponse == "1" {
						masterResponse = receiveResponse(conn)

						fmt.Println("--->> Your balance is : ", masterResponse)

					} else {
						fmt.Println(masterResponse)
					}
					fmt.Println("----------------------------------------------------------------------------")
				} else if userInput == "2" {
					sendRequest("deposit", conn)
					sendRequest(accNumb, conn)
					sendRequest(accPass, conn)
					fmt.Println("Please Enter enter the amount that you want to deposit.")
					fmt.Scan(&userInput)
					fmt.Println("----------------------------------------------------------------------------")
					fmt.Println("Counting cache please wait.")
					sendRequest(userInput, conn)
					fmt.Println("----------------------------------------------------------------------------")
					masterResponse = receiveResponse(conn)
					if masterResponse != "OK" {
						fmt.Println("Error happen please try again.")
						fmt.Println("----------------------------------------------------------------------------")
						continue
					}
					fmt.Println("Successful deposit")
					fmt.Println("----------------------------------------------------------------------------")
				} else if userInput == "3" {
					sendRequest("withdraw", conn)
					sendRequest(accNumb, conn)
					sendRequest(accPass, conn)
					fmt.Println("Please Enter enter the amount that you want to withdraw.")
					fmt.Scan(&userInput)
					fmt.Println("----------------------------------------------------------------------------")
					sendRequest(userInput, conn)
					masterResponse = receiveResponse(conn)
					if masterResponse != "OK" {
						fmt.Println(masterResponse)
						fmt.Println("----------------------------------------------------------------------------")
						continue
					}
					fmt.Println("Successful withdraw.")
					fmt.Println("----------------------------------------------------------------------------")
				} else {
					break
				}

				fmt.Println("Do you want to do something else?.")
				fmt.Println("Press 1 for YES.")
				fmt.Println("Or Press on any other key to quit.")
				fmt.Scan(&userInput)
				fmt.Println("----------------------------------------------------------------------------")
				if userInput != "1" {
					break
				}
			}
		} else if userInput == "2" {
			sendRequest("createAcc", conn)

			// acc num
			masterResponse = receiveResponse(conn)
			if masterResponse == "EOF" {
				fmt.Println("Error in connecting ")
				continue
			}
			fmt.Println(masterResponse)
			fmt.Scan(&userInput)
			fmt.Println("----------------------------------------------------------------------------")
			sendRequest(userInput, conn)

			// customer id
			masterResponse = receiveResponse(conn)
			if masterResponse == "EOF" {
				fmt.Println("Error in connecting ")
				continue
			}
			fmt.Println(masterResponse)
			fmt.Scan(&userInput)
			sendRequest(userInput, conn)
			fmt.Println("----------------------------------------------------------------------------")
			// password
			masterResponse = receiveResponse(conn)
			if masterResponse == "EOF" {
				fmt.Println("Error in connecting ")
				continue
			}
			fmt.Println(masterResponse)
			fmt.Scan(&userInput)
			sendRequest(userInput, conn)
			fmt.Println("----------------------------------------------------------------------------")
			// acc type
			masterResponse = receiveResponse(conn)
			if masterResponse == "EOF" {
				fmt.Println("Error in connecting ")
				continue
			}
			fmt.Println(masterResponse)
			fmt.Scan(&userInput)
			sendRequest(userInput, conn)
			fmt.Println("----------------------------------------------------------------------------")
			masterResponse = receiveResponse(conn)
			fmt.Println(masterResponse)
			fmt.Println("----------------------------------------------------------------------------")
		} else if userInput == "3" {
			sendRequest("quit", conn)
			fmt.Println("----------------------------------------------------------------------------")
			return
		}
	}
}

func sendRequest(str string, conn net.Conn) {
	var err error
	// Send the query to the master server
	_, err = conn.Write([]byte(str))
	if err != nil {
		fmt.Println("Error sending query to master server:", err)
		return
	}
}
func receiveResponse(conn net.Conn) string {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from slave:", err)
		return "EOF"
	}
	response := string(buf[:n])
	return response
}
