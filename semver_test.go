package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/blang/semver"
)

func TestBump(t *testing.T) {
	positiveTestTables := []struct {
		input    semver.Version
		options  bumpOptions
		expected semver.Version
	}{
		{
			getSemver("1.2.3"),
			bumpOptions{
				major: false,
				minor: false,
				patch: false,
				meta:  ""},
			getSemver("1.2.4")},
		{
			getSemver("1.2.3"),
			bumpOptions{
				major: false,
				minor: false,
				patch: true,
				meta:  ""},
			getSemver("1.2.4")},
		{
			getSemver("1.2.3"),
			bumpOptions{
				major: false,
				minor: true,
				patch: false,
				meta:  ""},
			getSemver("1.3.0")},
		{
			getSemver("1.2.3"),
			bumpOptions{
				major: true,
				minor: false,
				patch: false,
				meta:  ""},
			getSemver("2.0.0")},
		{
			getSemver("1.2.3+6.4"),
			bumpOptions{
				major: false,
				minor: false,
				patch: false,
				meta:  "7.3"},
			getSemver("1.2.4+7.3")},
		{
			getSemver("1.2.3+6.4"),
			bumpOptions{
				major: false,
				minor: false,
				patch: true,
				meta:  "7.3"},
			getSemver("1.2.4+7.3")},
		{
			getSemver("1.2.3+6.4"),
			bumpOptions{
				major: false,
				minor: true,
				patch: false,
				meta:  "7.3"},
			getSemver("1.3.0+7.3")},
		{
			getSemver("1.2.3+6.4"),
			bumpOptions{
				major: true,
				minor: false,
				patch: false,
				meta:  "7.3"},
			getSemver("2.0.0+7.3")},
	}

	for _, positiveTable := range positiveTestTables {
		newSemverErr := bump(&positiveTable.input, positiveTable.options)

		if newSemverErr != nil {
			t.Errorf("Could non bump [%q]", positiveTable.input)
		}

		if positiveTable.input.Compare(positiveTable.expected) != 0 {
			t.Errorf(`Expected: [%q], got: [%q], when:
            major: [%t],
            minor: [%t],
            patch: [%t],
            meta: [%s]`,
				positiveTable.expected,
				positiveTable.input,
				positiveTable.options.major,
				positiveTable.options.minor,
				positiveTable.options.patch,
				positiveTable.options.meta)
		}
	}

	negativeTestTables := []struct {
		input   semver.Version
		options bumpOptions
	}{
		{
			getSemver("1.2.3"),
			bumpOptions{
				major: false,
				minor: true,
				patch: true,
				meta:  ""}},
		{
			getSemver("1.2.3"),
			bumpOptions{
				major: true,
				minor: false,
				patch: true,
				meta:  ""}},
		{
			getSemver("1.2.3"),
			bumpOptions{
				major: true,
				minor: true,
				patch: false,
				meta:  ""}},
		{
			getSemver("1.2.3"),
			bumpOptions{
				major: true,
				minor: true,
				patch: true,
				meta:  ""}},
		{
			getSemver("1.2.3+6.4"),
			bumpOptions{
				major: false,
				minor: true,
				patch: true,
				meta:  "7.3"}},
		{
			getSemver("1.2.3+6.4"),
			bumpOptions{
				major: true,
				minor: false,
				patch: true,
				meta:  "7.3"}},
		{
			getSemver("1.2.3+6.4"),
			bumpOptions{
				major: true,
				minor: true,
				patch: false,
				meta:  "7.3"}},
		{
			getSemver("1.2.3+6.4"),
			bumpOptions{
				major: true,
				minor: true,
				patch: true,
				meta:  "7.3"}},
	}

	for _, negativeTable := range negativeTestTables {
		newSemverErr := bump(&negativeTable.input, negativeTable.options)

		if multiBumpError.Error() != newSemverErr.Error() {
			t.Errorf(`Expected error [%s] but got [%s], when:
            major: [%t],
            minor: [%t],
            patch: [%t],
            meta: [%s]`,
				multiBumpError.Error(),
				newSemverErr.Error(),
				negativeTable.options.major,
				negativeTable.options.minor,
				negativeTable.options.patch,
				negativeTable.options.meta)
		}
	}
}

func TestReadBumpTypesFromString(t *testing.T) {
	testTables := []struct {
		message        string
		expectedResult []bumpType
	}{
		{"This includes incompatible API changes. +major", []bumpType{bumpMajor}},
		{"Added new feature. +minor", []bumpType{bumpMinor}},
		{"Fixed a bug. +patch", []bumpType{bumpPatch}},
		{"I thought this was +patch. But the new commit shows +minor", []bumpType{bumpPatch, bumpMinor}},
		{"I thought this was +patch. But the new commit shows +major", []bumpType{bumpPatch, bumpMajor}},
		{"I thought this was +minor. But the new commit shows +major", []bumpType{bumpMinor, bumpMajor}},
		{"I thought this was +major or +minor. But the new commit shows +patch", []bumpType{bumpMajor, bumpMinor, bumpPatch}},
		{`Making a small bugfix, thought to be +patch.
        Adding functinality in a backwards compatible manner +minor.
        Made additional incompatible API changes. +major`, []bumpType{bumpPatch, bumpMinor, bumpMajor}},
	}

	for _, table := range testTables {
		bumpsFound := readBumpTypesFromString(table.message)

		if !reflect.DeepEqual(table.expectedResult, bumpsFound) {
			t.Errorf("Expected: [%v], got: [%v], when: message: [%s]", table.expectedResult, bumpsFound, table.message)
		}
	}
}

func TestFindMaxBump(t *testing.T) {
	testTables := []struct {
		bumpsPassed     []bumpType
		expectedMaxBump bumpType
	}{
		{[]bumpType{bumpMajor}, bumpMajor},
		{[]bumpType{bumpMinor}, bumpMinor},
		{[]bumpType{bumpPatch}, bumpPatch},
		{[]bumpType{bumpPatch, bumpMinor}, bumpMinor},
		{[]bumpType{bumpPatch, bumpMajor}, bumpMajor},
		{[]bumpType{bumpMinor, bumpMajor}, bumpMajor},
		{[]bumpType{bumpMajor, bumpMinor, bumpPatch}, bumpMajor},
		{[]bumpType{bumpPatch, bumpMinor, bumpMajor}, bumpMajor},
		{[]bumpType{bumpPatch, bumpMinor, bumpMajor, bumpMajor, bumpMinor, bumpMajor}, bumpMajor},
	}

	for _, table := range testTables {
		foundMaxBump := findMaxBump(table.bumpsPassed)

		if table.expectedMaxBump != foundMaxBump {
			t.Errorf("Expected: [%d], got: [%d], when: bumpsPassed: [%v]", table.expectedMaxBump, foundMaxBump, table.bumpsPassed)
		}
	}
}

func getSemver(version string) semver.Version {
	semVer, err := semver.Make(version)

	if err != nil {
		fmt.Printf("Could not make a semver for [%s]", version)
	}
	return semVer

}
