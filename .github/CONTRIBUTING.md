# Contributing

## Getting Started

We'd love to accept your contributions to this project! If you are a first time contributor, please review our [Contributing Guidelines](https://go-vela.github.io/docs/community/contributing_guidelines/) before proceeding.

### Prerequisites

* [Review the local development docs](../DOCS.md) - ensures you have the Vela application stack running locally
* [Review the commit guide we follow](https://chris.beams.io/posts/git-commit/#seven-rules) - ensure your commits follow our standards
* Review our [style guide](https://go-vela.github.io/docs/community/contributing_guidelines/#style-guide) to ensure your code is clean and consistent.

### Setup

* [Fork](/fork) this repository

* Clone this repository to your workstation:

```bash
# clone the project
git clone git@github.com:go-vela/server.git $HOME/go-vela/server
```

* Navigate to the repository code:

```bash
# change into the cloned project directory
cd $HOME/go-vela/server
```

* Point the original code at your fork:

```bash
# add a remote branch pointing to your fork
git remote add fork https://github.com/your_fork/server
```

### Development

**Please review the [local development documentation](../DOCS.md) for more information.**

* Navigate to the repository code:

```bash
# change into the cloned project directory
cd $HOME/go-vela/server
```

* Write your code and tests to implement the changes you desire.

* Run the repository code (ensures your changes perform as you desire):

```bash
# execute the `up` target with `make`
make up
```

* Test the repository code (ensures your changes don't break existing functionality):

```bash
# execute the `test` target with `make`
make test
```

* Clean the repository code (ensures your code meets the project standards):

```bash
# execute the `clean` target with `make`
make clean
```

* Push to your fork:

```bash
# push your code up to your fork
git push fork master
```

* Make sure to follow our [PR process](https://go-vela.github.io/docs/community/contributing_guidelines/#development-workflow) when opening a pull request

Thank you for your contribution!
