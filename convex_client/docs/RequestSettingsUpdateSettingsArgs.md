# RequestSettingsUpdateSettingsArgs

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AutoLayout** | Pointer to **bool** |  | [optional] 
**AutoSave** | Pointer to **bool** |  | [optional] 
**CurrentWorkspace** | Pointer to **string** |  | [optional] 
**Language** | Pointer to **string** |  | [optional] 
**LeftSidebarCollapsed** | Pointer to **bool** |  | [optional] 
**Notifications** | Pointer to **bool** |  | [optional] 
**RightSidebarCollapsed** | Pointer to **bool** |  | [optional] 
**SelectedNodeIds** | Pointer to **[]string** |  | [optional] 
**SelectedProjectId** | Pointer to **string** |  | [optional] 
**ShowBackground** | Pointer to **bool** |  | [optional] 
**ShowMiniMap** | Pointer to **bool** |  | [optional] 
**Theme** | Pointer to [**RequestSettingsUpdateSettingsArgsTheme**](RequestSettingsUpdateSettingsArgsTheme.md) |  | [optional] 
**UserId** | **string** |  | 

## Methods

### NewRequestSettingsUpdateSettingsArgs

`func NewRequestSettingsUpdateSettingsArgs(userId string, ) *RequestSettingsUpdateSettingsArgs`

NewRequestSettingsUpdateSettingsArgs instantiates a new RequestSettingsUpdateSettingsArgs object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRequestSettingsUpdateSettingsArgsWithDefaults

`func NewRequestSettingsUpdateSettingsArgsWithDefaults() *RequestSettingsUpdateSettingsArgs`

NewRequestSettingsUpdateSettingsArgsWithDefaults instantiates a new RequestSettingsUpdateSettingsArgs object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAutoLayout

`func (o *RequestSettingsUpdateSettingsArgs) GetAutoLayout() bool`

GetAutoLayout returns the AutoLayout field if non-nil, zero value otherwise.

### GetAutoLayoutOk

`func (o *RequestSettingsUpdateSettingsArgs) GetAutoLayoutOk() (*bool, bool)`

GetAutoLayoutOk returns a tuple with the AutoLayout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAutoLayout

`func (o *RequestSettingsUpdateSettingsArgs) SetAutoLayout(v bool)`

SetAutoLayout sets AutoLayout field to given value.

### HasAutoLayout

`func (o *RequestSettingsUpdateSettingsArgs) HasAutoLayout() bool`

HasAutoLayout returns a boolean if a field has been set.

### GetAutoSave

`func (o *RequestSettingsUpdateSettingsArgs) GetAutoSave() bool`

GetAutoSave returns the AutoSave field if non-nil, zero value otherwise.

### GetAutoSaveOk

`func (o *RequestSettingsUpdateSettingsArgs) GetAutoSaveOk() (*bool, bool)`

GetAutoSaveOk returns a tuple with the AutoSave field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAutoSave

`func (o *RequestSettingsUpdateSettingsArgs) SetAutoSave(v bool)`

SetAutoSave sets AutoSave field to given value.

### HasAutoSave

`func (o *RequestSettingsUpdateSettingsArgs) HasAutoSave() bool`

HasAutoSave returns a boolean if a field has been set.

### GetCurrentWorkspace

`func (o *RequestSettingsUpdateSettingsArgs) GetCurrentWorkspace() string`

GetCurrentWorkspace returns the CurrentWorkspace field if non-nil, zero value otherwise.

### GetCurrentWorkspaceOk

`func (o *RequestSettingsUpdateSettingsArgs) GetCurrentWorkspaceOk() (*string, bool)`

GetCurrentWorkspaceOk returns a tuple with the CurrentWorkspace field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrentWorkspace

`func (o *RequestSettingsUpdateSettingsArgs) SetCurrentWorkspace(v string)`

SetCurrentWorkspace sets CurrentWorkspace field to given value.

