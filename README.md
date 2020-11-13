# `paketo-buildpacks/microsoft-azure`
The Paketo Microsoft Azure Buildpack is a Cloud Native Buildpack that contributes Microsoft Azure agents and configures them to connec to their services.

## Behavior

* If `$BP_APPLICATION_INSIGHTS_ENABLED` is set to `true` and the application is Java
  * At build time, contributes an agent to a layer
  * At launch time, if credentials are available, configures the application to use the agent
* If `$BP_APPLICATION_INSIGHTS_ENABLED` is set to `true` and the application is NodeJS
  * At build time, contributes an agent to a layer
  * At launch time, if credentials are available, configures `$NODE_MODULES` with the agent path.  If the main module does not already require `applicationinsights`, prepends the main module with `require('applicationinsights').start();`.

### Credential Availability
If the applications runs within Microsoft Azure and the [Azure Metadata Service][m] is accessible, those credentials will be used.  If the application runs within any other environment, credentials must be provided with a service binding as described below.

[m]: https://docs.microsoft.com/en-us/azure/virtual-machines/windows/instance-metadata-service

## Configuration
| Environment Variable | Description
| -------------------- | -----------
| `$BP_AZURE_APPLICATION_INSIGHTS_ENABLED` | Whether to add Microsoft Azure Application Insights during build

## Bindings
The buildpack optionally accepts the following bindings:

### Type: `MicrosoftAzure`
|Key                  | Value   | Description
|---------------------|---------|------------
|`InstrumentationKey` | `<key>` | Azure Application Insights instrumentation key

### Type: `dependency-mapping`
|Key                   | Value   | Description
|----------------------|---------|------------
|`<dependency-digest>` | `<uri>` | If needed, the buildpack will fetch the dependency with digest `<dependency-digest>` from `<uri>`

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0
