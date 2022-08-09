# Prerequisites

* An Axway Platform user account that is assigned the AMPLIFY Central admin role
* A Layer7 Gateway on version 10.1 should be up and running and have APIs to be discovered and exposed in AMPLIFY Central

## Prepare AMPLIFY Central Environment

### Create an environment in Central

* Log into [Amplify Central](https://apicentral.axway.com)
* Navigate to "Topology" then "Environments"
* Click "+ Environment"
   * Select a name
   * Click "Save"
* To enable the viewing of the agent status in Amplify see [Visualize the agent status](https://docs.axway.com/bundle/amplify-central/page/docs/connect_manage_environ/environment_agent_resources/index.html#add-your-agent-resources-to-the-environment)

### Create a service account

* Create a public and private key pair locally using the openssl command

```sh
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits: 2048
openssl rsa -in private_key.pem -pubout -out public_key.pem
```

* Log into the [Amplify Platform](https://platform.axway.com)
* Navigate to "Organization" then "Service Accounts"
* Click "+ Service Account"
   * Select a name
   * Optionally add a description
   * Select "Client Certificate"
   * Select "Provide public key"
   * Select or paste the contents of the public_key.pem file
   * Select "Central admin"
   * Click "Save"
* Note the Client ID value, this and the key files will be needed for the agents

## Setup agent Environment Variables

The content below is an example environment variable file for starting the agent. Update the values as necessary to connect to Amplify Central and the Layer7 Gateway.
Copy the content into a file called `layer7.env`.

```ini
LAYER7_USERNAME=<Username>
LAYER7_PASSWORD=<Password>
LAYER7_HOST=<Gateway Host>
LAYER7_API=<Gateway REST API Endpoint>
LAYER7_POLLINTERVAL=30s

CENTRAL_GRPC_ENABLED=true
AGENTFEATURES_PERSISTCACHE=false
AGENTFEATURES_MARKETPLACEPROVISIONING=true

CENTRAL_PLATFORM_URL=https://platform.axway.com
CENTRAL_AUTH_URL=https://login.axway.com/auth
CENTRAL_URL=https://apicentral.axway.com

CENTRAL_AUTH_CLIENTID=<Client ID>
CENTRAL_AUTH_PRIVATEKEY=<File Path>
CENTRAL_AUTH_PUBLICKEY=<File Path
CENTRAL_AUTH_REALM=Broker

CENTRAL_ENVIRONMENT=<Environment Name>
CENTRAL_ORGANIZATIONID=426937327920148
CENTRAL_SSL_INSECURESKIPVERIFY=true
CENTRAL_SUBSCRIPTIONS_APPROVAL_MODE=auto
CENTRAL_VERSIONCHECKER=false
CENTRAL_TEAM="Default Team"

LOG_LEVEL=trace
LOG_FORMAT=line
```

# Layer7 Discovery Agent

The Discovery agent finds deployed API Proxies in the Layer7 API Gateway and sends them to Amplify Central.

## How to discover

The agent can discovery APIs in the Layer7 Gateway that have an OAS, Swagger, or WSDL file associated to the API. 
The API must have a spec file, it must be enabled, and it must not be an internal service.
If there is an API in the gateway, but does not meet the three requirements mentioned, then it will not be discovered.

### Swagger

To associate an API to a Swagger spec, a "Context Variable" is needed in the policy of the API.
To add a Context Variable to the API search for "Context Variable" in the "Assertions" tab in the top left of the Policy Manager.
Click on the assertion and add it to the policy. Set the "Variable Name" to `swagger.docUrl`. Set the "Expression" to the URL where the swagger can be found,
such as `https://petstore.swagger.io/v2/swagger.json`. Add the assertion to the policy when both fields are set.

### OpenAPI

To associate an API to an OpenAPI spec, a "Context Variable" is needed in the policy of the API.
To add a Context Variable to the API search for "Context Variable" in the "Assertions" tab in the top left of the Policy Manager.
Click on the assertion and add it to the policy. Set the "Variable Name" to `openapi.docUrl`. Set the "Expression" to the URL where the swagger can be found,
such as `https://petstore3.swagger.io/api/v3/openapi.json`. Add the assertion to the policy when both fields are set.

### WSDL

SOAP services published in the API Gateway are published from a WSDL or WSIL file. If the service is published from a WSDL file, then the service will be discovered by the agent.

### What is discovered?

The agent will find various values on an API such as
  * Service Name
  * Service Properties
  * Active Policy Version
  * API Endpoint

After an API is discovered there are two types of updates that can happen in Amplify Central regarding the API Service.

#### Major updates

  * If the API in the Gateway has a change to the "Resolution Path" found on the service in the Gateway, then a new revision is published in Central.
  * If a policy change of any kind is detected on the API, then a new revision is published to central.
  * If the active policy of the service changes, then a new revision is published to central

#### Minor updates

A minor update happens when a change is detected to the name of the service and to changes of the service properties.
A minor change does not trigger a new revision in central, and instead updates the latest revision.

## Build and run

The agent can be built by running `go build main.go -o ./bin/discovery`. The command `make build` can also be run to build the agent. This command will place the binary at `./bin/discovery`.

To run the agent run the command `./bin/discovery --envFile=layer7.env`

