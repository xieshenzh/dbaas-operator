# API Reference

## Packages
- [dbaas.redhat.com/v1beta1](#dbaasredhatcomv1beta1)


## dbaas.redhat.com/v1beta1

Package v1beta1 contains API Schema definitions for the dbaas v1beta1 API group

### Resource Types
- [DBaaSConnection](#dbaasconnection)
- [DBaaSInstance](#dbaasinstance)
- [DBaaSInventory](#dbaasinventory)
- [DBaaSPlatform](#dbaasplatform)
- [DBaaSPolicy](#dbaaspolicy)
- [DBaaSProvider](#dbaasprovider)



#### ConditionalProvisioningParameterData



ConditionalProvisioningParameterData defines the list of options with the corresponding default value available for a dropdown field, or the list of default values for an input text field in the UX based on the dependencies A provisioning parameter can have multiple option lists/default values depending on the dependent parameters. For instance, there are 4 different option lists for regions: one for dedicated cluster on GCP, one for dedicated on AWS, one for serverless on GCP, and one for serverless on AWS. If options lists are present, the field is displayed as a dropdown in the UX. Otherwise it is displayed as an input text.

_Appears in:_
- [ProvisioningParameter](#provisioningparameter)

| Field | Description |
| --- | --- |
| `dependencies` _[FieldDependency](#fielddependency) array_ | List of the dependent fields and their values |
| `options` _[Option](#option) array_ | Options displayed in the UX |
| `defaultValue` _string_ | Default value |


#### CredentialField



Defines the attributes.

_Appears in:_
- [DBaaSProviderSpec](#dbaasproviderspec)

| Field | Description |
| --- | --- |
| `key` _string_ | The name for this field. |
| `displayName` _string_ | A user-friendly name for this field. |
| `type` _string_ | The type of field: string, maskedstring, integer, or boolean. |
| `required` _boolean_ | Defines if the field is required or not. |
| `helpText` _string_ | Additional information about the field. |


#### DBaaSConnection



The schema for the DBaaSConnection API.



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `dbaas.redhat.com/v1beta1`
| `kind` _string_ | `DBaaSConnection`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[DBaaSConnectionSpec](#dbaasconnectionspec)_ |  |


#### DBaaSConnectionPolicy



DBaaSConnectionPolicy sets connection policy

_Appears in:_
- [DBaaSInventoryPolicy](#dbaasinventorypolicy)

| Field | Description |
| --- | --- |
| `namespaces` _string_ | Namespaces where DBaaSConnection and DBaaSInstance objects are only allowed to reference a policy's inventories. Using an asterisk surrounded by single quotes ('*'), allows all namespaces. If not set in the policy or by an inventory object, connections are only allowed in the inventory's namespace. |
| `nsSelector` _[LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#labelselector-v1-meta)_ | Use a label selector to determine the namespaces where DBaaSConnection and DBaaSInstance objects are only allowed to reference a policy's inventories. A label selector is a label query over a set of resources. Results use a logical AND from matchExpressions and matchLabels queries. An empty label selector matches all objects. A null label selector matches no objects. |


#### DBaaSConnectionSpec



Defines the desired state of a DBaaSConnection object.

_Appears in:_
- [DBaaSConnection](#dbaasconnection)
- [DBaaSProviderConnection](#dbaasproviderconnection)

| Field | Description |
| --- | --- |
| `inventoryRef` _[NamespacedName](#namespacedname)_ | A reference to the relevant DBaaSInventory custom resource (CR). |
| `instanceID` _string_ | The ID of the instance to connect to, as seen in the status of the referenced DBaaSInventory. |
| `instanceRef` _[NamespacedName](#namespacedname)_ | A reference to the DBaaSInstance CR used, if the InstanceID is not specified. |


#### DBaaSInstance



The schema for the DBaaSInstance API.



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `dbaas.redhat.com/v1beta1`
| `kind` _string_ | `DBaaSInstance`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[DBaaSInstanceSpec](#dbaasinstancespec)_ |  |


#### DBaaSInstanceSpec



Defines the desired state of a DBaaSInstance object.

_Appears in:_
- [DBaaSInstance](#dbaasinstance)
- [DBaaSProviderInstance](#dbaasproviderinstance)

| Field | Description |
| --- | --- |
| `inventoryRef` _[NamespacedName](#namespacedname)_ | A reference to the relevant DBaaSInventory custom resource (CR). |
| `provisioningParameters` _object (keys:[ProvisioningParameterType](#provisioningparametertype), values:string)_ | Parameters with values used for provisioning. |


#### DBaaSInventory



The schema for the DBaaSInventory API. Inventory objects must be created in a valid namespace, determined by the existence of a DBaaSPolicy object.



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `dbaas.redhat.com/v1beta1`
| `kind` _string_ | `DBaaSInventory`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[DBaaSOperatorInventorySpec](#dbaasoperatorinventoryspec)_ |  |


#### DBaaSInventoryPolicy



Sets the inventory policy.

_Appears in:_
- [DBaaSOperatorInventorySpec](#dbaasoperatorinventoryspec)
- [DBaaSPolicySpec](#dbaaspolicyspec)

| Field | Description |
| --- | --- |
| `disableProvisions` _boolean_ | Disables provisioning on inventory accounts. |
| `connections` _[DBaaSConnectionPolicy](#dbaasconnectionpolicy)_ | Namespaces where DBaaSConnection and DBaaSInstance objects are only allowed to reference a policy's inventories. |


#### DBaaSInventorySpec



DBaaSInventorySpec defines the Inventory Spec to be used by provider operators

_Appears in:_
- [DBaaSOperatorInventorySpec](#dbaasoperatorinventoryspec)
- [DBaaSProviderInventory](#dbaasproviderinventory)

| Field | Description |
| --- | --- |
| `credentialsRef` _[LocalObjectReference](#localobjectreference)_ | The secret containing the provider-specific connection credentials to use with the provider's API endpoint. The format specifies the secret in the provider’s operator for its DBaaSProvider custom resource (CR), such as the CredentialFields key. The secret must exist within the same namespace as the inventory. |


#### DBaaSOperatorInventorySpec



This object defines the desired state of a DBaaSInventory object.

_Appears in:_
- [DBaaSInventory](#dbaasinventory)

| Field | Description |
| --- | --- |
| `providerRef` _[NamespacedName](#namespacedname)_ | A reference to a DBaaSProvider custom resource (CR). |
| `DBaaSInventorySpec` _[DBaaSInventorySpec](#dbaasinventoryspec)_ | The properties that will be copied into the provider’s inventory. |
| `policy` _[DBaaSInventoryPolicy](#dbaasinventorypolicy)_ | The policy for this inventory. |


#### DBaaSPlatform



The schema for the DBaaSPlatform API.



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `dbaas.redhat.com/v1beta1`
| `kind` _string_ | `DBaaSPlatform`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[DBaaSPlatformSpec](#dbaasplatformspec)_ |  |


#### DBaaSPlatformSpec



Defines the desired state of a DBaaSPlatform object.

_Appears in:_
- [DBaaSPlatform](#dbaasplatform)

| Field | Description |
| --- | --- |
| `syncPeriod` _integer_ | The SyncPeriod set The minimum interval at which the provider operator controllers reconcile, the default value is 180 minutes. |


#### DBaaSPolicy



Enables administrative capabilities within a namespace, and sets a default inventory policy. Policy defaults can be overridden on a per-inventory basis.



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `dbaas.redhat.com/v1beta1`
| `kind` _string_ | `DBaaSPolicy`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[DBaaSPolicySpec](#dbaaspolicyspec)_ |  |


#### DBaaSPolicySpec



The specifications for a _DBaaSPolicy_ object. Enables administrative capabilities within a namespace, and sets a default inventory policy. Policy defaults can be overridden on a per-inventory basis.

_Appears in:_
- [DBaaSPolicy](#dbaaspolicy)

| Field | Description |
| --- | --- |
| `DBaaSInventoryPolicy` _[DBaaSInventoryPolicy](#dbaasinventorypolicy)_ |  |


#### DBaaSProvider



The schema for the DBaaSProvider API.



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `dbaas.redhat.com/v1beta1`
| `kind` _string_ | `DBaaSProvider`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[DBaaSProviderSpec](#dbaasproviderspec)_ |  |








#### DBaaSProviderSpec



Defines the desired state of a DBaaSProvider object.

_Appears in:_
- [DBaaSProvider](#dbaasprovider)

| Field | Description |
| --- | --- |
| `provider` _[DatabaseProviderInfo](#databaseproviderinfo)_ | Contains information about database provider and platform. |
| `groupVersion` _string_ | DBaaS API group version supported by the provider |
| `inventoryKind` _string_ | The name of the inventory custom resource definition (CRD) as defined by the database provider. |
| `connectionKind` _string_ | The name of the connection's custom resource definition (CRD) as defined by the provider. |
| `instanceKind` _string_ | The name of the instance's custom resource definition (CRD) as defined by the provider for provisioning. |
| `credentialFields` _[CredentialField](#credentialfield) array_ | Indicates what information to collect from the user interface and how to display fields in a form. |
| `allowsFreeTrial` _boolean_ | Indicates whether the provider offers free trials. |
| `externalProvisionURL` _string_ | The URL for provisioning instances by using the database provider's web portal. |
| `externalProvisionDescription` _string_ | Instructions on how to provision instances by using the database provider's web portal. |
| `provisioningParameters` _object (keys:[ProvisioningParameterType](#provisioningparametertype), values:[ProvisioningParameter](#provisioningparameter))_ | Parameter specs used by UX for provisioning a database instance |


#### DatabaseProviderInfo



Defines the information for a DBaaSProvider object.

_Appears in:_
- [DBaaSProviderSpec](#dbaasproviderspec)

| Field | Description |
| --- | --- |
| `name` _string_ | The name used to specify the service binding origin parameter. For example, 'Red Hat DBaaS / MongoDB Atlas'. |
| `displayName` _string_ | A user-friendly name for this database provider. For example, 'MongoDB Atlas'. |
| `displayDescription` _string_ | Indicates the description text shown for a database provider within the user interface. For example, the catalog tile description. |
| `icon` _[ProviderIcon](#providericon)_ | Indicates what icon to display on the catalog tile. |


#### FieldDependency



FieldDependency defines the name and value of a field used as a dependency

_Appears in:_
- [ConditionalProvisioningParameterData](#conditionalprovisioningparameterdata)

| Field | Description |
| --- | --- |
| `field` _[ProvisioningParameterType](#provisioningparametertype)_ | Name of the field used as a dependency |
| `value` _string_ | Value of the field used as a dependency |




#### LocalObjectReference



Contains enough information to locate the referenced object inside the same namespace.

_Appears in:_
- [DBaaSInventorySpec](#dbaasinventoryspec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name of the referent. |


#### NamespacedName



Defines the namespace and name of a k8s resource.

_Appears in:_
- [DBaaSConnectionSpec](#dbaasconnectionspec)
- [DBaaSInstanceSpec](#dbaasinstancespec)
- [DBaaSOperatorInventorySpec](#dbaasoperatorinventoryspec)

| Field | Description |
| --- | --- |
| `namespace` _string_ | The namespace where an object of a known type is stored. |
| `name` _string_ | The name for object of a known type. |




#### Option



Option defines the value and display value for an option in a dropdown, radio button or checkbox

_Appears in:_
- [ConditionalProvisioningParameterData](#conditionalprovisioningparameterdata)

| Field | Description |
| --- | --- |
| `value` _string_ | Value of the option |
| `displayValue` _string_ | Corresponding display value |






#### ProviderIcon



Follows the same field and naming formats as a comma-separated values (CSV) file.

_Appears in:_
- [DatabaseProviderInfo](#databaseproviderinfo)

| Field | Description |
| --- | --- |
| `base64data` _string_ |  |
| `mediatype` _string_ |  |


#### ProvisioningParameter



Information for a provisioning parameter

_Appears in:_
- [DBaaSProviderSpec](#dbaasproviderspec)

| Field | Description |
| --- | --- |
| `displayName` _string_ | A user-friendly name for this field. |
| `helpText` _string_ | Additional info about the field. |
| `conditionalData` _[ConditionalProvisioningParameterData](#conditionalprovisioningparameterdata) array_ | Lists of additional data containing the options or default values for the field. |


#### ProvisioningParameterType

_Underlying type:_ `string`



_Appears in:_
- [DBaaSInstanceSpec](#dbaasinstancespec)
- [DBaaSProviderSpec](#dbaasproviderspec)
- [FieldDependency](#fielddependency)



