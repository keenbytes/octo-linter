// Package dotgithub reads the contents of a .github directory, parsing all actions and workflows into structured data.
package dotgithub

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/keenbytes/octo-linter/v2/pkg/action"
	"github.com/keenbytes/octo-linter/v2/pkg/workflow"
)

// DotGithub represents contents of .github directory.
type DotGithub struct {
	Actions         map[string]*action.Action
	ExternalActions map[string]*action.Action
	Workflows       map[string]*workflow.Workflow
	Vars            map[string]bool
	Secrets         map[string]bool
}

// ReadDir scans the given directory and parses all GitHub Actions workflow and action YAML files into the struct.
func (d *DotGithub) ReadDir(p string) error {
	d.Actions = make(map[string]*action.Action)
	d.Workflows = make(map[string]*workflow.Workflow)

	err := d.getActionsFromDir(p)
	if err != nil {
		return fmt.Errorf("error getting actions from dir %s: %w", p, err)
	}

	err = d.getWorkflowsFromDir(p)
	if err != nil {
		return fmt.Errorf("error getting workflows from dir %s: %w", p, err)
	}

	// download all external actions used in actions' steps
	reExternal := regexp.MustCompile(`[a-zA-Z0-9\-\_]+\/[a-zA-Z0-9\-\_]+(\/[a-zA-Z0-9\-\_]){0,1}@[a-zA-Z0-9\.\-\_]+`)

	for _, a := range d.Actions {
		err := a.Unmarshal(false)
		if err != nil {
			return fmt.Errorf("error unmarshaling action: %w", err)
		}

		if a.Runs == nil || len(a.Runs.Steps) == 0 {
			continue
		}

		for i, step := range a.Runs.Steps {
			if !reExternal.MatchString(step.Uses) {
				continue
			}

			err := d.DownloadExternalAction(step.Uses)
			if err != nil {
				slog.Error(
					"error downloading external action",
					slog.String("action", a.DirName),
					slog.Int("step", i),
					slog.String("uses", step.Uses),
					slog.String("err", err.Error()),
				)
			}
		}
	}

	for _, w := range d.Workflows {
		err := w.Unmarshal(false)
		if err != nil {
			return fmt.Errorf("error unmarshaling workflow: %w", err)
		}

		for _, w := range d.Workflows {
			err := w.Unmarshal(false)
			if err != nil {
				return fmt.Errorf("error unmarshaling workflow: %w", err)
			}

			if len(w.Jobs) == 0 {
				continue
			}

			for _, job := range w.Jobs {
				if len(job.Steps) == 0 {
					continue
				}

				for i, step := range job.Steps {
					if !reExternal.MatchString(step.Uses) {
						continue
					}

					err := d.DownloadExternalAction(step.Uses)
					if err != nil {
						slog.Error(
							"error downloading external action",
							slog.String("workflow", w.FileName),
							slog.Int("step", i),
							slog.String("uses", step.Uses),
							slog.String("err", err.Error()),
						)
					}
				}
			}
		}
	}

	return nil
}

// ReadVars reads a file with GitHub Actions variables, parsing each line into the struct as a variable.
func (d *DotGithub) ReadVars(path string) error {
	if path == "" {
		return nil
	}

	d.Vars = make(map[string]bool)

	slog.Debug(
		"reading file with list of possible variable names...",
		slog.String("path", path),
	)

	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading vars file %s: %w", path, err)
	}

	l := strings.Fields(string(b))
	for _, v := range l {
		d.Vars[v] = true
	}

	return nil
}

// ReadSecrets reads a file with GitHub Actions secrets, parsing each line into the struct as a secret.
func (d *DotGithub) ReadSecrets(path string) error {
	if path == "" {
		return nil
	}

	d.Secrets = make(map[string]bool)

	slog.Debug(
		"reading file with list of possible secret names...",
		slog.String("path", path),
	)

	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading secrets file %s: %w", path, err)
	}

	l := strings.Fields(string(b))
	for _, s := range l {
		d.Secrets[s] = true
	}

	return nil
}

// GetAction returns an Action by its name.
func (d *DotGithub) GetAction(n string) *action.Action {
	return d.Actions[n]
}

// GetExternalAction returns an Action that is defined outside the current repository, by name.
func (d *DotGithub) GetExternalAction(n string) *action.Action {
	if d.ExternalActions == nil {
		d.ExternalActions = map[string]*action.Action{}
	}

	return d.ExternalActions[n]
}

