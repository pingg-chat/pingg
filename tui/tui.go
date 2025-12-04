package tui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/pingg-chat/pingg/models"
	"golang.org/x/term"
)

const (
	sidebarWidth  = 30
	minChatHeight = 10
)

// Terminal dimensions from environment (for SSH sessions)
var (
	envWidth  int
	envHeight int
)

type panelType int

const (
	channelsPanel panelType = iota
	dmsPanel
	messagesPanel
	inputPanel
)

type Model struct {
	user               *models.User
	width              int
	height             int
	activePanel        panelType
	selectedChannelIdx int
	selectedDMIdx      int
	textInput          textinput.Model
	viewport           viewport.Model
	messages           []Message
	ready              bool
	lastKey            string // debug: last key pressed
}

type Message struct {
	Time     string
	Author   string
	Icon     string
	Content  string
	IsLocked bool
}

func NewModel(user *models.User) Model {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.CharLimit = 500
	ti.Width = 50

	// Mock messages
	messages := []Message{
		{Time: "10:13", Author: "r2luna", Icon: "🔒", Content: "Olá, tudo bem?", IsLocked: true},
		{Time: "10:14", Author: "joedoe", Icon: "👁", Content: "Oi, tudo sim e você?", IsLocked: false},
		{Time: "10:15", Author: "r2luna", Icon: "🔒", Content: "Tudo ótimo, obrigado por perguntar!", IsLocked: true},
		{Time: "10:16", Author: "joedoe", Icon: "👁", Content: "Que bom! O que você tem feito ultimamente?", IsLocked: false},
	}

	return Model{
		user:               user,
		activePanel:        channelsPanel,
		selectedChannelIdx: 0,
		selectedDMIdx:      0,
		textInput:          ti,
		messages:           messages,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Debug: capture ALL messages
	m.lastKey = fmt.Sprintf("msg=%T", msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Use environment dimensions if available (for SSH), otherwise use detected size
		if envWidth > 0 && envHeight > 0 {
			m.width = envWidth
			m.height = envHeight
		} else {
			m.width = msg.Width
			m.height = msg.Height
		}

		if !m.ready {
			m.viewport = viewport.New(m.width-sidebarWidth-4, m.height-8)
			m.viewport.SetContent(m.renderMessages())
			m.ready = true
		} else {
			m.viewport.Width = m.width - sidebarWidth - 4
			m.viewport.Height = m.height - 8
		}

		return m, nil

	case tea.KeyMsg:
		// Debug: capture key details
		m.lastKey = fmt.Sprintf("K:t=%d s=%q r=%v", msg.Type, msg.String(), msg.Runes)

		// Global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// If in input panel, handle input first
		if m.activePanel == inputPanel {
			return m.updateInputPanel(msg)
		}

		// Panel navigation (lazygit style) - only when not in input
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "1":
			m.activePanel = channelsPanel
			return m, nil
		case "2":
			m.activePanel = dmsPanel
			return m, nil
		case "3":
			m.activePanel = messagesPanel
			return m, nil
		case "4":
			m.activePanel = inputPanel
			m.textInput.Focus()
			return m, nil
		}

		// Panel-specific controls
		switch m.activePanel {
		case channelsPanel:
			return m.updateChannelsPanel(msg)
		case dmsPanel:
			return m.updateDMsPanel(msg)
		case messagesPanel:
			return m.updateMessagesPanel(msg)
		}
	}

	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) updateChannelsPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	channels := m.getChannels()

	// Check for Enter key - multiple ways it can arrive via SSH
	isEnter := msg.Type == tea.KeyEnter ||
		msg.String() == "enter" ||
		(len(msg.Runes) == 1 && (msg.Runes[0] == '\r' || msg.Runes[0] == '\n'))

	if isEnter {
		m.activePanel = messagesPanel
		return m, nil
	}

	switch msg.String() {
	case "j", "down":
		if m.selectedChannelIdx < len(channels)-1 {
			m.selectedChannelIdx++
		}
	case "k", "up":
		if m.selectedChannelIdx > 0 {
			m.selectedChannelIdx--
		}
	}
	return m, nil
}

