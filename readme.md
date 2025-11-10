simplified implementation of a version control system, inspired by git.

implemented features:

- [x] initialize a repository (`init`)
- [x] add files to staging area (`add`)
- [x] commit changes (`commit`)
- [x] view commit history (`log`)
- [x] check repository status (`status`)
- [x] .gitignore support (respects .gitignore patterns in `add` and `status` commands)

todo:

- [ ] implement diff functionality
- [ ] add branching support
  - [x] create branches
  - [x] switch between branches
  - [ ] merge branches
- [ ] improve `add` command to support multiple files and directories
- [ ] add `reset` command to unstage changes or move head
