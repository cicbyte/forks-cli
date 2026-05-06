package clone

import (
	"fmt"
	"strings"

	"github.com/cicbyte/forks-cli/internal/utils"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CloneConfig clone 命令配置
type CloneConfig struct {
	Server    string
	Token     string
	RepoURL   string
	TargetDir string
	Force     bool
}

// CloneProcessor clone 处理器
type CloneProcessor struct {
	config *CloneConfig
}

func NewCloneProcessor(config *CloneConfig) *CloneProcessor {
	return &CloneProcessor{config: config}
}

type clonePhase int

const (
	phasePreparing clonePhase = iota
	phaseCloning
	phaseRemotes
	phaseDone
	phaseError
)

var (
	styleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	styleError   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	styleWarning = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	styleLabel   = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	styleURL     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	styleHint    = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	styleDim     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type cloneModel struct {
	phase      clonePhase
	spinner    spinner.Model
	info       *RepoInfo
	targetDir  string
	useMirror  bool
	prepareMsg string
	err        error
}

type phaseMsg struct{ phase clonePhase }
type prepareDoneMsg struct{ result PrepareResult }
type cloneDoneMsg struct {
	useMirror bool
	err       error
}
type remotesDoneMsg struct{ err error }

func newCloneModel(info *RepoInfo, targetDir string) cloneModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	return cloneModel{phase: phasePreparing, spinner: s, info: info, targetDir: targetDir}
}

func (m cloneModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg { return phaseMsg{phase: phasePreparing} })
}

func (m cloneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case phaseMsg:
		m.phase = msg.phase
		if m.phase == phaseDone || m.phase == phaseError {
			return m, tea.Quit
		}
		return m, nil
	case prepareDoneMsg:
		m.prepareMsg = msg.result.Message
		m.useMirror = msg.result.UseMirror
		m.phase = phaseCloning
		return m, nil
	case cloneDoneMsg:
		if msg.err != nil {
			m.err = msg.err
			m.phase = phaseError
			return m, tea.Quit
		}
		m.useMirror = msg.useMirror
		m.phase = phaseRemotes
		return m, nil
	case remotesDoneMsg:
		if msg.err != nil {
			m.err = msg.err
		}
		m.phase = phaseDone
		return m, tea.Quit
	}
	return m, nil
}

func (m cloneModel) View() string {
	switch m.phase {
	case phasePreparing:
		serverDisplay := m.info.ServerOrigin
		if serverDisplay == "" {
			serverDisplay = "(未配置)"
		}
		return fmt.Sprintf("  %s 正在联系镜像服务 %s ...", m.spinner.View(), serverDisplay)
	case phaseCloning:
		if m.useMirror {
			return fmt.Sprintf("  %s 正在克隆仓库（镜像加速）...", m.spinner.View())
		}
		return fmt.Sprintf("  %s 正在克隆仓库（直连）...", m.spinner.View())
	case phaseRemotes:
		return fmt.Sprintf("  %s 正在设置 remote...", m.spinner.View())
	case phaseDone, phaseError:
		return fmt.Sprintf("  %s 完成", m.spinner.View())
	}
	return ""
}

// Execute 执行 clone 流程（带 TUI 动画）
func (p *CloneProcessor) Execute() error {
	info, err := ResolveRepoInfo(p.config.RepoURL, p.config.Server)
	if err != nil {
		return err
	}

	targetDir := p.config.TargetDir
	if targetDir == "" {
		targetDir = strings.TrimSuffix(info.Repo, ".git")
	}

	type stepResult struct {
		prepareResult PrepareResult
		cloneErr      error
		remotesErr    error
	}
	result := &stepResult{}
	done := make(chan struct{})

	model := newCloneModel(info, targetDir)
	prog := tea.NewProgram(model)

	go func() {
		// 1. prepare
		var prepareResult PrepareResult
		if info.ServerOrigin != "" {
			prepareResult = PrepareFromServer(info.ServerOrigin, p.config.Token, info.Source, info.Author, info.Repo, p.config.Force)
		}
		result.prepareResult = prepareResult
		prog.Send(prepareDoneMsg{result: prepareResult})

		// 2. clone
		useMirror := prepareResult.UseMirror
		cloneURL := info.OriginalURL
		if useMirror {
			cloneURL = info.MirrorURL
		}
		var cloneErr error
		if useMirror {
			cloneErr = utils.RunGitNoProxy("clone", cloneURL, targetDir)
		} else {
			cloneErr = utils.RunGitSilent("clone", cloneURL, targetDir)
		}
		if cloneErr != nil {
			result.cloneErr = cloneErr
			prog.Send(cloneDoneMsg{useMirror: useMirror, err: cloneErr})
			close(done)
			return
		}
		prog.Send(cloneDoneMsg{useMirror: useMirror, err: nil})

		// 3. remotes
		cloneDir := utils.ResolveAbsDir(targetDir)
		var remotesErr error
		if err := utils.RunGitInDirSilent(cloneDir, "remote", "set-url", "origin", info.OriginalURL); err != nil {
			remotesErr = fmt.Errorf("设置 remote origin 失败: %w", err)
		}
		result.remotesErr = remotesErr
		prog.Send(remotesDoneMsg{err: remotesErr})
		close(done)
	}()

	_, _ = prog.Run()
	<-done

	if result.cloneErr != nil {
		printError(info, result.prepareResult.Message, result.cloneErr)
		return result.cloneErr
	}

	printDone(info)
	return nil
}

func printDone(info *RepoInfo) {
	fmt.Println()
	fmt.Println(styleSuccess.Render("  ✓ 克隆完成！"))
	fmt.Println()
	fmt.Printf("  %s  →  %s\n", styleLabel.Render("origin"), styleURL.Render(info.OriginalURL))
	fmt.Println()
}

func printError(info *RepoInfo, prepareMsg string, cloneErr error) {
	fmt.Println()
	fmt.Println(styleError.Render("  ✗ 克隆失败"))
	fmt.Println()
	if prepareMsg != "" {
		fmt.Println(styleWarning.Render(fmt.Sprintf("  %s", prepareMsg)))
		fmt.Println()
	}
	fmt.Println(styleDim.Render(fmt.Sprintf("  %v", cloneErr)))
	fmt.Println()
	if info.ServerOrigin == "" {
		fmt.Println(styleHint.Render("  提示: forks-cli config server <url>"))
		fmt.Println()
	}
}
