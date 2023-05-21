package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	listenAddr := ":9090"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	http.HandleFunc("/scrapper", scrape)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func scrape(w http.ResponseWriter, r *http.Request) {
	resp := InvokeResponse{
		Outputs:     map[string]resData{"res": {}},
		Logs:        []string{},
		ReturnValue: "",
	}

	sub := os.Getenv("AZURE_SUBSCRIPTION")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		resp.Logs = append(resp.Logs, fmt.Sprintf("failed to obtain a credential: %v", err))
		writeJSON(w, resp, http.StatusInternalServerError)
		return
	}

	scrapper, err := NewScrapper(cred, sub)
	if err != nil {
		resp.Logs = append(resp.Logs, fmt.Sprintf("unable to initialize scrapper: %v", err))
		writeJSON(w, resp, http.StatusInternalServerError)
		return
	}

	if err = scrapper.Run(); err != nil {
		resp.Logs = append(resp.Logs, fmt.Sprintf("scrapper failed: %v", err))
		writeJSON(w, resp, http.StatusInternalServerError)
		return
	}

	writeJSON(w, resp, http.StatusOK)
	return
}

func consoleHandler[T any](t *T) error {
	return json.NewEncoder(os.Stdout).Encode(t)
}

type resData = map[string]interface{}

type InvokeResponse struct {
	Outputs     map[string]resData
	Logs        []string
	ReturnValue interface{}
}

func writeJSON(w http.ResponseWriter, result InvokeResponse, code int) {
	result.Outputs["res"]["statuscode"] = code
	result.Outputs["res"]["headers"] = map[string]string{"Content-Type": "application/json"}
	result.Outputs["res"]["body"] = fmt.Sprintf("{\"status\":\"%v\"}", code)

	responseJson, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}
