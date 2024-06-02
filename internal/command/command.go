package command

import (
  "context"
  "os"
  "os/exec"
)

func InDir(dir string) *InDirCommand {
  return &InDirCommand{dir: dir}
}

type InDirCommand struct {
  dir string
}

func (i *InDirCommand) GoModTidy(ctx context.Context) error {
  cmd := exec.CommandContext(ctx, "go", "mod", "tidy")
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Dir = i.dir
  return cmd.Run()
}

func (i *InDirCommand) GoBuild(ctx context.Context) error {
  cmd := exec.CommandContext(ctx, "go", "build", "-o", "wombat")
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Dir = i.dir
  return cmd.Run()
}
