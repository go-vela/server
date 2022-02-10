# Contributing

### We'd love to accept your contributions to this project! If you are a first time contributor, please review our [Contributing Guidelines]() before proceeding.


## Bugs

Bug reports should be opened up as [issues](https://help.github.com/en/github/managing-your-work-on-github/about-issues) on the [go-vela/community](https://github.com/go-vela/community) repository!

## Feature Requests

Feature Requests should be opened up as [issues](https://help.github.com/en/github/managing-your-work-on-github/about-issues) on the [go-vela/community](https://github.com/go-vela/community) repository!

## Pull Requests

**NOTE: We recommend you start by opening a new issue describing the bug or feature you're intending to fix. Even if you think it's relatively minor, it's helpful to know what people are working on.**

We are always open to new PRs! You can follow the below guide for learning how you can contribute to the project!

## Getting Started

### Prerequisites

* [Review the commit guide we follow](https://chris.beams.io/posts/git-commit/#seven-rules) - ensure your commits follow our standards
* [Review the local development docs](../DOCS.md) - ensures you have the Vela application stack running locally

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
  * Be sure to follow our [style guide]() 

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

* Make sure to follow our [PR process]() when opening a pull request

Thank you for your contribution!
