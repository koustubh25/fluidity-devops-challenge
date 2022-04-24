package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/fluidity-money/devops-hiring-challenge/lib"
)

var regexpExtractComputeUnits = regexp.MustCompile(`consumed (\d+)`)

func extractComputes(websocketLogs devops_hiring_challenge.SolanaWebsocketLog) (units []uint64) {
	for _, log := range websocketLogs.Params.Result.Value.Logs {
		matches := regexpExtractComputeUnits.FindStringSubmatch(log)

		if len(matches) == 0 {
			continue
		}

		unitString := matches[0][9:]

		// we should be fine if the regex matched...

		unit, err := strconv.ParseUint(unitString, 10, 64)

		if err != nil {
			panic(fmt.Sprintf(
				"matched a value that was not a uint64! broken sexp! string is %#s, %v",
				unitString,
				err,
			))
		}

		units = append(units, unit)
	}

	return
}

func computeSummedUnits(units []uint64) (average uint64) {
	unitsLen := len(units)

	if unitsLen == 0 {
		return 0
	}

	for _, unit := range units {
		average += unit
	}

	return average / uint64(len(units))
}
