package importer

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ParsedTask struct {
	Content        string
	EstimatedHours float64
	ResourceURLs   []string
	SortOrder      int
	IsCheckpoint   bool
}

type ParsedDay struct {
	DayNumber int
	Title     string
	SortOrder int
	Tasks     []ParsedTask
}

type ParsedWeek struct {
	WeekNumber  int
	Title       string
	Subtitle    string
	Goal        string
	Deliverables string
	SortOrder   int
	Days        []ParsedDay
}

type ParsedPhase struct {
	PhaseNumber  int
	Title        string
	Subtitle     string
	Goal         string
	Deliverables string
	SortOrder    int
	Weeks        []ParsedWeek
}

type ParsedData struct {
	Phases []ParsedPhase
}

var (
	rePhaseH2  = regexp.MustCompile(`^## Phase (\d+)\s*[:：](.+)`)
	reWeekH3   = regexp.MustCompile(`^### 第 (\d+) 周\s*[:：](.+)`)
	reDayH4    = regexp.MustCompile(`^#### Day (\d+)\s*[:：(（].+[)）]\s*(.+)`)
	reDayH4Alt = regexp.MustCompile(`^#### Day (\d+)\s*[-–—](.+)`)
	reTask     = regexp.MustCompile(`^-\s*\[[ x]\]\s*(.+)`)
	reCheckpoint = regexp.MustCompile(`^\*\*自检清单[：:]*\*\*`)
	reURL      = regexp.MustCompile(`\((https?://[^)]+)\)`)
	reBracket  = regexp.MustCompile(`\[([^\]]+)\]$$`)
)

const (
	stateOutside = iota
	stateInPhase
	stateInWeek
	stateInDay
	stateInCheckpoint
)

type parser struct {
	state         int
	data          *ParsedData
	currentPhase  *ParsedPhase
	currentWeek   *ParsedWeek
	currentDay    *ParsedDay
	phaseGoal     strings.Builder
	weekGoal      strings.Builder
	taskOrder     int
	checkpointOrder int

	weekCount int
	dayCount  int
}

