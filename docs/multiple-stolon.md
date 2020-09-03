# Installing multiple stolon clusters into a Gravity cluster
This document provides instructions on how to build and install multiple stolon clusters into a Gravity cluster.

## Requirements

- tele binary: 5.5+
- stolon application: 1.13.0+

## Building
By default using provided [Makefile](../Makefile) stolon application is built as a Gravity application and could be installed in Gravity cluster with `gravity install` from resulting tarball.  
To add (or introduce) additional stolon instances into the cluster, use the [Application Catalog](https://gravitational.com/gravity/docs/catalog/). Application does not need to be pushed to catalog, but the process of making of that application image type applies here.

Stolon application contains a [Helm](https://helm.sh/) chart in `resources/charts/stolon`. Custom values are important for the installation into a Gravity cluster. [custom-values.yaml](../resources/custom-values.yaml) contains a set of values that can be used as a starting point for further customizations.

#### Building Docker images
There is a special target in [Makefile](../Makefile) to build all Docker images for the application:

```
$ make images
```

[version.sh](../version.sh) automatically calculates a tag name to use for Docker image(s) based on a git tag. Use the VERSION environment variable to override the version from command line.  
To build a Stolon application image with default custom-values file:

#### Changing version of Helm chart
Version of the stolon chart defaults to `0.1.0` but can be overridden by editing [Chart.yaml](../resources/charts/stolon/Chart.yaml).

#### Building Gravity application image
```
$ tele build --values resources/custom-values.yaml \
	--set registry="" --set tag=$(./version.sh) \
	--version=$(./version.sh) resources/charts/stolon 
```

*Note*: If you change the version (e.g. via VERSION environment variable), set the tag to the same value with `--set tag=<version>`  
As specified, the above command will create a tarball named `stolon-app-$(VERSION).tar` in the repository's root directory with Docker images and Helm charts bundled inside. To change the name of the output file, use `--output=<new-name>`.

## Installation
Resulted tarball with the application should be uploaded to Gravity cluster along with the [custom-values.yaml](../resources/custom-values.yaml) files because the last one needs to be used for setting custom Chart values.

``` shell
# gravity app install --values custom-values.yaml --set clusterName=kube-stolon2 --set keeperDataPath=/var/lib/data/stolon2 --name stolon2 stolon-app.tar
```
Each subsequent instance of stolon requires setting the values of `clusterName` and `keeperDataPath`.
- `clusterName` specifies a path in etcd that the stolon instance will use for its data. The path must be unique for each stolon cluster.
- `keeperDataPath` specifies a path on host each instance will use to store PostgreSQL data. The path should be unique for each stolon cluster.

## Upgrade
TODO
