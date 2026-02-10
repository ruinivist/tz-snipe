package main

import (
	"flag"
	"fmt"
	"os"
	"tz-snipe/internal/geodata"
	"tz-snipe/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// 1. arg parsing

	// go's flag package is the native arg parser for cmd line apps
	// --tz mode, default val, usage string
	tzFlag := flag.String("tz", "", "Timezone like +1000")
	// --github mode
	ghFlag := flag.String("github", "", "From a github username")
	// --time mode
	timeFlag := flag.String("time", "", "Time like 15:00")

	flag.Parse()

	// 2. load tz data
	db, err := geodata.Fetch()
	if err != nil {
		fmt.Printf("Error loading timezone data: %v\n", err)
		os.Exit(1) // 1 is error code
	}

	// 3. checks
	modes := 0
	for _, val := range []string{*tzFlag, *ghFlag, *timeFlag} {
		if val != "" {
			modes++
		}
	}

	if modes > 1 {
		fmt.Println("Cannot use multiple modes at once")
		os.Exit(1)
	}
	if modes == 0 {
		fmt.Println("No args provided... exiting")
		os.Exit(1)
	}

	// 4. exec
	p := tea.NewProgram(ui.NewModel(db, *ghFlag, *tzFlag, *timeFlag))
	if _, err := p.Run(); err != nil {
		// I keep on forgetting, fmt.ErrorF makes an error, does not print
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
