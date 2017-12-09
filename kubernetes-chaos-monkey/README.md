## Kubernetes Chaos Monkey

*WARNING:* This is strictly an example of how to use some of the kubernetes apis. If you wish to actually run a chaos monkey, please refer to either:

- https://github.com/asobti/kube-monkey
- https://github.com/Netflix/chaosmonkey

This example shows how to list and delete pods using the client-go kubernetes api

## Usage

```
./kubernetes-chaos-monkey
```

Every 10 seconds, this will delete a random pod from the cluster.