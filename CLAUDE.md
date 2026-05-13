# Go Template

## アーキテクチャ

- **エントリーポイント**: `cmd/gotemplate/main.go` → `run.go` の `Run()` を呼び出す
- **CLI**: `kong` でパース、`--config` と `--version` フラグ
- **設定**: `internal/config` でYAML設定のロード（デフォルト値あり、ファイル不存在は許容）
- **XDG**: 設定ファイルのデフォルトパスは `xdg.ConfigFile("go-template/config.yaml")`

```
cmd/gotemplate/  ─ CLIエントリーポイント
run.go           ─ アプリケーションロジック（CLIパース→設定ロード）
internal/config/ ─ 設定ファイルの読み込み
```

## テスト

- **Table driven test** を採用する
- `*_test.go` では `map[string]struct{}` パターンを使用
- **AAAパターン**（Arrange-Act-Assert）に従う
- `testify` の `assert` を使用する（fatalなエラーには `require`）
