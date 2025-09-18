# RequestTasksUpdateStateArgs

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**State** | [**RequestTasksUpdateStateArgsState**](RequestTasksUpdateStateArgsState.md) |  | 

## Methods

### NewRequestTasksUpdateStateArgs

`func NewRequestTasksUpdateStateArgs(id string, state RequestTasksUpdateStateArgsState, ) *RequestTasksUpdateStateArgs`

NewRequestTasksUpdateStateArgs instantiates a new RequestTasksUpdateStateArgs object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRequestTasksUpdateStateArgsWithDefaults

`func NewRequestTasksUpdateStateArgsWithDefaults() *RequestTasksUpdateStateArgs`

NewRequestTasksUpdateStateArgsWithDefaults instantiates a new RequestTasksUpdateStateArgs object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *RequestTasksUpdateStateArgs) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *RequestTasksUpdateStateArgs) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *RequestTasksUpdateStateArgs) SetId(v string)`

SetId sets Id field to given value.


### GetState

`func (o *RequestTasksUpdateStateArgs) GetState() RequestTasksUpdateStateArgsState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *RequestTasksUpdateStateArgs) GetStateOk() (*RequestTasksUpdateStateArgsState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *RequestTasksUpdateStateArgs) SetState(v RequestTasksUpdateStateArgsState)`

SetState sets State field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


