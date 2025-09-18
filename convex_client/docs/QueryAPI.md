# \QueryAPI

All URIs are relative to *http://127.0.0.1:3210*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiRunAgentsGetPost**](QueryAPI.md#ApiRunAgentsGetPost) | **Post** /api/run/agents/get | Calls a query at the path agents.js:get
[**ApiRunAgentsGetWithTasksPost**](QueryAPI.md#ApiRunAgentsGetWithTasksPost) | **Post** /api/run/agents/getWithTasks | Calls a query at the path agents.js:getWithTasks
[**ApiRunAgentsListByRolePost**](QueryAPI.md#ApiRunAgentsListByRolePost) | **Post** /api/run/agents/listByRole | Calls a query at the path agents.js:listByRole
[**ApiRunAgentsListByStatusPost**](QueryAPI.md#ApiRunAgentsListByStatusPost) | **Post** /api/run/agents/listByStatus | Calls a query at the path agents.js:listByStatus
[**ApiRunAgentsListPost**](QueryAPI.md#ApiRunAgentsListPost) | **Post** /api/run/agents/list | Calls a query at the path agents.js:list
[**ApiRunEventsGetRecentEventsPost**](QueryAPI.md#ApiRunEventsGetRecentEventsPost) | **Post** /api/run/events/getRecentEvents | Calls a query at the path events.js:getRecentEvents
[**ApiRunEventsSubscribeToProjectPost**](QueryAPI.md#ApiRunEventsSubscribeToProjectPost) | **Post** /api/run/events/subscribeToProject | Calls a query at the path events.js:subscribeToProject
[**ApiRunProjectsGetPost**](QueryAPI.md#ApiRunProjectsGetPost) | **Post** /api/run/projects/get | Calls a query at the path projects.js:get
[**ApiRunProjectsGetWithStatsPost**](QueryAPI.md#ApiRunProjectsGetWithStatsPost) | **Post** /api/run/projects/getWithStats | Calls a query at the path projects.js:getWithStats
[**ApiRunProjectsListPost**](QueryAPI.md#ApiRunProjectsListPost) | **Post** /api/run/projects/list | Calls a query at the path projects.js:list
[**ApiRunSettingsGetSettingsPost**](QueryAPI.md#ApiRunSettingsGetSettingsPost) | **Post** /api/run/settings/getSettings | Calls a query at the path settings.js:getSettings
[**ApiRunTasksGetPost**](QueryAPI.md#ApiRunTasksGetPost) | **Post** /api/run/tasks/get | Calls a query at the path tasks.js:get
[**ApiRunTasksListByAgentPost**](QueryAPI.md#ApiRunTasksListByAgentPost) | **Post** /api/run/tasks/listByAgent | Calls a query at the path tasks.js:listByAgent
[**ApiRunTasksListByProjectPost**](QueryAPI.md#ApiRunTasksListByProjectPost) | **Post** /api/run/tasks/listByProject | Calls a query at the path tasks.js:listByProject
[**ApiRunTasksListRootTasksPost**](QueryAPI.md#ApiRunTasksListRootTasksPost) | **Post** /api/run/tasks/listRootTasks | Calls a query at the path tasks.js:listRootTasks



## ApiRunAgentsGetPost

> ResponseAgentsGet ApiRunAgentsGetPost(ctx).RequestAgentsGet(requestAgentsGet).Execute()

Calls a query at the path agents.js:get

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
	requestAgentsGet := *openapiclient.NewRequestAgentsGet(*openapiclient.NewRequestAgentsGetArgs("Id_example")) // RequestAgentsGet | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunAgentsGetPost(context.Background()).RequestAgentsGet(requestAgentsGet).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunAgentsGetPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunAgentsGetPost`: ResponseAgentsGet
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunAgentsGetPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunAgentsGetPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestAgentsGet** | [**RequestAgentsGet**](RequestAgentsGet.md) |  | 

### Return type

[**ResponseAgentsGet**](ResponseAgentsGet.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunAgentsGetWithTasksPost

> ResponseAgentsGetWithTasks ApiRunAgentsGetWithTasksPost(ctx).RequestAgentsGetWithTasks(requestAgentsGetWithTasks).Execute()

Calls a query at the path agents.js:getWithTasks

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
	requestAgentsGetWithTasks := *openapiclient.NewRequestAgentsGetWithTasks(*openapiclient.NewRequestAgentsGetArgs("Id_example")) // RequestAgentsGetWithTasks | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunAgentsGetWithTasksPost(context.Background()).RequestAgentsGetWithTasks(requestAgentsGetWithTasks).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunAgentsGetWithTasksPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunAgentsGetWithTasksPost`: ResponseAgentsGetWithTasks
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunAgentsGetWithTasksPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunAgentsGetWithTasksPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestAgentsGetWithTasks** | [**RequestAgentsGetWithTasks**](RequestAgentsGetWithTasks.md) |  | 

### Return type

[**ResponseAgentsGetWithTasks**](ResponseAgentsGetWithTasks.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunAgentsListByRolePost

> ResponseAgentsListByRole ApiRunAgentsListByRolePost(ctx).RequestAgentsListByRole(requestAgentsListByRole).Execute()

Calls a query at the path agents.js:listByRole

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
	requestAgentsListByRole := *openapiclient.NewRequestAgentsListByRole(*openapiclient.NewRequestAgentsListByRoleArgs(openapiclient.Request_agents_listByRole_args_role{String: new(string)})) // RequestAgentsListByRole | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunAgentsListByRolePost(context.Background()).RequestAgentsListByRole(requestAgentsListByRole).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunAgentsListByRolePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunAgentsListByRolePost`: ResponseAgentsListByRole
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunAgentsListByRolePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunAgentsListByRolePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestAgentsListByRole** | [**RequestAgentsListByRole**](RequestAgentsListByRole.md) |  | 

### Return type

[**ResponseAgentsListByRole**](ResponseAgentsListByRole.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunAgentsListByStatusPost

> ResponseAgentsListByStatus ApiRunAgentsListByStatusPost(ctx).RequestAgentsListByStatus(requestAgentsListByStatus).Execute()

Calls a query at the path agents.js:listByStatus

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
	requestAgentsListByStatus := *openapiclient.NewRequestAgentsListByStatus(*openapiclient.NewRequestAgentsListByStatusArgs(openapiclient.Request_agents_listByStatus_args_status{String: new(string)})) // RequestAgentsListByStatus | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunAgentsListByStatusPost(context.Background()).RequestAgentsListByStatus(requestAgentsListByStatus).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunAgentsListByStatusPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunAgentsListByStatusPost`: ResponseAgentsListByStatus
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunAgentsListByStatusPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunAgentsListByStatusPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestAgentsListByStatus** | [**RequestAgentsListByStatus**](RequestAgentsListByStatus.md) |  | 

### Return type

[**ResponseAgentsListByStatus**](ResponseAgentsListByStatus.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunAgentsListPost

> ResponseAgentsList ApiRunAgentsListPost(ctx).RequestAgentsList(requestAgentsList).Execute()

Calls a query at the path agents.js:list

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
	requestAgentsList := *openapiclient.NewRequestAgentsList(map[string]interface{}(123)) // RequestAgentsList | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunAgentsListPost(context.Background()).RequestAgentsList(requestAgentsList).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunAgentsListPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunAgentsListPost`: ResponseAgentsList
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunAgentsListPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunAgentsListPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestAgentsList** | [**RequestAgentsList**](RequestAgentsList.md) |  | 

### Return type

[**ResponseAgentsList**](ResponseAgentsList.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunEventsGetRecentEventsPost

> ResponseEventsGetRecentEvents ApiRunEventsGetRecentEventsPost(ctx).RequestEventsGetRecentEvents(requestEventsGetRecentEvents).Execute()

Calls a query at the path events.js:getRecentEvents

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
	requestEventsGetRecentEvents := *openapiclient.NewRequestEventsGetRecentEvents(*openapiclient.NewRequestEventsGetRecentEventsArgs()) // RequestEventsGetRecentEvents | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunEventsGetRecentEventsPost(context.Background()).RequestEventsGetRecentEvents(requestEventsGetRecentEvents).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunEventsGetRecentEventsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunEventsGetRecentEventsPost`: ResponseEventsGetRecentEvents
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunEventsGetRecentEventsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunEventsGetRecentEventsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestEventsGetRecentEvents** | [**RequestEventsGetRecentEvents**](RequestEventsGetRecentEvents.md) |  | 

### Return type

[**ResponseEventsGetRecentEvents**](ResponseEventsGetRecentEvents.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunEventsSubscribeToProjectPost

> ResponseEventsSubscribeToProject ApiRunEventsSubscribeToProjectPost(ctx).RequestEventsSubscribeToProject(requestEventsSubscribeToProject).Execute()

Calls a query at the path events.js:subscribeToProject

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
	requestEventsSubscribeToProject := *openapiclient.NewRequestEventsSubscribeToProject(*openapiclient.NewRequestEventsSubscribeToProjectArgs("ProjectId_example")) // RequestEventsSubscribeToProject | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunEventsSubscribeToProjectPost(context.Background()).RequestEventsSubscribeToProject(requestEventsSubscribeToProject).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunEventsSubscribeToProjectPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunEventsSubscribeToProjectPost`: ResponseEventsSubscribeToProject
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunEventsSubscribeToProjectPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunEventsSubscribeToProjectPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestEventsSubscribeToProject** | [**RequestEventsSubscribeToProject**](RequestEventsSubscribeToProject.md) |  | 

### Return type

[**ResponseEventsSubscribeToProject**](ResponseEventsSubscribeToProject.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunProjectsGetPost

> ResponseProjectsGet ApiRunProjectsGetPost(ctx).RequestProjectsGet(requestProjectsGet).Execute()

Calls a query at the path projects.js:get

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
	requestProjectsGet := *openapiclient.NewRequestProjectsGet(*openapiclient.NewRequestAgentsGetArgs("Id_example")) // RequestProjectsGet | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunProjectsGetPost(context.Background()).RequestProjectsGet(requestProjectsGet).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunProjectsGetPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunProjectsGetPost`: ResponseProjectsGet
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunProjectsGetPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunProjectsGetPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestProjectsGet** | [**RequestProjectsGet**](RequestProjectsGet.md) |  | 

### Return type

[**ResponseProjectsGet**](ResponseProjectsGet.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunProjectsGetWithStatsPost

> ResponseProjectsGetWithStats ApiRunProjectsGetWithStatsPost(ctx).RequestProjectsGetWithStats(requestProjectsGetWithStats).Execute()

Calls a query at the path projects.js:getWithStats

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
	requestProjectsGetWithStats := *openapiclient.NewRequestProjectsGetWithStats(*openapiclient.NewRequestAgentsGetArgs("Id_example")) // RequestProjectsGetWithStats | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunProjectsGetWithStatsPost(context.Background()).RequestProjectsGetWithStats(requestProjectsGetWithStats).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunProjectsGetWithStatsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunProjectsGetWithStatsPost`: ResponseProjectsGetWithStats
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunProjectsGetWithStatsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunProjectsGetWithStatsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestProjectsGetWithStats** | [**RequestProjectsGetWithStats**](RequestProjectsGetWithStats.md) |  | 

### Return type

[**ResponseProjectsGetWithStats**](ResponseProjectsGetWithStats.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunProjectsListPost

> ResponseProjectsList ApiRunProjectsListPost(ctx).RequestProjectsList(requestProjectsList).Execute()

Calls a query at the path projects.js:list

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
	requestProjectsList := *openapiclient.NewRequestProjectsList(map[string]interface{}(123)) // RequestProjectsList | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunProjectsListPost(context.Background()).RequestProjectsList(requestProjectsList).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunProjectsListPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunProjectsListPost`: ResponseProjectsList
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunProjectsListPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunProjectsListPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestProjectsList** | [**RequestProjectsList**](RequestProjectsList.md) |  | 

### Return type

[**ResponseProjectsList**](ResponseProjectsList.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunSettingsGetSettingsPost

> ResponseSettingsGetSettings ApiRunSettingsGetSettingsPost(ctx).RequestSettingsGetSettings(requestSettingsGetSettings).Execute()

Calls a query at the path settings.js:getSettings

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
	requestSettingsGetSettings := *openapiclient.NewRequestSettingsGetSettings(*openapiclient.NewRequestSettingsGetSettingsArgs("UserId_example")) // RequestSettingsGetSettings | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunSettingsGetSettingsPost(context.Background()).RequestSettingsGetSettings(requestSettingsGetSettings).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunSettingsGetSettingsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunSettingsGetSettingsPost`: ResponseSettingsGetSettings
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunSettingsGetSettingsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunSettingsGetSettingsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestSettingsGetSettings** | [**RequestSettingsGetSettings**](RequestSettingsGetSettings.md) |  | 

### Return type

[**ResponseSettingsGetSettings**](ResponseSettingsGetSettings.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunTasksGetPost

> ResponseTasksGet ApiRunTasksGetPost(ctx).RequestTasksGet(requestTasksGet).Execute()

Calls a query at the path tasks.js:get

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
	requestTasksGet := *openapiclient.NewRequestTasksGet(*openapiclient.NewRequestAgentsGetArgs("Id_example")) // RequestTasksGet | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunTasksGetPost(context.Background()).RequestTasksGet(requestTasksGet).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunTasksGetPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunTasksGetPost`: ResponseTasksGet
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunTasksGetPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunTasksGetPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestTasksGet** | [**RequestTasksGet**](RequestTasksGet.md) |  | 

### Return type

[**ResponseTasksGet**](ResponseTasksGet.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunTasksListByAgentPost

> ResponseTasksListByAgent ApiRunTasksListByAgentPost(ctx).RequestTasksListByAgent(requestTasksListByAgent).Execute()

Calls a query at the path tasks.js:listByAgent

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
	requestTasksListByAgent := *openapiclient.NewRequestTasksListByAgent(*openapiclient.NewRequestTasksListByAgentArgs("AgentId_example")) // RequestTasksListByAgent | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunTasksListByAgentPost(context.Background()).RequestTasksListByAgent(requestTasksListByAgent).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunTasksListByAgentPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunTasksListByAgentPost`: ResponseTasksListByAgent
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunTasksListByAgentPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunTasksListByAgentPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestTasksListByAgent** | [**RequestTasksListByAgent**](RequestTasksListByAgent.md) |  | 

### Return type

[**ResponseTasksListByAgent**](ResponseTasksListByAgent.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunTasksListByProjectPost

> ResponseTasksListByProject ApiRunTasksListByProjectPost(ctx).RequestTasksListByProject(requestTasksListByProject).Execute()

Calls a query at the path tasks.js:listByProject

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
	requestTasksListByProject := *openapiclient.NewRequestTasksListByProject(*openapiclient.NewRequestTasksListByProjectArgs("ProjectId_example")) // RequestTasksListByProject | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunTasksListByProjectPost(context.Background()).RequestTasksListByProject(requestTasksListByProject).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunTasksListByProjectPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunTasksListByProjectPost`: ResponseTasksListByProject
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunTasksListByProjectPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunTasksListByProjectPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestTasksListByProject** | [**RequestTasksListByProject**](RequestTasksListByProject.md) |  | 

### Return type

[**ResponseTasksListByProject**](ResponseTasksListByProject.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRunTasksListRootTasksPost

> ResponseTasksListRootTasks ApiRunTasksListRootTasksPost(ctx).RequestTasksListRootTasks(requestTasksListRootTasks).Execute()

Calls a query at the path tasks.js:listRootTasks

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
	requestTasksListRootTasks := *openapiclient.NewRequestTasksListRootTasks(*openapiclient.NewRequestTasksListRootTasksArgs("ProjectId_example")) // RequestTasksListRootTasks | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.QueryAPI.ApiRunTasksListRootTasksPost(context.Background()).RequestTasksListRootTasks(requestTasksListRootTasks).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `QueryAPI.ApiRunTasksListRootTasksPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRunTasksListRootTasksPost`: ResponseTasksListRootTasks
	fmt.Fprintf(os.Stdout, "Response from `QueryAPI.ApiRunTasksListRootTasksPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRunTasksListRootTasksPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestTasksListRootTasks** | [**RequestTasksListRootTasks**](RequestTasksListRootTasks.md) |  | 

### Return type

[**ResponseTasksListRootTasks**](ResponseTasksListRootTasks.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

