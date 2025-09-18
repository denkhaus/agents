# RequestAgentsCreateArgs

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Capabilities** | **[]string** |  | 
**Description** | **string** |  | 
**Name** | **string** |  | 
**Role** | [**RequestAgentsListByRoleArgsRole**](RequestAgentsListByRoleArgsRole.md) |  | 

## Methods

### NewRequestAgentsCreateArgs

`func NewRequestAgentsCreateArgs(capabilities []string, description string, name string, role RequestAgentsListByRoleArgsRole, ) *RequestAgentsCreateArgs`

NewRequestAgentsCreateArgs instantiates a new RequestAgentsCreateArgs object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRequestAgentsCreateArgsWithDefaults

`func NewRequestAgentsCreateArgsWithDefaults() *RequestAgentsCreateArgs`

NewRequestAgentsCreateArgsWithDefaults instantiates a new RequestAgentsCreateArgs object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCapabilities

`func (o *RequestAgentsCreateArgs) GetCapabilities() []string`

GetCapabilities returns the Capabilities field if non-nil, zero value otherwise.

### GetCapabilitiesOk

`func (o *RequestAgentsCreateArgs) GetCapabilitiesOk() (*[]string, bool)`

GetCapabilitiesOk returns a tuple with the Capabilities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapabilities

`func (o *RequestAgentsCreateArgs) SetCapabilities(v []string)`

SetCapabilities sets Capabilities field to given value.


### GetDescription

`func (o *RequestAgentsCreateArgs) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *RequestAgentsCreateArgs) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *RequestAgentsCreateArgs) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetName

`func (o *RequestAgentsCreateArgs) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *RequestAgentsCreateArgs) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *RequestAgentsCreateArgs) SetName(v string)`

SetName sets Name field to given value.


### GetRole

`func (o *RequestAgentsCreateArgs) GetRole() RequestAgentsListByRoleArgsRole`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *RequestAgentsCreateArgs) GetRoleOk() (*RequestAgentsListByRoleArgsRole, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *RequestAgentsCreateArgs) SetRole(v RequestAgentsListByRoleArgsRole)`

SetRole sets Role field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


