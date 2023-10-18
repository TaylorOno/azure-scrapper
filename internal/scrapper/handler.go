package scrapper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func Handle(w http.ResponseWriter, _ *http.Request) {
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
	_, err := w.Write(responseJson)
	if err != nil {
		log.Printf("failed to write response: %v", err)
		return
	}
}
