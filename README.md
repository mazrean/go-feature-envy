# go-feature-envy

GoのFeature Envyを検出する静的解析ツール。
Goにはクラスの概念はないので、構造体のメソッドで異なる型のフィールドやメソッドを基準を超えて呼び出していた場合Feature Envyと判定している。

## Usage
```bash
$ go build cmd/gofeatureenvy
$ go vet -vettool=./gofeatureenvy {{対象ソースコードのパス}}
```
