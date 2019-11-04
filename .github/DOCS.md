# Documentation

For installation and usage, please [visit our docs](https://go-vela.github.io/docs).

## Services

If you followed [the instructions from the contributing guide](CONTRIBUTING.md/#getting-started), you should have 3 services running as Docker containers on your machine:

* vela
* redis
* postgres

### Vela

The `vela` service is running the actual Vela server and API. The [docker-compose](../docker-compose.yml) file is already setup to connect will the other services (`redis` and `postgres`) as well as the OAuth app you created from the [getting started section](CONTRIBUTING.md/#getting-started).

### Redis

The `redis` service hosts the [Redis](https://redis.io/) database used for Vela's queue implementation.

### Postgres

The `postgres` service hosts the [Postgresql](https://www.postgresql.org/) database used for storing Vela's data at rest.

## API

With the `vela` service running, you can login to Vela @ http://localhost:8080/login. This should have you authorize the application and in return provide a token via URL query parameter.

After you've logged in to Vela, you can now integrate with the Vela API. Currently the Vela API supports two methods of authentication:

* `Authorization` Header (Bearer Token)
  * curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/v1/users
* `access_token` Query Parameter
  * curl http://localhost:8080/api/v1/users?access_token=<TOKEN>

## CLI

Coming soon!

## Executing Builds

In order to execute builds on your local machine, you'll also need to create a Vela worker to process the workloads pushed to the `redis` queue.

To create a worker, you can follow the  [documentation](https://github.com/go-vela/worker/blob/master/.github/DOCS.md) found in the [go-vela/worker](https://github.com/go-vela/worker) repository.