### HasCurrentWorkspace

`func (o *RequestSettingsUpdateSettingsArgs) HasCurrentWorkspace() bool`

HasCurrentWorkspace returns a boolean if a field has been set.

### GetLanguage

`func (o *RequestSettingsUpdateSettingsArgs) GetLanguage() string`

GetLanguage returns the Language field if non-nil, zero value otherwise.

### GetLanguageOk

`func (o *RequestSettingsUpdateSettingsArgs) GetLanguageOk() (*string, bool)`

GetLanguageOk returns a tuple with the Language field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLanguage

`func (o *RequestSettingsUpdateSettingsArgs) SetLanguage(v string)`

SetLanguage sets Language field to given value.

### HasLanguage

`func (o *RequestSettingsUpdateSettingsArgs) HasLanguage() bool`

HasLanguage returns a boolean if a field has been set.

### GetLeftSidebarCollapsed

`func (o *RequestSettingsUpdateSettingsArgs) GetLeftSidebarCollapsed() bool`

GetLeftSidebarCollapsed returns the LeftSidebarCollapsed field if non-nil, zero value otherwise.

### GetLeftSidebarCollapsedOk

`func (o *RequestSettingsUpdateSettingsArgs) GetLeftSidebarCollapsedOk() (*bool, bool)`

GetLeftSidebarCollapsedOk returns a tuple with the LeftSidebarCollapsed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeftSidebarCollapsed

`func (o *RequestSettingsUpdateSettingsArgs) SetLeftSidebarCollapsed(v bool)`

SetLeftSidebarCollapsed sets LeftSidebarCollapsed field to given value.

### HasLeftSidebarCollapsed

`func (o *RequestSettingsUpdateSettingsArgs) HasLeftSidebarCollapsed() bool`

HasLeftSidebarCollapsed returns a boolean if a field has been set.

### GetNotifications

`func (o *RequestSettingsUpdateSettingsArgs) GetNotifications() bool`

GetNotifications returns the Notifications field if non-nil, zero value otherwise.

### GetNotificationsOk

`func (o *RequestSettingsUpdateSettingsArgs) GetNotificationsOk() (*bool, bool)`

GetNotificationsOk returns a tuple with the Notifications field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotifications

`func (o *RequestSettingsUpdateSettingsArgs) SetNotifications(v bool)`

SetNotifications sets Notifications field to given value.

### HasNotifications

`func (o *RequestSettingsUpdateSettingsArgs) HasNotifications() bool`

HasNotifications returns a boolean if a field has been set.

### GetRightSidebarCollapsed

`func (o *RequestSettingsUpdateSettingsArgs) GetRightSidebarCollapsed() bool`

GetRightSidebarCollapsed returns the RightSidebarCollapsed field if non-nil, zero value otherwise.

### GetRightSidebarCollapsedOk

`func (o *RequestSettingsUpdateSettingsArgs) GetRightSidebarCollapsedOk() (*bool, bool)`

GetRightSidebarCollapsedOk returns a tuple with the RightSidebarCollapsed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRightSidebarCollapsed

`func (o *RequestSettingsUpdateSettingsArgs) SetRightSidebarCollapsed(v bool)`

SetRightSidebarCollapsed sets RightSidebarCollapsed field to given value.

### HasRightSidebarCollapsed

`func (o *RequestSettingsUpdateSettingsArgs) HasRightSidebarCollapsed() bool`

HasRightSidebarCollapsed returns a boolean if a field has been set.

### GetSelectedNodeIds

`func (o *RequestSettingsUpdateSettingsArgs) GetSelectedNodeIds() []string`

GetSelectedNodeIds returns the SelectedNodeIds field if non-nil, zero value otherwise.

### GetSelectedNodeIdsOk

`func (o *RequestSettingsUpdateSettingsArgs) GetSelectedNodeIdsOk() (*[]string, bool)`

