package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Service struct {
	Name string `json:"Name"`
	URL  string `json:"URL"`
}

func main() {

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help", "--help", "-h", "-help":
			printUsage()
			os.Exit(0)
		}
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	serverAddr := os.Getenv("CHECKER_ADDR")
	if serverAddr == "" {
		serverAddr = "http://localhost:8081/services"
	}

	switch cmd {
	case "list":
		listServices(serverAddr)

	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		addName := addCmd.String("name", "", "Service name")
		addURL := addCmd.String("url", "", "Service URL")

		addCmd.Parse(os.Args[2:])
		if *addName == "" || *addURL == "" {
			printUsage()
			os.Exit(1)
		}
		addService(serverAddr, *addName, *addURL)

	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		deleteName := deleteCmd.String("name", "", "Service name")

		deleteCmd.Parse(os.Args[2:])
		if *deleteName == "" {
			fmt.Println("Usage: checker-cli delete -name <name>")
			os.Exit(1)
		}
		deleteService(serverAddr, *deleteName)
	case "help":
		printUsage()
	default:
		fmt.Println("Unrecognized command: " + cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: checker-cli [command]")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  list                    List all monitored services")
	fmt.Println("  add -name <name> -url <url>  Add a new service to monitor")
	fmt.Println("  delete -name <name>     Delete a service from monitoring")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  - You can set server address via CHECKER_ADDR environment variable")
	fmt.Println("  - Default address: http://localhost:8080/services")
}
func listServices(
	addr string,
) {
	resp, err := http.Get(addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var raw interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		fmt.Println("Error:", err)
		return
	}

	pretty, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}
	fmt.Println(string(pretty))
}

func addService(
	addr string,
	name string,
	url string,
) {
	s := Service{
		Name: name,
		URL:  url,
	}
	data, _ := json.Marshal(s)
	resp, err := http.Post(
		addr,
		"application/json",
		bytes.NewReader(data),
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Service added", name)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(
			"Service failed",
			resp.StatusCode,
			string(body),
		)
	}

}

func deleteService(
	addr string,
	name string,
) {
	s := Service{
		Name: name,
	}
	data, _ := json.Marshal(s)
	req, err := http.NewRequest(
		"DELETE",
		addr,
		bytes.NewReader(data),
	)
	if err != nil {
		fmt.Println("Error create delete request:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error while do http delete:", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Service deleted", name)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Delete failed", resp.StatusCode, string(body))
	}

}
