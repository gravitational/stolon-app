# stolon-app

Opinionated [stolon](https://github.com/sorintlab/stolon) gravity app.

## Provided resources

Once deployed or installed, this app will provide:

 * HA Postgres via stolon proxy servers, as `svc/stolon-postgres`

## Gravity app

You can import as a gravity app, by running the following:

```
gravity app import --vendor --state-dir=/var/lib/gravity . gravitational.io/stolon-app:0.0.1
```

## Development

There are several development Makefile targets to simplify your workflow:

 * `dev-push` push images to `apiserver:5000`
 * `dev-deploy` deploy the bootstrap with `kubectl`
 * `dev-clean` destroy all cluster resources
 * `dev-redeploy` clean and then deploy the cluster

