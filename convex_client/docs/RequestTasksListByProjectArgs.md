# RequestTasksListByProjectArgs

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IncludeCompleted** | Pointer to **bool** |  | [optional] 
**ProjectId** | **string** |  | 

## Methods

### NewRequestTasksListByProjectArgs

`func NewRequestTasksListByProjectArgs(projectId string, ) *RequestTasksListByProjectArgs`

NewRequestTasksListByProjectArgs instantiates a new RequestTasksListByProjectArgs object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRequestTasksListByProjectArgsWithDefaults

`func NewRequestTasksListByProjectArgsWithDefaults() *RequestTasksListByProjectArgs`

NewRequestTasksListByProjectArgsWithDefaults instantiates a new RequestTasksListByProjectArgs object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIncludeCompleted

`func (o *RequestTasksListByProjectArgs) GetIncludeCompleted() bool`

GetIncludeCompleted returns the IncludeCompleted field if non-nil, zero value otherwise.

### GetIncludeCompletedOk

`func (o *RequestTasksListByProjectArgs) GetIncludeCompletedOk() (*bool, bool)`

GetIncludeCompletedOk returns a tuple with the IncludeCompleted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncludeCompleted

`func (o *RequestTasksListByProjectArgs) SetIncludeCompleted(v bool)`

SetIncludeCompleted sets IncludeCompleted field to given value.

### HasIncludeCompleted

`func (o *RequestTasksListByProjectArgs) HasIncludeCompleted() bool`

HasIncludeCompleted returns a boolean if a field has been set.

### GetProjectId

`func (o *RequestTasksListByProjectArgs) GetProjectId() string`

GetProjectId returns the ProjectId field if non-nil, zero value otherwise.

### GetProjectIdOk

`func (o *RequestTasksListByProjectArgs) GetProjectIdOk() (*string, bool)`

GetProjectIdOk returns a tuple with the ProjectId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProjectId

`func (o *RequestTasksListByProjectArgs) SetProjectId(v string)`

SetProjectId sets ProjectId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


