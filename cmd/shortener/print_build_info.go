package main

import (
	"fmt"
	"os"
)

const buildInfoDefaultValue = "N/A"

var (
	buildVersion = buildInfoDefaultValue
	buildDate    = buildInfoDefaultValue
	buildCommit  = buildInfoDefaultValue
)

func printBuildInfo() {
	if buildVersion == buildInfoDefaultValue {
		if buildVersionEnv, ok := os.LookupEnv("VERSION"); ok {
			buildVersion = buildVersionEnv
		}
	}

	if buildDate == buildInfoDefaultValue {
		if buildDateEnv, ok := os.LookupEnv("BUILD_DATE"); ok {
			buildDate = buildDateEnv
		}
	}

	if buildCommit == buildInfoDefaultValue {
		if buildCommitEnv, ok := os.LookupEnv("COMMIT_HASH"); ok {
			buildCommit = buildCommitEnv
		}
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
