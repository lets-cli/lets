package downloadprogress

import (
	"fmt"
	"io"
	"math"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	bubblesprogress "charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
	"github.com/lets-cli/lets/internal/fetch"
)

const (
	defaultWidth = 80
	minBarWidth  = 10
	maxBarWidth  = 30
)

type Observer struct {
	writer     io.Writer
	width      int
	noColor    bool
	animate    bool
	throttle   time.Duration
	finalPause time.Duration
	now        func() time.Time
}

type Option func(*Observer)

func WithWidth(width int) Option {
	return func(observer *Observer) {
		observer.width = width
	}
}

func WithNoColor(noColor bool) Option {
	return func(observer *Observer) {
		observer.noColor = noColor
	}
}

func WithThrottle(throttle time.Duration) Option {
	return func(observer *Observer) {
		observer.throttle = throttle
	}
}

func WithFinalPause(finalPause time.Duration) Option {
	return func(observer *Observer) {
		observer.finalPause = finalPause
	}
}

func WithNow(now func() time.Time) Option {
	return func(observer *Observer) {
		observer.now = now
	}
}

func New(writer io.Writer, options ...Option) *Observer {
	observer := &Observer{
		writer:     writer,
		width:      detectWidth(writer),
		throttle:   100 * time.Millisecond,
		finalPause: 750 * time.Millisecond,
		animate:    isTerminal(writer),
		now:        time.Now,
	}

	for _, option := range options {
		option(observer)
	}

	if observer.width <= 0 {
		observer.width = defaultWidth
	}

	if observer.now == nil {
		observer.now = time.Now
	}

	return observer
}

func (o *Observer) Start(info fetch.ProgressInfo) fetch.ProgressTracker { //nolint:ireturn // Implements fetch.ProgressObserver.
	if info.TotalBytes > 0 && o.animate {
		return newAnimatedTracker(o, info)
	}

	tracker := &manualTracker{
		observer: o,
		info:     info,
		label:    downloadLabel(info.URL),
	}
	tracker.render(false)

	return tracker
}

type manualTracker struct {
	observer   *Observer
	info       fetch.ProgressInfo
	label      string
	read       int64
	started    bool
	lastRender time.Time
	lastWidth  int
}

func (t *manualTracker) Add(n int64) {
	t.read += n

	now := t.observer.now()
	if t.observer.throttle > 0 && !t.lastRender.IsZero() && now.Sub(t.lastRender) < t.observer.throttle {
		return
	}

	t.render(false)
}

func (t *manualTracker) Done(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(t.observer.writer)
		return
	}

	t.render(true)
}

func (t *manualTracker) render(done bool) {
	if !t.started {
		_, _ = fmt.Fprintln(t.observer.writer, t.labelLine())
		t.started = true
	}

	line := t.progressLine()
	if done {
		_, _ = fmt.Fprintf(t.observer.writer, "\r\033[2K%s\n", line)
		t.lastWidth = 0
		t.lastRender = t.observer.now()

		return
	}

	lineWidth := lipgloss.Width(line)

	padding := ""
	if t.lastWidth > lineWidth {
		padding = strings.Repeat(" ", t.lastWidth-lineWidth)
	}

	_, _ = fmt.Fprintf(t.observer.writer, "\r%s%s", line, padding)

	t.lastWidth = lineWidth
	t.lastRender = t.observer.now()
}

func (t *manualTracker) labelLine() string {
	return labelLine("Downloading", t.label, t.observer.width)
}

func (t *manualTracker) progressLine() string {
	if t.info.TotalBytes > 0 {
		return t.knownSizeLine()
	}

	return t.unknownSizeLine()
}

func (t *manualTracker) knownSizeLine() string {
	percent := clamp(float64(t.read)/float64(t.info.TotalBytes), 0, 1)
	suffix := fmt.Sprintf("%3.0f%% %s/%s", percent*100, formatBytes(t.read), formatBytes(t.info.TotalBytes))

	bar := ""

	barWidth := t.barWidth(suffix)
	if barWidth >= minBarWidth {
		model := t.progressModel(barWidth)
		bar = model.ViewAs(percent)
	}

	if bar == "" {
		return suffix
	}

	return fmt.Sprintf("%s %s", bar, suffix)
}

