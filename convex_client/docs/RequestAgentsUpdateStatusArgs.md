# RequestAgentsUpdateStatusArgs

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**IsStreaming** | Pointer to **bool** |  | [optional] 
**Status** | [**RequestAgentsListByStatusArgsStatus**](RequestAgentsListByStatusArgsStatus.md) |  | 

## Methods

### NewRequestAgentsUpdateStatusArgs

`func NewRequestAgentsUpdateStatusArgs(id string, status RequestAgentsListByStatusArgsStatus, ) *RequestAgentsUpdateStatusArgs`

NewRequestAgentsUpdateStatusArgs instantiates a new RequestAgentsUpdateStatusArgs object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRequestAgentsUpdateStatusArgsWithDefaults

`func NewRequestAgentsUpdateStatusArgsWithDefaults() *RequestAgentsUpdateStatusArgs`

NewRequestAgentsUpdateStatusArgsWithDefaults instantiates a new RequestAgentsUpdateStatusArgs object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *RequestAgentsUpdateStatusArgs) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *RequestAgentsUpdateStatusArgs) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *RequestAgentsUpdateStatusArgs) SetId(v string)`

SetId sets Id field to given value.


### GetIsStreaming

`func (o *RequestAgentsUpdateStatusArgs) GetIsStreaming() bool`

GetIsStreaming returns the IsStreaming field if non-nil, zero value otherwise.

### GetIsStreamingOk

`func (o *RequestAgentsUpdateStatusArgs) GetIsStreamingOk() (*bool, bool)`

GetIsStreamingOk returns a tuple with the IsStreaming field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsStreaming

`func (o *RequestAgentsUpdateStatusArgs) SetIsStreaming(v bool)`

SetIsStreaming sets IsStreaming field to given value.

### HasIsStreaming

`func (o *RequestAgentsUpdateStatusArgs) HasIsStreaming() bool`

HasIsStreaming returns a boolean if a field has been set.

### GetStatus

`func (o *RequestAgentsUpdateStatusArgs) GetStatus() RequestAgentsListByStatusArgsStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *RequestAgentsUpdateStatusArgs) GetStatusOk() (*RequestAgentsListByStatusArgsStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *RequestAgentsUpdateStatusArgs) SetStatus(v RequestAgentsListByStatusArgsStatus)`

SetStatus sets Status field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


