package test

import (
	"os"
	"testing"

	"github.com/xyd/web3-learning-tracker/internal/importer"
)

func TestParse_RealFile(t *testing.T) {
	data, err := importer.Parse("../../sources/web3_infra_3month_plan.md")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(data.Phases) != 3 {
		t.Errorf("expected 3 phases, got %d", len(data.Phases))
	}

	totalWeeks := 0
	totalDays := 0
	totalTasks := 0
	for _, p := range data.Phases {
		if p.PhaseNumber == 0 {
			t.Error("phase number should not be 0")
		}
		if p.Title == "" {
			t.Error("phase title should not be empty")
		}
		totalWeeks += len(p.Weeks)
		for _, w := range p.Weeks {
			if w.WeekNumber == 0 {
				t.Error("week number should not be 0")
			}
			totalDays += len(w.Days)
			for _, d := range w.Days {
				totalTasks += len(d.Tasks)
			}
		}
	}
	if totalWeeks != 12 {
		t.Errorf("expected 12 weeks, got %d", totalWeeks)
	}
	if totalDays != 84 {
		t.Errorf("expected 84 days, got %d", totalDays)
	}
	if totalTasks != 232 {
		t.Errorf("expected 232 tasks, got %d", totalTasks)
	}
}

func TestParse_Fixture(t *testing.T) {
	fixture := `# Test Plan

## Phase 1：测试阶段

> **阶段目标：** 测试目标描述
> **阶段产出：** 测试产出

### 第 1 周：测试周

> **本周目标：** 本周测试目标

#### Day 1（周一）— 测试日
- [ ] 任务一
- [ ] 任务二 [参考链接](https://example.com)
- [ ] 任务三

#### Day 2（周二）— 第二日
- [ ] 任务四

**自检清单：**
- [ ] 检查项一
- [ ] 检查项二
`

	tmpFile := "/tmp/test_fixture.md"
	os.WriteFile(tmpFile, []byte(fixture), 0644)
	defer os.Remove(tmpFile)

	data, err := importer.Parse(tmpFile)
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	if len(data.Phases) != 1 {
		t.Fatalf("expected 1 phase, got %d", len(data.Phases))
	}
	p := data.Phases[0]
	if p.PhaseNumber != 1 {
		t.Errorf("expected phase 1, got %d", p.PhaseNumber)
	}
	if p.Goal != "测试目标描述" {
		t.Errorf("expected goal, got %q", p.Goal)
	}
	if len(p.Weeks) != 1 {
		t.Fatalf("expected 1 week, got %d", len(p.Weeks))
	}
	w := p.Weeks[0]
	if w.WeekNumber != 1 {
		t.Errorf("expected week 1, got %d", w.WeekNumber)
	}
	if len(w.Days) != 2 {
		t.Fatalf("expected 2 days, got %d", len(w.Days))
	}

	d1 := w.Days[0]
	if len(d1.Tasks) != 3 {
		t.Errorf("expected 3 tasks in day 1, got %d", len(d1.Tasks))
	}
	if len(d1.Tasks[1].ResourceURLs) != 1 {
		t.Error("task 2 should have resource URL")
	}

	d2 := w.Days[1]
	if len(d2.Tasks) != 3 {
		t.Errorf("expected 3 tasks in day 2 (1 regular + 2 checkpoints), got %d", len(d2.Tasks))
	}
	if !d2.Tasks[1].IsCheckpoint {
		t.Error("task after self-check should be checkpoint")
	}
	if !d2.Tasks[2].IsCheckpoint {
		t.Error("second checkpoint task should be checkpoint")
	}
}
