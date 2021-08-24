package bar

import "fmt"

type Bar struct {
	percent int64 // progress percentage
	cur int64 // current progress
	total int64 // total value for progress
	rate string // the actual progress bar to be printed
	graph string // the fill value for progress bar
	Description string
}

func (bar *Bar) NewOption(start, total int64) {
	bar.cur = start
	bar.total = total
	if bar.graph == "" {
		bar.graph = "█"
	}
	if bar.total == -1 {
		bar.graph = ""
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph // initial progress position
	}
}

func (bar *Bar) getPercent() int64 {
	return int64((float32(bar.cur) / float32(bar.total))*50)
}


func (bar *Bar) NewOptionWithGraph(start, total int64, graph string) {
	bar.graph = graph
	bar.NewOption(start, total)
}

func (bar *Bar) Play(cur int64, description string) {
	if len(description) > 0 {
		bar.Description = description
	}
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last {
		var i int64 = 0
		for ; i < bar.percent - last; i++ {
			bar.rate += bar.graph
		}
		if bar.total == -1 {
			fmt.Printf("\r%s [%8d]", bar.Description, bar.cur)
		} else {
			fmt.Printf("\r%s [%-50s]%3d%% %8d/%d ", bar.Description, bar.rate, bar.percent*2, bar.cur, bar.total)
		}

	}
}

func (bar *Bar) Increment()  {
	bar.cur += 1
	bar.Play(bar.cur, bar.Description)

}

func (bar *Bar) Finish(){
	fmt.Println()
}