func (t *manualTracker) unknownSizeLine() string {
	return formatBytes(t.read)
}

func (t *manualTracker) barWidth(suffix string) int {
	spaceForBar := t.observer.width - 1 - lipgloss.Width(suffix)
	if spaceForBar < minBarWidth {
		return 0
	}

	return min(maxBarWidth, max(minBarWidth, spaceForBar))
}

func (t *manualTracker) progressModel(width int) bubblesprogress.Model {
	model := bubblesprogress.New(
		bubblesprogress.WithWidth(width),
		bubblesprogress.WithoutPercentage(),
		bubblesprogress.WithFillCharacters('#', '-'),
	)
	if t.observer.noColor {
		model.FullColor = nil
		model.EmptyColor = nil
	}

	return model
}

type animatedTracker struct {
	observer   *Observer
	program    *tea.Program
	done       chan struct{}
	label      string
	read       int64
	total      int64
	lastUpdate time.Time
}

func newAnimatedTracker(observer *Observer, info fetch.ProgressInfo) *animatedTracker {
	ready := make(chan struct{})
	done := make(chan struct{})
	label := downloadLabel(info.URL)
	_, _ = fmt.Fprintln(observer.writer, labelLine("Downloading", label, observer.width))
	model := newProgressModel(observer, label, info.TotalBytes, ready)
	program := tea.NewProgram(
		model,
		tea.WithInput(nil),
		tea.WithOutput(observer.writer),
		tea.WithoutSignals(),
	)

	tracker := &animatedTracker{
		observer: observer,
		program:  program,
		done:     done,
		label:    label,
		total:    info.TotalBytes,
	}

	go func() {
		_, _ = program.Run()

		close(done)
	}()

	<-ready

	return tracker
}

func (t *animatedTracker) Add(n int64) {
	t.read += n

	now := t.observer.now()
	if t.observer.throttle > 0 && !t.lastUpdate.IsZero() && now.Sub(t.lastUpdate) < t.observer.throttle {
		return
	}

	t.program.Send(progressMsg{read: t.read, total: t.total})
	t.lastUpdate = now
}

func (t *animatedTracker) Done(err error) {
	if err != nil {
		t.program.Send(progressErrMsg{})
	} else {
		t.program.Send(progressDoneMsg{read: t.read, total: t.total})
	}

	<-t.done

	if err == nil {
		_, _ = fmt.Fprintf(t.observer.writer, "%s\n", t.progressLine())
	}
}

func (t *animatedTracker) progressLine() string {
	bar := t.progressModel().ViewAs(1)
	return fmt.Sprintf("%s 100%% %s/%s", bar, formatBytes(t.total), formatBytes(t.total))
}

func (t *animatedTracker) progressModel() bubblesprogress.Model {
	model := bubblesprogress.New(
		bubblesprogress.WithWidth(barWidthForTerminal(t.observer.width)),
		bubblesprogress.WithDefaultBlend(),
		bubblesprogress.WithoutPercentage(),
		bubblesprogress.WithFillCharacters('#', '-'),
	)
	if t.observer.noColor {
		model.FullColor = nil
		model.EmptyColor = nil
	}

	return model
}

type progressMsg struct {
	read  int64
	total int64
}

type progressDoneMsg struct {
	read  int64
	total int64
}

type progressErrMsg struct{}

type progressQuitMsg struct{}

type progressModel struct {
	label      string
	read       int64
	total      int64
	width      int
	finalPause time.Duration
	ready      chan struct{}
	readyOnce  *sync.Once
	progress   bubblesprogress.Model
}

func newProgressModel(observer *Observer, label string, total int64, ready chan struct{}) progressModel {
	model := bubblesprogress.New(
		bubblesprogress.WithWidth(barWidthForTerminal(observer.width)),
		bubblesprogress.WithDefaultBlend(),
		bubblesprogress.WithoutPercentage(),
		bubblesprogress.WithFillCharacters('#', '-'),
	)
	if observer.noColor {
		model.FullColor = nil
		model.EmptyColor = nil
	}

	return progressModel{
		label:      label,
		total:      total,
		width:      observer.width,
		finalPause: observer.finalPause,
		ready:      ready,
		readyOnce:  &sync.Once{},
		progress:   model,
	}
}

