# raas-be-go-sample

Go向けBackendソースのサンプルです

## 1.起動方法
### 1.1 プライベート・レポジトリへのアクセスを許可する環境変数の設定

SuTech社提供のraas-client-goはSuTech社のプライベート・レポジトリに格納されています。

取得に際し、環境変数の設定が必要です。環境変数は以下の通りです。
Linux環境下でのコマンドになっていますので、Windowsの場合は読み替えてください。
`export GOPRIVATE="github.com/SuTech-JP/raas-client-go"`

### 1.2 ビルド

依存関係の解決のために、以下のコマンドを実行します

`go build`

### 1.3 application.yamlを準備する
（SuTech社より取得した値をXXXX部分に記載）
```
raasConfig:
  application: XXXXX
  landscape: dev
  token: XXXXX
```

### 1.4 実行

`go run main.go`

## 2.サンプル概要
本サンプルは別途SuTech社が提供するFrontend用サンプルを結合して動作する想定のサンプルとなっています。
サンプルでは以下３つのAPIを提供しています。

### 2.1 Frontend用コンポーネントを表示するためのsession作成用API
- /raas/report/session
- /raas/datatraveler/session

### 2.2 帳票作成結果（PDF/JSON/CSV）を取得する
- /raas/report/result/{targetId}

### 2.3 帳票レイアウト一覧を取得する
- /raas/report/layout/{application}/{schema}

## 3.組込方法
WIP