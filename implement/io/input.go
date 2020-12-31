package iio

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/spf13/pflag"

	sio "aman/struct/io"
)

type InputStruct sio.InputStruct

const (
	optionText = "options arguments\n"
	usageText  = "usage: aman <command> [arguments]\n" +
		"  <command>: command which has man result (ls, git,...)\n" +
		" [arguments]: arguments which command have (`status` for `git status`)\n"
)

/*
 * @description コンストラクタ
 */
func NewInput(version string) *InputStruct {
	input := &InputStruct{
		Commands:   []string{},
		Options:    []string{},
		Query:      "",
		CursorPosX: 2,
	}
	input.Version = version
	return input
}

/*
 * @description コマンドライン引数を取得
 */
func (myself *InputStruct) Parse() error {
	// オプションのセット
	var versionFlag *bool = pflag.BoolP("version", "v", false, "show version")
	var helpFlag *bool = pflag.BoolP("help", "h", false, "show help")
	pflag.Parse()

	// バージョン表示
	if *versionFlag {
		fmt.Println(myself.Version)
		return errors.New("")
	}

	args := pflag.Args()

	// 引数がない場合はヘルプ表示
	if len(args) < 1 || *helpFlag {
		fmt.Fprintf(os.Stderr, optionText)
		pflag.PrintDefaults()
		fmt.Fprintf(os.Stderr, usageText)
		return errors.New("")
	}

	myself.Commands = args
	return nil
}

/*
 * @description 入力を削除する
 */
func (myself *InputStruct) DeleteInput() {
	if 0 < len(myself.Query) {
		myself.CursorPosX -= runewidth.RuneWidth([]rune(myself.Query)[utf8.RuneCountInString(myself.Query)-1])
		myself.Query = string([]rune(myself.Query)[:utf8.RuneCountInString(myself.Query)-1])
	}
	if myself.CursorPosX < 2 {
		myself.CursorPosX = 2
	}
}

/*
 * @description 選択した行のオプションを抽出する
 * @param line オプション説明文
 */
func (myself *InputStruct) ExtractOption(line string) {
	// 文字列を空白区切で区切ったものの先頭がオプションのはずなのでそれを取得
	var selectedOption string = strings.Split(line, " ")[0]
	// 末端の改行を削除する
	selectedOption = strings.TrimRight(selectedOption, "\n")
	// 末尾のカンマを削除する
	selectedOption = strings.TrimRight(selectedOption, ",")

	// 重複選択を制限する
	for _, option := range myself.Options {
		// 一致するオプションが見つかったら追加処理を行わず、returnする
		if option == selectedOption {
			return
		}
	}

	// ストック
	myself.Options = append(myself.Options, selectedOption)
}

/*
 * @description スペース入力
 */
func (myself *InputStruct) PutSpace() {
	myself.Query += " "
	myself.CursorPosX++
}

/*
 * @description キー入力
 * @param キーイベント
 */
func (myself *InputStruct) PutKey(ev termbox.Event) {
	myself.Query += string(ev.Ch)
	for _, r := range string(ev.Ch) {
		myself.CursorPosX += runewidth.RuneWidth(r)
	}
}
