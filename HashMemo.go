package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/lxn/walk"
	decl "github.com/lxn/walk/declarative"
	"golang.org/x/crypto/sha3"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	runMainTable()
}

func runMainTable() error {
	var mw *walk.MainWindow
	var tv *walk.TableView
	var cm *colModel = newColModel("./")
	_, err := decl.MainWindow{

		ContextMenuItems: []decl.MenuItem{
			decl.Action{
				Text:        "🗘refresh",
				OnTriggered: func() {},
			},
			decl.Action{
				Text: "copy",
				OnTriggered: func() {
					if err := walk.Clipboard().SetText(cm.items[tv.CurrentIndex()].Sha3_256); err != nil {
						log.Print("Copy: ", err)
					}
				},
			},
		},

		AssignTo: &mw,
		Title:    "HashMemo",
		Layout:   decl.VBox{MarginsZero: true, SpacingZero: true},
		Children: []decl.Widget{
			decl.PushButton{
				Text: "Toggle",
				OnClicked: func() {
					fmt.Println("menu")
					mw.Menu()
				},
			},
			decl.TableView{
				Name:             "테이블뷰얌",
				AssignTo:         &tv,
				AlternatingRowBG: true,
				ColumnsOrderable: true,
				MultiSelection:   true,
				Columns:          newColView(),
				Model:            cm,

				OnItemActivated: func() {
					i := tv.CurrentIndex()
					item := cm.items[i]
					defer tv.UpdateItem(i)
					fmt.Println("OnItemActivated", i, item.Name)
					if newMemo, isOK := runSetMemoDialog(mw, item.Name, item.Memo); isOK {
						item.Memo = newMemo
					}

					//value := tv.Items[tv.CurrentIndex()].value
					//walk.MsgBox(mw, "Value", value, walk.MsgBoxIconInformation)
				},
				//OnSelectedIndexesChanged: func() {
				//	fmt.Println("OnSelectedIndexesChanged")
				//},
				//OnBoundsChanged: func() {
				//	fmt.Println("OnBoundsChanged")
				//},
				//OnKeyDown: func(key walk.Key) {
				//	fmt.Println("OnKeyDown")
				//},
				//OnCurrentIndexChanged: func() {
				//	fmt.Println("OnCurrentIndexChanged")
				//},
				//OnKeyPress: func(key walk.Key) {
				//	fmt.Println("OnKeyPress")
				//},
				//OnKeyUp: func(key walk.Key) {
				//	fmt.Println("OnKeyUp")
				//},
				//OnMouseDown: func(x, y int, button walk.MouseButton) {
				//	fmt.Println("OnMouseDown")
				//},
				//OnMouseMove: func(x, y int, button walk.MouseButton) {
				//	fmt.Println("OnMouseMove", x, y, tv.Name())
				//},
				//OnMouseUp: func(x, y int, button walk.MouseButton) {
				//	fmt.Println("OnMouseUp")
				//},
				//OnSizeChanged: func() {
				//	fmt.Println("OnSizeChanged")
				//},
			},
		},
	}.Run()
	if err != nil {
		log.Fatal(err)
	}
	return err
}

type col struct {
	Name     string
	Sha3_256 string
	Memo     string
	IsDir    bool
}
type colModel struct {
	walk.SortedReflectTableModelBase
	items []*col
}

func (m *colModel) Items() interface{} {
	return m.items
}

func newColView() []decl.TableViewColumn {
	return []decl.TableViewColumn{
		{Name: "Name"},
		{Name: "Sha3_256"},
		{Name: "Memo"},
		{Name: "IsDir"},
	}
}
func newColModel(dirRoot string) *colModel {
	files, err := ioutil.ReadDir(dirRoot)
	if err != nil {
		log.Fatal(err)
	}

	m := &colModel{items: make([]*col, len(files))}
	for i, file := range files {
		name := file.Name()
		m.items[i] = &col{
			Name:     name,
			Sha3_256: checksumSha256(name),
			Memo:     "없음",
			IsDir:    file.IsDir(),
		}
	}
	return m
}

//체크섬 값이 제대로 나온건지 아래 사이트에서 교차검증 필요
//https://www.virustotal.com/
//https://emn178.github.io/online-tools/sha256_checksum.html
func checksumSha256(fileName string) string {
	f, err := os.Open(fileName)
	if err != nil {
		return err.Error()
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func runSetMemoDialog(owner walk.Form, title, old string) (string, bool) {
	var dlg *walk.Dialog
	var edit *walk.LineEdit
	var acceptPB, cancelPB *walk.PushButton
	result, _ := decl.Dialog{
		AssignTo:      &dlg,
		Title:         title,
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		Layout:        decl.VBox{},
		Children: []decl.Widget{
			decl.LineEdit{
				AssignTo: &edit,
				Text:     old,
			},
			decl.Composite{
				Layout: decl.HBox{},
				Children: []decl.Widget{
					decl.PushButton{
						AssignTo:  &acceptPB,
						Text:      "OK",
						OnClicked: func() { dlg.Accept() },
					},
					decl.PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
	return edit.Text(), result == walk.DlgCmdOK
}
