# Documentation

This document intends to provide information on how to get the Vela application running locally.

For more information, please see our [installation docs](https://go-vela.github.io/docs/install/).

## Prerequisites

* [Docker](https://docs.docker.com/install/) - building block for local development
* [Docker Compose](https://docs.docker.com/compose/install/) - start up local development
* [Github OAuth Client](https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/) - building block for local Development
* [Golang](https://golang.org/dl/) - for source code and [dependency management](https://github.com/golang/go/wiki/Modules)
* [Make](https://www.gnu.org/software/make/) - start up local development

## Setup

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

* Create an [OAuth App](https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/) and obtain secrets for local development:
  * `Application name` = `Vela - local` (name of the OAuth application shouldn't matter)
  * `Homepage URL` = `http://localhost:8080` (base URL of the server)
  * `Authorization callback URL` = `http://localhost:8888/account/authenticate` (authenticate endpoint of the base URL of the UI)

**NOTE: This will work for GitHub or GitHub Enterprise.**

* Add OAuth client secrets to a local `secrets.env` file:

```bash
# add Github Client ID to local secrets file for `docker-compose`
echo "VELA_SOURCE_CLIENT=<Github OAuth Client ID>" >> secrets.env

# add Github Client Secret to local secrets file for `docker-compose`
echo "VELA_SOURCE_SECRET=<Github OAuth Client Secret>" >> secrets.env
```

## Start

**NOTE: Please review the [setup section](#setup) before starting.

* Navigate to the repository code:

```bash
# change into the cloned project directory
cd $HOME/go-vela/server
```

* Run the repository code:

```bash
# execute the `up` target with `make`
make up
```

* Navigate to the Web UI @ http://localhost:8888
