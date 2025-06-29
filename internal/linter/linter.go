package linter

import (
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"time"

	"github.com/keenbytes/octo-linter/v2/internal/linter/rule"
	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

const (
	HasNoErrorsOrWarnings = iota
	HasErrors
	HasOnlyWarnings
)

type Linter struct {
	Config *Config
}

func (l *Linter) Lint(d *dotgithub.DotGithub) (uint8, error) {
	if l.Config == nil {
		panic("Config cannot be nil")
	}
	if d == nil {
		panic("DotGithub cannot be empty")
	}

	summary := newSummary()
	numCPU := runtime.NumCPU()

	chJobs := make(chan Job)
	chWarnings := make(chan string)
	chErrors := make(chan string)

	wg := sync.WaitGroup{}
	wg.Add(numCPU)

	go func() {
		for _, action := range d.Actions {
			for ruleIdx, ruleEntry := range l.Config.Rules {
				if ruleEntry.FileType() & rule.DotGithubFileTypeAction == 0 {
					continue
				}
				isError := l.Config.IsError(ruleEntry.ConfigName(rule.DotGithubFileTypeAction))
				chJobs <- Job{
					rule:      ruleEntry,
					file:      action,
					dotGithub: d,
					isError:   isError,
					value:     l.Config.Values[ruleIdx],
				}
				summary.numJob.Add(1)
			}
		}

		for _, workflow := range d.Workflows {
			for ruleIdx, ruleEntry := range l.Config.Rules {
				if ruleEntry.FileType() & rule.DotGithubFileTypeWorkflow == 0 {
					continue
				}
				isError := l.Config.IsError(ruleEntry.ConfigName(rule.DotGithubFileTypeWorkflow))
				chJobs <- Job{
					rule:      ruleEntry,
					file:      workflow,
					dotGithub: d,
					isError:   isError,
					value:     l.Config.Values[ruleIdx],
				}
				summary.numJob.Add(1)
			}
		}

		close(chJobs)
		wg.Done()
	}()

	go func() {
		for {
			job, more := <-chJobs
			if more {
				compliant, err := job.Run(chWarnings, chErrors)
				if err != nil {
					slog.Error(fmt.Sprintf("%s\n", err.Error()))
					summary.numError.Add(1)
					continue
				}
				if !compliant {
					if job.isError {
						summary.numError.Add(1)
					} else {
						summary.numWarning.Add(1)
					}
				}
				summary.numProcessed.Add(1)
				continue
			}

			close(chWarnings)
			close(chErrors)

			wg.Done()
			return
		}
	}()

	for range numCPU - 2 {
		go func() {
			chWarningsClosed := false
			chErrorsClosed := false

			ticker := time.NewTicker(500 * time.Millisecond)

			for {
				select {
				case s, more := <-chWarnings:
					if more {
						if s != "" {
							slog.Warn(s)
						}
					} else {
						chWarningsClosed = true
					}
				case s, more := <-chErrors:
					if more {
						if s != "" {
							slog.Error(s)
						}
					} else {
						chErrorsClosed = true
					}
				case <-ticker.C:
					if chWarningsClosed && chErrorsClosed {
						wg.Done()
						return
					}
				}
			}
		}()
	}

	wg.Wait()

	finalStatus := HasNoErrorsOrWarnings

	if summary.numError.Load() > 0 {
		finalStatus = HasErrors
	} else {
		if summary.numWarning.Load() > 0 {
			finalStatus = HasOnlyWarnings
		}
	}

	slog.Debug(fmt.Sprintf("number of rules returning errors: %d", summary.numError.Load()))
	slog.Debug(fmt.Sprintf("number of rules returning warnings: %d", summary.numWarning.Load()))
	slog.Debug(fmt.Sprintf("number of rules processed in total: %d", summary.numProcessed.Load()))

	return uint8(finalStatus), nil
}
