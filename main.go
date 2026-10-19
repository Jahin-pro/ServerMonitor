package main

import (
	"fmt"

	"golang.org/x/sys/windows/svc/mgr"

	"golang.org/x/sys/windows/svc"
)

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
	serviceNames := []string{"CodexSandboxService.OpenAI.Codex", "LxssManager", "DefinitelyNotInstalled_Test_12345", "Spooler", "PostgreSQL"}

	for _, serviceName := range serviceNames {
		isRunning, err := checkService(serviceName)
		if err == nil {
			fmt.Println("Service is running:", isRunning)
		} else {
			fmt.Printf("Error checking service %s: %v\n", serviceName, err)
		}
	}
}

func checkServer(name string) (bool, error) {
	if name == "" {
		return false, fmt.Errorf("server name  is empty")
	}

	fmt.Println("Checking server:", name)
	return true, nil
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
	fmt.Println("Connected to SCM")

	services, err := m.OpenService(serviceName)
	if err != nil {
		return false, fmt.Errorf("failed to open service: %w", err)
	}
	fmt.Println("Service:", services.Name)
	defer services.Close()

	status, err := services.Query()
	if err != nil {
		return false, fmt.Errorf("failed to query service status: %w", err)
	}
	fmt.Println("Service status:", status.State)

	return status.State == svc.Running, nil
}
