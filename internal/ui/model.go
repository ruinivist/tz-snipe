package ui

import (
	"fmt"
	"strings"
	"tz-snipe/internal/core"
	"tz-snipe/internal/geodata"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	db              geodata.TimezoneMap
	username        string
	manualTz        string
	targetTime      string
	preds           []core.Prediction
	loadingGithubTz bool
	err             error // this is how error handling is to be done in bubbletea
	// just pass the error to ui layer
	displayTz string
	table     table.Model
}

// go quirk: struct trick to get lsp to highlight errors
var _ tea.Model = Model{}

// REMINDER REMINDER REMINDER
// Model, NewModel, Init, Commands, Update, View is all you need

func NewModel(db geodata.TimezoneMap, ghUser, manualTz, targetTime string) Model {
	m := Model{
		db:              db,
		username:        ghUser,
		manualTz:        manualTz,
		targetTime:      targetTime,
		loadingGithubTz: false,
	}
	switch {
	case manualTz != "":
		preds, err := core.GetStatsForTZs([]string{manualTz})
		if err != nil {
			m.err = err
		} else {
			m.preds = preds
			m.table = newTable(m.preds)
			m.displayTz = m.manualTz
		}
	case targetTime != "":
		tzs, err := core.GetTZsMatchingTime(db, targetTime)
		if err != nil {
			m.err = err
		} else if len(tzs) == 0 {
			m.err = fmt.Errorf("no timezones found matching %s. Is time in hh:mm format?", targetTime)
		} else {
			preds, err := core.GetStatsForTZs(tzs)
			if err != nil {
				m.err = err
			} else {
				m.preds = preds
				m.table = newTable(m.preds)
				m.displayTz = strings.Join(tzs, ", ")
			}
		}
	case ghUser != "":
		m.loadingGithubTz = true
	}

	return m
}

func (m Model) Init() tea.Cmd {
	if m.err != nil {
		return tea.Quit
	}
	if m.loadingGithubTz {
		return fetchGithubCmd(m.username)
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// tbf this interrupt handling should be default
	// go's type assertion trick aghain
	// you want it handled first and globally so that
	// fall throughs work nicely
	if asKey, casted := msg.(tea.KeyMsg); casted {
		if asKey.Type == tea.KeyCtrlC || asKey.String() == "q" {
			return m, tea.Quit
		}
	}

	// .(type) is special go syntax
	// defiendonly for interfaces and must be used in a switch
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
		return m, tea.Quit
		// seems like bubbletea will render once and then quit
		// as I get my error msg printed

	case tzsMsg:
		preds, err := core.GetStatsForTZs(msg)
		if err != nil {
			m.err = err
		} else {
			m.preds = preds
			m.table = newTable(m.preds)
			m.displayTz = strings.Join(msg, ", ")
			m.loadingGithubTz = false
		}
		return m, nil
	}

	// anuthing else, to the table

	// go quirk: non name on left side of ...
	// I can't use := as then it tries to declare it
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.err != nil {
		return AppStyle.Render(ErrStyle.Render("Error: " + m.err.Error()))
	}

	if m.loadingGithubTz {
		return AppStyle.Render(SubtleStyle.Render("Fetching github commits..."))
	}

	if len(m.preds) == 0 {
		return AppStyle.Render(SubtleStyle.Render("No valid tzs found. Are you sure tz is correct/user is active?"))
	}

	var header string
	if m.manualTz != "" {
		header += SubtleStyle.Render(fmt.Sprintf("Tz : %s\n", m.displayTz))
	} else if m.username != "" {
		// you need Width to split, MaxWidth just truncates
		header += SubtleStyle.Width(48).Render(fmt.Sprintf("User : %s\nTz   : %s\n", m.username, m.displayTz))
	} else {
		header += SubtleStyle.Render(fmt.Sprintf("Time : %s\nTz   : %s\n", m.targetTime, m.displayTz))
	}

	table := m.table.View()

	view := lipgloss.JoinVertical(lipgloss.Left, header, table)
	return AppStyle.Render(view)
}
