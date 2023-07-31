package main

import (
	"domgolonka/blockparty/database"
	"domgolonka/blockparty/model"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

const address = "https://blockpartyplatform.mypinata.cloud/ipfs/"

func readFromCSF(file string) ([][]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func scrap(cid string) (*model.Ipfs, error) {
	resp, err := http.Get(address + cid)
	if err != nil {
		// skip if we can't get the cid (Invalid url)
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data *model.Ipfs
	err = json.Unmarshal(b, &data)
	if err != nil {
		// skip if the format is invalid
		return nil, err
	}
	return data, nil
}

func scrapeIpfs(db *database.Client) {
	var ipfs []*model.Ipfs
	var wg sync.WaitGroup

	list, err := readFromCSF("ipfs_cids.csv")
	if err != nil {
		log.Fatalln(err)
	}
	for _, c := range list {
		wg.Add(1)
		go func(listCids []string) {
			for _, cid := range listCids {

				defer wg.Done()
				data, err := scrap(cid)
				if err != nil {
					return
				}
				data.Cid = cid
				ipfs = append(ipfs, data)
			}
		}(c)

	}
	wg.Wait()

	err = db.InsertIpfs(ipfs)
	if err != nil {
		return // skip if we can't insert the data
	}
}
func healthHandler(db *database.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check the health of the server and return a status code accordingly
		if serverIsHealthy(db) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Server is healthy")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Server is not healthy")
		}
	}
}

func serverIsHealthy(db *database.Client) bool {
	err := db.Ping()
	if err != nil {
		return false
	}
	return true
}

func handleRequests(db *database.Client) {
	http.HandleFunc("/health", healthHandler(db))
	log.Fatal(http.ListenAndServe(":80", nil))
}

func main() {
	db, err := database.NewClient()
	if err != nil {
		panic(err)
	}
	scrapeIpfs(db)
	handleRequests(db)
}
