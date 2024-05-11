# raas-be-go-sample

Go向けBackendソースのサンプルです

## How to setup
### 1.1 プライベート・レポジトリへのアクセスを許可する環境変数の設定

RestClient自体はSuTech社で管理されているプライベート・レポジトリに格納されています。
取得に際し、環境変数の設定が必要です。環境変数は以下の通りです。
Linux環境下でのコマンドになっていますので、Windowsの場合は読み替えてください。

`export GOPRIVATE="github.com/SuTech-JP/raas-client-go"`

### 1.2 依存関係の追加

環境変数の設定が完了したら、以下のコマンドで依存関係を追加します。
`raas`というモジュールの名前空間でRestClientなどが利用できるようになります。

`go get github.com/SuTech-JP/raas-client-go`

### 1.3 ビルド

依存関係の解消のために、以下のコマンドを実行します

`go build`

### 1.4 application.yamlを準備する
（SuTech社より取得した値をXXXX部分に記載）
```
raasConfig:
  application: XXXXX
  landscape: dev
  token: XXXXX
```

### 1.5 実行

`go run main.go`

## 概要
feのsampleをSuTech社より取得して結合する

### 2.1 FEを初期化する為のsession作成を行う
/raas/report/session
/raas/datatraveler/session
にアクセスがあった際に該当の処理が起動する

### 2.2 帳票作成結果（PDF/JSON/CSV）を取得する
/raas/report/result/{targetId}
にアクセスがあった際に該当の処理が起動する


## 組込方法
WIP