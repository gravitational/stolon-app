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


### Recovering Using a Continuous Archive Backup
Stolon application uses [Continuous Archive Backup](https://www.postgresql.org/docs/9.4/static/continuous-archiving.html) feature of PostgreSQL database.
You can restore Stolon database to the latest archived checkpoint.
1. Create new pod, where all restoration steps above[2-8] will perform.

``` shell
kubectl create -f /var/lib/gravity/local/packages/unpacked/gravitational.io/stolon-app/1.2.3/resources/restore.yaml
```
2. Delete old stolon resources. Important: **Do not create new stolon keepers before restoring data from WAL's in step 5.**

``` shell
kubectl delete daemonset stolon-keeper
kubectl delete deployment stolon-sentinel
```
3. Clean stolon data directories on host. You can check on which nodes stolon daemonset is exists with command below.

``` shell
kubectl get nodes -L stolon-keeper

rm -rf /var/lib/data/stolon/*
```
4. Create new empty database and restore latest WAL's.

``` shell
su - postgres -c "rm -rf /var/lib/postgresql/9.4/main/"
su - postgres -c "envdir /etc/wal-e.d/env /usr/local/bin/wal-e backup-fetch /var/lib/postgresql/9.4/main LATEST"
echo "restore_command = '/usr/bin/envdir /etc/wal-e.d/env /usr/local/bin/wal-e wal-fetch %f %p'" > /var/lib/postgresql/9.4/main/recovery.conf
chown postgres:postgres -R /var/lib/postgresql/9.4/main/
service postgresql start
```
5. Create dump of restored data.

``` shell
PGHOST=localhost PGUSER=stolon pg_dumpall > /root/dump.sql
```
6. Clear stolon cluster data in etcd.

``` shell
etcdctl --endpoint "https://${NODE_NAME}:2379" rm /stolon/cluster/kube-stolon --recursive
```
7. Create new sentinels and keepers.

``` shell
kubectl create -f /var/lib/gravity/resources/1.2.3/resources/sentinel.yaml
kubectl create -f /var/lib/gravity/resources/1.2.3/resources/keeper.yaml
```
8. Restore data into freshly created stolon cluster.

``` shell
psql -h stolon-postgres.default.svc -U stolon -d postgres < /tmp/dump.sql
```
9. Delete `stolon-restore` deployment after checking data in stolon database.

``` shell
kubectl delete deployment stolon-restore
```

### Development

There are several development `Makefile` targets to simplify your workflow:

 * `dev-push` push images to `apiserver:5000`
 * `dev-deploy` deploy the bootstrap with `kubectl`
 * `dev-clean` destroy all cluster resources
 * `dev-redeploy` clean and then deploy the cluster
 * `dev-hatest` run integration test.
