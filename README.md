# raas-be-go-sample

Go向けBackendソースのサンプルです

## 1.起動方法

### 1.1 raas-client-goを取得できるようにする
SuTech社提供のraas-client-goは、GitHubのprivateリポジトリ（`github.com/SuTech-JP/raas-client-go`）にあります。
SuTech社に登録したdeploy key（このリポジトリ専用の、読み取り専用のSSH鍵）を使って取得します。

#### 1.1.1 SSH鍵を作成し、公開鍵をSuTech社へ送る
```
ssh-keygen -t ed25519 -f ~/.ssh/raas-client-go -C "raas-client-go deploy key"
```
- SuTech社へは**公開鍵（`~/.ssh/raas-client-go.pub`）だけ**を送ってください。秘密鍵は送らず、社外に出さないでください
- SuTech社が読み取り専用で登録し、完了をご連絡します。連絡を受けてから次の手順へ進んでください
- 開発用とCI用など、用途ごとに鍵を分けることもできます。その場合はそれぞれの公開鍵を送ってください
- 1つの鍵を他のGitHubリポジトリと共用することはできません（GitHubの制約です）
- 鍵が不要になった場合や、漏えいの疑いがある場合はSuTech社へご連絡ください。登録を削除します

#### 1.1.2 SSHの接続先を設定する
`~/.ssh/config` に次を追加します。他のGitHubの利用と鍵が混ざらないよう、専用のホスト名（`github-raas-client-go`）を使います。
```
Host github-raas-client-go
    HostName github.com
    User git
    IdentityFile ~/.ssh/raas-client-go
    IdentitiesOnly yes
```

接続を確認します。`SuTech-JP/raas-client-go` への認証に成功した旨が表示されれば設定できています。
```
ssh -T git@github-raas-client-go
```

#### 1.1.3 GitとGoを設定する
```
git config --global url."ssh://git@github-raas-client-go/SuTech-JP/raas-client-go".insteadOf "https://github.com/SuTech-JP/raas-client-go"
go env -w GOPRIVATE=github.com/SuTech-JP/raas-client-go
```
- `insteadOf` は、置き換えの対象をraas-client-goのリポジトリまでに絞ってください。`https://github.com/SuTech-JP` 全体を置き換えると、この鍵で読めない他のリポジトリにも鍵が使われます
- `GOPRIVATE` を設定しないと、Goは公開のプロキシ（proxy.golang.org）とチェックサムのデータベース（sum.golang.org）に問い合わせて取得に失敗します

#### 1.1.4 CIで取得する場合
- 秘密鍵はCIのsecretに登録し、ジョブの中で一時ファイルに書き出して権限を`600`にしてください
- 1.1.2と1.1.3の設定はジョブの中で行い、共有のランナーのglobal設定を書き換えないでください
- GitHubのホスト鍵は、[GitHubが公開しているフィンガープリント](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/githubs-ssh-key-fingerprints)と照合してから`known_hosts`に登録してください。`StrictHostKeyChecking=no` は使わないでください

### 1.2 ビルド
依存ライブラリを取得・検証してからビルドします。
```
go mod download
go mod verify
go build ./...
```
- raas-client-goのバージョンは `go get github.com/SuTech-JP/raas-client-go@v1.0.8` のように明示して更新してください（`@latest` 任せにしない）
- `go.mod` と `go.sum` はリポジトリにコミットしてください。privateなモジュールは公開のチェックサムのデータベースで検証されないため、`go.sum` が取得物の改ざんを検出する唯一の手段になります

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
