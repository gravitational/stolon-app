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
### Building images
Execute `all` make target(could be omiited as it is default target).
```sh
$ make
```

### Building self-sufficient gravity image(a.k.a `Cluster Image`)
Download gravity and tele binaries
```
make download-binaries
```

Dowload and unpack dependent application packages into state directory(`./state` by default)
```
make install-dependent-packages
```

Build cluster image
```
export PATH=$(pwd)/bin:$PATH
make build-app
```

*Optional*: Build cluster image with intermediate runtime
```
export PATH=$(pwd)/bin:$PATH
make build-app INTERMEDIATE_RUNTIME_VERSION=5.2.18
```

## Prerequisites

* docker >=  1.8
* golang >= 1.13
* GNU make
* kubectl >= 1.13
