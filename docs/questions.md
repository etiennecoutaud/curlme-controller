### How would you deploy your controller to a Kubernetes cluster ?

This controller is very simple so there is no need of a template configuration manager like kustomize or helm, I consider that tools must be very simple to install to get value immediately this is why I provide an all in one YAML file to deploy the controller easily

### Why did k/k decide to keep the spec and status away from each other?

`spec` is the desired state and `status` is the observed state. It allow the reconcilation loop to attempt to make `status` matches with `spec`

### The content returned when curling URLs may be always different. How is it going to affect your controllers?

As the result is not predictible it's going to affect testing, because we need predictive data to compare the expect and the result during test.

Controller must also manage errors which can be raised during execution.
