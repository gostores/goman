package cmd

import (
	"strings"
	"time"

	"github.com/govenue/configurator"
)

// Licenses contains all possible licenses a user can choose from.
var Licenses = make(map[string]License)

// License represents a software license agreement, containing the Name of
// the license, its possible matches (on the command line as given to goman),
// the header to be used with each file on the file's creating, and the text
// of the license
type License struct {
	Name            string   // The type of license in use
	PossibleMatches []string // Similar names to guess
	Text            string   // License text data
	Header          string   // License header for source files
}

func init() {
	// Allows a user to not use a license.
	Licenses["none"] = License{"None", []string{"none", "false"}, "", ""}

	initApache2()
	initMit()
	initBsdClause3()
	initBsdClause2()
	initGpl2()
	initGpl3()
	initLgpl()
	initAgpl()
}

// getLicense returns license specified by user in flag or in config.
// If user didn't specify the license, it returns Apache License 2.0.
//
// TODO: Inspect project for existing license
func getLicense() License {
	// If explicitly flagged, use that.
	if userLicense != "" {
		return findLicense(userLicense)
	}

	// If user wants to have custom license, use that.
	if configurator.IsSet("license.header") || configurator.IsSet("license.text") {
		return License{Header: configurator.GetString("license.header"),
			Text: configurator.GetString("license.text")}
	}

	// If user wants to have built-in license, use that.
	if configurator.IsSet("license") {
		return findLicense(configurator.GetString("license"))
	}

	// If user didn't set any license, use Apache 2.0 by default.
	return Licenses["apache"]
}

func copyrightLine() string {
	author := configurator.GetString("author")

	year := configurator.GetString("year") // For tests.
	if year == "" {
		year = time.Now().Format("2006")
	}

	return "Copyright © " + year + " " + author
}

// findLicense looks for License object of built-in licenses.
// If it didn't find license, then the app will be terminated and
// error will be printed.
func findLicense(name string) License {
	found := matchLicense(name)
	if found == "" {
		er("unknown license: " + name)
	}
	return Licenses[found]
}

// matchLicense compares the given a license name
// to PossibleMatches of all built-in licenses.
// It returns blank string, if name is blank string or it didn't find
// then appropriate match to name.
func matchLicense(name string) string {
	if name == "" {
		return ""
	}

	for key, lic := range Licenses {
		for _, match := range lic.PossibleMatches {
			if strings.EqualFold(name, match) {
				return key
			}
		}
	}

	return ""
}
