# provider-nats

`provider-nats` is a [Crossplane](https://crossplane.io/) Provider
that implements [NATS Jetstream](https://docs.nats.io/nats-concepts/jetstream) managed resources

Currently the provider supports the following resources:
- Stream: https://docs.nats.io/nats-concepts/jetstream/streams
- Consumer: https://docs.nats.io/nats-concepts/jetstream/consumers

Future releases might implement the key/value store and the object store as well.

## Examples

You might find the [examples](examples) directory helpful. Every example in this directory is deployable in a `make dev` environment.

For a full spec of possible options look either at [apis/consumer/v1alpha1/consumer/consumer_types.go](apis/consumer/v1alpha1/consumer/consumer_types.go), [apis/stream/v1alpha1/stream/stream_types.go](apis/stream/v1alpha1/stream/stream_types.go)
or use the `kubectl explain` command.

```bash
$ kubectl explain streams.nats.crossplane.io.spec.forProvider
$ kubectl explain consumers.nats.crossplane.io.spec.forProvider
```

You can also follow the [NATS Jetstream configuration docs](https://docs.nats.io/nats-concepts/jetstream/streams#configuration) and [NATS Jetstream consumer configuration docs](https://docs.nats.io/nats-concepts/jetstream/consumers#configuration) for more information as the managed resources implement basically the same configuration options.

### Example stream resource

```yaml
apiVersion: nats.crossplane.io/v1alpha1
kind: Stream
metadata:
  name: foo
spec:
  forProvider:
    domain: foo
    config:
      subjects:
        - foo.>
      retention: Limits
      storage: File
      maxBytes: 102400
      discard: Old
  providerConfigRef:
    name: default
```

### Example minimal pull consumer resource

```yaml
apiVersion: nats.crossplane.io/v1alpha1
kind: Consumer
metadata:
  name: pull
spec:
  forProvider:
    stream: foo
    config:
      pull: {}
  providerConfigRef:
    name: default

```
## Developing locally

Start a local development environment using `kind` with crossplane and a complete NATS environment. Ensure that you can reach `nats.nats.svc` on 127.0.0.1

```bash
# The provider needs to know how to reach the NATS server. We use kubectl port-forward to expose the NATS server on localhost
echo "nats.nats.svc" >> /etc/hosts
# Create the cluster
make dev
# Please ensure you reach the right cluster before we start the port-forward. Usually it should be the right one after `make dev`
kubectl port-forward -n nats svc/nats 4222:4222
```

Once the environment is up and running you can deploy the provider and the managed resources. Have a look at the files used to e2e tests ([cluster/local/e2e/manifests/](cluster/local/e2e/manifests/)).

### Debugging

Just start the debugger of your choice to debug `cmd/provider/main.go`.
The only thing that is important is, that your KUBECONFIG points to a dev cluster with the CRDs deployed.

## E2E Tests

Simply run `make e2e` for running the e2e tests. They spin up a `kind` cluster with a fully provisioned NATS environment including leaf nats servers and JWT/NKEY based authentication.
The e2e tests ensure that managing streams and consumers works as expected.

If you want to develop or debug e2e tests you have some options for the creation of the test cluster.

Provide the following env variables for `make e2e`:
- `E2E_DEV`: if `true`, the e2e cluster is not deleted after the tests
- `E2E_SKIP_TESTS`: if `true`, the e2e tests are skipped

To create a cluster that is capable of running the e2e tests and keep it running during development of the e2e tests you can combine both:

```bash
E2E_DEV=true E2E_SKIP_TESTS=true make e2e
```