func (m progressModel) Init() tea.Cmd {
	m.readyOnce.Do(func() {
		close(m.ready)
	})

	return nil
}

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn // Required by Bubble Tea's model interface.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.progress.SetWidth(barWidthForTerminal(msg.Width))

		return m, nil

	case progressMsg:
		m.read = msg.read
		m.total = msg.total

		return m, m.progress.SetPercent(m.percent())

	case progressDoneMsg:
		m.read = msg.read
		m.total = msg.total

		return m, tea.Batch(m.progress.SetPercent(1), m.quitAfterFinalPause())

	case progressErrMsg:
		return m, tea.Quit

	case progressQuitMsg:
		return m, tea.Quit

	case bubblesprogress.FrameMsg:
		var cmd tea.Cmd

		m.progress, cmd = m.progress.Update(msg)

		return m, cmd

	default:
		return m, nil
	}
}

func (m progressModel) View() tea.View {
	return tea.NewView(m.progressLine())
}

func (m progressModel) progressLine() string {
	return fmt.Sprintf("%s %3.0f%% %s/%s", m.progress.View(), m.percent()*100, formatBytes(m.read), formatBytes(m.total))
}

func (m progressModel) percent() float64 {
	if m.total <= 0 {
		return 0
	}

	return clamp(float64(m.read)/float64(m.total), 0, 1)
}

func (m progressModel) quitAfterFinalPause() tea.Cmd {
	return tea.Tick(m.finalPause, func(time.Time) tea.Msg {
		return progressQuitMsg{}
	})
}

func barWidthForTerminal(width int) int {
	suffixWidth := lipgloss.Width(" 100% 1023.9 KiB/1023.9 KiB")

	spaceForBar := width - 1 - suffixWidth
	if spaceForBar < minBarWidth {
		return minBarWidth
	}

	return min(maxBarWidth, max(minBarWidth, spaceForBar))
}

func detectWidth(writer io.Writer) int {
	file, ok := writer.(term.File)
	if !ok || !isTerminal(writer) {
		return defaultWidth
	}

	width, _, err := term.GetSize(file.Fd())
	if err != nil || width <= 0 {
		return defaultWidth
	}

	return width
}

func isTerminal(writer io.Writer) bool {
	file, ok := writer.(term.File)
	return ok && term.IsTerminal(file.Fd())
}

func labelLine(verb, label string, width int) string {
	return fmt.Sprintf("%s %s", verb, truncateMiddle(label, width-lipgloss.Width(verb)-1))
}

func downloadLabel(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return path.Base(stripSecretParts(rawURL))
	}

	filename := path.Base(parsed.Path)
	if filename != "." && filename != "/" && filename != "" {
		return filename
	}

	if parsed.Host != "" {
		return parsed.Host
	}

	return redactURL(rawURL)
}

func redactURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return stripSecretParts(rawURL)
	}

	parsed.User = nil
	parsed.RawQuery = ""
	parsed.Fragment = ""

	return parsed.String()
}

func stripSecretParts(rawURL string) string {
	idx := len(rawURL)
	if queryIdx := strings.Index(rawURL, "?"); queryIdx >= 0 && queryIdx < idx {
		idx = queryIdx
	}

	if fragmentIdx := strings.Index(rawURL, "#"); fragmentIdx >= 0 && fragmentIdx < idx {
		idx = fragmentIdx
	}

	return rawURL[:idx]
}

func truncateMiddle(s string, width int) string {
	if width <= 0 {
		return ""
	}

	if lipgloss.Width(s) <= width {
		return s
	}

	if width <= 3 {
		return strings.Repeat(".", width)
	}

	runes := []rune(s)
	keep := width - 3
	left := keep / 2
	right := keep - left

	if left+right >= len(runes) {
		return s
	}

	return string(runes[:left]) + "..." + string(runes[len(runes)-right:])
}

func formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}

	units := []string{"KiB", "MiB", "GiB", "TiB"}
	value := float64(bytes) / 1024

	unit := units[0]
	for _, nextUnit := range units[1:] {
		if value < 1024 {
			break
		}

		value /= 1024
		unit = nextUnit
	}

	return fmt.Sprintf("%.1f %s", value, unit)
}

func clamp(value, minValue, maxValue float64) float64 {
	return math.Max(minValue, math.Min(maxValue, value))
}
