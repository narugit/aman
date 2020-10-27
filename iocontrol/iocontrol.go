package iocontrol

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aman/modules"
	"github.com/nsf/termbox-go"
)

/*
 * height  ウィンドウの高さ
 * page    現在のページ番号 定義域は[0, maxPage]
 * maxPage 最大ページ番号
 */
type IoController struct {
	height  int
	page    int
	maxPage int
}

/*
 * @param manLists オプションとオプション説明が格納された文字列と、各オプション説明の行数の配列
 * @description IoControllerのコンストラクタ
 */
func NewIoController(manLists []modules.ManData) *IoController {
	_, height := termbox.Size()
	iocontroller := IoController{
		height:  height,
		page:    0,
		maxPage: 0,
	}
	return &iocontroller
}

func DeleteInput(inputs *string) {
	var space = ""
	for i := 0; i < len(*inputs); i++ {
		space += " "
	}
	fmt.Printf("\r%s", space)
	if 0 < len(*inputs) {
		*inputs = (*inputs)[:len(*inputs)-1]
	}
}

func (iocontroller *IoController) ReceiveKeys(inputs *string, selectedPos *int) int {
	var ev termbox.Event = termbox.PollEvent()

	if ev.Type != termbox.EventKey {
		return 0
	}

	switch ev.Key {
	case termbox.KeyEsc:
		return -1
	case termbox.KeyArrowUp:
		return 90
	case termbox.KeyArrowDown:
		return 91
	case termbox.KeyArrowRight:
		iocontroller.page++
		if iocontroller.maxPage < iocontroller.page {
			iocontroller.page = iocontroller.maxPage
		}
	case termbox.KeyArrowLeft:
		iocontroller.page--
		if iocontroller.page < 0 {
			iocontroller.page = 0
		}
	case termbox.KeySpace:
		*inputs += " "
		break
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		DeleteInput(inputs)
		break
	case termbox.KeyEnter:
		return 99
	default:
		iocontroller.page = 0
		*selectedPos = 0
		*inputs += string(ev.Ch)
		break
	}
	return 1
}

func RenderQuery(inputs *string) {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
	fmt.Printf("\r> %s\n", *inputs)
}

func (iocontroller *IoController) RenderResult(selectedPos int, result []modules.ManData, pageList []int) {
	const SEPARATOR = "----------"
	var row = 0
	fmt.Printf("%d/%d", iocontroller.page+1, iocontroller.maxPage+1)
	fmt.Println(SEPARATOR)
	row++
	if iocontroller.height <= row {
		return
	}

	if len(result) == 0 {
		return
	}

	for i := pageList[iocontroller.page]; i < pageList[iocontroller.page+1]; i++ {
		row += strings.Count(result[i].Contents, "\n") + 2
		if iocontroller.height <= row {
			return
		}
		var state string = "\r%s\n"
		if selectedPos == i {
			// 選択行だけ赤色に変更
			state = "\r\x1b[31m%s\x1b[0m\n"
		}
		fmt.Printf(state, result[i].Contents)
		fmt.Println(SEPARATOR)
	}
}

/*
 * @param manLists オプションとオプション説明が格納された文字列と、各オプション説明の行数の配列
 * @description 各ページの先頭となるオプション配列manListsのindex番号が格納された配列を生成する
 */
func (iocontroller *IoController) LocatePages(manLists []modules.ManData) []int {
	var maxLineNumber = -1
	pageList := []int{0}
	// >行と---の2行
	var lineCount = 2
	var page = 0
	iocontroller.maxPage = 0

	for i := 0; i < len(manLists); i++ {
		// for文を抜けた後に、ウィンドウの高さが低すぎて描画できないかを判定するために、
		// 一番行数の多いオプション説明文の行数を求める
		if maxLineNumber < manLists[i].LineNumber {
			maxLineNumber = manLists[i].LineNumber
		}

		// ウィンドウの高さをオーバーしてしまう場合、次のページにオプション説明を表示する
		if iocontroller.height < lineCount+manLists[i].LineNumber {
			lineCount = 2
			page++
			pageList = append(pageList, i)
			if iocontroller.maxPage < page {
				iocontroller.maxPage = page
			}
		}

		lineCount += manLists[i].LineNumber

		if i == len(manLists)-1 {
			pageList = append(pageList, i+1)
		}
	}

	if iocontroller.height < maxLineNumber {
		panic(errors.New("Window height is too small"))
	}

	return pageList
}
