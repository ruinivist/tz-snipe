package ui

import (
	"tz-snipe/internal/github"

	tea "github.com/charmbracelet/bubbletea"
)

type tzsMsg []string
type errMsg error

func fetchGithubCmd(username string) tea.Cmd {
	return func() tea.Msg {
		tzs, err := github.GetTzsForUser(username)
		if err != nil {
			return errMsg(err)
		}
		return tzsMsg(tzs)
	}
}
