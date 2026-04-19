package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyou/post2bsky/bsky"
)

const maxChars = 300

type Model struct {
	client    *bsky.Client
	textarea  textarea.Model
	posts     []bsky.Post
	status    string
	statusErr bool
	posting   bool
	authed    bool
	width     int
}

func NewModel(client *bsky.Client) Model {
	ta := textarea.New()
	ta.Placeholder = "What's on your mind?"
	ta.SetHeight(5)
	ta.SetWidth(50)
	ta.CharLimit = 0
	ta.ShowLineNumbers = false
	ta.Focus()

	return Model{
		client:   client,
		textarea: ta,
		width:    60,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(authCmd(m.client), textarea.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.textarea.SetWidth(msg.Width - 6)

	case AuthDoneMsg:
		m.authed = true
		m.status = "Authenticated"
		return m, loadFeedCmd(m.client)

	case AuthErrMsg:
		m.status = msg.Err.Error()
		m.statusErr = true

	case FeedLoadedMsg:
		m.posts = msg.Posts

	case FeedErrMsg:
		// フィード取得失敗は無視

	case PostDoneMsg:
		m.posting = false
		m.textarea.SetValue("")
		m.status = "Posted!"
		m.statusErr = false
		return m, loadFeedCmd(m.client)

	case PostErrMsg:
		m.posting = false
		m.status = msg.Err.Error()
		m.statusErr = true

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if msg.Type == tea.KeyCtrlS {
			return m, m.submitPost()
		}

		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) submitPost() tea.Cmd {
	text := strings.TrimSpace(m.textarea.Value())
	if len([]rune(text)) == 0 || len([]rune(text)) > maxChars || m.posting || !m.authed {
		return nil
	}
	m.posting = true
	return postCmd(m.client, text)
}

func (m Model) View() string {
	w := m.width - 4
	if w < 40 {
		w = 40
	}

	divider := dividerStyle.Render(strings.Repeat("─", w))

	title := titleStyle.Render("post2bsky")

	charCount := len([]rune(m.textarea.Value()))
	var counter string
	if charCount > maxChars {
		counter = counterOverStyle.Render(fmt.Sprintf("%d/%d", charCount, maxChars))
	} else {
		counter = counterStyle.Render(fmt.Sprintf("%d/%d", charCount, maxChars))
	}
	counterLine := lipgloss.NewStyle().Width(w).Align(lipgloss.Right).Render(counter)

	hintLine := lipgloss.NewStyle().Width(w).Align(lipgloss.Right).Render(hintStyle.Render("ctrl+s: Post   ctrl+c: Quit"))

	var statusLine string
	if m.posting {
		statusLine = hintStyle.Render("Posting...")
	} else if m.status != "" {
		if m.statusErr {
			statusLine = statusErrStyle.Render(m.status)
		} else {
			statusLine = statusOkStyle.Render(m.status)
		}
	}

	recentParts := []string{recentHeaderStyle.Render("Recent posts (last 5):")}
	if len(m.posts) == 0 {
		recentParts = append(recentParts, hintStyle.Render("  (none)"))
	}
	for _, p := range m.posts {
		text := p.Text
		runes := []rune(text)
		if len(runes) > 60 {
			text = string(runes[:57]) + "..."
		}
		recentParts = append(recentParts, recentItemStyle.Render("• "+text))
	}

	body := strings.Join([]string{
		title,
		divider,
		m.textarea.View(),
		counterLine,
		divider,
		hintLine,
	}, "\n")

	if statusLine != "" {
		body += "\n" + statusLine
	}

	body += "\n" + divider + "\n" + strings.Join(recentParts, "\n")

	return borderStyle.Render(body)
}

func authCmd(client *bsky.Client) tea.Cmd {
	return func() tea.Msg {
		if err := client.Authenticate(); err != nil {
			return AuthErrMsg{Err: err}
		}
		return AuthDoneMsg{}
	}
}

func loadFeedCmd(client *bsky.Client) tea.Cmd {
	return func() tea.Msg {
		posts, err := client.GetOwnRecentPosts(5)
		if err != nil {
			return FeedErrMsg{Err: err}
		}
		return FeedLoadedMsg{Posts: posts}
	}
}

func postCmd(client *bsky.Client, text string) tea.Cmd {
	return func() tea.Msg {
		if err := client.CreatePost(text); err != nil {
			return PostErrMsg{Err: err}
		}
		return PostDoneMsg{}
	}
}
