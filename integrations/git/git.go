package git

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-cmd/cmd"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/vortex14/gotyphoon/environment"
	"github.com/vortex14/gotyphoon/interfaces"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Git struct {
	settings *environment.Settings
	repo     *git.Repository
	Path     string
	workTree *git.Worktree
	Project  interfaces.Project
}

var fileStatusMapping = map[git.StatusCode]string{
	git.Unmodified:         "",
	git.Untracked:          "Untracked",
	git.Modified:           "Modified",
	git.Added:              "Added",
	git.Deleted:            "Deleted",
	git.Renamed:            "Renamed",
	git.Copied:             "Copied",
	git.UpdatedButUnmerged: "Updated",
}

func (g *Git) TestGetDir(projectName string) error {
	if _, err := os.Stat(g.settings.Projects + "/" + projectName + "/.git"); os.IsNotExist(err) {
		return err
	}
	return nil
}

func (g *Git) Commit(message string, opt *git.CommitOptions) error {
	g.LoadRepo()

	commit, err := g.workTree.Commit(message, opt)
	obj, err := g.repo.CommitObject(commit)
	if err != nil {
		color.Red("%s", err.Error())
		os.Exit(1)
	}
	fmt.Println(obj)

	return err

}

func (g *Git) Push(remote string, branch string)  {
	g.LoadRepo()

	pushCommand := cmd.NewCmd("git", "push", remote, branch)
	pushCommand.Dir = g.Path
	s := <-pushCommand.Start()

	for _, l := range s.Stdout {
		color.Yellow("%s", l)
	}



	//TO DO: research AUTH method

	//err := g.repo.Push(&git.PushOptions{
	//	RemoteName: remote,
	//	Auth: nil,
	//})
	//if err != nil {
	//	color.Red("%s", err.Error())
	//	return
	//}
}

func (g *Git) RemovePyCacheFiles()  {
	//HOW TO RUN  THIS sample ???? ::::
	//find . -name "__pycache__" -exec rm -r "{}" \;
	//removeCommand := cmd.NewCmd("find", ".", "-name", "__pycache__", "-exec", "rm", "-r")
	//removeCommand.Dir = g.Path
	//s := <-removeCommand.Start()
	//for _, l := range s.Stdout {
	//	color.Yellow("%s", l)
	//}

	g.LoadRepo()
	pyignore := "__pycache__"
	_ = filepath.Walk(g.Path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				color.Red("%s", err.Error())
			}
			if strings.Contains(path, pyignore) {

				if _, err := os.Stat(path); err == nil {
					pyIgnoreDir := strings.Split(path, pyignore)[0] + pyignore

					color.Yellow("%s",pyIgnoreDir)
					err := os.RemoveAll(pyIgnoreDir)
					if err != nil {
						color.Red("%s", err.Error())
					}

				}

			}

			return nil
		})


}

func (g *Git) ShowPyCacheFiles()  {

	g.LoadRepo()
	removeCommand := cmd.NewCmd("find", ".", "-name", "__pycache__")
	removeCommand.Dir = g.Path
	s := <-removeCommand.Start()
	for _, l := range s.Stdout {
		color.Yellow("%s", l)
	}
}

func (g Git) AddAll()  {
	g.LoadRepo()
	//g.AddAndCommit()

}

func (g *Git) RemoveAllUnTrackingFiles()  {
	g.LoadRepo()
	cleanCommand := cmd.NewCmd("git", "clean", "-fd")
	cleanCommand.Dir = g.Path
	s := <-cleanCommand.Start()
	for _, l := range s.Stdout {
		color.Yellow("%s", l)
	}
}

func (g *Git) GetStoreBranchName() string {
	currTime := time.Now()
	name := fmt.Sprintf("backup-store-branch-%d-%d-%d", currTime.Day(), int(currTime.Month()), currTime.Year())
	return name
}

func (g *Git) SaveLocalChanging(){
	newBranchName := g.GetStoreBranchName()

	g.AddAndCommit("save backup changes")

	err := g.SwitchBranch(newBranchName)
	if err != nil {
		color.Red("%s", err.Error())
		return
	}
}

func (g *Git) CreateBranch(name string)  {
	headRef, err := g.repo.Head()
	if err != nil {
		color.Red("%s", err.Error())
	}
	ref := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/"+name), headRef.Hash())
	err = g.repo.Storer.SetReference(ref)
	if err != nil {
		color.Red("%s", err.Error())
	}
}

func (g *Git) SwitchBranch(branch string) error {
	refs := fmt.Sprintf("refs/heads/%s", branch)
	err := g.workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(refs),
		//Create: false,
		//Force:  false,
		//Keep:   false,
	})
	if err != nil {
		color.Red("%s", err.Error())

	}
	return err
}

func (g *Git) CreateBranchAndCommit(message string, branch string)  {
	g.LoadRepo()

	g.AddAndCommit(message)

	err := g.SwitchBranch(branch)

	if err != nil {
		g.CreateBranch(branch)
		_ = g.SwitchBranch(branch)
	}


}


func (g *Git) LocalResetLikeRemote(remote string, branch string, backup bool)  {
	g.LoadRepo()
	g.SwitchBranch(branch)

	if backup {
		g.SaveLocalChanging()
	}
	fetchCommand := cmd.NewCmd("git", "fetch", remote)
	fetchCommand.Dir = g.Path
	<-fetchCommand.Start()
	resetOpt := fmt.Sprintf("%s/%s", remote, branch)
	color.Yellow(resetOpt)
	resetCommand := cmd.NewCmd("git", "reset", "--hard" , resetOpt)
	resetCommand.Dir = g.Path
	out := <-resetCommand.Start()
	for _, l := range out.Stdout {
		color.Yellow("%s", l)
	}
}

func (g *Git) AddAndCommit(message string)  {
	g.LoadRepo()

	AddCommand := cmd.NewCmd("git", "add", ".")
	AddCommand.Dir = g.Path
	<-AddCommand.Start()
	err := g.Commit(message, &git.CommitOptions{

	})
	if err != nil {
		color.Red("%s", err.Error())
		return
	}
}

func (g *Git) LoadRepo() {
	if g.repo == nil {
		repo, err := git.PlainOpen(g.Path)
		if err != nil {
			color.Red("%s", err.Error())
			os.Exit(1)
		}
		g.repo = repo
		w, err := g.repo.Worktree()
		if err != nil {
			color.Red("%s", err.Error())
			os.Exit(1)
		}
		g.workTree = w

	}
}


func (g *Git) RepoStatus() () {
	g.LoadRepo()
	s, err := g.workTree.Status()
	if err != nil {
		color.Red("%s", err.Error())
		os.Exit(1)
	}
	if s != nil {
		for filename := range s {
			var untracked bool

			if s.IsUntracked(filename) {
				untracked = true
			}
			fileStatus := s.File(filename)
			if !untracked && fileStatus.Staging == git.Untracked && fileStatus.Worktree == git.Untracked {
				fileStatus.Staging = git.Unmodified
				fileStatus.Worktree = git.Unmodified
			}

			color.Yellow("%s -> %s", fileStatusMapping[fileStatus.Worktree], filename)

		}
	}
}