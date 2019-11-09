package main

import (
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/morikuni/aec"
)

func main() {
	statusMap = map[int]string{
		Started: "started",
		Stopped: "stopped",
	}

	maxWorkers := 2
	workItems := []Work{
		{Command: "sleep 1.1", Name: "Func 1"},
		{Command: "sleep 0.1", Name: "Func 2"},
		{Command: "sleep 0.5", Name: "Func 3"},
		{Command: "sleep 1.4", Name: "Func 4"},
		{Command: "sleep 2.1", Name: "Func 5"},
		{Command: "sleep 0.1", Name: "Func 6"},
	}

	parallelRun(workItems, maxWorkers)
}

type Work struct {
	Command  string
	Name     string
	Started  *time.Time
	Finished *time.Time
}

func (w *Work) Status() int {
	if w.Started == nil {
		return NotStarted
	} else if w.Finished == nil {
		return Started
	}
	return Stopped
}

type StatusUpdate struct {
	Name   string
	Status int
}

const Started int = 1
const Stopped int = 2
const NotStarted int = 0

var statusMap map[int]string

type ConsolePainter struct {
	WorkMap *map[string]*Work
}

type ByCreated []Work

func (a ByCreated) Len() int      { return len(a) }
func (a ByCreated) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCreated) Less(i, j int) bool {
	return a[i].Status() > a[j].Status()
}

func (c *ConsolePainter) Paint() {
	workMapOrder := []Work{}
	for _, v := range *c.WorkMap {
		workMapOrder = append(workMapOrder, *v)
	}
	sort.Sort(ByCreated(workMapOrder))

	done := 0
	numLines := len(workMapOrder)

	upN := aec.Up(uint(numLines + 1))
	fmt.Print(upN)

	var color aec.ANSI
	for _, v := range workMapOrder {
		var line string
		if v.Started == nil {
			color = aec.CyanF
			line = fmt.Sprintf("[%-40s] - %-20s", color.Apply(v.Name), "stopped")
		} else if v.Finished == nil {
			color = aec.YellowF
			runTime := time.Since(*v.Started).Seconds()
			line = fmt.Sprintf("[%-40s] - %-20s %-.1fs", color.Apply(v.Name), "running", float64(runTime))

		} else {
			color = aec.GreenF

			runTime := (*v.Finished).Sub(*v.Started).Seconds()

			// runTime := time.Since(*v.Started).Seconds()
			line = fmt.Sprintf("[%-40s] - %-20s %-.1fs", color.Apply(v.Name), "done", float64(runTime))

			done++
		}

		fmt.Printf("%-100s\n", line)
	}

	fmt.Printf("Done: %d/%d\n", done, numLines)
}

func parallelRun(workItems []Work, maxWorkers int) {
	wg := sync.WaitGroup{}
	statusWg := sync.WaitGroup{}

	painter := ConsolePainter{}
	workMap := make(map[string]*Work)
	for _, item := range workItems {
		copied := item
		workMap[item.Name] = &copied
	}
	numLines := 0
	for range workMap {
		numLines++
	}

	for i := 0; i < (numLines + 1); i++ {
		fmt.Println("")
	}

	painter.WorkMap = &workMap

	workChannel := make(chan Work)
	statusUpdate := make(chan StatusUpdate)

	for i := 0; i < maxWorkers; i++ {
		go func(index int) {
			wg.Add(1)

			for work := range workChannel {
				statusUpdate <- StatusUpdate{Name: work.Name, Status: Started}

				workParts := strings.Split(work.Command, " ")
				targetCmd := exec.Command(workParts[0], workParts[1:]...)
				out, _ := targetCmd.CombinedOutput()
				strings.Contains(string(out), "")

				statusUpdate <- StatusUpdate{Name: work.Name, Status: Stopped}
			}

			wg.Done()
		}(i)
	}

	go func() {
		statusWg.Add(1)

		for update := range statusUpdate {
			mapItem := workMap[update.Name]
			if mapItem != nil {
				if update.Status == Started {
					var v time.Time
					v = time.Now()
					mapItem.Started = &v
				} else if update.Status == Stopped {
					if mapItem.Finished != nil {
						log.Fatalln("Setting finished status again??")
					}
					var finished time.Time
					finished = time.Now()
					mapItem.Finished = &finished
				}
				painter.Paint()
			}
		}
		statusWg.Done()
	}()

	for i := 0; i < len(workItems); i++ {
		workChannel <- workItems[i]
	}

	close(workChannel)

	wg.Wait()

	go func() {
		close(statusUpdate)
	}()
	statusWg.Wait()
}
