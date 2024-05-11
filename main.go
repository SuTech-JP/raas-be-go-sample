package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/SuTech-JP/raas-client-go"
	"github.com/go-yaml/yaml"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Arguments struct {
	BackUrl   string
	SubUrl    string
	SubDomain string
}

type RaasConfig struct {
	Application string `json:"application" yaml:"application"`
	Landscape   string `json:"landscape" yaml:"landscape"`
	Token       string `json:"token" yaml:"token"`
}

type AppConfig struct {
	RaasConfig RaasConfig `json:"raasConfig" yaml:"raasConfig"`
}

// ---------------
// 実装サンプル：RaaSRestClientを用いたRaaSとの通信
// ---------------
func main() {
	r := mux.NewRouter()

	// 設定ファイルからRaasへの接続情報を取得
	appConfig, appConfigError := loadConfigForYaml()
	if appConfigError != nil {
		log.Fatal("failed to load appConfig...", appConfigError)
	}
	config, configErr := raas.NewRaaSConnectConfig(appConfig.RaasConfig.Application, appConfig.RaasConfig.Landscape, appConfig.RaasConfig.Token)
	if configErr != nil {
		log.Fatal("Config is invalid...", configErr)
	}

	// サンプル：Sessionの発行
	r.HandleFunc("/raas/{msa}/session", func(w http.ResponseWriter, r *http.Request) {
		// RaasUserContextの作成
		context, contextErr := raas.NewRaasUserContext("tenant", "sub")
		if contextErr != nil {
			log.Fatal("Context is invalid...")
		}

		// 引数を準備
		vars := mux.Vars(r)
		msa := vars["msa"]
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var args Arguments
		if err := json.Unmarshal(body, &args); err != nil {
			http.Error(w, "Error parsing JSON data", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// RaasのSession発行APIを実行
		result, err := raas.RaaSRestClient(*config, *context).CreateExternalSession(msa, args.BackUrl, args.SubUrl, args.SubDomain)
		if err != nil {
			log.Fatalf("Errors: %v", err.Error())
		}

		// resultはバイト配列なので、APIドキュメントを参考に適切な形に変換する
		var jsonObject map[string]any //キャストはご自由に
		encodeError := json.Unmarshal(result, &jsonObject)
		if encodeError != nil {
			log.Fatalf("Cannot encode to JSON")
		}
		// JSONに変換してレスポンスに書き出し
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonObject)
	})

	// サンプル：レイアウト一覧の取得
	r.HandleFunc("/raas/report/layout/{application}/{schema}", func(w http.ResponseWriter, r *http.Request) {
		// RaasUserContextの作成
		context, contextErr := raas.NewRaasUserContext("tenant", "sub")
		if contextErr != nil {
			log.Fatal("Context is invalid...")
		}

		// 引数を準備
		vars := mux.Vars(r)
		application := vars["application"]
		schema := vars["schema"]

		// レイアウト取得APIの実行
		requestUrl := fmt.Sprintf("/report/layouts/%s/%s", application, schema)
		result, err := raas.RaaSRestClient(*config, *context).Get(requestUrl, nil)
		if err != nil {
			log.Fatalf("Errors: %v", err.Error())
		}
		// resultはバイト配列なので、APIドキュメントを参考に適切な形に変換する
		var jsonObject []map[string]any //キャストはご自由に
		encodeError := json.Unmarshal(result, &jsonObject)
		if encodeError != nil {
			log.Fatalf("Cannot encode to JSON")
		}
		// JSONに変換してレスポンスに書き出し
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonObject)
	})

	// サンプル：CSVインポートにより作成されたログデータの取得
	r.HandleFunc("/raas/report/result/{targetId}", func(w http.ResponseWriter, r *http.Request) {
		// RaasUserContextの作成
		context, contextErr := raas.NewRaasUserContext("tenant", "sub")
		if contextErr != nil {
			log.Fatal("Config or Context is invalid...")
		}

		// 引数を準備
		vars := mux.Vars(r)
		targetId := vars["targetId"]

		// ログデータ取得APIを実行
		requestUrl := fmt.Sprintf("/datatraveler/import/logs/%s", targetId)
		result, err := raas.RaaSRestClient(*config, *context).Get(requestUrl, nil)
		if err != nil {
			log.Fatalf("Errors: %v", err.Error())
		}

		// resultはバイト配列なので、APIドキュメントを参考に適切な形に変換する
		var jsonObject map[string]any //キャストはご自由に
		encodeError := json.Unmarshal(result, &jsonObject)
		if encodeError != nil {
			log.Fatalf("Cannot encode to JSON")
		}
		// JSONに変換してレスポンスに書き出し
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonObject)
	})

	//サーバの設定
	log.Println("Server starting on http://localhost:8080")
	handler := cors.Default().Handler(r)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// application.yamlの読み込み
func loadConfigForYaml() (*AppConfig, error) {
	f, err := os.Open("application.yaml")
	if err != nil {
		log.Fatal("loadConfigForYaml os.Open err:", err)
		return nil, err
	}
	defer f.Close()

	var cfg AppConfig
	err = yaml.NewDecoder(f).Decode(&cfg)
	return &cfg, err
}
