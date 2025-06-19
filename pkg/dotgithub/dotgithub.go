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

	"github.com/keenbytes/octo-linter/pkg/action"
	"github.com/keenbytes/octo-linter/pkg/workflow"
)

type DotGithub struct {
	Actions         map[string]*action.Action
	ExternalActions map[string]*action.Action
	Workflows       map[string]*workflow.Workflow
	Vars            map[string]bool
	Secrets         map[string]bool
}

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
			return err
		}
		if a.Runs != nil && a.Runs.Steps != nil && len(a.Runs.Steps) > 0 {
			for i, step := range a.Runs.Steps {
				if reExternal.MatchString(step.Uses) {
					err := d.DownloadExternalAction(step.Uses)
					if err != nil {
						slog.Error(fmt.Sprintf("action '%s' step %d: error downloading external action '%s': %s", a.DirName, i, step.Uses, err.Error()))
					}
				}
			}
		}
	}
	for _, w := range d.Workflows {
		err := w.Unmarshal(false)
		if err != nil {
			return err
		}
		for _, w := range d.Workflows {
			err := w.Unmarshal(false)
			if err != nil {
				return err
			}
			if w.Jobs != nil && len(w.Jobs) > 0 {
				for _, job := range w.Jobs {
					if job.Steps == nil || len(job.Steps) == 0 {
						continue
					}
					for i, step := range job.Steps {
						if reExternal.MatchString(step.Uses) {
							err := d.DownloadExternalAction(step.Uses)
							if err != nil {
								slog.Error(fmt.Sprintf("workflow '%s' step %d: error downloading external action '%s': %s", w.FileName, i, step.Uses, err.Error()))
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (d *DotGithub) ReadVars(path string) error {
	if path == "" {
		return nil
	}

	d.Vars = make(map[string]bool)

	slog.Debug(fmt.Sprintf("reading file with list of possible variable names %s ...", path))

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

func (d *DotGithub) ReadSecrets(path string) error {
	if path == "" {
		return nil
	}

	d.Secrets = make(map[string]bool)

	slog.Debug(fmt.Sprintf("reading file with list of possible secret names %s ...", path))

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

func (d *DotGithub) GetAction(n string) *action.Action {
	return d.Actions[n]
}

func (d *DotGithub) GetExternalAction(n string) *action.Action {
	if d.ExternalActions == nil {
		d.ExternalActions = map[string]*action.Action{}
	}
	return d.ExternalActions[n]
}

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

	slog.Debug(fmt.Sprintf("downloading %s ...", actionURLPrefix+directory+"/action.yml"))

	req, err := http.NewRequest("GET", actionURLPrefix+directory+"/action.yml", strings.NewReader(""))
	if err != nil {
		return err
	}
	c := &http.Client{}
	resp, err := c.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		slog.Debug(fmt.Sprintf("downloading %s ...", actionURLPrefix+directory+"/action.yaml"))

		req, err = http.NewRequest("GET", actionURLPrefix+directory+"/action.yaml", strings.NewReader(""))
		if err != nil {
			return err
		}
		resp, err = c.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
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
		return err
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

func (d *DotGithub) IsVarExist(n string) bool {
	_, ok := d.Vars[n]
	return ok
}

func (d *DotGithub) IsSecretExist(n string) bool {
	_, ok := d.Secrets[n]
	return ok
}
