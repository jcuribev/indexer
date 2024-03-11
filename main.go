package main

import (
	"Indexer/Indexer"
	"Indexer/file"
	"Indexer/profiler"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

var (
	errMissingFileFlag = errors.New("missing file argument")
)

func main() {
	if err := runAndProfileIndexer(); err != nil {
		fmt.Printf("\nError: %v\n", err)
		os.Exit(1)
	}
}

func runAndProfileIndexer() error {
	sourceFile, cpuProfileName, memProfileName, err := parseFlags()
	if err != nil {
		return err
	}

	cpuFile, err := profiler.HandleProfilerFile(*cpuProfileName, "cpu")
	if err != nil {
		return err
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		return errors.New("could not start CPU profile: " + err.Error())
	}
	defer pprof.StopCPUProfile()

	memoryFile, err := profiler.HandleProfilerFile(*memProfileName, "memory")
	if err != nil {
		return err
	}

	defer memoryFile.Close()

	if err := processSourceFile(*sourceFile); err != nil {
		return err
	}

	runtime.GC()
	if err := pprof.WriteHeapProfile(memoryFile); err != nil {
		return errors.New("could not write memory profile: " + err.Error())
	}

	return nil
}

func processSourceFile(sourceFile string) error {
	tgzFile, err := file.OpenSourceFile(sourceFile)
	if err != nil {
		return err
	}
	defer tgzFile.Close()

	tarReader, err := file.GetTgzReader(tgzFile)
	if err != nil {
		return err
	}

	if err := Indexer.IterateTarReader(tarReader); err != nil {
		return err
	}

	if err := Indexer.IndexEmailsToDatabase(); err != nil {
		return err
	}

	return nil
}

func parseFlags() (*string, *string, *string, error) {
	sourceFile := flag.String("file", "", "read file")
	cpuProfileName := flag.String("cpuprofile", "", "write cpu profile to `file`")
	memProfileName := flag.String("memprofile", "", "write memory profile to `file`")

	flag.Parse()

	if *sourceFile == "" {
		return nil, nil, nil, errMissingFileFlag
	}
	return sourceFile, cpuProfileName, memProfileName, nil
}