func (m Model) updateDMsPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	dms := m.getDMs()

	// Check for Enter key - multiple ways it can arrive via SSH
	isEnter := msg.Type == tea.KeyEnter ||
		msg.String() == "enter" ||
		(len(msg.Runes) == 1 && (msg.Runes[0] == '\r' || msg.Runes[0] == '\n'))

	if isEnter {
		m.activePanel = messagesPanel
		return m, nil
	}

	switch msg.String() {
	case "j", "down":
		if m.selectedDMIdx < len(dms)-1 {
			m.selectedDMIdx++
		}
	case "k", "up":
		if m.selectedDMIdx > 0 {
			m.selectedDMIdx--
		}
	}
	return m, nil
}

func (m Model) updateMessagesPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "i":
		m.activePanel = inputPanel
		m.textInput.Focus()
		return m, nil
	case "j", "down":
		m.viewport.LineDown(1)
	case "k", "up":
		m.viewport.LineUp(1)
	case "d":
		m.viewport.HalfViewDown()
	case "u":
		m.viewport.HalfViewUp()
	case "g":
		m.viewport.GotoTop()
	case "G":
		m.viewport.GotoBottom()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) updateInputPanel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Debug: store last key info
	m.lastKey = fmt.Sprintf("type=%d str=%q runes=%v", msg.Type, msg.String(), msg.Runes)

	// Check for Enter key - multiple ways it can arrive via SSH
	isEnter := msg.Type == tea.KeyEnter ||
		msg.String() == "enter" ||
		msg.String() == "ctrl+j" ||
		(len(msg.Runes) == 1 && (msg.Runes[0] == '\r' || msg.Runes[0] == '\n'))

	// Check for Esc key
	isEsc := msg.Type == tea.KeyEsc || msg.String() == "esc"

	if isEsc {
		m.textInput.Blur()
		m.activePanel = messagesPanel
		return m, nil
	}

	if isEnter {
		if m.textInput.Value() != "" {
			// Add message
			newMsg := Message{
				Time:     time.Now().Format("15:04"),
				Author:   m.user.Username,
				Icon:     "🔒",
				Content:  m.textInput.Value(),
				IsLocked: true,
			}
			m.messages = append(m.messages, newMsg)
			m.viewport.SetContent(m.renderMessages())
			m.viewport.GotoBottom()
			m.textInput.SetValue("")
		}
		return m, nil
	}

	// Pass all other keys to textInput
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	sidebar := m.renderSidebar()
	chat := m.renderChat()

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		chat,
	)
}

func (m Model) renderSidebar() string {
	var boxes []string

	// Header box (no navigation)
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(0, 1)

	headerText := fmt.Sprintf("Welcome: %s\nto pingg.me", m.user.Username)
	headerContent := headerStyle.Render(headerText)

	headerBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238"))

	headerBox := headerBoxStyle.Render(headerContent)
	boxes = append(boxes, headerBox)

	// Channels box
	channels := m.getChannels()
	channelsContent := m.renderChannelsList(channels)

	channelsBorderColor := lipgloss.Color("238")
	if m.activePanel == channelsPanel {
		channelsBorderColor = lipgloss.Color("86")
	}

	channelsBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(channelsBorderColor).
		Padding(0, 1)

	channelsBox := channelsBoxStyle.Render(channelsContent)
	boxes = append(boxes, channelsBox)

	// DMs box
	dms := m.getDMs()
	dmsContent := m.renderDMsList(dms)

	dmsBorderColor := lipgloss.Color("238")
	if m.activePanel == dmsPanel {
		dmsBorderColor = lipgloss.Color("86")
	}

	dmsBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(dmsBorderColor).
		Padding(0, 1)

	dmsBox := dmsBoxStyle.Render(dmsContent)
	boxes = append(boxes, dmsBox)

	sidebarStyle := lipgloss.NewStyle().
		Width(sidebarWidth).
		Height(m.height).
		Padding(0, 1)

	return sidebarStyle.Render(lipgloss.JoinVertical(lipgloss.Left, boxes...))
}

