package linter

import (
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"sync"

	"github.com/keenbytes/octo-linter/pkg/dotgithub"
)

const (
	HasNoErrorsOrWarnings = iota
	HasErrors
	HasOnlyWarnings
)

type Linter struct {
	Config   *Config
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
	chDoneProcessing := make(chan struct{})

	defer close(chDoneProcessing)
	defer close(chErrors)
	defer close(chWarnings)
	defer close(chJobs)

	wg := sync.WaitGroup{}
	wg.Add(numCPU)

	go func() {
		for _, action := range d.Actions {
			for _, rule := range l.Config.Rules {
				if !strings.HasPrefix(rule.GetConfigName(), "action_") {
					continue
				}
				isError := l.Config.IsError(rule.GetConfigName())
				chJobs <- Job{
					rule:      rule,
					file:      action,
					dotGithub: d,
					isError:   isError,
				}
				summary.numJob.Add(1)
			}
		}

		for _, workflow := range d.Workflows {
			for _, rule := range l.Config.Rules {
				if !strings.HasPrefix(rule.GetConfigName(), "workflow_") {
					continue
				}
				isError := l.Config.IsError(rule.GetConfigName())
				chJobs <- Job{
					rule:      rule,
					file:      workflow,
					dotGithub: d,
					isError:   isError,
				}
				summary.numJob.Add(1)
			}
		}

		for {
			if summary.numJob.Load() == summary.numProcessed.Load() {
				for range numCPU {
					chDoneProcessing <- struct{}{}
				}
				break
			}
		}
	}()

	go func() {
		for {
			select {
			case s := <-chWarnings:
				slog.Warn(s)
			case s := <-chErrors:
				slog.Error(s)
			case <-chDoneProcessing:
				wg.Done()
			}
		}
	}()

	for range numCPU - 2 {
		go func() {
			for {
				select {
				case job := <-chJobs:
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
				case <-chDoneProcessing:
					wg.Done()
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
