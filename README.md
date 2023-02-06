# Permalinker

Generate permalink URLs (Github, Gitlab, etc.) from a local file path 

## Install

```sh
# Install by building from source
➜ go install github.com/denis-engcom/permalinker@latest

# Run without arguments or with `-h` or with `--help` to print usage
➜ permalinker
```

## Example usage

```sh
# Navigate to a git repo
# ➜ cd <some-git-repo>
➜ pwd
/Users/denis/code/github.com/denis-engcom/permalinker

➜ git log -1 --oneline
4bc40a1 (HEAD -> main, origin/main) Initial commit

➜ ls
README.md go.mod    go.sum    main.go   notes.md

➜ permalinker main.go 20
https://github.com/denis-engcom/permalinker/blob/4bc40a14289d35d6d4ffed99dd1ef27559146b05/main.go#L20
```

## Next feature(s) to implement

* If we're behind origin, abort.
    * Perform a fetch, get current branch, check whether we're behind branch's origin.
	* Turn off fetch/check using `--no-fetch`
