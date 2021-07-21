package typhoon

import "github.com/rivo/tview"

func RunUI()  {
	//box := tview.NewBox().SetBorder(true).SetTitle("Typhoon cli dashboard")
	//if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
	//	panic(err)
	//}

	//app := tview.NewApplication()
	//list := tview.NewList().
	//	AddItem("List item 1", "Some explanatory text", 'a', nil).
	//	AddItem("List item 2", "Some explanatory text", 'b', nil).
	//	AddItem("List item 3", "Some explanatory text", 'c', nil).
	//	AddItem("List item 4", "Some explanatory text", 'd', nil).
	//	AddItem("Quit", "Press to exit", 'q', func() {
	//		app.Stop()
	//	})
	//if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
	//	panic(err)
	//}

	//time.Sleep(time.Second * 20)


	//app := tview.NewApplication()
	//button := tview.NewButton("Hit Enter to close").SetSelectedFunc(func() {
	//	app.Stop()
	//})
	//button.SetBorder(true).SetRect(0, 0, 22, 3)
	//if err := app.SetRoot(button, false).SetFocus(button).Run(); err != nil {
	//	panic(err)
	//}


	//app := tview.NewApplication()
	//textView := tview.NewTextView().
	//	SetDynamicColors(true).
	//	SetChangedFunc(func() {
	//		app.Draw()
	//	})
	//textView.SetBorder(true).SetTitle("Stdin")
	//go func() {
	//	w := tview.ANSIWriter(textView)
	//	if _, err := io.Copy(w, os.Stdin); err != nil {
	//		panic(err)
	//	}
	//}()
	//if err := app.SetRoot(textView, true).Run(); err != nil {
	//	panic(err)
	//}
	//pageCount := 2
	//
	//app := tview.NewApplication()
	//pages := tview.NewPages()
	//for page := 0; page < pageCount; page++ {
	//	func(page int) {
	//		pages.AddPage(fmt.Sprintf("page-%d", page),
	//			tview.NewModal().
	//				SetText(fmt.Sprintf("This is page %d. Choose where to go next.", page+1)).
	//				AddButtons([]string{"Next", "Quit"}).
	//				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
	//					if buttonIndex == 0 {
	//						pages.SwitchToPage(fmt.Sprintf("page-%d", (page+1)%pageCount))
	//					} else {
	//						app.Stop()
	//					}
	//				}),
	//			false,
	//			page == 0)
	//	}(page)
	//}
	//if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
	//	panic(err)
	//}


	//app := tview.NewApplication()
	//dropdown := tview.NewDropDown().
	//	SetLabel("Select an option (hit Enter): ").
	//	SetOptions([]string{"First", "Second", "Third", "Fourth", "Fifth"}, nil)
	//if err := app.SetRoot(dropdown, true).SetFocus(dropdown).Run(); err != nil {
	//	panic(err)
	//}

	app := tview.NewApplication()
	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Left (1/2 x width of Top)"), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (20 cols)"), 20, 1, false)
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}



}