GetSelectedNodeIdsOk returns a tuple with the SelectedNodeIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSelectedNodeIds

`func (o *RequestSettingsUpdateSettingsArgs) SetSelectedNodeIds(v []string)`

SetSelectedNodeIds sets SelectedNodeIds field to given value.

### HasSelectedNodeIds

`func (o *RequestSettingsUpdateSettingsArgs) HasSelectedNodeIds() bool`

HasSelectedNodeIds returns a boolean if a field has been set.

### GetSelectedProjectId

`func (o *RequestSettingsUpdateSettingsArgs) GetSelectedProjectId() string`

GetSelectedProjectId returns the SelectedProjectId field if non-nil, zero value otherwise.

### GetSelectedProjectIdOk

`func (o *RequestSettingsUpdateSettingsArgs) GetSelectedProjectIdOk() (*string, bool)`

GetSelectedProjectIdOk returns a tuple with the SelectedProjectId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSelectedProjectId

`func (o *RequestSettingsUpdateSettingsArgs) SetSelectedProjectId(v string)`

SetSelectedProjectId sets SelectedProjectId field to given value.

### HasSelectedProjectId

`func (o *RequestSettingsUpdateSettingsArgs) HasSelectedProjectId() bool`

HasSelectedProjectId returns a boolean if a field has been set.

### GetShowBackground

`func (o *RequestSettingsUpdateSettingsArgs) GetShowBackground() bool`

GetShowBackground returns the ShowBackground field if non-nil, zero value otherwise.

### GetShowBackgroundOk

`func (o *RequestSettingsUpdateSettingsArgs) GetShowBackgroundOk() (*bool, bool)`

GetShowBackgroundOk returns a tuple with the ShowBackground field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShowBackground

`func (o *RequestSettingsUpdateSettingsArgs) SetShowBackground(v bool)`

SetShowBackground sets ShowBackground field to given value.

### HasShowBackground

`func (o *RequestSettingsUpdateSettingsArgs) HasShowBackground() bool`

HasShowBackground returns a boolean if a field has been set.

### GetShowMiniMap

`func (o *RequestSettingsUpdateSettingsArgs) GetShowMiniMap() bool`

GetShowMiniMap returns the ShowMiniMap field if non-nil, zero value otherwise.

### GetShowMiniMapOk

`func (o *RequestSettingsUpdateSettingsArgs) GetShowMiniMapOk() (*bool, bool)`

GetShowMiniMapOk returns a tuple with the ShowMiniMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShowMiniMap

`func (o *RequestSettingsUpdateSettingsArgs) SetShowMiniMap(v bool)`

SetShowMiniMap sets ShowMiniMap field to given value.

### HasShowMiniMap

`func (o *RequestSettingsUpdateSettingsArgs) HasShowMiniMap() bool`

HasShowMiniMap returns a boolean if a field has been set.

### GetTheme

`func (o *RequestSettingsUpdateSettingsArgs) GetTheme() RequestSettingsUpdateSettingsArgsTheme`

GetTheme returns the Theme field if non-nil, zero value otherwise.

### GetThemeOk

`func (o *RequestSettingsUpdateSettingsArgs) GetThemeOk() (*RequestSettingsUpdateSettingsArgsTheme, bool)`

GetThemeOk returns a tuple with the Theme field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTheme

`func (o *RequestSettingsUpdateSettingsArgs) SetTheme(v RequestSettingsUpdateSettingsArgsTheme)`

SetTheme sets Theme field to given value.

### HasTheme

`func (o *RequestSettingsUpdateSettingsArgs) HasTheme() bool`

HasTheme returns a boolean if a field has been set.

### GetUserId

`func (o *RequestSettingsUpdateSettingsArgs) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *RequestSettingsUpdateSettingsArgs) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *RequestSettingsUpdateSettingsArgs) SetUserId(v string)`

SetUserId sets UserId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


