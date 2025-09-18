# RequestEventsEmitEventArgs

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | **interface{}** |  | 
**EntityId** | **string** |  | 
**Type** | **string** |  | 
**UserId** | Pointer to **string** |  | [optional] 

## Methods

### NewRequestEventsEmitEventArgs

`func NewRequestEventsEmitEventArgs(data interface{}, entityId string, type_ string, ) *RequestEventsEmitEventArgs`

NewRequestEventsEmitEventArgs instantiates a new RequestEventsEmitEventArgs object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRequestEventsEmitEventArgsWithDefaults

`func NewRequestEventsEmitEventArgsWithDefaults() *RequestEventsEmitEventArgs`

NewRequestEventsEmitEventArgsWithDefaults instantiates a new RequestEventsEmitEventArgs object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *RequestEventsEmitEventArgs) GetData() interface{}`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *RequestEventsEmitEventArgs) GetDataOk() (*interface{}, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *RequestEventsEmitEventArgs) SetData(v interface{})`

SetData sets Data field to given value.


### SetDataNil

`func (o *RequestEventsEmitEventArgs) SetDataNil(b bool)`

 SetDataNil sets the value for Data to be an explicit nil

### UnsetData
`func (o *RequestEventsEmitEventArgs) UnsetData()`

UnsetData ensures that no value is present for Data, not even an explicit nil
### GetEntityId

`func (o *RequestEventsEmitEventArgs) GetEntityId() string`

GetEntityId returns the EntityId field if non-nil, zero value otherwise.

### GetEntityIdOk

`func (o *RequestEventsEmitEventArgs) GetEntityIdOk() (*string, bool)`

GetEntityIdOk returns a tuple with the EntityId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEntityId

`func (o *RequestEventsEmitEventArgs) SetEntityId(v string)`

SetEntityId sets EntityId field to given value.


### GetType

`func (o *RequestEventsEmitEventArgs) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *RequestEventsEmitEventArgs) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *RequestEventsEmitEventArgs) SetType(v string)`

SetType sets Type field to given value.


### GetUserId

`func (o *RequestEventsEmitEventArgs) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *RequestEventsEmitEventArgs) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *RequestEventsEmitEventArgs) SetUserId(v string)`

SetUserId sets UserId field to given value.

### HasUserId

`func (o *RequestEventsEmitEventArgs) HasUserId() bool`

HasUserId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


