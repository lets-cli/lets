package progressbar

import (
	"fmt"
	"image/color"
	"io"
	"math"
	"net/url"
	"path"
	"strings"
	"time"

	bubblesprogress "charm.land/bubbles/v2/progress"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
	"github.com/lets-cli/lets/internal/fetch"
	"github.com/lets-cli/lets/internal/theme"
	"github.com/lets-cli/lets/internal/util"
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
	throttle   time.Duration
	fillColor  color.Color
	emptyColor color.Color
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

func WithTheme(themeName string) Option {
	return func(observer *Observer) {
		if observer.noColor {
			return
		}

		file, ok := observer.writer.(term.File)
		if !ok || !term.IsTerminal(file.Fd()) {
			return
		}

		observer.fillColor, observer.emptyColor = theme.ProgressColorsByName(themeName, file)
	}
}

func WithThrottle(throttle time.Duration) Option {
	return func(observer *Observer) {
		observer.throttle = throttle
	}
}

func WithNow(now func() time.Time) Option {
	return func(observer *Observer) {
		observer.now = now
	}
}

func New(writer io.Writer, options ...Option) *Observer {
	observer := &Observer{
		writer:   writer,
		width:    detectWidth(writer),
		throttle: 100 * time.Millisecond,
		now:      time.Now,
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
	} else {
		applyProgressColors(&model, t.observer.fillColor, t.observer.emptyColor)
	}

	return model
}

func detectWidth(writer io.Writer) int {
	file, ok := writer.(term.File)
	if !ok || !util.IsTerminalWriter(writer) {
		return defaultWidth
	}

	width, _, err := term.GetSize(file.Fd())
	if err != nil || width <= 0 {
		return defaultWidth
	}

	return width
}

func applyProgressColors(model *bubblesprogress.Model, fill, empty color.Color) {
	if fill != nil {
		model.FullColor = fill
	}

	if empty != nil {
		model.EmptyColor = empty
	}
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
