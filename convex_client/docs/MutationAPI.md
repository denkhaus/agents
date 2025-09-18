# \MutationAPI

All URIs are relative to *http://127.0.0.1:3210*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiRunAgentsAssignTaskPost**](MutationAPI.md#ApiRunAgentsAssignTaskPost) | **Post** /api/run/agents/assignTask | Calls a mutation at the path agents.js:assignTask
[**ApiRunAgentsCreatePost**](MutationAPI.md#ApiRunAgentsCreatePost) | **Post** /api/run/agents/create | Calls a mutation at the path agents.js:create
[**ApiRunAgentsUnassignTaskPost**](MutationAPI.md#ApiRunAgentsUnassignTaskPost) | **Post** /api/run/agents/unassignTask | Calls a mutation at the path agents.js:unassignTask
[**ApiRunAgentsUpdateStatusPost**](MutationAPI.md#ApiRunAgentsUpdateStatusPost) | **Post** /api/run/agents/updateStatus | Calls a mutation at the path agents.js:updateStatus
[**ApiRunClearDatabaseClearDatabasePost**](MutationAPI.md#ApiRunClearDatabaseClearDatabasePost) | **Post** /api/run/clearDatabase/clearDatabase | Calls a mutation at the path clearDatabase.js:clearDatabase
[**ApiRunEventsCleanupOldEventsPost**](MutationAPI.md#ApiRunEventsCleanupOldEventsPost) | **Post** /api/run/events/cleanupOldEvents | Calls a mutation at the path events.js:cleanupOldEvents
[**ApiRunEventsEmitEventPost**](MutationAPI.md#ApiRunEventsEmitEventPost) | **Post** /api/run/events/emitEvent | Calls a mutation at the path events.js:emitEvent
[**ApiRunProjectsCreatePost**](MutationAPI.md#ApiRunProjectsCreatePost) | **Post** /api/run/projects/create | Calls a mutation at the path projects.js:create
[**ApiRunProjectsRemovePost**](MutationAPI.md#ApiRunProjectsRemovePost) | **Post** /api/run/projects/remove | Calls a mutation at the path projects.js:remove
[**ApiRunProjectsUpdateEditableFieldsPost**](MutationAPI.md#ApiRunProjectsUpdateEditableFieldsPost) | **Post** /api/run/projects/updateEditableFields | Calls a mutation at the path projects.js:updateEditableFields
[**ApiRunProjectsUpdatePositionPost**](MutationAPI.md#ApiRunProjectsUpdatePositionPost) | **Post** /api/run/projects/updatePosition | Calls a mutation at the path projects.js:updatePosition
[**ApiRunProjectsUpdateStatsPost**](MutationAPI.md#ApiRunProjectsUpdateStatsPost) | **Post** /api/run/projects/updateStats | Calls a mutation at the path projects.js:updateStats
[**ApiRunSeedSeedDatabasePost**](MutationAPI.md#ApiRunSeedSeedDatabasePost) | **Post** /api/run/seed/seedDatabase | Calls a mutation at the path seed.js:seedDatabase
[**ApiRunSettingsUpdateSettingsPost**](MutationAPI.md#ApiRunSettingsUpdateSettingsPost) | **Post** /api/run/settings/updateSettings | Calls a mutation at the path settings.js:updateSettings
[**ApiRunSettingsUpdateThemePost**](MutationAPI.md#ApiRunSettingsUpdateThemePost) | **Post** /api/run/settings/updateTheme | Calls a mutation at the path settings.js:updateTheme
[**ApiRunTasksCreatePost**](MutationAPI.md#ApiRunTasksCreatePost) | **Post** /api/run/tasks/create | Calls a mutation at the path tasks.js:create
[**ApiRunTasksRemovePost**](MutationAPI.md#ApiRunTasksRemovePost) | **Post** /api/run/tasks/remove | Calls a mutation at the path tasks.js:remove
[**ApiRunTasksUpdateEditableFieldsPost**](MutationAPI.md#ApiRunTasksUpdateEditableFieldsPost) | **Post** /api/run/tasks/updateEditableFields | Calls a mutation at the path tasks.js:updateEditableFields
[**ApiRunTasksUpdatePositionPost**](MutationAPI.md#ApiRunTasksUpdatePositionPost) | **Post** /api/run/tasks/updatePosition | Calls a mutation at the path tasks.js:updatePosition
[**ApiRunTasksUpdateStatePost**](MutationAPI.md#ApiRunTasksUpdateStatePost) | **Post** /api/run/tasks/updateState | Calls a mutation at the path tasks.js:updateState



## ApiRunAgentsAssignTaskPost

> ResponseAgentsAssignTask ApiRunAgentsAssignTaskPost(ctx).RequestAgentsAssignTask(requestAgentsAssignTask).Execute()

Calls a mutation at the path agents.js:assignTask

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestAgentsAssignTask := *openapiclient.NewRequestAgentsAssignTask(*openapiclient.NewRequestAgentsAssignTaskArgs("AgentId_example", "TaskId_example")) // RequestAgentsAssignTask | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunAgentsAssignTaskPost(context.Background()).RequestAgentsAssignTask(requestAgentsAssignTask).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunAgentsAssignTaskPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunAgentsAssignTaskPost`: ResponseAgentsAssignTask
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunAgentsAssignTaskPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunAgentsAssignTaskPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestAgentsAssignTask** | [**RequestAgentsAssignTask**](RequestAgentsAssignTask.md) |  | 

### Return type

[**ResponseAgentsAssignTask**](ResponseAgentsAssignTask.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunAgentsCreatePost

> ResponseAgentsCreate ApiRunAgentsCreatePost(ctx).RequestAgentsCreate(requestAgentsCreate).Execute()

Calls a mutation at the path agents.js:create

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestAgentsCreate := *openapiclient.NewRequestAgentsCreate(*openapiclient.NewRequestAgentsCreateArgs([]string{"Capabilities_example"}, "Description_example", "Name_example", openapiclient.Request_agents_listByRole_args_role{String: new(string)})) // RequestAgentsCreate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunAgentsCreatePost(context.Background()).RequestAgentsCreate(requestAgentsCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunAgentsCreatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunAgentsCreatePost`: ResponseAgentsCreate
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunAgentsCreatePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunAgentsCreatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestAgentsCreate** | [**RequestAgentsCreate**](RequestAgentsCreate.md) |  | 

### Return type

[**ResponseAgentsCreate**](ResponseAgentsCreate.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunAgentsUnassignTaskPost

> ResponseAgentsUnassignTask ApiRunAgentsUnassignTaskPost(ctx).RequestAgentsUnassignTask(requestAgentsUnassignTask).Execute()

Calls a mutation at the path agents.js:unassignTask

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestAgentsUnassignTask := *openapiclient.NewRequestAgentsUnassignTask(*openapiclient.NewRequestAgentsAssignTaskArgs("AgentId_example", "TaskId_example")) // RequestAgentsUnassignTask | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunAgentsUnassignTaskPost(context.Background()).RequestAgentsUnassignTask(requestAgentsUnassignTask).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunAgentsUnassignTaskPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunAgentsUnassignTaskPost`: ResponseAgentsUnassignTask
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunAgentsUnassignTaskPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunAgentsUnassignTaskPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestAgentsUnassignTask** | [**RequestAgentsUnassignTask**](RequestAgentsUnassignTask.md) |  | 

### Return type

[**ResponseAgentsUnassignTask**](ResponseAgentsUnassignTask.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunAgentsUpdateStatusPost

> ResponseAgentsUpdateStatus ApiRunAgentsUpdateStatusPost(ctx).RequestAgentsUpdateStatus(requestAgentsUpdateStatus).Execute()

Calls a mutation at the path agents.js:updateStatus

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestAgentsUpdateStatus := *openapiclient.NewRequestAgentsUpdateStatus(*openapiclient.NewRequestAgentsUpdateStatusArgs("Id_example", openapiclient.Request_agents_listByStatus_args_status{String: new(string)})) // RequestAgentsUpdateStatus | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunAgentsUpdateStatusPost(context.Background()).RequestAgentsUpdateStatus(requestAgentsUpdateStatus).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunAgentsUpdateStatusPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunAgentsUpdateStatusPost`: ResponseAgentsUpdateStatus
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunAgentsUpdateStatusPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunAgentsUpdateStatusPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestAgentsUpdateStatus** | [**RequestAgentsUpdateStatus**](RequestAgentsUpdateStatus.md) |  | 

### Return type

[**ResponseAgentsUpdateStatus**](ResponseAgentsUpdateStatus.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunClearDatabaseClearDatabasePost

> ResponseClearDatabaseClearDatabase ApiRunClearDatabaseClearDatabasePost(ctx).RequestClearDatabaseClearDatabase(requestClearDatabaseClearDatabase).Execute()

Calls a mutation at the path clearDatabase.js:clearDatabase

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestClearDatabaseClearDatabase := *openapiclient.NewRequestClearDatabaseClearDatabase(map[string]interface{}(123)) // RequestClearDatabaseClearDatabase | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunClearDatabaseClearDatabasePost(context.Background()).RequestClearDatabaseClearDatabase(requestClearDatabaseClearDatabase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunClearDatabaseClearDatabasePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunClearDatabaseClearDatabasePost`: ResponseClearDatabaseClearDatabase
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunClearDatabaseClearDatabasePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunClearDatabaseClearDatabasePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestClearDatabaseClearDatabase** | [**RequestClearDatabaseClearDatabase**](RequestClearDatabaseClearDatabase.md) |  | 

### Return type

[**ResponseClearDatabaseClearDatabase**](ResponseClearDatabaseClearDatabase.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunEventsCleanupOldEventsPost

> ResponseEventsCleanupOldEvents ApiRunEventsCleanupOldEventsPost(ctx).RequestEventsCleanupOldEvents(requestEventsCleanupOldEvents).Execute()

Calls a mutation at the path events.js:cleanupOldEvents

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestEventsCleanupOldEvents := *openapiclient.NewRequestEventsCleanupOldEvents(*openapiclient.NewRequestEventsCleanupOldEventsArgs(float32(123))) // RequestEventsCleanupOldEvents | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunEventsCleanupOldEventsPost(context.Background()).RequestEventsCleanupOldEvents(requestEventsCleanupOldEvents).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunEventsCleanupOldEventsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunEventsCleanupOldEventsPost`: ResponseEventsCleanupOldEvents
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunEventsCleanupOldEventsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunEventsCleanupOldEventsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestEventsCleanupOldEvents** | [**RequestEventsCleanupOldEvents**](RequestEventsCleanupOldEvents.md) |  | 

### Return type

[**ResponseEventsCleanupOldEvents**](ResponseEventsCleanupOldEvents.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunEventsEmitEventPost

> ResponseEventsEmitEvent ApiRunEventsEmitEventPost(ctx).RequestEventsEmitEvent(requestEventsEmitEvent).Execute()

Calls a mutation at the path events.js:emitEvent

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestEventsEmitEvent := *openapiclient.NewRequestEventsEmitEvent(*openapiclient.NewRequestEventsEmitEventArgs(interface{}(123), "EntityId_example", "Type_example")) // RequestEventsEmitEvent | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunEventsEmitEventPost(context.Background()).RequestEventsEmitEvent(requestEventsEmitEvent).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunEventsEmitEventPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunEventsEmitEventPost`: ResponseEventsEmitEvent
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunEventsEmitEventPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunEventsEmitEventPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestEventsEmitEvent** | [**RequestEventsEmitEvent**](RequestEventsEmitEvent.md) |  | 

### Return type

[**ResponseEventsEmitEvent**](ResponseEventsEmitEvent.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunProjectsCreatePost

> ResponseProjectsCreate ApiRunProjectsCreatePost(ctx).RequestProjectsCreate(requestProjectsCreate).Execute()

Calls a mutation at the path projects.js:create

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestProjectsCreate := *openapiclient.NewRequestProjectsCreate(*openapiclient.NewRequestProjectsCreateArgs("Description_example", "Title_example")) // RequestProjectsCreate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunProjectsCreatePost(context.Background()).RequestProjectsCreate(requestProjectsCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunProjectsCreatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunProjectsCreatePost`: ResponseProjectsCreate
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunProjectsCreatePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunProjectsCreatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestProjectsCreate** | [**RequestProjectsCreate**](RequestProjectsCreate.md) |  | 

### Return type

[**ResponseProjectsCreate**](ResponseProjectsCreate.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunProjectsRemovePost

> ResponseProjectsRemove ApiRunProjectsRemovePost(ctx).RequestProjectsRemove(requestProjectsRemove).Execute()

Calls a mutation at the path projects.js:remove

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestProjectsRemove := *openapiclient.NewRequestProjectsRemove(*openapiclient.NewRequestAgentsGetArgs("Id_example")) // RequestProjectsRemove | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunProjectsRemovePost(context.Background()).RequestProjectsRemove(requestProjectsRemove).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunProjectsRemovePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunProjectsRemovePost`: ResponseProjectsRemove
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunProjectsRemovePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunProjectsRemovePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestProjectsRemove** | [**RequestProjectsRemove**](RequestProjectsRemove.md) |  | 

### Return type

[**ResponseProjectsRemove**](ResponseProjectsRemove.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunProjectsUpdateEditableFieldsPost

> ResponseProjectsUpdateEditableFields ApiRunProjectsUpdateEditableFieldsPost(ctx).RequestProjectsUpdateEditableFields(requestProjectsUpdateEditableFields).Execute()

Calls a mutation at the path projects.js:updateEditableFields

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestProjectsUpdateEditableFields := *openapiclient.NewRequestProjectsUpdateEditableFields(*openapiclient.NewRequestProjectsUpdateEditableFieldsArgs("Id_example")) // RequestProjectsUpdateEditableFields | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunProjectsUpdateEditableFieldsPost(context.Background()).RequestProjectsUpdateEditableFields(requestProjectsUpdateEditableFields).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunProjectsUpdateEditableFieldsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunProjectsUpdateEditableFieldsPost`: ResponseProjectsUpdateEditableFields
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunProjectsUpdateEditableFieldsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunProjectsUpdateEditableFieldsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestProjectsUpdateEditableFields** | [**RequestProjectsUpdateEditableFields**](RequestProjectsUpdateEditableFields.md) |  | 

### Return type

[**ResponseProjectsUpdateEditableFields**](ResponseProjectsUpdateEditableFields.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunProjectsUpdatePositionPost

> ResponseProjectsUpdatePosition ApiRunProjectsUpdatePositionPost(ctx).RequestProjectsUpdatePosition(requestProjectsUpdatePosition).Execute()

Calls a mutation at the path projects.js:updatePosition

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestProjectsUpdatePosition := *openapiclient.NewRequestProjectsUpdatePosition(*openapiclient.NewRequestProjectsUpdatePositionArgs("Id_example", float32(123), float32(123))) // RequestProjectsUpdatePosition | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunProjectsUpdatePositionPost(context.Background()).RequestProjectsUpdatePosition(requestProjectsUpdatePosition).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunProjectsUpdatePositionPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunProjectsUpdatePositionPost`: ResponseProjectsUpdatePosition
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunProjectsUpdatePositionPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunProjectsUpdatePositionPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestProjectsUpdatePosition** | [**RequestProjectsUpdatePosition**](RequestProjectsUpdatePosition.md) |  | 

### Return type

[**ResponseProjectsUpdatePosition**](ResponseProjectsUpdatePosition.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunProjectsUpdateStatsPost

> ResponseProjectsUpdateStats ApiRunProjectsUpdateStatsPost(ctx).RequestProjectsUpdateStats(requestProjectsUpdateStats).Execute()

Calls a mutation at the path projects.js:updateStats

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestProjectsUpdateStats := *openapiclient.NewRequestProjectsUpdateStats(*openapiclient.NewRequestProjectsUpdateStatsArgs(float32(123), "Id_example", float32(123))) // RequestProjectsUpdateStats | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunProjectsUpdateStatsPost(context.Background()).RequestProjectsUpdateStats(requestProjectsUpdateStats).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunProjectsUpdateStatsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunProjectsUpdateStatsPost`: ResponseProjectsUpdateStats
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunProjectsUpdateStatsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunProjectsUpdateStatsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestProjectsUpdateStats** | [**RequestProjectsUpdateStats**](RequestProjectsUpdateStats.md) |  | 

### Return type

[**ResponseProjectsUpdateStats**](ResponseProjectsUpdateStats.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunSeedSeedDatabasePost

> ResponseSeedSeedDatabase ApiRunSeedSeedDatabasePost(ctx).RequestSeedSeedDatabase(requestSeedSeedDatabase).Execute()

Calls a mutation at the path seed.js:seedDatabase

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestSeedSeedDatabase := *openapiclient.NewRequestSeedSeedDatabase(map[string]interface{}(123)) // RequestSeedSeedDatabase | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunSeedSeedDatabasePost(context.Background()).RequestSeedSeedDatabase(requestSeedSeedDatabase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunSeedSeedDatabasePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunSeedSeedDatabasePost`: ResponseSeedSeedDatabase
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunSeedSeedDatabasePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunSeedSeedDatabasePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestSeedSeedDatabase** | [**RequestSeedSeedDatabase**](RequestSeedSeedDatabase.md) |  | 

### Return type

[**ResponseSeedSeedDatabase**](ResponseSeedSeedDatabase.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunSettingsUpdateSettingsPost

> ResponseSettingsUpdateSettings ApiRunSettingsUpdateSettingsPost(ctx).RequestSettingsUpdateSettings(requestSettingsUpdateSettings).Execute()

Calls a mutation at the path settings.js:updateSettings

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestSettingsUpdateSettings := *openapiclient.NewRequestSettingsUpdateSettings(*openapiclient.NewRequestSettingsUpdateSettingsArgs("UserId_example")) // RequestSettingsUpdateSettings | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunSettingsUpdateSettingsPost(context.Background()).RequestSettingsUpdateSettings(requestSettingsUpdateSettings).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunSettingsUpdateSettingsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunSettingsUpdateSettingsPost`: ResponseSettingsUpdateSettings
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunSettingsUpdateSettingsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunSettingsUpdateSettingsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestSettingsUpdateSettings** | [**RequestSettingsUpdateSettings**](RequestSettingsUpdateSettings.md) |  | 

### Return type

[**ResponseSettingsUpdateSettings**](ResponseSettingsUpdateSettings.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunSettingsUpdateThemePost

> ResponseSettingsUpdateTheme ApiRunSettingsUpdateThemePost(ctx).RequestSettingsUpdateTheme(requestSettingsUpdateTheme).Execute()

Calls a mutation at the path settings.js:updateTheme

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestSettingsUpdateTheme := *openapiclient.NewRequestSettingsUpdateTheme(*openapiclient.NewRequestSettingsUpdateThemeArgs(openapiclient.Request_settings_updateSettings_args_theme{String: new(string)}, "UserId_example")) // RequestSettingsUpdateTheme | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunSettingsUpdateThemePost(context.Background()).RequestSettingsUpdateTheme(requestSettingsUpdateTheme).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunSettingsUpdateThemePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunSettingsUpdateThemePost`: ResponseSettingsUpdateTheme
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunSettingsUpdateThemePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunSettingsUpdateThemePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestSettingsUpdateTheme** | [**RequestSettingsUpdateTheme**](RequestSettingsUpdateTheme.md) |  | 

### Return type

[**ResponseSettingsUpdateTheme**](ResponseSettingsUpdateTheme.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunTasksCreatePost

> ResponseTasksCreate ApiRunTasksCreatePost(ctx).RequestTasksCreate(requestTasksCreate).Execute()

Calls a mutation at the path tasks.js:create

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestTasksCreate := *openapiclient.NewRequestTasksCreate(*openapiclient.NewRequestTasksCreateArgs(float32(123), float32(123), "Description_example", "ProjectId_example", "Title_example")) // RequestTasksCreate | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunTasksCreatePost(context.Background()).RequestTasksCreate(requestTasksCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunTasksCreatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunTasksCreatePost`: ResponseTasksCreate
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunTasksCreatePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunTasksCreatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestTasksCreate** | [**RequestTasksCreate**](RequestTasksCreate.md) |  | 

### Return type

[**ResponseTasksCreate**](ResponseTasksCreate.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunTasksRemovePost

> ResponseTasksRemove ApiRunTasksRemovePost(ctx).RequestTasksRemove(requestTasksRemove).Execute()

Calls a mutation at the path tasks.js:remove

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestTasksRemove := *openapiclient.NewRequestTasksRemove(*openapiclient.NewRequestAgentsGetArgs("Id_example")) // RequestTasksRemove | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunTasksRemovePost(context.Background()).RequestTasksRemove(requestTasksRemove).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunTasksRemovePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunTasksRemovePost`: ResponseTasksRemove
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunTasksRemovePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunTasksRemovePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestTasksRemove** | [**RequestTasksRemove**](RequestTasksRemove.md) |  | 

### Return type

[**ResponseTasksRemove**](ResponseTasksRemove.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunTasksUpdateEditableFieldsPost

> ResponseTasksUpdateEditableFields ApiRunTasksUpdateEditableFieldsPost(ctx).RequestTasksUpdateEditableFields(requestTasksUpdateEditableFields).Execute()

Calls a mutation at the path tasks.js:updateEditableFields

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestTasksUpdateEditableFields := *openapiclient.NewRequestTasksUpdateEditableFields(*openapiclient.NewRequestProjectsUpdateEditableFieldsArgs("Id_example")) // RequestTasksUpdateEditableFields | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunTasksUpdateEditableFieldsPost(context.Background()).RequestTasksUpdateEditableFields(requestTasksUpdateEditableFields).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunTasksUpdateEditableFieldsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunTasksUpdateEditableFieldsPost`: ResponseTasksUpdateEditableFields
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunTasksUpdateEditableFieldsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunTasksUpdateEditableFieldsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestTasksUpdateEditableFields** | [**RequestTasksUpdateEditableFields**](RequestTasksUpdateEditableFields.md) |  | 

### Return type

[**ResponseTasksUpdateEditableFields**](ResponseTasksUpdateEditableFields.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunTasksUpdatePositionPost

> ResponseTasksUpdatePosition ApiRunTasksUpdatePositionPost(ctx).RequestTasksUpdatePosition(requestTasksUpdatePosition).Execute()

Calls a mutation at the path tasks.js:updatePosition

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestTasksUpdatePosition := *openapiclient.NewRequestTasksUpdatePosition(*openapiclient.NewRequestProjectsUpdatePositionArgs("Id_example", float32(123), float32(123))) // RequestTasksUpdatePosition | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunTasksUpdatePositionPost(context.Background()).RequestTasksUpdatePosition(requestTasksUpdatePosition).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunTasksUpdatePositionPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunTasksUpdatePositionPost`: ResponseTasksUpdatePosition
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunTasksUpdatePositionPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunTasksUpdatePositionPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestTasksUpdatePosition** | [**RequestTasksUpdatePosition**](RequestTasksUpdatePosition.md) |  | 

### Return type

[**ResponseTasksUpdatePosition**](ResponseTasksUpdatePosition.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunTasksUpdateStatePost

> ResponseTasksUpdateState ApiRunTasksUpdateStatePost(ctx).RequestTasksUpdateState(requestTasksUpdateState).Execute()

Calls a mutation at the path tasks.js:updateState

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/denkhaus/agents/convex"
)

func main() {
	requestTasksUpdateState := *openapiclient.NewRequestTasksUpdateState(*openapiclient.NewRequestTasksUpdateStateArgs("Id_example", openapiclient.Request_tasks_updateState_args_state{String: new(string)})) // RequestTasksUpdateState | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MutationAPI.ApiRunTasksUpdateStatePost(context.Background()).RequestTasksUpdateState(requestTasksUpdateState).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MutationAPI.ApiRunTasksUpdateStatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunTasksUpdateStatePost`: ResponseTasksUpdateState
	fmt.Fprintf(os.Stdout, "Response from `MutationAPI.ApiRunTasksUpdateStatePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunTasksUpdateStatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestTasksUpdateState** | [**RequestTasksUpdateState**](RequestTasksUpdateState.md) |  | 

### Return type

[**ResponseTasksUpdateState**](ResponseTasksUpdateState.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

