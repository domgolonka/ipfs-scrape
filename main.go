package main

import (
	"domgolonka/blockparty/database"
	"domgolonka/blockparty/model"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

const address = "https://blockpartyplatform.mypinata.cloud/ipfs/"

var listCids = []string{
	"bafkreifovhtvvrx5jmo2b4ne2hoyk4t3c276jc7weva5s57ilupiiuqg2y",
	"bafkreifwbyviygstiqmjcijju33r6scctuyxciqiepcrs2ym2bbpf7c7rq",
	"bafkreicpcdl32e5l4kusphczsswo3wcrjo7fyt4iktgnhctdnhim6o3xwe",
	"bafkreicjwnp5xostqloqtfbwawhzii7nekddbzppi77rqf4r4nw4wrsd4q",
	"bafkreib7xfz4cyv7gp2enc3zcxj6us42dz3v6m6aan4abueqzudcfjvmby",
	"bafkreih3ofrl5cdne3o7bmabuugi2scqbvk6pqxaskbl4irgi2i7345vvy",
	"bafkreihztk3kzqfhasvow6elspu6vqurhc4xgxew6efa45cbftx5vo3li4",
	"bafkreifzz3xqlmxuk2mapgrv532jryigp5aqgoas6r7eksfl7ccvfsc3ge",
	"bafkreiekw2qbsfx7s7lyzcoabyskdpheteab3ogzzcrxboqzydv6fqhix4",
	"bafkreieizajnzclz3fda52pmsdb7yxaf73aszq6kakqntpmawfa5bhvfau",
	"bafkreicuerb3ixr7hni77pr6pnmsavades3wxcuwugy4umvsjqstvpoelu",
	"bafkreieuhkm2v4gzhuf52cfi2ob47ve4p6lo25rlbk7ckrc4sywivc2npi",
	"bafkreigfin2zgkd45bdelxxmizpa7dvuvgelvxd5oodexc52uruu3cg4pi",
	"bafkreihjkq3w2gkplz2fclroceqco7nit4ekquhnvfbj7iikt6izuswlva",
	"bafkreibwqdysbhjc4uwratubu5dpsolcnyx7mo3pwey7gfnbhm2wo2eoze",
	"bafkreicuerb3ixr7hni77pr6pnmsavadeshawnuwugy4umvswqstv55el9",
	"bafkreigt5pduy2qtxou3nvfuwq4xyfziawgse7bwzj5pjfesnhg6kvaj7y",
	"bafkreig4qyrpmt2gmvfcxnz4gqbodxjp2dsy6sgllkupnup4gys7wqnt2y",
	"QmamdCZfHLy18hAix12h7ntu41vNjjMGQCj2gEqijiEcRs",
	"bafkreibrdju7aievsss6lq2iem3kzqvskz6ig3zk2psdfy7bejhdp72qzy",
	"bafkreiailk323slyxvcrc7k3gofhtnsp7d6k5s5c2li3fa3uscur5yznlm",
	"QmT3Xz6QBLLre9KfCLthUWVjZAHhMzFDx34MiSKaoUh3pi/1464.json",
	"bafybeia67q6eabx2rzu6datbh3rnsoj7cpupudckijgc5vtxf46zpnk2t4/3885",
}

func scrap(cid string) (*model.Ipfs, error) {
	resp, err := http.Get(address + cid)
	if err != nil {
		// skip if we can't get the cid (Invalid url)
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
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

	for _, cid := range listCids {
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			data, err := scrap(cid)
			if err != nil {
				return
			}
			data.Cid = cid
			ipfs = append(ipfs, data)
		}(cid)

	}
	wg.Wait()

	err := db.InsertIpfs(ipfs)
	if err != nil {
		log.Fatalln(err)
	}
}

func tokens(db *database.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ipfs, err := db.SelectIpfs()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		bodyMarshal(w, http.StatusOK, ipfs)

	}
}

func cidTokens(db *database.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cid := mux.Vars(r)["cid"]
		ipfs, err := db.SelectIpfsByCid(cid)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid cid")
			return
		}
		bodyMarshal(w, http.StatusOK, ipfs)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	bodyMarshal(w, code, map[string]string{"error": message})
}

func bodyMarshal(w http.ResponseWriter, code int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(payload)
	if err != nil {
		w.Write([]byte(`{"success":false,"message":"Some internal error occurred"}`))
		return errors.New("some internal error occurred")
	}
	w.WriteHeader(code)
	w.Write(resp)
	return nil
}

func handleRequests(db *database.Client) {
	http.HandleFunc("/tokens", tokens(db))
	http.HandleFunc("/tokens/{cid}", cidTokens(db))
	log.Fatal(http.ListenAndServe(":8084", nil))
}
func handleDatabase() (*database.Client, error) {
	godotenv.Load()

	dbuser := os.Getenv("DB_USER")
	dbpass := os.Getenv("DB_PASS")
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	return database.NewClient(dbhost, dbport, dbuser, dbpass)
}

func main() {
	db, err := handleDatabase()
	if err != nil {
		panic(err)
	}
	scrapeIpfs(db)
	handleRequests(db)
}
