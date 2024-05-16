# raas-be-go-sample

Go向けBackendソースのサンプルです

## 1.起動方法

### 1.1 githubのuser/tokenを環境変数に登録する
SuTech社提供のraas-client-goをGitHub.Packageから取得出来るようにする
(SuTech社による権限付与が必要)
```
export RAAS_GITHUB_USERNAME=XXXXXXX
export RAAS_GITHUB_TOKEN=XXXXXXX
```

### 1.2 SuTech社のGitHub.Packageリポジトリに、提供されたuser/tokenでアクセスするための設定を行う
```
git config --global url."https://${RAAS_GITHUB_USERNAME}:${RAAS_GITHUB_TOKEN}@github.com/SuTech-JP".insteadOf "https://github.com/SuTech-JP"
```

### 1.2 ビルド
依存ライブラリの取得のため、以下のコマンドを実行します

`go install`

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
本サンプルは別途SuTech社が提供するFrontend用サンプルと結合して動作する想定のサンプルとなっています。
本サンプルでは以下３つのAPIを提供しています。

### 2.1 Frontend用コンポーネントを表示するためのsession作成用API
- /raas/report/session
- /raas/datatraveler/session

### 2.2 帳票作成結果（PDF/JSON/CSV）を取得する
- /raas/report/result/{targetId}

### 2.3 帳票レイアウト一覧を取得する
- /raas/report/layout/{application}/{schema}

## 3.組込方法
### 3.1 raas-client-goを利用するための設定行う
1.1, 1.2と同じ手順を実施する

### 3.2 FE用のsession関数を作成する
2.1のsession作成用APIと同等の処理を作成する

### 3.3 データ連携処理を作成する（DataImportLogIdを保存する）
CSVインポートを実行した後にDataImportLogIdを保存する
（tenant , sub も一緒に保存することを推奨する）

### 3.4 データ連携処理を作成する（PDF作成処理が終わったDataImportLogIdの結果を取得する）
3.3のデータを元にBEにてデータ連携処理を実装する
データ連携処理は、別途SuTech社から提供されるシーケンス図を参考に実装する
