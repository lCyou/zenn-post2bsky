package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
	"github.com/kyou/post2bsky/bsky"
	"github.com/kyou/post2bsky/config"
	"github.com/kyou/post2bsky/ui"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Config error:", err)
		os.Exit(1)
	}

	if cfg.NeedsLogin() {
		if err := promptLogin(cfg); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := cfg.Save(); err != nil {
			fmt.Fprintln(os.Stderr, "Warning: could not save config:", err)
		}
	}

	client := bsky.NewClient(cfg.PDS)
	client.SetCredentials(cfg.Handle, cfg.AppPassword)

	model := ui.NewModel(client)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func promptLogin(cfg *config.Config) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Bluesky Handle: ")
	if !scanner.Scan() {
		return fmt.Errorf("aborted")
	}
	cfg.Handle = strings.TrimSpace(scanner.Text())

	fmt.Print("App Password: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	cfg.AppPassword = strings.TrimSpace(string(password))

	if cfg.Handle == "" || cfg.AppPassword == "" {
		return fmt.Errorf("handle and app password are required")
	}
	return nil
}
