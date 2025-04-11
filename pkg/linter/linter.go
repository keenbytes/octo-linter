package linter

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/loglevel"
)

const (
	HasNoErrorsOrWarnings = iota
	HasErrors
	HasOnlyWarnings
)

type Linter struct {
	Config   *Config
	LogLevel int
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
				_, isError := l.Config.Errors[rule.GetConfigName()]
				chJobs <- Job{
					rule: rule,
					file: action,
					dotGithub: d,
					isError: isError,
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
				if l.LogLevel == loglevel.LogLevelDebug || l.LogLevel == loglevel.LogLevelErrorsAndWarnings {
					fmt.Fprintf(os.Stdout, "wrn: %s\n", s)
				}
			case s := <-chErrors:
				if l.LogLevel != loglevel.LogLevelNone {
					fmt.Fprintf(os.Stderr, "err: %s\n", s)
				}
			case <-chDoneProcessing:
				wg.Done()
			}
		}
	}()

	for range numCPU-2 {
		go func() {
			for {
				select {
				case job := <- chJobs:
					ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10 * time.Second))
					compliant, err := job.Run(ctx, chWarnings, chErrors)
					cancel()
					if err != nil {
						if l.LogLevel != loglevel.LogLevelNone {
							fmt.Fprintf(os.Stderr, "!!!: %s\n", err.Error())
						}
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
				case <- chDoneProcessing:
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

	if l.LogLevel == loglevel.LogLevelDebug {
		fmt.Fprintf(os.Stderr, "dbg: number of errors: %d\n", summary.numError.Load())
		fmt.Fprintf(os.Stderr, "dbg: number of warnings: %d\n", summary.numWarning.Load())
		fmt.Fprintf(os.Stderr, "dbg: number of processed: %d\n", summary.numProcessed.Load())
	}

	return uint8(finalStatus), nil
}
