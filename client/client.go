package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	p := tea.NewProgram(initialModel(conn))

	go getResponse(p, conn)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg    error
	serverMsg string
)

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
	conn        net.Conn
	ready       bool
}

func initialModel(conn net.Conn) model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "> "
	ta.CharLimit = 256

	ta.SetHeight(1)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messages:    []string{},
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		conn:        conn,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.Height = msg.Height - 4
			m.viewport.SetContent("Connected to EchoTea Server. Press return to send message.")
			m.ready = true
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.messages = append(m.messages, m.senderStyle.Render("Client: ")+m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))

			_, err := m.conn.Write([]byte(m.textarea.Value()))
			if err != nil {
				log.Fatal(err)
			}

			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	case serverMsg:
		m.messages = append(m.messages, m.senderStyle.Render("Server: ")+strings.TrimRight(string(msg), "\n"))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.textarea.Reset()
		m.viewport.GotoBottom()

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func getResponse(p *tea.Program, conn net.Conn) {
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		p.Send(serverMsg(msg))
	}
}