func Parse(filePath string) (*ParsedData, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	p := &parser{
		data:  &ParsedData{},
		state: stateOutside,
	}
	p.resetGoal()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		p.processLine(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	p.flushGoal()
	p.flushDay()
	p.flushWeek()
	p.flushPhase()
	return p.data, nil
}

func (p *parser) resetGoal() {
	p.phaseGoal.Reset()
	p.weekGoal.Reset()
}

func (p *parser) processLine(line string) {
	switch {
	case rePhaseH2.MatchString(line):
		p.flushDay()
		p.flushWeek()
		p.flushPhase()
		matches := rePhaseH2.FindStringSubmatch(line)
		num, _ := strconv.Atoi(matches[1])
		p.currentPhase = &ParsedPhase{
			PhaseNumber: num,
			Title:       strings.TrimSpace(matches[2]),
			SortOrder:   num,
		}
		p.state = stateInPhase
		p.phaseGoal.Reset()
		p.weekCount = 0
		p.currentWeek = nil

	case reWeekH3.MatchString(line):
		p.flushDay()
		p.flushWeek()
		matches := reWeekH3.FindStringSubmatch(line)
		num, _ := strconv.Atoi(matches[1])
		p.weekCount++
		p.currentWeek = &ParsedWeek{
			WeekNumber: num,
			Title:      strings.TrimSpace(matches[2]),
			SortOrder:  num,
		}
		p.state = stateInWeek
		p.weekGoal.Reset()
		p.dayCount = 0
		p.currentDay = nil

	case reDayH4.MatchString(line):
		p.flushDay()
		p.processDayHeader(reDayH4.FindStringSubmatch(line))
	case reDayH4Alt.MatchString(line):
		p.flushDay()
		p.processDayHeader(reDayH4Alt.FindStringSubmatch(line))

	case reCheckpoint.MatchString(line):
		p.flushDay()
		p.state = stateInCheckpoint

	case reTask.MatchString(line):
		if p.state == stateInDay || p.state == stateInCheckpoint {
			p.processTask(reTask.FindStringSubmatch(line))
		}

	default:
		p.collectBlockquote(line)
	}
}

func (p *parser) processDayHeader(matches []string) {
	num, _ := strconv.Atoi(matches[1])
	title := strings.TrimSpace(matches[2])
	p.dayCount++
	p.currentDay = &ParsedDay{
		DayNumber: num,
		Title:     title,
		SortOrder: num,
	}
	p.state = stateInDay
	p.taskOrder = 0
	p.checkpointOrder = 0
}

func (p *parser) processTask(matches []string) {
	content := strings.TrimSpace(matches[1])
	urls := extractURLs(content)
	content = cleanContent(content)

	task := ParsedTask{
		Content:      content,
		ResourceURLs: urls,
		IsCheckpoint: p.state == stateInCheckpoint,
	}
	if p.state == stateInDay {
		p.taskOrder++
		task.SortOrder = p.taskOrder
	} else {
		p.checkpointOrder++
		task.SortOrder = p.checkpointOrder
	}
	p.currentDay.Tasks = append(p.currentDay.Tasks, task)
}

func (p *parser) collectBlockquote(line string) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || trimmed == "---" || trimmed == "[TOC]" {
		return
	}

	switch p.state {
	case stateInPhase:
		if p.currentWeek == nil {
			if strings.HasPrefix(trimmed, "> **阶段目标：**") {
				p.phaseGoal.WriteString(strings.TrimPrefix(trimmed, "> **阶段目标：**"))
			} else if strings.HasPrefix(trimmed, "> **阶段产出：**") {
				p.currentPhase.Deliverables = strings.TrimPrefix(trimmed, "> **阶段产出：**")
			} else if strings.HasPrefix(trimmed, "> ") {
				line := strings.TrimPrefix(trimmed, "> ")
				if p.currentPhase.Subtitle == "" {
					p.currentPhase.Subtitle = line
				}
				if p.currentWeek == nil {
					p.phaseGoal.WriteString(line)
				}
			}
		}
	case stateInWeek:
		if p.currentDay == nil {
			if strings.HasPrefix(trimmed, "> **本周目标：**") {
				p.weekGoal.WriteString(strings.TrimPrefix(trimmed, "> **本周目标：**"))
			} else if strings.HasPrefix(trimmed, "> **本周产出：**") {
				p.currentWeek.Deliverables = strings.TrimPrefix(trimmed, "> **本周产出：**")
			} else if strings.HasPrefix(trimmed, "> ") {
				line := strings.TrimPrefix(trimmed, "> ")
				if p.currentDay == nil {
					p.weekGoal.WriteString(line)
				}
			}
		}
	}
}

func (p *parser) flushGoal() {
	if p.currentPhase != nil && p.phaseGoal.Len() > 0 {
		p.currentPhase.Goal = strings.TrimSpace(p.phaseGoal.String())
	}
	if p.currentWeek != nil && p.weekGoal.Len() > 0 {
		p.currentWeek.Goal = strings.TrimSpace(p.weekGoal.String())
	}
}

func (p *parser) flushDay() {
	if p.currentDay != nil && len(p.currentDay.Tasks) > 0 && p.currentWeek != nil {
		p.currentWeek.Days = append(p.currentWeek.Days, *p.currentDay)
	}
	p.currentDay = nil
}

func (p *parser) flushWeek() {
	if p.currentWeek != nil && p.currentPhase != nil {
		if len(p.currentWeek.Days) > 0 {
			p.currentPhase.Weeks = append(p.currentPhase.Weeks, *p.currentWeek)
		}
	}
	p.currentWeek = nil
}

func (p *parser) flushPhase() {
	if p.currentPhase != nil {
		p.currentPhase.Goal = strings.TrimSpace(p.phaseGoal.String())
		if len(p.currentPhase.Weeks) > 0 {
			p.data.Phases = append(p.data.Phases, *p.currentPhase)
		}
	}
	p.currentPhase = nil
}

func extractURLs(content string) []string {
	matches := reURL.FindAllStringSubmatch(content, -1)
	var urls []string
	for _, m := range matches {
		urls = append(urls, m[1])
	}
	return urls
}

func cleanContent(content string) string {
	content = reURL.ReplaceAllString(content, "")
	content = reBracket.ReplaceAllString(content, "")
	content = strings.TrimSpace(content)
	runes := []rune(content)
	if len(runes) > 0 && runes[len(runes)-1] == '(' {
		runes = runes[:len(runes)-1]
	}
	return strings.TrimSpace(string(runes))
}
