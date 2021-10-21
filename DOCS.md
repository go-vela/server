# Documentation

This document intends to provide information on how to get the Vela application running locally.

For more information, please see our [installation docs](https://go-vela.github.io/docs/install/).

## Prerequisites

This section covers the dependencies required to get the Vela application running locally.

* [Docker](https://docs.docker.com/install/) - building block for local development
* [Docker Compose](https://docs.docker.com/compose/install/) - start up local development
* [Github OAuth Client](https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/) - building block for local development
* [Golang](https://golang.org/dl/) - for source code and [dependency management](https://github.com/golang/go/wiki/Modules)
* [Make](https://www.gnu.org/software/make/) - start up local development

## Setup

**NOTE: Please review the [prerequisites section](#prerequisites) before moving forward.**

This section covers the configuration required to get the Vela application running locally.

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

* If using GitHub Enterprise (default: `https://github.com/`), add the Web URL to a local `.env` file:

```bash
# add Github Enterprise Web URL to local `.env` file for `docker-compose`
echo "VELA_SOURCE_ADDR=<GitHub Enterprise Web URL>" >> .env
```

* Create an [OAuth App](https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/) and obtain secrets for local development:
  * `Application name` = `Vela - local` (name of the OAuth application shouldn't matter)
  * `Homepage URL` = `http://localhost:8080` (base URL of the server)
  * `Authorization callback URL` = `http://localhost:8080/authenticate` (authenticate endpoint of the base URL of the server)

* Add OAuth client secrets to a local `.env` file:

```bash
# add Github Client ID to local `.env` file for `docker-compose`
echo "VELA_SOURCE_CLIENT=<Github OAuth Client ID>" >> .env

# add Github Client Secret to local `.env` file for `docker-compose`
echo "VELA_SOURCE_SECRET=<Github OAuth Client Secret>" >> .env
```

## Start

**NOTE: Please review the [setup section](#setup) before moving forward.**

This section covers the commands required to get the Vela application running locally.

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

* Navigate to the web UI:

```bash
# opens a browser window to the page http://localhost:8888
open http://localhost:8888
```

## Repo

**NOTE: Please review the [start section](#start) before moving forward.**

In order to run a build in Vela, you'll need to add a repo to the locally running system.

<details><summary>click to reveal content</summary>
<p>

1. Navigate to the `Source Repositories` page @ http://localhost:8888/account/source-repos
  * For conveinence, you can reference our documentation to [learn how to enable a repo](https://go-vela.github.io/docs/usage/enable_repo/).

2. Click the blue drop down arrow on the left side next to the org that contains the repo you want to enable.

3. Find the repo you want to enable in the drop down list and click the blue `Enable` button on the right side.
  * You should received a `success` message telling you `<org>/<repo> enabled.`

4. Click the blue `View` button to navigate directly to the repo.
  * You should be redirected to http://localhost:8888/<org>/<repo>

</p>
</details>

## Pipeline

**NOTE: Please review the [repo section](#repo) before moving forward.**

In order to run a build in Vela, you'll need to add a pipeline to the repo that has been added to the locally running system.

<details><summary>click to reveal content</summary>
<p>

1. Create a Vela [pipeline](https://go-vela.github.io/docs/tour/) to define a workflow for Vela to run.
  * For convenience, you can reference our documentation to use [one of our example pipelines](https://go-vela.github.io/docs/usage/examples/).

2. Add the pipeline to the repo that was enabled above.

</p>
</details>

## Build

**NOTE: Please review the [pipeline section](#pipeline) before moving forward.**

In order to run a build in Vela, you'll need to capture a valid webhook payload to mock sending it to the locally running system.

<details><summary>click to reveal content</summary>
<p>

1. Review GitHub's [documentation on webhooks](https://developer.github.com/webhooks/)

2. Find the [recent delivery](https://developer.github.com/webhooks/testing/#listing-recent-deliveries) for the pipeline that was added to your repo.

3. Create a request locally for http://localhost:8080/webhook and replicate all parts from the recent delivery.
  * You should use whatever tool feels most comfortable and natural to you (`curl`, `Postman`, `Insomnia` etc.).
  * You should replicate all the request headers and the request body from the recent delivery.

4. Send the request and navigate directly to the repo (http://localhost:8888/<org>/<repo>) to watch the build run live.

</p>
</details>

## Services

This section covers the different services in the stack when the Vela application is running locally.

<details><summary>click to reveal content</summary>
<p>

### Server

The `server` Docker compose service hosts the Vela server and API.

This component is used for processing web requests and managing resources in the database and publishing builds to the FIFO queue.

### Worker

The `worker` Docker compose service hosts the Vela build daemon.

This component is used for pulling builds from the FIFO queue and executing them based off their configuration.

### UI

The `ui` Docker compose service hosts the Vela UI.

This component is used for providing a user-friendly interface for triggering actions in the Vela system.

### Redis

The `redis` Docker compose service hosts the Redis database.

This component is used for publishing builds to a FIFO queue.

For more information, please review [the official documentation](https://redis.io/).

### Postgres

The `postgres` Docker compose service hosts the Postgresql database.

This component is used for storing data at rest.

For more information, please review [the official documentation](https://www.postgresql.org/).

### Vault

The `vault` Docker compose service hosts the HashiCorp Vault instance.

This component is used for storing sensitive data like secrets.

For more information, please review [the official documentation](https://www.vaultproject.io/).

</p>
</details>
