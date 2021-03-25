# jobber
Run a list of commands in a queue from a text file or stdin using all cores.

## Building
**requires [golang](https://golang.org/) to build!**
1. Clone the repository with `git clone https://github.com/zalefin/jobber`
2. Change directory into the newly cloned repository with `cd jobber`
3. Build jobber with `go build`

## Usage
`./jobber [TARGETS PATH]`

where a targets file looks something like

```
program1 arg1 arg2
program2 arg1 arg2
program3 arg1 arg2
program4 arg1 arg2
```

additionally, jobber can take a list from stdin and do things like `cat targets.txt | ./jobber -`

