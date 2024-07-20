package main

import (
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/Scrin/scanmytesla-influxdb/common/logging"
	"github.com/Scrin/scanmytesla-influxdb/common/version"
	"github.com/Scrin/scanmytesla-influxdb/config"
	"github.com/Scrin/scanmytesla-influxdb/influxdb"
	"github.com/Scrin/scanmytesla-influxdb/reader"
	"github.com/rs/zerolog/log"
)

func main() {
	config, err := config.ReadConfig()
	logging.Setup() // logging should be set up with logging config before logging a possible error in the config, weird, I know
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	if len(os.Args) < 2 {
		log.Fatal().Msgf("Usage: go run %s <scan my tesla logfile>", os.Args[0])
	}

	log.Info().Str("version", version.Version).Msg("Setting up")

	influxdb.Setup(config.InfluxDB)

	file := os.Args[1]

	r, err := regexp.Compile(`.* ([0-9]{4}-[0-9][0-9]?-[0-9][0-9]? [0-9][0-9]?-[0-9][0-9]?-[0-9][0-9]?)\.csv`)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to compile regex")
	}
	matches := r.FindStringSubmatch(file)
	if len(matches) < 2 {
		log.Fatal().Msg("Failed to extract timestamp from file name")
	}

	log.Info().Str("file", file).Str("timestamp", matches[1]).Msg("Reading data")

	data, err := reader.ReadCSV(file)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read CSV")
	}

	log.Info().Msg("Processing")

	fileStartTs, err := time.Parse("2006-01-02 15-04-05", matches[1])
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse timestamp from file name")
	}

	var columns []string
	valueCount := 0
	for lineNumber, line := range data {
		if lineNumber == 0 {
			columns = line
			continue
		}
		var timestamp time.Time
		var values []influxdb.Value
		for columnIndex, column := range line {

			switch columns[columnIndex] {
			case "Time":
				val, err := strconv.ParseInt(column, 10, 64)
				if err != nil {
					log.Fatal().Err(err).Str("timestamp", column).Msg("Malformed timestamp")
				}
				timestamp = fileStartTs.Add(time.Duration(val) * time.Millisecond)
			default:
				if column != "" && column != "Infinity" {
					value, err := strconv.ParseFloat(column, 64)
					if err != nil {
						log.Fatal().Err(err).Str("value", column).Msg("Malformed value")
					}
					values = append(values, influxdb.Value{Field: columns[columnIndex], Value: value})
					valueCount++
				}
			}
		}

		if timestamp.IsZero() {
			log.Warn().Int("line_number", lineNumber).Msg("Encountered a line with no timestamp")
			continue
		} else if len(values) == 0 {
			log.Warn().Int("line_number", lineNumber).Msg("Encountered a line with no values")
			continue
		} else {
			err = influxdb.Send(timestamp, values)
			if err != nil {
				log.Error().Err(err).Int("line_number", lineNumber).Msg("Failed to write data to InfluxDB")
			}
		}
	}

	log.Info().Int("value_count", valueCount).Msg("Processed data, finishing writing to InfluxDB")

	influxdb.Close()

	log.Info().Msg("Done")
}
