# RequestTasksCreateArgs

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AssignedAgent** | Pointer to **string** |  | [optional] 
**Complexity** | **float32** |  | 
**Dependencies** | Pointer to **[]string** |  | [optional] 
**Depth** | **float32** |  | 
**Description** | **string** |  | 
**Estimate** | Pointer to **float32** |  | [optional] 
**ParentId** | Pointer to **string** |  | [optional] 
**ProjectId** | **string** |  | 
**Title** | **string** |  | 

## Methods

### NewRequestTasksCreateArgs

`func NewRequestTasksCreateArgs(complexity float32, depth float32, description string, projectId string, title string, ) *RequestTasksCreateArgs`

NewRequestTasksCreateArgs instantiates a new RequestTasksCreateArgs object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRequestTasksCreateArgsWithDefaults

`func NewRequestTasksCreateArgsWithDefaults() *RequestTasksCreateArgs`

NewRequestTasksCreateArgsWithDefaults instantiates a new RequestTasksCreateArgs object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAssignedAgent

`func (o *RequestTasksCreateArgs) GetAssignedAgent() string`

GetAssignedAgent returns the AssignedAgent field if non-nil, zero value otherwise.

### GetAssignedAgentOk

`func (o *RequestTasksCreateArgs) GetAssignedAgentOk() (*string, bool)`

GetAssignedAgentOk returns a tuple with the AssignedAgent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAssignedAgent

`func (o *RequestTasksCreateArgs) SetAssignedAgent(v string)`

SetAssignedAgent sets AssignedAgent field to given value.

### HasAssignedAgent

`func (o *RequestTasksCreateArgs) HasAssignedAgent() bool`

HasAssignedAgent returns a boolean if a field has been set.

### GetComplexity

`func (o *RequestTasksCreateArgs) GetComplexity() float32`

GetComplexity returns the Complexity field if non-nil, zero value otherwise.

### GetComplexityOk

`func (o *RequestTasksCreateArgs) GetComplexityOk() (*float32, bool)`

GetComplexityOk returns a tuple with the Complexity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetComplexity

`func (o *RequestTasksCreateArgs) SetComplexity(v float32)`

SetComplexity sets Complexity field to given value.


### GetDependencies

`func (o *RequestTasksCreateArgs) GetDependencies() []string`

GetDependencies returns the Dependencies field if non-nil, zero value otherwise.

### GetDependenciesOk

`func (o *RequestTasksCreateArgs) GetDependenciesOk() (*[]string, bool)`

GetDependenciesOk returns a tuple with the Dependencies field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDependencies

`func (o *RequestTasksCreateArgs) SetDependencies(v []string)`

SetDependencies sets Dependencies field to given value.

### HasDependencies

`func (o *RequestTasksCreateArgs) HasDependencies() bool`

HasDependencies returns a boolean if a field has been set.

### GetDepth

`func (o *RequestTasksCreateArgs) GetDepth() float32`

GetDepth returns the Depth field if non-nil, zero value otherwise.

### GetDepthOk

`func (o *RequestTasksCreateArgs) GetDepthOk() (*float32, bool)`

GetDepthOk returns a tuple with the Depth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDepth

`func (o *RequestTasksCreateArgs) SetDepth(v float32)`

SetDepth sets Depth field to given value.


### GetDescription

`func (o *RequestTasksCreateArgs) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *RequestTasksCreateArgs) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *RequestTasksCreateArgs) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetEstimate

`func (o *RequestTasksCreateArgs) GetEstimate() float32`

GetEstimate returns the Estimate field if non-nil, zero value otherwise.

### GetEstimateOk

`func (o *RequestTasksCreateArgs) GetEstimateOk() (*float32, bool)`

GetEstimateOk returns a tuple with the Estimate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEstimate

`func (o *RequestTasksCreateArgs) SetEstimate(v float32)`

SetEstimate sets Estimate field to given value.

### HasEstimate

`func (o *RequestTasksCreateArgs) HasEstimate() bool`

HasEstimate returns a boolean if a field has been set.

### GetParentId

`func (o *RequestTasksCreateArgs) GetParentId() string`

GetParentId returns the ParentId field if non-nil, zero value otherwise.

### GetParentIdOk

`func (o *RequestTasksCreateArgs) GetParentIdOk() (*string, bool)`

GetParentIdOk returns a tuple with the ParentId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParentId

`func (o *RequestTasksCreateArgs) SetParentId(v string)`

SetParentId sets ParentId field to given value.

### HasParentId

`func (o *RequestTasksCreateArgs) HasParentId() bool`

HasParentId returns a boolean if a field has been set.

### GetProjectId

`func (o *RequestTasksCreateArgs) GetProjectId() string`

GetProjectId returns the ProjectId field if non-nil, zero value otherwise.

### GetProjectIdOk

`func (o *RequestTasksCreateArgs) GetProjectIdOk() (*string, bool)`

GetProjectIdOk returns a tuple with the ProjectId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProjectId

`func (o *RequestTasksCreateArgs) SetProjectId(v string)`

SetProjectId sets ProjectId field to given value.


### GetTitle

`func (o *RequestTasksCreateArgs) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *RequestTasksCreateArgs) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *RequestTasksCreateArgs) SetTitle(v string)`

SetTitle sets Title field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


