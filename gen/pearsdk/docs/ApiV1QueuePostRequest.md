# ApiV1QueuePostRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**VideoId** | **string** |  | 
**InsertPosition** | Pointer to **string** |  | [optional] [default to "INSERT_AT_END"]

## Methods

### NewApiV1QueuePostRequest

`func NewApiV1QueuePostRequest(videoId string, ) *ApiV1QueuePostRequest`

NewApiV1QueuePostRequest instantiates a new ApiV1QueuePostRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiV1QueuePostRequestWithDefaults

`func NewApiV1QueuePostRequestWithDefaults() *ApiV1QueuePostRequest`

NewApiV1QueuePostRequestWithDefaults instantiates a new ApiV1QueuePostRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVideoId

`func (o *ApiV1QueuePostRequest) GetVideoId() string`

GetVideoId returns the VideoId field if non-nil, zero value otherwise.

### GetVideoIdOk

`func (o *ApiV1QueuePostRequest) GetVideoIdOk() (*string, bool)`

GetVideoIdOk returns a tuple with the VideoId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVideoId

`func (o *ApiV1QueuePostRequest) SetVideoId(v string)`

SetVideoId sets VideoId field to given value.


### GetInsertPosition

`func (o *ApiV1QueuePostRequest) GetInsertPosition() string`

GetInsertPosition returns the InsertPosition field if non-nil, zero value otherwise.

### GetInsertPositionOk

`func (o *ApiV1QueuePostRequest) GetInsertPositionOk() (*string, bool)`

GetInsertPositionOk returns a tuple with the InsertPosition field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInsertPosition

`func (o *ApiV1QueuePostRequest) SetInsertPosition(v string)`

SetInsertPosition sets InsertPosition field to given value.

### HasInsertPosition

`func (o *ApiV1QueuePostRequest) HasInsertPosition() bool`

HasInsertPosition returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


