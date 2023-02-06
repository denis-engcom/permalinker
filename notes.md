permalinker <path>[:<line-number>] ...
Composes a permalink URL using the repository, current commit, file path, and line number. Provide at least one file, each with (optional) colon and line number

Options:
  --rev          Set URL commit to the following revision. [Default: HEAD] (the current commit in the local repo)
  --remote       Set the repo config remove verified to derive the repo HTTPS URL. [Default: origin]
  --url          Set the repo URL directly (ignores --remote)
  -v, --verbose  Print permalinker debug logs
  -h, --help     Print usage text

Idea:

File searching can be handled outside via `fzf` or other command line tools.
Line number searching can be handled outside via `cat -n`, `less -n`, `ag`, etc.

For each file...

Using go-git, use https://pkg.go.dev/github.com/go-git/go-git/v5#PlainOpenWithOptions with given file path and DetectDotGit = true
Then, get abs path of path provided: `filepath.Abs(...)`
Get abs path of repo https://pkg.go.dev/github.com/go-git/go-git/v5@v5.5.2/config#Config: `repo.Config().Core.Worktree`
Subtract abs repo path from file path using https://pkg.go.dev/path/filepath#Rel `filepath.Rel(repo path, abs file path)`
Get URL of repo via `repo.Config().Core.Remotes["origin"]`
* Use https://pkg.go.dev/github.com/go-git/go-git/v5@v5.5.2/plumbing/transport `NewEndpoint(...), ep.Protocol = "https", ep.String()`
Get the latest commit from https://pkg.go.dev/github.com/go-git/go-git/v5@v5.5.2#Repository.Log, `repo.Log(), iter.Next(), iter.Close()`
Print [repo https url]/blob/[commit][relative file path][#L123 (line-number)]

---

Future plan: add ability to provide ref (instead of using current commit), https://pkg.go.dev/github.com/go-git/go-git/v5#Repository.ResolveRevision, then use that hash

(maybe?)
Error if file is not unmodified: https://pkg.go.dev/github.com/go-git/go-git/v5#Status.File: `repo.Worktree().Status()[path].Worktree and .Staging == Unmodified`
* "Error: Cannot guarantee accuracy. File is not unmodified. Please unstage and revert the file related to the current commit."
