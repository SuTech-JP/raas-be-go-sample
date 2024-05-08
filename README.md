# raas-be-go-sample
RaaSへのアクセスを代替するGo言語用のRestClientです。

# How to setup
## 1.1 プライベート・レポジトリへのアクセスを許可する環境変数の設定

RestClient自体はSuTech社で管理されているプライベート・レポジトリに格納されています。
取得に際し、環境変数の設定が必要です。環境変数は以下の通りです。Linux環境下でのコマンドに
なっていますの。Windowsの場合は読み替えてください。

`export GOPRIVATE="github.com/SuTech-JP/raas-client-go"`

## 1.2 依存関係の追加

環境変数の設定が完了したら、以下のコマンドで依存関係を追加します。
`raas`というモジュールの名前空間でRestClientなどが利用できるようになります。

`go get github.com/SuTech-JP/raas-client-go`

## 1.3 ビルド

依存関係の解消のために、以下のコマンドを実行します

`go build`

## 1.4 RaaSRestClientを利用する

RaaSRestClientには、必須のパラメータを格納する構造体が２つあります。
`RaaSConnectConfig`と`RaaSUserContext`になります。
それぞれについて見ていきましょう

### 1.4.1 RaaSConnectConfig

`RaaSConnectConfig`は、バックエンドサーバに接続する際に必要な情報を格納する構造体です。
`NewRaaSConnectConfig`というファクトリーメソッドがあります。こちらを使って初期化を行います。
渡すパラメータは全て**必須**です。Tokenなどの大事な情報でもありますので、環境変数を利用することを推奨します。これらの情報はSuTech社から提供されます。以下が、実装例です。

```go
	config, configErr := raas.NewRaaSConnectConfig(
			os.Getenv("RaaS-Application"), 
			os.Getenv("RaaS-Landscape"), 
			os.Getenv("RaaS-Token"))

```

### 1.4.2 RaaSUserContext

`RaaSUserContext`は、ご利用のサーバ等の情報を格納します。必須の初期化パラメータと
任意で設定ができるSetterとの組み合わせによって構成されています。必須パラメータは、tenantとsubです。提供された情報に基づいて設定を行ってください。ここでも、`NewRaasUserContext`というファクトリーメソッドが用意されています。任意で設定できるSetterに関しても、下記の実装例を参考にしてください。

```go
    //tenantとsubは必須
	context, contextErr := raas.NewRaasUserContext("tenant", "sub")
    //以下は任意設定
    context.SetSubAlias("subAlias")
    context.SetSubDomain("subDomain")
    context.SetTenantAlias("subTenantAlias")

```

### 1.4.3 RaaSRestClient

1.4.1と1.4.2で準備したものを渡して初期化を行います。

```go
raas.RaaSRestClient(*config, *context)
```

このクライアントには、公開メソッドとして。RaaSとのセッションを確立するメソッドと
GET,POST,DELETE,PUTのHTTPメソッドを利用してRaaSにアクセスするメソッドが用意されています。

##　1.5 RaaSとのセッションを確立する場合

RaaSとのセッションを確立し、新たな処理を行う場合には以下のようにします

```go
    //パラメータは自身の環境に置き換えてください
    const backUrl = "http://localhost:3000/send-data"
	const subUrl = "subUrl"
	const subDomain = "subDomain" //override
	//msa must be [report] or [datatraveler]
	extSession, err := raas.RaaSRestClient(*config, *context).CreateExternalSession(msa, backUrl, subUrl, subDomain)
	if err == nil {
		//extSessionはJSON形式での返却が望ましいです。
        //JSONへのエンコードをここで行います
	} else {
        //エラー処理を確実に行ってください
	}

```

## 1.6 RaaSに問い合わせをする場合

ここではGETを利用した方法を紹介します

```go

    //GET, POST, DELETE, PUTの結果は、map[string]anyで返却されます
	resultMap, err := raas.RaaSRestClient(*config, *context).Get(requestUrl, nil)
```

## 1.7 サンプルについて

サンプルはURLの加工しやすさから、`github.com/gorilla/mux` を利用しています。
適宜、ご自身の環境に置き換えてください

## 1.8 権利関係

サンプルを含む、RaaSに関する全ての著作権は株式会社SuTechが保有しています