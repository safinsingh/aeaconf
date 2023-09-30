package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
)

func fatalErr(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	iniFile, err := ini.Load("testing/scoring.ini")
	if err != nil {
		fatalErr(errors.Wrap(err, "failed to read config"))
	}

	var config Config
	err = iniFile.MapTo(&config)
	if err != nil {
		fatalErr(errors.Wrap(err, "failed to map ini config to struct"))
	}

	customConditionsSection := iniFile.Section("custom_conditions")
	customConditions := make(map[string]string)
	for key := range customConditionsSection.KeysHash() {
		value := customConditionsSection.Key(key).String()
		customConditions[key] = value
	}

	checksSection := iniFile.Section("checks")
	for key := range checksSection.KeysHash() {
		value := checksSection.Key(key).String()

		// if : at the very least before the last character in the string...
		if idx := strings.IndexByte(value, ':'); idx >= 0 {
			points, err := strconv.Atoi(strings.TrimSpace(value[:idx]))
			if err != nil {
				fatalErr(errors.Wrapf(err, "invalid point value for vuln %s", key))
			}

			// portion after delimeter = condition
			cond := ParseConditionFromStringWith(value[idx+1:], customConditions)
			config.Checks = append(config.Checks, Check{Message: key, Points: points, Cond: cond})
		}
	}

	spew.Dump(config)
}
