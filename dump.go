package Macho

import (
	"justAnotherDev/machodump-mobile/entitlements"
	"justAnotherDev/machodump-mobile/helpers"
	"os"
	"fmt"

	"github.com/blacktop/go-macho"
)

func DumpFile(input string)(string) {
	output := ""
	if _, err := os.Stat(input); os.IsNotExist(err) {
		output += fmt.Sprintf("Fatal: input does not exist: %s", input)
		return output
	}

	// process file details
	output += fmt.Sprintf("Parsing file %q.", input)

	file, err := os.Open(input)

	if err != nil {
		output += err.Error()
		return output
	}

	machoFile, err := macho.NewFile(file)

	if err != nil {
		output += fmt.Sprintf("Error parsing file as Macho-O: %q. Exiting.", err.Error())
		return output
	}

	// get code signature
	cd := machoFile.CodeSignature()

	if cd == nil {
		// print file details
		output += helpers.GetFileDetails(machoFile)
		output += helpers.GetLibs(machoFile.ImportedLibraries())
		output += helpers.GetLoads(machoFile.Loads)

		output += "No code signing section in binary, exiting"
		return output
	}

	// print file details
	output += helpers.GetFileDetails(machoFile)
	output += helpers.GetLibs(machoFile.ImportedLibraries())
	output += helpers.GetLoads(machoFile.Loads)

	// print the details
	output += helpers.GetCDs(cd.CodeDirectories)

	// parse the CMS sig, if it's there
	if len(cd.CMSSignature) > 0 {
		helpers.ParseCMSSig(cd.CMSSignature)
	}

	output += helpers.GetRequirements(cd.Requirements)


	// get array of entitlements
	ents, err := entitlements.GetEntsFromXMLString(cd.Entitlements)

	if err != nil {
		output += fmt.Sprintf("Errored when trying to extract ents: %s", err.Error())
	}
	output += helpers.GetEnts(ents)

	output += "Fin."
	return output
}