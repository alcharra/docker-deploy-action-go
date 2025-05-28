package client

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func streamComposeOutput(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	printedPull := false
	printedStop := false
	printedStart := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			fmt.Println("      \u21B3 ...")
			continue
		}

		switch {
		case strings.Contains(line, "Pulling") || strings.Contains(line, "Pulled"):
			if !printedPull {
				fmt.Println("   \U0001F4E5 Pulling images...")
				printedPull = true
			}
		case strings.Contains(line, "Stopping") || strings.Contains(line, "Stopped") ||
			strings.Contains(line, "Removing") || strings.Contains(line, "Removed"):
			if !printedStop {
				fmt.Println("   \U0001F4E6 Stopping services...")
				printedStop = true
			}
		case strings.Contains(line, "Creating") || strings.Contains(line, "Created") ||
			strings.Contains(line, "Starting") || strings.Contains(line, "Started"):
			if !printedStart {
				fmt.Println("   \U0001F4E6 Starting services...")
				printedStart = true
			}
		}
		fmt.Printf("      \u21B3 %s\n", line)
	}
}

func streamStackOutput(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	serviceMap := make(map[string]string)
	var serviceOrder []string
	convergedSet := make(map[string]bool)

	var currentID string
	var printedVerifyingFor string
	var lastCountdown string
	var rollbackInProgress bool
	var rollbackService string
	printedUpdateHeader := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			fmt.Println("      \u21B3 ...")
			continue
		}

		if strings.HasPrefix(line, "Updating service ") {
			if !printedUpdateHeader {
				fmt.Println("   \U0001F527 Updating services...")
				printedUpdateHeader = true
			}

			start := strings.Index(line, "service ") + len("service ")
			mid := strings.Index(line, " (id: ")
			end := strings.LastIndex(line, ")")

			if start > 0 && mid > start && end > mid {
				name := line[start:mid]
				id := line[mid+6 : end]
				serviceMap[id] = name
				serviceOrder = append(serviceOrder, id)
			}

			fmt.Printf("      \u21B3 %s\n", line)
			continue
		}

		if strings.Contains(line, "rollback: manually requested rollback") {
			rollbackInProgress = true
			if printedVerifyingFor != "" {
				rollbackService = printedVerifyingFor
			}
			fmt.Printf("   \U0001F501 Rolling back %s\n", rollbackService)
			continue
		}

		if rollbackInProgress && strings.Contains(line, "rolling back update:") {
			fmt.Printf("      \u21B3 %s\n", line)
			continue
		}

		if rollbackInProgress && strings.Contains(line, "converged") && strings.Contains(line, "verify: Service") {
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				key := parts[len(parts)-2]
				name := serviceMap[key]
				if name == "" {
					name = key
				}

				fmt.Printf("   \u2705 Service '%s' convergence complete\n", name)
				fmt.Printf("      \u21B3 %s\n", line)

				convergedSet[key] = true
				currentID = ""
				printedVerifyingFor = ""
				rollbackService = ""
				rollbackInProgress = false
			}
			continue
		}

		if strings.Contains(line, "verify: Waiting") {
			if currentID == "" {
				for _, id := range serviceOrder {
					if !convergedSet[id] && id != printedVerifyingFor {
						currentID = id
						break
					}
				}
			}

			name := ""
			if currentID != "" && currentID != printedVerifyingFor {
				name = serviceMap[currentID]
				printedVerifyingFor = currentID
			} else if rollbackInProgress && rollbackService != "" {
				name = rollbackService
				printedVerifyingFor = rollbackService
			}

			if name != "" {
				fmt.Printf("   \U0001F9EA Verifying service %s...\n", name)
			}

			if line != lastCountdown {
				lastCountdown = line
				fmt.Printf("      \u21B3 %s\n", line)
			}
			continue
		}

		if strings.Contains(line, "verify: Service") && strings.Contains(line, "converged") {
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				key := parts[len(parts)-2]
				name := serviceMap[key]
				if name == "" {
					name = key
				}

				fmt.Printf("   \u2705 Service '%s' convergence complete\n", name)
				fmt.Printf("      \u21B3 %s\n", line)

				convergedSet[key] = true
				currentID = ""
				printedVerifyingFor = ""
			}
			continue
		}

		fmt.Printf("      \u21B3 %s\n", line)
	}

	if currentID != "" && !convergedSet[currentID] {
		name := serviceMap[currentID]
		fmt.Printf("   \u2705 Service '%s' convergence complete\n", name)
	}
}