// DownloadExternalAction downloads a GitHub Action from its “uses” path (e.g., "actions/checkout@v4").
func (d *DotGithub) DownloadExternalAction(path string) error {
	if d.ExternalActions == nil {
		d.ExternalActions = map[string]*action.Action{}
	}

	if d.ExternalActions[path] != nil {
		return nil
	}

	repoVersion := strings.Split(path, "@")
	ownerRepoDir := strings.SplitN(repoVersion[0], "/", 3)

	directory := ""
	if len(ownerRepoDir) > 2 {
		directory = "/" + ownerRepoDir[2]
	}

	actionURLPrefix := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", ownerRepoDir[0], ownerRepoDir[1], repoVersion[1])

	urlYML := actionURLPrefix + directory + "/action.yml"
	slog.Debug(
		"downloading external action yaml",
		slog.String("url", urlYML),
	)

	req, err := http.NewRequest("GET", urlYML, strings.NewReader(""))
	if err != nil {
		return fmt.Errorf("error creating http request for action yml: %w", err)
	}

	c := &http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("error doing http request for action yml: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		urlYAML := actionURLPrefix + directory + "/action.yaml"
		slog.Debug(
			"downloading external action yaml",
			slog.String("url", urlYAML),
		)

		req, err = http.NewRequest("GET", urlYAML, strings.NewReader(""))
		if err != nil {
			return fmt.Errorf("error creating http request for action yaml: %w", err)
		}

		resp, err = c.Do(req)
		if err != nil {
			return fmt.Errorf("error doing http request for action yaml: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil
		}
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)

	d.ExternalActions[path] = &action.Action{
		Path:    path,
		DirName: "",
		Raw:     b,
	}

	err = d.ExternalActions[path].Unmarshal(true)
	if err != nil {
		return fmt.Errorf("error unmarshaling external action: %w", err)
	}

	return nil
}

func (d *DotGithub) getActionsFromDir(p string) error {
	dirActions := filepath.Join(p, "actions")

	entries, err := os.ReadDir(dirActions)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error reading actions directory: %w", err)
		}
	}

	for _, e := range entries {
		dirAction := filepath.Join(dirActions, e.Name())

		// only directories
		fileInfo, err := os.Stat(dirAction)
		if err != nil {
			return fmt.Errorf("error getting os.Stat on %s: %w", dirAction, err)
		}

		if !fileInfo.IsDir() {
			continue
		}

		// search for action.yml or action.yaml file
		ymlAction := filepath.Join(dirAction, "action.yml")
		_, err = os.Stat(ymlAction)

		ymlNotFound := os.IsNotExist(err)
		if err != nil && !ymlNotFound {
			return fmt.Errorf("error getting os.Stat on %s: %w", ymlAction, err)
		}

		if ymlNotFound {
			yamlAction := filepath.Join(dirAction, "action.yaml")
			_, err = os.Stat(yamlAction)

			yamlNotFound := os.IsNotExist(err)
			if err != nil && !yamlNotFound {
				return fmt.Errorf("error getting os.Stat on %s: %w", yamlAction, err)
			}

			if !yamlNotFound {
				ymlAction = yamlAction
			} else {
				continue
			}
		}

		d.Actions[e.Name()] = &action.Action{
			Path:    ymlAction,
			DirName: e.Name(),
		}
	}

	return nil
}

func (d *DotGithub) getWorkflowsFromDir(p string) error {
	dirWorkflows := filepath.Join(p, "workflows")

	entries, err := os.ReadDir(dirWorkflows)
	if err != nil {
		return fmt.Errorf("error reading workflows directory %s: %w", dirWorkflows, err)
	}

	nameRegex := regexp.MustCompile(`\.y[a]{0,1}ml$`)
	for _, e := range entries {
		m := nameRegex.MatchString(e.Name())
		if !m {
			continue
		}

		ymlWorkflow := filepath.Join(dirWorkflows, e.Name())

		fileInfo, err := os.Stat(ymlWorkflow)
		if err != nil {
			return fmt.Errorf("error getting os.Stat on %s: %w", ymlWorkflow, err)
		}

		if !fileInfo.Mode().IsRegular() {
			continue
		}

		d.Workflows[e.Name()] = &workflow.Workflow{
			Path: ymlWorkflow,
		}
	}

	return nil
}

// IsVarExist checks whether the variable has been loaded from the variables file.
func (d *DotGithub) IsVarExist(n string) bool {
	_, ok := d.Vars[n]
	return ok
}

// IsSecretExist checks whether the secret has been loaded from the secrets file.
func (d *DotGithub) IsSecretExist(n string) bool {
	_, ok := d.Secrets[n]
	return ok
}
