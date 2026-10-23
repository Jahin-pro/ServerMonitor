package main

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc/mgr"

	"golang.org/x/sys/windows/svc"

	"net"
)

type ServiceResult struct {
	Name      string
	IsRunning bool
	Err       error
}

type Server struct {
	Name        string
	Address     string
	IsAvailable bool
}

func (s Server) PrintInfo() {
	fmt.Println("Server:", s.Name)
	fmt.Println("Address:", s.Address)
	fmt.Println("Is Available:", s.IsAvailable)
}

func (s *Server) SetAvailability(available bool) {
	s.IsAvailable = available
}

func main() {
	fmt.Println("ServerMonitor")
	/*serviceNames := []string{"CodexSandboxService.OpenAI.Codex", "LxssManager", "DefinitelyNotInstalled_Test_12345", "Spooler", "PostgreSQL"}
	serviceResults := collectResults(serviceNames)
	printReport(serviceResults)*/

	connectToLocalSQL, err := checkTCP("localhost:5432", 2*time.Second)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Connet to localhost:5432: %v\n", connectToLocalSQL)
}

func checkServer(name string) (bool, error) {
	if name == "" {
		return false, fmt.Errorf("server name  is empty")
	}

	fmt.Println("Checking server:", name)
	return true, nil
}

func collectResults(names []string) []ServiceResult {
	var result []ServiceResult

	for _, serviceName := range names {
		isRunning, err := checkService(serviceName)
		result = append(result, ServiceResult{Name: serviceName, IsRunning: isRunning, Err: err})
	}
	return result
}

func checkService(serviceName string) (bool, error) {
	if serviceName == "" {
		return false, fmt.Errorf("service name is empty")
	}

	m, err := mgr.Connect()
	if err != nil {
		return false, fmt.Errorf("failed to connect to SCM: %w", err)
	}
	defer m.Disconnect()

	services, err := m.OpenService(serviceName)
	if err != nil {
		return false, fmt.Errorf("failed to open service: %w", err)
	}
	defer services.Close()

	status, err := services.Query()
	if err != nil {
		return false, fmt.Errorf("failed to query service status: %w", err)
	}

	return status.State == svc.Running, nil
}

func printReport(results []ServiceResult) {
	for _, result := range results {
		if result.Err != nil {
			fmt.Printf("Service: %s, Error: %v\n", result.Name, result.Err)
		} else {
			fmt.Printf("Service: %s, Is Running: %v\n", result.Name, result.IsRunning)
		}
	}
}

func checkTCP(address string, timeout time.Duration) (bool, error) {
	if address == "" {
		return false, fmt.Errorf("address is empty")
	}
	connection, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false, fmt.Errorf("failed to connect to %s: %w", address, err)
	}
	defer connection.Close()
	return true, nil
}
