package main

import (
	"encoding/json"
	"fmt"
	"io"
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

type Arguments struct {
	BackUrl   string
	SubUrl    string
	SubDomain string
}

// ---------------
// 実装サンプル：RaaSRestClientを用いたRaaSとの通信
// ---------------
func main() {
	r := mux.NewRouter()
	// 以下の３つは、環境変数からセットすることを推奨します。本サンプルでも環境変数から取得しています
	// os.Setenv(KEY_APP, "アプリケーション名")
	// os.Setenv(KEY_LND, "ランドスケープ名")
	// os.Setenv(KEY_TKN, "付与されたトークン")
	// １：ExternalSessionの確立
	r.HandleFunc("/raas/{msa}/session", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		msa := vars["msa"]
		config, configErr := raas.NewRaaSConnectConfig(os.Getenv(KEY_APP), os.Getenv(KEY_LND), os.Getenv(KEY_TKN))
		context, contextErr := raas.NewRaasUserContext("tenant", "sub")
		if configErr == nil && contextErr == nil {
			//bodyから値を取得
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}
			//ここではArgumentsにマッピング
			var args Arguments
			if err := json.Unmarshal(body, &args); err != nil {
				http.Error(w, "Error parsing JSON data", http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()
			//引数
			backUrl := args.BackUrl
			subUrl := args.SubUrl
			subDomain := args.SubDomain
			//msaは、"report" か　"datatraveler"
			//resultMapはmap[string]anyで返却されます
			resultMap, err := raas.RaaSRestClient[map[string]any](*config, *context).CreateExternalSession(msa, backUrl, subUrl, subDomain)
			if err == nil {
				//JSONに変換して書き出し
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resultMap)
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
			requestUrl := fmt.Sprintf("/report/layouts/%s/%s", application, schema)
			//Get,Put,Post,DeleteのHTTPメソッドは、map[string]anyで返却されます
			resultMap, err := raas.RaaSRestClient[[]map[string]any](*config, *context).Get(requestUrl, nil)
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
			requestUrl := fmt.Sprintf("/datatraveler/import/logs/%s", targetId)
			//Get,Put,Post,DeleteのHTTPメソッドは、map[string]anyで返却されます
			resultMap, err := raas.RaaSRestClient[map[string]any](*config, *context).Get(requestUrl, nil)
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
