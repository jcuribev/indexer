package file

import (
	"Indexer/json_manager"
	"archive/tar"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const malformedFilesDir = "./MalformedEmails/"
const indexEmailsDir = "./IndexEmails/emails"

func WriteEmailToFile(json []byte, jsonFile *os.File) error {

	if _, err := jsonFile.Write(json); err != nil {
		fmt.Println("Write email to file failed")
		return err
	}

	return nil
}

func SeparateNewEntryWithComma(jsonFile *os.File) error {
	if _, err := jsonFile.Write([]byte(",\n")); err != nil {
		fmt.Println("Add new line to file failed")
		return err
	}

	return nil
}

func CreateJsonFile(fileNumber int, jsonFile *os.File) (*os.File, error) {
	if jsonFile != nil {
		return jsonFile, nil
	}

	fileDir, err := filepath.Abs(indexEmailsDir + strconv.Itoa(fileNumber) + ".ndjson")

	if err != nil {
		fmt.Println("Find JSON filepath failed")
		return nil, err
	}

	jsonFile, err = os.Create(fileDir)

	if err != nil {
		fmt.Println("Create JSON failed")
		return nil, err
	}

	if err := json_manager.InitFile(jsonFile); err != nil {
		return nil, err
	}

	return jsonFile, nil
}

func StoreMalformedFile(tarHeader *tar.Header, fileContent []byte) error {

	fileDir, err := filepath.Abs(malformedFilesDir + tarHeader.FileInfo().Name())

	if err != nil {
		fmt.Println("Find malformed emails filepath failed")
		return err
	}

	file, err := os.Create(fileDir)

	if err != nil {
		fmt.Println("Create malformed email file failed")
		return err
	}

	_, err = file.Write(fileContent)

	if err != nil {
		fmt.Println("Write content of malformed file failed")
		return err
	}

	file.Close()

	return nil
}

func CreateProfileFile(fileName string) (*os.File, error) {
	profileFile, err := os.Create(fileName)

	if err != nil {
		return nil, err
	}

	return profileFile, nil
}
