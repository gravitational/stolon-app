# stolon-app

[Stolon](https://github.com/sorintlab/stolon) packaged as gravity app.

## Provides

Kubernetes resouces in `resources` folder.

Once deployed or installed, this app will provide HA Postgres via stolon proxy servers, you can find endpoint using:

 ```sh
$  kubectl describe svc/stolon-postgres
 ```

## Based on

[This example](https://github.com/sorintlab/stolon/tree/master/examples/kubernetes) included with [stolon](https://github.com/sorintlab/stolon)


## Building

```sh
$ make
```

## Prerequisites

* docker >=  1.8
* golang 1.5.4
* GNU make
* kubectl >= 1.3

## Gravity app

You can import it in the OpsCenter as a gravity app, by running the following:

```sh
$ make reimport
```

### Creating site

```sh
$ gravity site create --app="gravitational.io/stolon-app:0.0.5"`
```

**Note**: you might want to deploy it on kubernetes manually but it's not recommended.

### Development

There are several development `Makefile` targets to simplify your workflow:

 * `dev-push` push images to `apiserver:5000`
 * `dev-deploy` deploy the bootstrap with `kubectl`
 * `dev-clean` destroy all cluster resources
 * `dev-redeploy` clean and then deploy the cluster
