# fls-interpreter
FSL interpreter is an application which allows to process the fictitious scripting language (FSL) that is written in JSON.
The interpreter receives FSL as input. The FSL defines variables and functions.
The following operations are supported: 
- create 
- delete
- update
- add
- subtract
- multiply
- divide
- print
- function call 

Variables are always numeric.

### Script example
```
{
  "var1":1,
  "var2":2,
  
  "init": [
    {"cmd" : "#setup" }
  ],
  
  "setup": [
    {"cmd":"update", "id": "var1", "value":3.5},
    {"cmd":"print", "value": "#var1"},
    {"cmd":"#sum", "id": "var1", "value1":"#var1", "value2":"#var2"},
    {"cmd":"print", "value": "#var1"},
    {"cmd":"create", "id": "var3", "value":5},
    {"cmd":"delete", "id": "var1"},
    {"cmd":"#printAll"}
  ],
  
  "sum": [
      {"cmd":"add", "id": "$id", "operand1":"$value1", "operand2":"$value2"}
  ],

  "printAll":
  [
    {"cmd":"print", "value": "#var1"},
    {"cmd":"print", "value": "#var2"},
    {"cmd":"print", "value": "#var3"}
  ]
}
```
The interpreter could be used either as CLI tool or as microservice deployed in kubernetes cluster.
Script could be passed to service through REST endpoints as JSON payload or as a file.
also, the script can be executed via gPRC call.
### 0. Dependencies
In order to be able to build interpreter app following apps/tools should be installed.

- go toolchain > 1.16 https://go.dev/dl/
- docker https://docs.docker.com/engine/install/ubuntu/
- helm https://helm.sh/docs/intro/install/
- kind https://kind.sigs.k8s.io/
- jq - optional

### 1. Clone repo

1. `git clone https://github.com/alavrovinfb/fls-interpreter.git`

### 2. Build and run as CLI tool locally

1. `cd fls-interpreter`
2. `make build-local`
3. `bin/fls-interpreter --script.files="example/sample-script.txt,example/sample-script2.txt"`

example output:
```
2023/02/03 09:54:13 executing script example/sample-script.txt
3.5
5.5
undefined
2
5
2023/02/03 09:54:13 script example/sample-script.txt done.
2023/02/03 09:54:13 executing script example/sample-script2.txt
3.5
5.5
undefined
2
2.5
2023/02/03 09:54:13 script example/sample-script2.txt done.
```

### 3. Run as service
1. install kind if needed `make install-kind`
2. `make run-in-kind`
3. Wait until deployment will be rolled out
```
Waiting for deployment "fls-interpreter" rollout to finish: 0 of 1 updated replicas are available...
deployment "fls-interpreter" successfully rolled out
```
4. `kubectl -n fls port-forward svc/fls-interpreter 8080`
5. Pass script as payload: `curl -XPOST http://localhost:8080/fls-interpreter/v1/scripts -d '@./example/sample-script-curl.txt' |jq`

expected result:
```
{
  "result": [
    "3.5",
    "5.5",
    "undefined",
    "2",
    "5"
  ]
}
```
6. Pass script as file: `curl -XPOST http://localhost:8080/fls-interpreter/v1/files -F 'script=@./example/sample-script.txt' |jq`
Due to variables and functions are retained after previous execution, we need a way to reset them if needed.
It can be accomplished through `rest` endpoint.
`curl http://localhost:8080/fls-interpreter/v1/reset |jq`
7. To destroy test cluster run: `make kind-destroy`

### API doc (swagger)

Service provides `/swagger` endpoint to observe REST API scheme.
To browse API documentation following steps should be done:
1. `kubectl -n fls port-forward svc/fls-interpreter 8080`
2. Open in browser `http://localhost:8080/swagger`

### Monitoring and profiling

For monitoring and profiling purposes service provides:
- Prometheus' metrics endpoint `/metrics`
- Profiling endpoints:
```
/debug/pprof/
/debug/pprof/cmdline
/debug/pprof/profile
/debug/pprof/symbol
/debug/pprof/trace
```
To use these endpoints, pod port should be exposed.
1. `export POD_NAME=$(kubectl get pods --namespace fls -l "app.kubernetes.io/name=fls-interpreter,app.kubernetes.io/instance=fls-interpreter" -o jsonpath="{.items[0].metadata.name}")`
2. `kubectl port-forward $POD_NAME 8081:8081`
3. `curl http://localhost:8081/metrics`

links for useful articles about profiling:
- https://pkg.go.dev/net/http/pprof
- https://gist.github.com/slok/33dad1d0d0bae07977e6d32bcc010188
