# Permalinker

```sh
# Install using the go runtime
go install github.com/denis-engcom/permalinker@latest

# Run without arguments to print usage
permalinker
```

## Next feature(s) to implement

* Perform a fetch, get current branch, check whether we're behind branch's origin.
	* If we're behind origin, abort 
	* Turn off fetch/check using `--no-fetch`
