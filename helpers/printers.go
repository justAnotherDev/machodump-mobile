package helpers

import (
	"fmt"
	"justAnotherDev/machodump-mobile/entitlements"
	"strings"
	"unicode"

	"github.com/blacktop/go-macho"
	"github.com/blacktop/go-macho/types"
	"github.com/dustin/go-humanize/english"

	ctypes "github.com/blacktop/go-macho/pkg/codesign/types"
)

// GetFileDetails gets the general file details to console
func GetFileDetails(m *macho.File)(string) {
	return fmt.Sprintf("File Details:\n"+
		"\tMagic: %s\n"+
		"\tType: %s\n"+
		"\tCPU: %s, %s\n"+
		"\tCommands: %d (Size: %d)\n"+
		"\tFlags: %s\n"+
		"\tUUID: %s\n",
		m.FileHeader.Magic,
		m.FileHeader.Type,
		m.FileHeader.CPU, m.FileHeader.SubCPU.String(m.FileHeader.CPU),
		m.FileHeader.NCommands,
		m.FileHeader.SizeCommands,
		m.FileHeader.Flags.Flags(),
		m.UUID())
}

// GetLibs gets the loaded libraries to console
func GetLibs(libs []string)(string) {
	output := fmt.Sprintf("File imports %s\n", english.Plural(len(libs), "library:", "libraries:"))

	for i, lib := range libs {
		output += fmt.Sprintf("\t%d: %q\n", i, lib)
	}

	return output
}

// GetLoads gets the interesting load commands to console
func GetLoads(loads []macho.Load)(string) {
	output := fmt.Sprintf("File has %s. Interesting %s:\n", english.Plural(len(loads), "load command", "load commands"), english.PluralWord(len(loads), "command", "commands"))

	for i, load := range loads {
		switch load.Command() {
		case types.LC_VERSION_MIN_IPHONEOS:
			fallthrough
		case types.LC_ENCRYPTION_INFO:
			fallthrough
		case types.LC_ENCRYPTION_INFO_64:
			fallthrough
		case types.LC_SOURCE_VERSION:
			output += fmt.Sprintf("\tLoad %d (%s): %s\n", i, load.Command(), load.String())
		}
	}

	return output
}

// GetCDs gets the code directory details to console
func GetCDs(CDs []ctypes.CodeDirectory)(string) {
	output := fmt.Sprintf("Binary has %s:\n", english.Plural(len(CDs), "Code Directory", "Code Directories"))

	for i, dir := range CDs {
		output += fmt.Sprintf("\tCodeDirectory %d:\n", i)

		output += fmt.Sprintf("\t\tIdent: \"%s\"\n", dir.ID)

		if len(dir.TeamID) > 0 && isASCII(dir.TeamID) {
			output += fmt.Sprintf("\t\tTeam ID: %q\n", dir.TeamID)
		}

		output += fmt.Sprintf("\t\tCD Hash: %s\n", dir.CDHash)
		output += fmt.Sprintf("\t\tCode slots: %d\n", len(dir.CodeSlots))
		output += fmt.Sprintf("\t\tSpecial slots: %d\n", len(dir.SpecialSlots))

		for _, slot := range dir.SpecialSlots {
			output += fmt.Sprintf("\t\t\t%s\n", slot.Desc)
		}
	}

	return output
}

// GetRequirements gets the requirement sections to console
func GetRequirements(reqs []ctypes.Requirement)(string) {
	output := fmt.Sprintf("Binary has %s:\n", english.Plural(len(reqs), "requirement", "requirements"))

	for i, req := range reqs {
		output += fmt.Sprintf("\tRequirement %d (%s): %s\n", i, req.Type, req.Detail)
	}

	return output
}

// GetEnts gets the entitlements to console
func GetEnts(ents *entitlements.EntsStruct)(string) {

	if ents == nil {
		return "Binary has no entitlements\n"
	}

	entries := false
	output := ""
	// print the boolean entries
	if ents.BooleanValues != nil && len(ents.BooleanValues) > 0 {
		output += fmt.Sprintf("Binary has %s:\n", english.Plural(len(ents.BooleanValues), "boolean entitlement", "boolean entitlements"))

		for _, ent := range ents.BooleanValues {
			output += fmt.Sprintf("\t%s: %t\n", ent.Name, ent.Value)
		}

		entries = true
	}

	// print the string entries
	if ents.StringValues != nil && len(ents.StringValues) > 0 {
		output += fmt.Sprintf("Binary has %s:\n", english.Plural(len(ents.StringValues), "string entitlement", "string entitlements"))

		for i, ent := range ents.StringValues {
			output += fmt.Sprintf("\t%d %s: %q\n", i, ent.Name, ent.Value)
		}

		entries = true
	}

	// print the integer entries
	if ents.IntegerValues != nil && len(ents.IntegerValues) > 0 {
		output += fmt.Sprintf("Binary has %s:\n", english.Plural(len(ents.IntegerValues), "integer entitlement", "integer entitlements"))

		for i, ent := range ents.IntegerValues {
			output += fmt.Sprintf("\t%d %s: %d\n", i, ent.Name, ent.Value)
		}

		entries = true
	}

	// print the string array entries
	if ents.StringArrayValues != nil && len(ents.StringArrayValues) > 0 {
		output += fmt.Sprintf("Binary has %s:\n", english.Plural(len(ents.StringArrayValues), "string array entitlement", "string array entitlements"))

		for i, ent := range ents.StringArrayValues {

			valueList := ""

			for _, str := range ent.Values {
				valueList = valueList + str + ", "
			}

			valueList = strings.TrimSuffix(valueList, ", ")

			output += fmt.Sprintf("\t%d %s: [%q]\n", i, ent.Name, valueList)
		}

		entries = true
	}

	if !entries {
		output += fmt.Sprintf("Binary has no entitlements\n")
	}

	return output
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
