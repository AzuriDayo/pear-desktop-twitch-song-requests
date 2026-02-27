# ApiV1SearchPostRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Query** | **string** |  | 
**Params** | Pointer to **string** |  | [optional] 
**Continuation** | Pointer to **string** |  | [optional] 

## Methods

### NewApiV1SearchPostRequest

`func NewApiV1SearchPostRequest(query string, ) *ApiV1SearchPostRequest`

NewApiV1SearchPostRequest instantiates a new ApiV1SearchPostRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiV1SearchPostRequestWithDefaults

`func NewApiV1SearchPostRequestWithDefaults() *ApiV1SearchPostRequest`

NewApiV1SearchPostRequestWithDefaults instantiates a new ApiV1SearchPostRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetQuery

`func (o *ApiV1SearchPostRequest) GetQuery() string`

GetQuery returns the Query field if non-nil, zero value otherwise.

### GetQueryOk

`func (o *ApiV1SearchPostRequest) GetQueryOk() (*string, bool)`

GetQueryOk returns a tuple with the Query field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuery

`func (o *ApiV1SearchPostRequest) SetQuery(v string)`

SetQuery sets Query field to given value.


### GetParams

`func (o *ApiV1SearchPostRequest) GetParams() string`

GetParams returns the Params field if non-nil, zero value otherwise.

### GetParamsOk

`func (o *ApiV1SearchPostRequest) GetParamsOk() (*string, bool)`

GetParamsOk returns a tuple with the Params field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParams

`func (o *ApiV1SearchPostRequest) SetParams(v string)`

SetParams sets Params field to given value.

### HasParams

`func (o *ApiV1SearchPostRequest) HasParams() bool`

HasParams returns a boolean if a field has been set.

### GetContinuation

`func (o *ApiV1SearchPostRequest) GetContinuation() string`

GetContinuation returns the Continuation field if non-nil, zero value otherwise.

### GetContinuationOk

`func (o *ApiV1SearchPostRequest) GetContinuationOk() (*string, bool)`

GetContinuationOk returns a tuple with the Continuation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContinuation

`func (o *ApiV1SearchPostRequest) SetContinuation(v string)`

SetContinuation sets Continuation field to given value.

### HasContinuation

`func (o *ApiV1SearchPostRequest) HasContinuation() bool`

HasContinuation returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


