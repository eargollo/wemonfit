package fitbit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TimeSeriesEntry struct {
	DateTime string `json:"dateTime"`
	Value    string `json:"value"`
}

type TimeSeries struct {
	Entries []TimeSeriesEntry `json:"body-weight"`
}

type Weight struct {
	Bmi    float64 `json:"bmi"`
	Date   string  `json:"date"`
	Fat    float64 `json:"fat"`
	Logid  uint    `json:"logid"`
	Source string  `json:"source"`
	Time   string  `json:"time"`
	Weight float64 `json:"weight"`
}

type Weights struct {
	Weights []Weight `json:"weight"`
}

func (cli *Client) AllWeights() (weights []Weight, err error) {
	// First get the time series that has an entry for each day of the month
	// and averages the measurements.
	// The first entry in time series has the date for the first measurement.

	ts, err := cli.TimeSeries()
	if err != nil {
		return
	}

	// Get initial date to iterate from it to today reading all the weights
	initialDate, err := time.Parse("2006-01-02", ts.Entries[0].DateTime)
	if err != nil {
		return weights, fmt.Errorf(`could not parse timeseries date value "%s": %v`,
			ts.Entries[0].DateTime,
			err)
	}

	// Subtract 24 hours to make sure that no value is missed.
	initialDate = initialDate.Add(-24 * time.Hour)

	return cli.Weights(initialDate)
}

func (cli *Client) Weights(initialDate time.Time) (weights []Weight, err error) {
	var we Weights

	for initialDate.Before(time.Now()) {
		initialDate = initialDate.Add(30 * 24 * time.Hour)
		requestUrl := fmt.Sprintf(
			"https://api.fitbit.com/1/user/-/body/log/weight/date/%s/30d.json",
			initialDate.Format("2006-01-02"))

		response, err := cli.client.Get(requestUrl)
		if err != nil {
			return weights, fmt.Errorf("could not retrieve weiths for month ending at %s, error %v", initialDate.Format("2006-01-02"), err)
		}

		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&we)
		if err != nil {
			return weights, fmt.Errorf("could not decode payload for month ending at %s, error %v", initialDate.Format("2006-01-02"), err)
		}
		weights = append(weights, we.Weights...)
	}

	return weights, nil
}

func (cli *Client) TimeSeries() (ts TimeSeries, err error) {
	var response *http.Response

	response, err = cli.client.Get(
		"https://api.fitbit.com/1/user/-/body/weight/date/today/max.json")
	if err != nil {
		return
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&ts); err != nil {
		return ts, fmt.Errorf("could not decode time series: %v", err)
	}

	return
}
