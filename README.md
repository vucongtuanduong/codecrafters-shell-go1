# Shell implementation in Go
A bash-inspired shell built from scratch in Go with the help of CodeCrafters
---


## Features
- [x] **Builtin commands:**
  - [x] `exit` — exit the shell
  - [x] `echo` — print the arguments
  - [x] `type` — show whether a command is a builtin or where an executable lives
  - [x] `pwd`
  - [x] `cd` - change directory (absolute path, relative path, home directory)
- [x] **Execute external programs** available in the `PATH`
- [x] Navigation
  - [x] pwd builtin
  - [x] cd builtin
- [ ] Quoting
  - [x] Single quotes
  - [x] Double quotes
- [ ] Redirection
- [ ] Autocompletion
- [ ] Pipelines
---

## Quick start 

Clone/use your private workspace and then:

- Run:
```bash
./your_program.sh 
```

---

## Usage examples 🧪

Example REPL session:
```
$ echo Hello, Codecrafters!
Hello, Codecrafters!
$ type echo
echo is a shell builtin
$ type ls
ls is /bin/ls
$ ls -la
# (output from /bin/ls)
$ exit
```

---

