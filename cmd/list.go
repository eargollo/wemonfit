/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"log"
	"time"

	"github.com/eargollo/wemonfit/pkg/fitbit"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Fitbit weight entries as CSV",
	Long: `List all Fitbit weight entries as CSV. Format:
ID, date, time, weight, bmi, fat%`,
	Run: func(cmd *cobra.Command, args []string) {
		cli, err := fitbit.New(secret)
		if err != nil {
			log.Fatalf("Could not initialize Fitbit client. Error: %v", err)
		}

		var ws []fitbit.Weight

		from, _ := cmd.Flags().GetString("from")
		if from == "" {
			// All weights
			ws, err = cli.AllWeights()
		} else {
			// Get initial date
			initialDate, err := time.Parse("2006/01/02", from)
			if err != nil {
				log.Fatalf("Could not parse from argument '%s' as a date. Make sure you format it as YYYY/MM/DD.", from)
			}
			log.Printf("List weights from %s", initialDate)
			ws, err = cli.Weights(initialDate)
		}

		if err != nil {
			log.Fatalf("Error getting all weights. Error: %v", err)
		}

		if len(ws) == 0 {
			fmt.Println("There are no entries for weight on FitBit!")
		}

		// Print CSV header
		fmt.Printf("%s, %s, %s, %s, %s, %s\n",
			"Logid",
			"Date",
			"Time",
			"Weight",
			"BMI",
			"Fat%",
		)
		for _, entry := range ws {
			fmt.Printf("%d, %s, %s, %f, %f, %f\n",
				entry.Logid,
				entry.Date,
				entry.Time,
				entry.Weight,
				entry.Bmi,
				entry.Fat,
			)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().String("from", "", "Initial date in format YYYY/MM/DD from which to list Fitbit weights")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
