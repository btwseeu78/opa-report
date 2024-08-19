/*
Copyright © 2024 Arpan Chatterjee btwseeu78@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"opa-report/opa"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
		saveConstraint()
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func saveConstraint() {
	incluster := false

	constraintList, err := opa.GetConstraints(&incluster)

	if err != nil {
		fmt.Printf("Unable to Get the data")
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
		"Kind", "Name", "Enforcement Action", "Violation Kind",
		"Violation Name", "Violation Namespace", "Violation Message", "Violation Enforcement Action",
	}
	if err := writer.Write(headers); err != nil {
		log.Fatalf("Failed to write headers to CSV file: %v", err)
	}

	// Write data rows

	for _, c := range constraintList {
		var hashMap []string
		for _, v := range c.Status.Violations {
			name := v.Name
			if v.Kind == "Pod" {
				splittedString := strings.Split(v.Name, "-")

				if len(splittedString) >= 3 {
					name = strings.Join(splittedString[0:len(splittedString)-2], "-")
				}

			}
			uniqueNameSpacedName := v.Namespace + "-" + name
			if len(hashMap) == 0 || !slices.Contains(hashMap, uniqueNameSpacedName) {
				hashMap = append(hashMap, uniqueNameSpacedName)
			} else if slices.Contains(hashMap, uniqueNameSpacedName) {
				continue
			}

			row := []string{
				c.Meta.Kind,
				c.Meta.Name,
				c.Spec.EnforcementAction,
				v.Kind,
				name,
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
