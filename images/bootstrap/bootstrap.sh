#!/bin/sh

cd /root/cfssl

if kubectl get secret/cluster-ca ; then
	echo "secret/cluster-ca already exists"
else
	cfssl gencert -initca ca-csr.json|cfssljson -bare ca -

	kubectl create secret generic cluster-ca \
		--from-file=ca.pem=ca.pem \
		--from-file=ca-key=ca-key.pem \
		--from-file=ca.csr=ca.csr
fi

if kubectl get secret/cluster-default-ssl ; then
	echo "secret/cluster-ca already exists"
else
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json \
		-profile=server default-server-csr.json | cfssljson -bare default-server
	cp default-server.pem default-server-with-chain.pem
	cat ca.pem >> default-server-with-chain.pem

	kubectl create secret generic cluster-default-ssl \
		--from-file=default-server.pem=default-server.pem \
		--from-file=default-server-with-chain.pem=default-server-with-chain.pem \
		--from-file=default-server-key.pem=default-server-key.pem \
		--from-file=default-server.csr=default-server.csr
fi

if kubectl get secret/cluster-kube-system-ssl ; then
	echo "secret/cluster-kube-system-ssl already exists"
else
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json \
		-profile=server kube-system-server-csr.json | cfssljson -bare kube-system-server
	cp kube-system-server.pem kube-system-server-with-chain.pem
	cat ca.pem >> kube-system-server-with-chain.pem

	kubectl create secret generic cluster-kube-system-ssl \
		--from-file=default-server.pem=default-server.pem \
		--from-file=kube-system-server-with-chain.pem=kube-system-server-with-chain.pem \
		--from-file=default-server-key.pem=default-server-key.pem \
		--from-file=default-server.csr=default-server.csr
fi

/usr/local/bin/stolonboot -sentinels 1 -rpc 1
