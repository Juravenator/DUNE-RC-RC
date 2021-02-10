# DUNE-RC-RC

```
$ minikube start --mount-string ~/git/DUNE-RC-RC:/DUNE-RC-RC --mount
$ minikube node add
$ minikube node add
$ kubectl taint nodes minikube-m02 dedicated=APA:NoSchedule
$ kubectl taint nodes minikube-m03 dedicated=APA:NoSchedule
$ kubectl proxy --port=8080 &
$ curl --header "Content-Type: application/json-patch+json" --request PATCH --data '[{"op": "add", "path": "/status/capacity/rc.ccm~1ru-42-apa-tpc-a", "value": "1"}]' http://localhost:8080/api/v1/nodes/minikube-m02/status
$ curl --header "Content-Type: application/json-patch+json" --request PATCH --data '[{"op": "add", "path": "/status/capacity/rc.ccm~1ru-43-apa-tpc-a", "value": "1"}]' http://localhost:8080/api/v1/nodes/minikube-m03/status
$ curl --header "Content-Type: application/json-patch+json" --request PATCH --data '[{"op": "add", "path": "/status/capacity/rc.ccm~1ru-42-51-apa-pds-a", "value": "1"}]' http://localhost:8080/api/v1/nodes/minikube-m03/status
$ make
$ make manifests
$ kubectl apply -f config/crd/bases
$ kubectl apply -f config/samples/something.yaml
$ make run
```