# `gcr.io/paketo-buildpacks/azure-application-insights`
The Paketo Azure Application Insights Buildpack is a Cloud Native Buildpack that contributes the Application Insights Agent and configures it to connect to the service.

## Behavior
This buildpack will participate if all the following conditions are met

* A binding exists with `type` of `ApplicationInsights`

The buildpack will do the following for Java applications:

* Contributes a Java agent to a layer and configures `$JAVA_OPTS` to use it
* Transforms the contents of the binding secret to environment variables with the pattern `APPINSIGHTS_<KEY>=<VALUE>`

The buildpack will do the following NodeJS applications:

* Contributes a NodeJS agent to a layer and configures `$NODE_MODULES` to use it
* If main module does not already require `appinsights` module, prepends the main module with `require('applicationinsights').start();`
* Transforms the contents of the binding secret to environment variables with the pattern `APPINSIGHTS_<KEY>=<VALUE>`

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0
