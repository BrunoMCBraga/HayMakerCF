package main

import (
	"errors"
	"fmt"

	"github.com/BrunoMCBraga/HayMakerCF/commandlinegenerators"
	"github.com/BrunoMCBraga/HayMakerCF/commandlineprocessors"
	"github.com/BrunoMCBraga/HayMakerCF/globalstringsproviders"
)

func main() {

	commandlinegenerators.PrepareCommandLineProcessing()

	fmt.Println(globalstringsproviders.GetMenuPictureString())

	commandlinegenerators.ParseCommandLine()
	parameters := commandlinegenerators.GetParametersDict()
	processCommandLineProcessorError := commandlineprocessors.ProcessCommandLine(parameters)

	if processCommandLineProcessorError != nil {
		fmt.Println(errors.New("HayMakerCF->main->commandlineprocessors.ProcessCommandLine:" + processCommandLineProcessorError.Error()))
	}

}
