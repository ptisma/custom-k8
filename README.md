# custom-k8

## Goal of this project is to learn more about the k8 CRD-s and Operators. This will be achieved by writing our own Operator/CRD.

#### sample-api
Codebase for our simple app (Dockerfile, docker-compose and go files). Written in Go.

#### sample-api-operator
Codebase for our Operator which is going to "watch over" for our deployed sample-api app on k8. Generated and written using go/operator-sdk.
operator-sdk version: "v1.31.0", commit: "e67da35ef4fff3e471a208904b2a142b27ae32b1", kubernetes version: "1.26.0", go version: "go1.19.11", GOOS: "linux", GOARCH: "amd64"

#### sample-api-app-crd.yaml contains the configuration for our CRD sample-api app
Example k8 yaml for our CRD simple-api app.

##### Instructions

Run Operator locally:
1. Create empty folder sample-api-operator and navigate to it
2. operator-sdk init --domain example.com --repo https://github.com/ptisma/custom-k8 sample-api-operator
3. operator-sdk create api --group=app --version=v1 --kind=SampleAPIApp --resource --controller
4. go mod tidy
5. make generate
6. make manifests
7. make install
8. make run

Create CRD:
kubectl apply -f sample-api-app-crd.yaml