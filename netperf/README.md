# netperf tests

These tests were performed with the imac always acting as the server, in the following setups:

## bare metal
```
[gdirkx@glenn-imac09 ~]$ netserver -D
[glenn-macbook ~]# bash client.sh 10.69.0.3 macbook-to-imac 100
```

## docker swarm
```
[gdirkx@glenn-imac09 ~]$ docker swarm init --advertise-addr 10.69.0.3
[gdirkx@glenn-imac09 ~]$ docker network create -d overlay --attachable my_net
[glenn-macbook ~]$ docker swarm join --advertise-addr 10.69.0.2 --token abc 10.69.0.3:2377
[gdirkx@glenn-imac09 ~]$ docker run -d --network my_net tailoredcloud/netperf:v2.7 netserver -D
[glenn-macbook ~]$ docker run -it --network my_net tailoredcloud/netperf:v2.7 sh
$ apk update
$ apk add bash
$ bash client.sh 10.0.1.8 macbook-to-imac 100
```

## docker host net
Same as docker swarm, except using `--network host` and `10.69.0.3` as IP.

## Kubernetes with Calico
```
$ kubectl apply -f pods.yaml
$ kubectl get pods -owide
NAME   READY   STATUS    RESTARTS   AGE   IP             NODE            NOMINATED NODE   READINESS GATES
alp1   1/1     Running   0          17s   10.42.248.72   glenn-imac09    <none>           <none>
alp2   1/1     Running   0          17s   10.42.254.5    glenn-macbook   <none>           <none>
```

```
[gdirkx@glenn-imac09 ~]$ kubectl exec -it alp1 -- sh
$ netserver -D
[gdirkx@glenn-imac09 ~]$ kubectl exec -it alp2 -- sh
$ apk update
$ apk add bash
$ bash client.sh 10.42.248.72 macbook-to-imac 100
```

For good measure, a check to see if we're leveraging the kernel and not relying on software overlay:
```
[glenn-imac09 ~]# ip route
default via 10.69.0.2 dev enp0s10 proto static metric 100 
10.42.248.64 dev cali071a07b47bc scope link 
blackhole 10.42.248.64/26 proto bird 
10.42.248.65 dev cali0b82c630646 scope link 
10.42.248.66 dev calicf568764cd4 scope link 
10.42.248.70 dev calidc0cea093f8 scope link 
10.42.254.0/26 via 10.69.0.2 dev enp0s10 proto bird 
10.69.0.0/24 dev enp0s10 proto kernel scope link src 10.69.0.3 metric 100 
10.85.0.0/16 dev cni0 proto kernel scope link src 10.85.0.1 linkdown 
172.17.0.0/16 dev docker0 proto kernel scope link src 172.17.0.1 linkdown 
172.18.0.0/16 dev docker_gwbridge proto kernel scope link src 172.18.0.1
[glenn-macbook ~]# ip route
default via 192.168.0.1 dev wlp3s0 proto dhcp metric 600 
10.42.248.64/26 via 10.69.0.3 dev ens9 proto bird 
10.42.254.0 dev cali4793d1165f6 scope link 
blackhole 10.42.254.0/26 proto bird 
10.69.0.0/24 dev ens9 proto kernel scope link src 10.69.0.2 metric 100 
172.17.0.0/16 dev docker0 proto kernel scope link src 172.17.0.1 linkdown 
172.18.0.0/16 dev docker_gwbridge proto kernel scope link src 172.18.0.1 
192.168.0.0/24 dev wlp3s0 proto kernel scope link src 192.168.0.231 metric 600
```