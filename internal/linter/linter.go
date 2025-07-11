package linter

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/keenbytes/octo-linter/v2/internal/linter/glitch"
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

func (l *Linter) Lint(d *dotgithub.DotGithub, output string, outputLimit int) (int, error) {
	if l.Config == nil {
		panic("Config cannot be nil")
	}

	if d == nil {
		panic("DotGithub cannot be empty")
	}

	summary := newSummary()
	numCPU := runtime.NumCPU()

	chJobs := make(chan Job)
	chWarnings := make(chan glitch.Glitch)
	chErrors := make(chan glitch.Glitch)

	wg := sync.WaitGroup{}
	wg.Add(numCPU)

	go func() {
		for _, action := range d.Actions {
			for ruleIdx, ruleEntry := range l.Config.Rules {
				if ruleEntry.FileType()&rule.DotGithubFileTypeAction == 0 {
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
				if ruleEntry.FileType()&rule.DotGithubFileTypeWorkflow == 0 {
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
					slog.Error(
						"error running job",
						slog.String("err", err.Error()),
					)
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
				case glitchInstance, more := <-chWarnings:
					if more {
						slog.Warn(
							glitchInstance.ErrText,
							slog.String("path", glitchInstance.Path),
							slog.String("rule", glitchInstance.RuleName),
						)

						glitchInstance.IsError = false
						summary.addGlitch(&glitchInstance)
					} else {
						chWarningsClosed = true
					}
				case glitchInstance, more := <-chErrors:
					if more {
						slog.Error(
							glitchInstance.ErrText,
							slog.String("path", glitchInstance.Path),
							slog.String("rule", glitchInstance.RuleName),
						)

						glitchInstance.IsError = true
						summary.addGlitch(&glitchInstance)
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

	slog.Debug(
		"summary",
		slog.Int("rules_returning_errors", int(summary.numError.Load())),
		slog.Int("rules_processed", int(summary.numProcessed.Load())),
		slog.Int("glitches", len(summary.glitches)),
	)

	if output != "" {
		outputMd := filepath.Join(output, "output.md")
		slog.Debug(
			"writing markdown output",
			slog.String("path", outputMd),
		)

		_ = os.Remove(outputMd)

		if outputLimit < 0 {
			outputLimit = 0
		}

		md := summary.markdown("octo-linter summary", outputLimit)

		err := os.WriteFile(outputMd, []byte(md), 0644)
		if err != nil {
			return finalStatus, fmt.Errorf("error writing markdown output: %w", err)
		}
	}

	return finalStatus, nil
}
