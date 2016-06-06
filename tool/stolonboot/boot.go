package main

import (
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/gravitational/trace"
)

func bootCluster(sentinels int, keepers int, proxies int, password string) error {
	err := createEtcd()
	if err != nil {
		return trace.Wrap(err)
	}

	err = createSentinels(sentinels)
	if err != nil {
		return trace.Wrap(err)
	}

	err = createSecret(password)
	if err != nil {
		return trace.Wrap(err)
	}

	err = createKeepers(keepers)
	if err != nil {
		return trace.Wrap(err)
	}

	err = createProxies(proxies)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

func createEtcd() error {
	log.Infof("creating etcd")
	cmd := exec.Command("/usr/local/bin/kubectl", "create", "-f", "/resources/etcd.yml")
	out, err := cmd.Output()
	if err != nil {
		return trace.Wrap(err)
	}
	log.Infof("cmd output: %s", string(out))
	return nil
}

func createSentinels(sentinels int) error {
	log.Infof("creating sentinels")
	cmd := exec.Command("/usr/local/bin/kubectl", "create", "-f", "/resources/sentinel.yml")
	out, err := cmd.Output()
	if err != nil {
		return trace.Wrap(err)
	}
	log.Infof("cmd output: %s", string(out))
	return nil
}

func createSecret(password string) error {
	log.Infof("creating secret")
	cmd := exec.Command("/usr/local/bin/kubectl", "create", "-f", "/resources/secret.yml")
	out, err := cmd.Output()
	if err != nil {
		return trace.Wrap(err)
	}
	log.Infof("cmd output: %s", string(out))
	return nil
}

func createKeepers(keepers int) error {
	log.Infof("creating keepers")
	cmd := exec.Command("/usr/local/bin/kubectl", "create", "-f", "/resources/keeper.yml")
	out, err := cmd.Output()
	if err != nil {
		return trace.Wrap(err)
	}
	log.Infof("cmd output: %s", string(out))
	return nil
}

func createProxies(proxies int) error {
	log.Infof("creating proxies")
	cmd := exec.Command("/usr/local/bin/kubectl", "create", "-f", "/resources/proxy.yml")
	out, err := cmd.Output()
	if err != nil {
		return trace.Wrap(err)
	}
	log.Infof("cmd output: %s", string(out))
	return nil
}