func (m Model) renderChannelsList(channels []*models.Channel) string {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Bold(true)

	sb.WriteString(titleStyle.Render("(1) Channels"))
	sb.WriteString("\n\n")

	for i, channel := range channels {
		prefix := "- "
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))

		if i == m.selectedChannelIdx {
			if m.activePanel == channelsPanel {
				prefix = "◇ "
				style = style.Bold(true).Foreground(lipgloss.Color("86"))
			} else {
				prefix = "◇ "
				style = style.Foreground(lipgloss.Color("250"))
			}
		}

		channelIcon := "#"
		if channel.IsPrivate {
			channelIcon = "🔒"
		}

		line := fmt.Sprintf("%s%s %s", prefix, channelIcon, channel.Name)
		sb.WriteString(style.Render(line))
		if i < len(channels)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (m Model) renderDMsList(dms []*models.Channel) string {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Bold(true)

	sb.WriteString(titleStyle.Render("(2) Direct Messages"))
	sb.WriteString("\n\n")

	for i, dm := range dms {
		prefix := "- "
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))

		if i == m.selectedDMIdx {
			if m.activePanel == dmsPanel {
				prefix = "◇ "
				style = style.Bold(true).Foreground(lipgloss.Color("86"))
			} else {
				prefix = "◇ "
				style = style.Foreground(lipgloss.Color("250"))
			}
		}

		line := fmt.Sprintf("%s◇ %s", prefix, dm.Name)
		sb.WriteString(style.Render(line))
		if i < len(dms)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (m Model) renderChat() string {
	// Header box (no navigation)
	selectedChannel := m.getSelectedChannel()
	var headerText string
	if selectedChannel != nil {
		prefix := "#"
		if selectedChannel.IsDM {
			prefix = "💬"
		} else if selectedChannel.IsPrivate {
			prefix = "🔒"
		}
		headerText = fmt.Sprintf("%s %s", prefix, selectedChannel.Name)
		if selectedChannel.Description != "" {
			headerText = fmt.Sprintf("%s\n%s", headerText, selectedChannel.Description)
		} else {
			headerText = fmt.Sprintf("%s\n ", headerText)
		}
	} else {
		headerText = "Select a channel\n "
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Padding(0, 1)
	headerContent := headerStyle.Render(headerText)

	headerBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238"))

	headerBox := headerBoxStyle.Render(headerContent)

	// Messages box
	messagesBorderColor := lipgloss.Color("238")
	if m.activePanel == messagesPanel {
		messagesBorderColor = lipgloss.Color("86") // Green when focused
	}

	var messagesContent string
	if m.ready {
		messagesContent = m.viewport.View()
	} else {
		messagesContent = m.renderMessages()
	}

	messagesInnerStyle := lipgloss.NewStyle().
		Padding(0, 1)
	messagesInner := messagesInnerStyle.Render(messagesContent)

	messagesStyle := lipgloss.NewStyle().
		Height(m.height - 14).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(messagesBorderColor)

	messagesBox := messagesStyle.Render(messagesInner)

	// Input box
	inputBorderColor := lipgloss.Color("238")
	if m.activePanel == inputPanel {
		inputBorderColor = lipgloss.Color("86") // Green when focused
	}

	inputContent := fmt.Sprintf("> %s", m.textInput.View())

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(inputBorderColor).
		Padding(0, 1)

	inputBox := inputStyle.Render(inputContent)

	// Help with debug info
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(0, 1)

	debugInfo := ""
	if m.lastKey != "" {
		debugInfo = fmt.Sprintf(" | DEBUG: %s", m.lastKey)
	}
	help := helpStyle.Render(fmt.Sprintf("1-4:panels | j/k:nav | Enter:send | esc:back | w=%d h=%d%s", m.width, m.height, debugInfo))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerBox,
		messagesBox,
		inputBox,
		help,
	)
}

