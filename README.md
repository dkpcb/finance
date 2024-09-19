## 概要

このアプリケーションは、投資信託の取引履歴と基準価額を管理し、評価する。主要な機能としては、取引と価格データのインポート、データ処理、およびユーザー固有の取引数、資産評価、損益計算を取得するためのAPIを提供する。

## 機能

### 機能1: 取引履歴と基準価額データのインポート
- **説明**: アプリケーションは `trade_history.csv` と `reference_prices.csv` の2つのCSVファイルを読み込み、その内容をMySQL 8.0データベースに保存する。
- **ファイル構成**:
  - `trade_history.csv`: `user_id`、`fund_id`、`quantity`、`trade_date` の列を含む取引履歴を記録。
  - `reference_prices.csv`: `fund_id`、`reference_price_date`、`reference_price` の列を含む基準価額履歴を記録。

### 機能2: Dockerを使用したアプリケーションのコンテナ化
- **説明**: アプリケーションはDockerを使用してコンテナ化され、実行される。
- **コマンド**:
  - `make dev/run/import`: データベースにデータをインポートする。
  - `make dev/run/server`: アプリケーションサーバーを起動し、`localhost:8080` でアクセス可能にする。

### 機能3: ユーザーの取引回数を取得するAPI
- **説明**: 特定の `user_id` に対する取引回数を取得する。
- **APIエンドポイント**: `/{user_id}/trades`
- **リクエスト**: `http://localhost:8080/A1B2C3D4E5/trades`　
- **レスポンス**:
  ```json
  { 
    "count": 197 
  }
  ```

### 機能4: 現在の資産評価額と評価損益を取得するAPI
- **説明**: 特定の `user_id` に対する現在の資産評価額と評価損益を取得する。
- **APIエンドポイント**: `/{user_id}/assets`
- **リクエスト**: `http://localhost:8080/A1B2C3D4E5/assets`　
- **レスポンス**:
  ```json
  { 
    "date": "2024-06-01", 
    "current_value": 17661, 
    "current_pl": 62 
  }
  ```

### 機能5: 特定の日付における資産評価額と評価損益を取得するAPI
- **説明**: 特定の `user_id` と `date` に対する資産評価額と評価損益を取得する。
- **APIエンドポイント**: `/{user_id}/assets?date={date}`
- **リクエスト**: `http://localhost:8080/A1B2C3D4E5/assets?date=2023-12-08`　
- **レスポンス**:
  ```json
  {
    "date": "2023-12-08",
    "current_value": 5418,
    "current_pl": 16
  }
  ```

### 機能6: 年ごとの資産評価額と評価損益を取得するAPI
- **説明**: 特定の `user_id` における現在時点の資産評価額と評価損益を、購入年ごとに集計して返す。
- **APIエンドポイント**: `/{user_id}/assets/byYear`
- **リクエスト**: `http://localhost:8080/A1B2C3D4E5/assets/byYear`
- **レスポンス**:
  ```json
  { 
    "date": "2024-06-01", 
    "assets": [ 
      { 
        "year": 2024, 
        "current_value": 10945, 
        "current_pl": 38 
      }, 
      { 
        "year": 2023, 
        "current_value": 6715, 
        "current_pl": 23 
      }
    ] 
  }
  ```

## 実装の詳細

- **言語**: Go
- **データベース**: MySQL 8.0
- **コンテナ**: Docker