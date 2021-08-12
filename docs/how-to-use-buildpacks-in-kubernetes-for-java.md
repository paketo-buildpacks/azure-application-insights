## How to use buildpacks in Kubernetes for Java

The Paketo Azure Application Insights Buildpack is a Cloud Native Buildpack
that contributes the Application Insights Agent and configures it to connect
to the service.

This article describes how to use this buildpack for Java applications from
the Kubernetes environment.

### Preconditions

* Kubernetes
* [kpack](https://github.com/pivotal/kpack) or [Tanzu Build Service](https://network.pivotal.io/products/build-service/)

### Build Phase

You need to prepare one `ConfigMap` from Kubernetes when build. It may be similar
to the example below.

> Note: the type should be `ApplicationInsights`.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: sample-cm
  namespace: build-service
data:
  provider: Azure
  type: ApplicationInsights
```

Then consume the `Confimap` from an Image resource. For example `image.yaml` seen below.

> Note: the `metadataRef` should point to the name of `ConfigMap` above.

```yaml
apiVersion: kpack.io/v1alpha1
kind: Image
metadata:
  name: sample-image
  namespace: build-service
spec:
  tag: xxxx.example.com/sample-repo:sample-tag
  serviceAccount: sample-build-service
  builder:
    name: sample-builder
    kind: Builder
  source:
    git:
      url: https://github.com/xxxx/xxxx
      revision: dev 
  build:
    bindings:
    - name: appinsights
      metadataRef:
        name: sample-cm
```

Now you can use below command to build the image.

```shell
kubectl apply -f image.yaml
```

If everything goes well, the output of build will contain the information below.

```
Paketo Azure Application Insights Buildpack 4.3.0
  https://github.com/paketo-buildpacks/azure-application-insights
  Azure Application Insights Java Agent 3.0.3: Contributing to layer
    Reusing cached download from buildpack
    Copying to /layers/paketo-buildpacks_azure-application-insights/azure-application-insights-java
    Writing env.launch/JAVA_TOOL_OPTIONS.append
    Writing env.launch/JAVA_TOOL_OPTIONS.delim
  Launch Helper: Contributing to layer
    Creating /layers/paketo-buildpacks_azure-application-insights/helper/exec.d/properties
```

### Runtime Phase

It is important that the binding be present at runtime because the binding data is not embedded into the image, so you need to prepare one `Secret` in Kubernetes before the Java application bootup,
it may be similar as below.

> Note: the type should be `ApplicationInsights`.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: sample-secret
  namespace: default
stringData:
  type: ApplicationInsights
  connection-string: "xxxx"
  sampling-percentage: "100.0"
```

Additional settings may be added as key/value pairs in the secret.

You also need to prepare at least 2 required items for Kubernetes deployment. For example `deployment.yaml` seen below.

* Mount the `Secret` as volume.
* Point the environment variable `CNB_BINDINGS` to the path of mounted `Secret`.

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      namespace: default
    spec:
      containers:
        env:
        - name: CNB_BINDINGS
          value: /bindings
        image: xxxx.example.com/sample-repo:sample-tag
        volumeMounts:
        - mountPath: /bindings/application-insights-settings
          name: sample-secret-volume
      volumes:
      - name: sample-secret-volume
        secret:
          secretName: sample-secret
```

Apply the configuration with the following command:

```shell
kubectl apply -f deployment.yaml
```

Finally, the Java application will bootup with the agent and consume the `Secret`
content as environment variables.

