package profiler

import (
	"Indexer/file"
	"fmt"
	"os"
)

func HandleProfilerFile(fileName string, target string) (*os.File, error) {
	if fileName == "" {
		fmt.Printf("No %v Profile name provided, proceeding.\n", target)
		return nil, nil
	}

	profileFile, err := file.CreateProfileFile(fileName)
	if err != nil {
		return nil, err
	}

	return profileFile, nil
}
