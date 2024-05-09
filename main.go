package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SuTech-JP/raas-client-go"
	"github.com/gorilla/mux"
    "github.com/rs/cors"
)

const KEY_APP = "RaaS-Application"
const KEY_LND = "RaaS-Landscape"
const KEY_TKN = "RaaS-Token"

// ---------------
// 実装サンプル：RaaSRestClientを用いたRaaSとの通信
// ---------------
func main() {
	r := mux.NewRouter()
	// only for test
	os.Setenv(KEY_APP, "アプリケーション名")
	os.Setenv(KEY_LND, "ランドスケープ名")
	os.Setenv(KEY_TKN, "付与されたトークン")
	// １：ExternalSessionの確立
	r.HandleFunc("/raas/{msa}/session", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		msa := vars["msa"]
		config, configErr := raas.NewRaaSConnectConfig(os.Getenv(KEY_APP), os.Getenv(KEY_LND), os.Getenv(KEY_TKN))
		context, contextErr := raas.NewRaasUserContext("tenant", "sub")
		if configErr == nil && contextErr == nil {
			const backUrl = "返却時のURL"
			const subUrl = "subUrl"
			const subDomain = "subDomain" //override
			//msa must be [report] or [datatraveler]
			//extSessionはExtSessionという特殊な構造体が定義されています。
			extSession, err := raas.RaaSRestClient(*config, *context).CreateExternalSession(msa, backUrl, subUrl, subDomain)
			if err == nil {
				//encoding into JSON
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(extSession)
			} else {
				log.Fatalf("Errors: %v", err.Error())
			}
		} else {
			log.Fatal("Config or Context is invalid...")
		}
	})

	//GETの送信例１
	r.HandleFunc("/raas/report/layout/{application}/{schema}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		application := vars["application"]
		schema := vars["schema"]
		config, configErr := raas.NewRaaSConnectConfig(os.Getenv(KEY_APP), os.Getenv(KEY_LND), os.Getenv(KEY_TKN))
		context, contextErr := raas.NewRaasUserContext("tenant", "sub")
		if configErr == nil && contextErr == nil {
			requestUrl := fmt.Sprintf("/raas/report/layout/%s/%s", application, schema)
			//Get,Put,Post,DeleteのHTTPメソッドは、map[string]anyで返却されます
			resultMap, err := raas.RaaSRestClient(*config, *context).Get(requestUrl, nil)
			if err == nil {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resultMap)
			} else {
				log.Fatalf("Errors: %v", err.Error())
			}
		} else {
			log.Fatal("Config or Context is invalid...")
		}
	})

	//GETの送信例２
	r.HandleFunc("/raas/report/result/{targetId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		targetId := vars["targetId"]
		config, configErr := raas.NewRaaSConnectConfig(os.Getenv(KEY_APP), os.Getenv(KEY_LND), os.Getenv(KEY_TKN))
		context, contextErr := raas.NewRaasUserContext("tenant", "sub")
		if configErr == nil && contextErr == nil {
			requestUrl := fmt.Sprintf("/raas/report/result/%s", targetId)
			//Get,Put,Post,DeleteのHTTPメソッドは、map[string]anyで返却されます
			resultMap, err := raas.RaaSRestClient(*config, *context).Get(requestUrl, nil)
			if err == nil {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resultMap)
			} else {
				log.Fatalf("Errors: %v", err.Error())
			}
		} else {
			log.Fatal("Config or Context is invalid...")
		}
	})

	//サーバの設定
	log.Println("Server starting on http://localhost:8080")
    handler := cors.Default().Handler(r)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
