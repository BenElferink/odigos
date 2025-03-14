# Odigos Chaos Testing

Another kind of tests we would like to run, are chaos tests.
The purpose of those tests is to tests the odigos platform with different fault and errors on the cluster level.

## Tools

- [Kubernetes In Docker (Kind)](https://kind.sigs.k8s.io/) - a tool for running local Kubernetes clusters using Docker container “nodes”.
- [Chainsaw](https://kyverno.github.io/chainsaw/) - To orchestrate the different Kubernetes actions.
- [Chaos-mesh](https://github.com/chaos-mesh/chaos-mesh) - In order to simulate faults in the cluster.
- [Tempo](https://github.com/grafana/tempo) - Distributed tracing backend. Chosen due to its query language that allows for easy querying of traces.

## Running chaos tests locally

### Prerequisites

Install these tools once when setting up your local testing environment the first time.

- [Kubernetes In Docker (KinD)](https://kind.sigs.k8s.io/) - a tool for running local Kubernetes clusters using Docker container “nodes”.

- [Chainsaw](https://kyverno.github.io/chainsaw/) - To orchestrate the different Kubernetes actions.
  - Hombrew:

  ```bash
  brew tap kyverno/chainsaw https://github.com/kyverno/chainsaw
  brew install kyverno/chainsaw/chainsaw
  ```

  - Go:

  ```bash
  go install github.com/kyverno/chainsaw@latest
  ```

- helm, yq and jq installed:

```bash
brew install yq
brew install jq
brew install helm
```

### Preparing

You can run all the below steps with `make dev-tests-setup`.

- Fresh Kubernetes cluster in kubectl context. For local development, you can use KinD but also managed clusters like EKS will work. you can create the cluster with `make dev-tests-kind-cluster`.
- Odigos CLI compiled at the `cli` directory in odigos OSS repo (which is expected to be cloned as sibling of the current repo). To compile the cli executable, go to the OSS repository and run: `make cli-build`.

### Running the Tests

To run specific scenarios, for example `network-latency/leader-election` run from Odigos root directory:

```bash
chainsaw test tests/chaos/network-latency/leader-election
```

## Writing new scenarios

Every scenario should include some/all of the following:

- Install destination (usually Tempo)
- Install test applications
- Install Odigos
- Select apps for instrumentation and configure destination
- Generate traffic
- Validate traces

Scenarios are written in yaml files called `chainsaw-test.yaml` according to the Chainsaw schema.

See the [following document](https://kyverno.github.io/chainsaw/latest/test/) for more information on how to write scenarios.

Scenarios should be placed in the `tests/chaos/<fault-created>/<scenario-name>` directory.

Tests that expect to validate traces should have TraceQL validations and have them be placed in the `tests/chaos/<fault-created>/<scenraio-name>/traceql` directory.

After writing and testing new scenario, you should also add it to the GitHub Action file location at:
`.github/workflows/chaos.yaml` to run it on every pull request.

## Working with TraceQL

TraceQL is a query language that allows you to query traces in Tempo.
It is used in the end-to-end tests to validate the traces generated by Odigos.

### Connecting to Tempo

In order to run TraceQL queries, you need to connect to Tempo.
Tempo is installed automatically in the chaos test, so if you ran a scenario you can connect to it.
You can do this by port-forwarding the Tempo service:

```bash
kubectl port-forward svc/chaos-tests-tempo 3100:3100 -n traces
```

### Querying traces

Then you can execute TraceQL queries via:

```bash
curl -G -s http://localhost:3100/api/search --data-urlencode 'q={ resource.odigos.version = "chaos-test"}'
```

To get full individual trace you can use the following command:

```bash
curl -G -s http://localhost:3100/api/traces/3debdffae5920741a53d1bd015c62b29
```

For both APIs it is recommended to pipe the results to `jq` for better readability and `less` for paging.

### Writing queries

See [the following document](https://grafana.com/docs/tempo/latest/traceql/) for more information on how to write queries.
In order to add new traceql test, you need to add a new yaml file in the following schema:

```yaml
apiVersion: chaos.tests.odigos.io/v1
kind: TraceTest
description: <Description>
query: |
    <TraceQL Query>
expected:
  count: <Number of expected spans>
```

Once you have the file, you can run the test via:

```bash
tests/chaos/common/traceql_runner.sh <path-to-yaml-file>
```

## Debugging

When tests fail, and it's related to some traceql query not succeeding, it can be useful to setup a grafana ui to commit queries and see the traces that are stored in tempo.

### TL;DR

```bash
make dev-tests-grafana
make dev-tests-grafana-port-forward
```

Then browse to `http://localhost:3080/explore`.

### Detailed Steps

- Install grafana with helm once:

```bash
make dev-tests-grafana
```

or manually:

```bash
helm install -n traces grafana grafana/grafana \
  --set "env.GF_AUTH_ANONYMOUS_ENABLED=true" \
  --set "env.GF_AUTH_ANONYMOUS_ORG_ROLE=Admin" \
  --set "datasources.datasources\.yaml.apiVersion=1" \
  --set "datasources.datasources\.yaml.datasources[0].name=Tempo" \
  --set "datasources.datasources\.yaml.datasources[0].type=tempo" \
  --set "datasources.datasources\.yaml.datasources[0].url=http://chaos-tests-tempo:3100" \
  --set "datasources.datasources\.yaml.datasources[0].access=proxy" \
  --set "datasources.datasources\.yaml.datasources[0].isDefault=true"
```

- Port forward to the grafana service:

```bash
make dev-tests-grafana-port-forward
```

or manually:

```bash
kubectl port-forward svc/grafana 3080:80 -n traces
```

- Browse to `http://localhost:3080/explore`.

- You can now write queries, run them, and see the traces that are stored in tempo to troubleshoot your test issues. example query:

```
{resource.service.name = "coupon"}
```
