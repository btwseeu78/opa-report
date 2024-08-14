package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"opa-report/opa"
	"os"
	"strings"
)

func main() {
	// var incluster bool
	incluster := false

	constraintList, err := opa.GetConstraints(&incluster)

	if err != nil {
		fmt.Printf("Unable to Get the data")
	}
	for _, cons := range constraintList {
		for _, vls := range cons.Status.Violations {
			name := vls.Name
			if vls.Kind == "Pod" {
				splittedString := strings.Split(vls.Name, "-")

				if len(splittedString) >= 3 {
					name := strings.Join(splittedString[0:len(splittedString)-2], "-")
					fmt.Println(name)
				}

			}
			fmt.Printf("violation kind: %s , Name: %s , Constraint Name: %s, Enforcemnet: %s", vls.Kind, name, cons.Meta.Name, cons.Spec.EnforcementAction)
			fmt.Println("\n----------------------------------")
		}
	}

	// Create a new CSV file
	file, err := os.Create("constraints.csv")
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{
		"Kind", "Name", "Enforcement Action", "Total Violations", "Violation Kind",
		"Violation Name", "Violation Namespace", "Violation Message", "Violation Enforcement Action",
	}
	if err := writer.Write(headers); err != nil {
		log.Fatalf("Failed to write headers to CSV file: %v", err)
	}

	// Write data rows
	for _, c := range constraintList {
		for _, v := range c.Status.Violations {
			row := []string{
				c.Meta.Kind,
				c.Meta.Name,
				c.Spec.EnforcementAction,
				fmt.Sprintf("%v", c.Status.TotalViolations),
				v.Kind,
				v.Name,
				v.Namespace,
				v.Message,
				v.EnforcementAction,
			}
			if err := writer.Write(row); err != nil {
				fmt.Printf("Failed to write row to CSV file: %v", err)
			}
		}
	}

	fmt.Println("CSV file 'constraints.csv' created successfully!")
}
