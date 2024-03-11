package Indexer

import (
	"Indexer/file"
	"fmt"
	"net/http"
	"os"
)

const bulkAddress = "http://localhost:4080/api/_bulkv2"
const directory = "./IndexEmails/"

func IndexEmailsToDatabase() error {

	filesInfo, err := file.ReadFilesFromDirectory(directory)

	if err != nil {
		return err
	}

	for i := range filesInfo {

		file, err := os.Open(directory + filesInfo[i].Name())
		if err != nil {
			println("couldn't open file:" + file.Name())
			fmt.Printf("err: %v\n", err)
			continue
		}
		defer file.Close()

		request, err := http.NewRequest("POST", bulkAddress, file)

		if err != nil {
			return err
		}

		request.SetBasicAuth(os.Getenv("ZINC_FIRST_ADMIN_USER"), os.Getenv("ZINC_FIRST_ADMIN_PASSWORD"))

		response, err := http.DefaultClient.Do(request)

		if err != nil {
			println("couldn't POST file:" + file.Name())
			fmt.Printf("err: %v\n", err)
		}
		defer response.Body.Close()
	}

	return nil
}
