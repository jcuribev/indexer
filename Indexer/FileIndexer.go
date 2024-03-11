package Indexer

import (
	"Indexer/email"
	"Indexer/file"
	"Indexer/json_manager"
	"archive/tar"
	"fmt"
	"io"
	"os"
)

type FileInformation struct {
	entriesWritten int
	fileNumber     int
	jsonFile       *os.File
	isFirstEntry   bool
}

const maxEntriesPerJson = 50000

func IterateTarReader(tarReader *tar.Reader) error {

	fileInfo := newFileInformation()

	for {
		tarHeader, err := tarReader.Next()
		if err == io.EOF {
			endJsonFile(fileInfo.jsonFile)
			break
		}

		if err != nil {
			fmt.Println("Getting next TarHaeader failed")
			return err
		}

		if tarHeader.Typeflag != tar.TypeReg {
			continue
		}

		fileInfo.jsonFile, err = file.CreateJsonFile(fileInfo.fileNumber, fileInfo.jsonFile)
		if err != nil {
			return err
		}

		err = createJsonEmailAndWriteToFile(tarReader, &fileInfo)
		if err != nil {
			err := handleEmailError(tarReader, tarHeader)
			if err != nil {
				return err
			}
			continue
		}

		fileInfo = registerNewFileEntry(fileInfo)

		if !hasReachedMaxEntries(fileInfo.entriesWritten) {
			continue
		}

		fileInfo, err = changeFileForNextEntries(fileInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func registerNewFileEntry(fileInfo FileInformation) FileInformation {
	return updateFileInformation(fileInfo.entriesWritten+1, fileInfo.fileNumber, fileInfo.jsonFile, false)
}

func hasReachedMaxEntries(entriesWritten int) bool {
	return entriesWritten >= maxEntriesPerJson
}

func changeFileForNextEntries(fileInfo FileInformation) (FileInformation, error) {
	if err := endJsonFile(fileInfo.jsonFile); err != nil {
		return fileInfo, err
	}

	return updateFileInformation(0, fileInfo.fileNumber+1, nil, true), nil
}
func endJsonFile(file *os.File) error {
	defer file.Close()

	return json_manager.FinishFile(file)
}

func newFileInformation() FileInformation {
	return FileInformation{
		entriesWritten: 0,
		fileNumber:     0,
		jsonFile:       nil,
		isFirstEntry:   true,
	}
}

func updateFileInformation(entriesWritten int, fileNumber int, jsonFile *os.File, isFirstEntry bool) FileInformation {
	return FileInformation{
		entriesWritten: entriesWritten,
		fileNumber:     fileNumber,
		jsonFile:       jsonFile,
		isFirstEntry:   isFirstEntry,
	}
}

func createJsonEmailAndWriteToFile(tarReader *tar.Reader, fileInfo *FileInformation) error {
	jsonEmail, err := createJsonEmail(tarReader)
	if err != nil {
		return err
	}

	if fileInfo.isFirstEntry == false {
		if err = file.SeparateNewEntryWithComma(fileInfo.jsonFile); err != nil {
			return err
		}
	}

	return file.WriteEmailToFile(jsonEmail, fileInfo.jsonFile)

}

func createJsonEmail(tarReader *tar.Reader) (jsonEmail []byte, err error) {
	fileContent, err := file.ReadFileContent(tarReader)
	if err != nil {
		return nil, err
	}

	email, err := email.FileContentToEmail(string(fileContent))
	if err != nil {
		return nil, err
	}

	return json_manager.EmailToJson(email)
}

func handleEmailError(tarReader *tar.Reader, tarHeader *tar.Header) error {
	fileContent, err := file.ReadFileContent(tarReader)

	if err != nil {
		return err
	}

	return file.StoreMalformedFile(tarHeader, fileContent)
}
