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

	if len(os.Args) < 2 {
		fmt.Println("Usage: checker-cli [list|add|delete]")
		os.Exit(1)
	}

	cmd := os.Args[1]

	serverAddr := os.Getenv("CHECKER_ADDR")
	if serverAddr == "" {
		serverAddr = "http://localhost:8080/services"
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
			fmt.Println("Usage: checker-cli add -name <name> -url <url>")
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

	default:
		fmt.Println("Unrecognized command: " + cmd)
		os.Exit(1)
	}
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