func (m Model) renderMessages() string {
	var sb strings.Builder

	for i, msg := range m.messages {
		msgStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))

		line := fmt.Sprintf("%s %s %s: %s",
			msg.Time,
			msg.Icon,
			msg.Author,
			msg.Content,
		)
		sb.WriteString(msgStyle.Render(line))

		if i < len(m.messages)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (m Model) getChannels() []*models.Channel {
	var channels []*models.Channel
	for _, ch := range m.user.Channels {
		if !ch.IsDM {
			channels = append(channels, ch)
		}
	}
	return channels
}

func (m Model) getDMs() []*models.Channel {
	var dms []*models.Channel
	for _, ch := range m.user.Channels {
		if ch.IsDM {
			dms = append(dms, ch)
		}
	}
	return dms
}

func (m Model) getSelectedChannel() *models.Channel {
	if m.activePanel == channelsPanel || (m.activePanel != dmsPanel && m.selectedChannelIdx >= 0) {
		channels := m.getChannels()
		if m.selectedChannelIdx >= 0 && m.selectedChannelIdx < len(channels) {
			return channels[m.selectedChannelIdx]
		}
	} else if m.activePanel == dmsPanel || m.selectedDMIdx >= 0 {
		dms := m.getDMs()
		if m.selectedDMIdx >= 0 && m.selectedDMIdx < len(dms) {
			return dms[m.selectedDMIdx]
		}
	}
	return nil
}

// SetTerminalSize sets the terminal dimensions from environment variables
// This is useful when running via SSH where the terminal size might not be detected
func SetTerminalSize(width, height int) {
	envWidth = width
	envHeight = height
}

// LoadTerminalSizeFromEnv loads terminal dimensions from environment variables
// Supports both WHISP_COLS/WHISP_ROWS (from Whisp SSH) and TERM_WIDTH/TERM_HEIGHT
func LoadTerminalSizeFromEnv() {
	// Try WHISP variables first (from SSH)
	if w := os.Getenv("WHISP_COLS"); w != "" {
		if width, err := strconv.Atoi(w); err == nil && width > 0 {
			envWidth = width
		}
	} else if w := os.Getenv("TERM_WIDTH"); w != "" {
		if width, err := strconv.Atoi(w); err == nil && width > 0 {
			envWidth = width
		}
	}

	if h := os.Getenv("WHISP_ROWS"); h != "" {
		if height, err := strconv.Atoi(h); err == nil && height > 0 {
			envHeight = height
		}
	} else if h := os.Getenv("TERM_HEIGHT"); h != "" {
		if height, err := strconv.Atoi(h); err == nil && height > 0 {
			envHeight = height
		}
	}
}

func Run(user *models.User) error {
	// Load terminal size from environment if available
	LoadTerminalSizeFromEnv()

	// Force color profile for SSH sessions
	// This ensures colors work correctly when running via SSH
	lipgloss.SetColorProfile(termenv.TrueColor)

	// Put terminal in raw mode for SSH sessions
	// This is necessary because Whisp sets canonical mode by default
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		oldState, err := term.MakeRaw(fd)
		if err == nil {
			defer term.Restore(fd, oldState)
		}
	}

	model := NewModel(user)

	// If we already have dimensions from env, set them
	if envWidth > 0 && envHeight > 0 {
		model.width = envWidth
		model.height = envHeight
	}

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		// Enable explicit input/output for SSH sessions
		tea.WithInput(os.Stdin),
		tea.WithOutput(os.Stdout),
	)
	_, err := p.Run()
	return err
}
