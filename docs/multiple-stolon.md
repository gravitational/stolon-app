# Installing multiple stolon clusters in one Gravity cluster
This document provides an instruction on how to build and install multiple stolon clusters 
in one Gravity cluster.

## Requirements

- tele binary: 5.5+
- stolon application: 1.13.0+

## Building
By default using provided [Makefile](../Makefile) stolon application is built as a Gravity application and could be installed in Gravity cluster with `gravity install` from resulting tarball.  
For the additional stolon(s) in cluster [Application Catalog](https://gravitational.com/gravity/docs/catalog/) feature could be reused. Application does not need to be pushed to catalog, but the process of making of that application image type applies here.

Stolon application contains [Helm](https://helm.sh/) chart in `resources/charts/stolon`. Custom values are stored in [custom-values.yaml](../resources/custom-values.yaml). They are important for installation in Gravity clusters and should be used as a starting point for your customization.

#### Building Docker images
There is a special target in [Makefile](../Makefile) to build Docker images:

```
$ make images
```

Tag of the resulting images is automatically calculated from git by [version.sh](../version.sh) script and could be overriden by `VERSION` environment variable.  
To build a Stolon application image with default custom-values file:

#### Changing version of Helm chart
Version of the stolon chart is hardcoded to `0.1.0` in [Chart.yaml](../resources/charts/stolon/Chart.yaml). You could replace it with command `sed -i "s/0.1.0/$(./version.sh)/g" resources/charts/stolon/Chart.yaml` or with any version you want to use.

#### Building Gravity application image
```
$ tele build --values resources/custom-values.yaml \
	--set registry="" --set tag=$(./version.sh) \
	--version=$(./version.sh) resources/charts/stolon 
```

*Note*: If `VERSION` was overriden during the building of images then `--set tag=` should be set with the same value.  
Command produces tarball with vendored images and Helm chart inside. Default name of the tarball will be stolon-app-$(VERSION).tar in the root directory of the repository. Flag `-o(--output) filename.tar` could be added to `tele build` command to specify the name of a resulting tarball.

## Installation
Resulted tarball with the application should be uploaded to Gravity cluster along with the [custom-values.yaml](../resources/custom-values.yaml) files because the last one needs to be used for setting custom Chart values.

``` shell
# gravity app install --values custom-values.yaml --set clusterName=kube-stolon2 --set keeperDataPath=/var/lib/data/stolon2 --name stolon2 stolon-app.tar
```
There are two important values that must be set for the second(or more) stolon in Gravity cluster: `clusterName` and `keeperDataPath`.
- `clusterName` specifies the path in Etcd to store data about the stolon cluster. Name should be unique for each stolon cluster.
- `keeperDataPath` specifies hostPath on the hosts for storing stolon data. Path should be unique for each stolon cluster.

## Upgrade
TODO
