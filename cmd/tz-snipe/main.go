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

	flag.Parse()

	// 2. load tz data
	db, err := geodata.Fetch()
	if err != nil {
		fmt.Printf("Error loading timezone data: %v\n", err)
		os.Exit(1) // 1 is error code
	}

	// 3. checks
	if *tzFlag != "" && *ghFlag != "" {
		fmt.Println("Cannot use both --tz and --github modes")
		os.Exit(1)
	}
	if *tzFlag == "" && *ghFlag == "" {
		fmt.Println("Both args empty... exiting")
		os.Exit(1)
	}

	// 4. exec
	p := tea.NewProgram(ui.NewModel(db, *ghFlag, *tzFlag))
	if _, err := p.Run(); err != nil {
		// I keep on forgetting, fmt.ErrorF makes an error, does not print
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
