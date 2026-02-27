# ApiV1SongInfoGet200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Title** | **string** |  | 
**Artist** | **string** |  | 
**Views** | **float32** |  | 
**UploadDate** | Pointer to **string** |  | [optional] 
**ImageSrc** | Pointer to **string** |  | [optional] 
**IsPaused** | Pointer to **bool** |  | [optional] 
**SongDuration** | **float32** |  | 
**ElapsedSeconds** | Pointer to **float32** |  | [optional] 
**Url** | Pointer to **string** |  | [optional] 
**Album** | Pointer to **string** |  | [optional] 
**VideoId** | **string** |  | 
**PlaylistId** | Pointer to **string** |  | [optional] 
**MediaType** | **string** |  | 

## Methods

### NewApiV1SongInfoGet200Response

`func NewApiV1SongInfoGet200Response(title string, artist string, views float32, songDuration float32, videoId string, mediaType string, ) *ApiV1SongInfoGet200Response`

NewApiV1SongInfoGet200Response instantiates a new ApiV1SongInfoGet200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiV1SongInfoGet200ResponseWithDefaults

`func NewApiV1SongInfoGet200ResponseWithDefaults() *ApiV1SongInfoGet200Response`

NewApiV1SongInfoGet200ResponseWithDefaults instantiates a new ApiV1SongInfoGet200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTitle

`func (o *ApiV1SongInfoGet200Response) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *ApiV1SongInfoGet200Response) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *ApiV1SongInfoGet200Response) SetTitle(v string)`

SetTitle sets Title field to given value.


### GetArtist

`func (o *ApiV1SongInfoGet200Response) GetArtist() string`

GetArtist returns the Artist field if non-nil, zero value otherwise.

### GetArtistOk

`func (o *ApiV1SongInfoGet200Response) GetArtistOk() (*string, bool)`

GetArtistOk returns a tuple with the Artist field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetArtist

`func (o *ApiV1SongInfoGet200Response) SetArtist(v string)`

SetArtist sets Artist field to given value.


### GetViews

`func (o *ApiV1SongInfoGet200Response) GetViews() float32`

GetViews returns the Views field if non-nil, zero value otherwise.

### GetViewsOk

`func (o *ApiV1SongInfoGet200Response) GetViewsOk() (*float32, bool)`

GetViewsOk returns a tuple with the Views field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetViews

`func (o *ApiV1SongInfoGet200Response) SetViews(v float32)`

SetViews sets Views field to given value.


### GetUploadDate

`func (o *ApiV1SongInfoGet200Response) GetUploadDate() string`

GetUploadDate returns the UploadDate field if non-nil, zero value otherwise.

### GetUploadDateOk

`func (o *ApiV1SongInfoGet200Response) GetUploadDateOk() (*string, bool)`

GetUploadDateOk returns a tuple with the UploadDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUploadDate

`func (o *ApiV1SongInfoGet200Response) SetUploadDate(v string)`

SetUploadDate sets UploadDate field to given value.

### HasUploadDate

`func (o *ApiV1SongInfoGet200Response) HasUploadDate() bool`

HasUploadDate returns a boolean if a field has been set.

### GetImageSrc

`func (o *ApiV1SongInfoGet200Response) GetImageSrc() string`

GetImageSrc returns the ImageSrc field if non-nil, zero value otherwise.

### GetImageSrcOk

`func (o *ApiV1SongInfoGet200Response) GetImageSrcOk() (*string, bool)`

GetImageSrcOk returns a tuple with the ImageSrc field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageSrc

`func (o *ApiV1SongInfoGet200Response) SetImageSrc(v string)`

SetImageSrc sets ImageSrc field to given value.

### HasImageSrc

`func (o *ApiV1SongInfoGet200Response) HasImageSrc() bool`

HasImageSrc returns a boolean if a field has been set.

### GetIsPaused

`func (o *ApiV1SongInfoGet200Response) GetIsPaused() bool`

GetIsPaused returns the IsPaused field if non-nil, zero value otherwise.

### GetIsPausedOk

`func (o *ApiV1SongInfoGet200Response) GetIsPausedOk() (*bool, bool)`

GetIsPausedOk returns a tuple with the IsPaused field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPaused

`func (o *ApiV1SongInfoGet200Response) SetIsPaused(v bool)`

SetIsPaused sets IsPaused field to given value.

### HasIsPaused

`func (o *ApiV1SongInfoGet200Response) HasIsPaused() bool`

HasIsPaused returns a boolean if a field has been set.

### GetSongDuration

`func (o *ApiV1SongInfoGet200Response) GetSongDuration() float32`

GetSongDuration returns the SongDuration field if non-nil, zero value otherwise.

### GetSongDurationOk

`func (o *ApiV1SongInfoGet200Response) GetSongDurationOk() (*float32, bool)`

GetSongDurationOk returns a tuple with the SongDuration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSongDuration

`func (o *ApiV1SongInfoGet200Response) SetSongDuration(v float32)`

SetSongDuration sets SongDuration field to given value.


### GetElapsedSeconds

`func (o *ApiV1SongInfoGet200Response) GetElapsedSeconds() float32`

GetElapsedSeconds returns the ElapsedSeconds field if non-nil, zero value otherwise.

### GetElapsedSecondsOk

`func (o *ApiV1SongInfoGet200Response) GetElapsedSecondsOk() (*float32, bool)`

GetElapsedSecondsOk returns a tuple with the ElapsedSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetElapsedSeconds

`func (o *ApiV1SongInfoGet200Response) SetElapsedSeconds(v float32)`

SetElapsedSeconds sets ElapsedSeconds field to given value.

### HasElapsedSeconds

`func (o *ApiV1SongInfoGet200Response) HasElapsedSeconds() bool`

HasElapsedSeconds returns a boolean if a field has been set.

### GetUrl

`func (o *ApiV1SongInfoGet200Response) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *ApiV1SongInfoGet200Response) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *ApiV1SongInfoGet200Response) SetUrl(v string)`

SetUrl sets Url field to given value.

### HasUrl

`func (o *ApiV1SongInfoGet200Response) HasUrl() bool`

HasUrl returns a boolean if a field has been set.

### GetAlbum

`func (o *ApiV1SongInfoGet200Response) GetAlbum() string`

GetAlbum returns the Album field if non-nil, zero value otherwise.

### GetAlbumOk

`func (o *ApiV1SongInfoGet200Response) GetAlbumOk() (*string, bool)`

GetAlbumOk returns a tuple with the Album field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlbum

`func (o *ApiV1SongInfoGet200Response) SetAlbum(v string)`

SetAlbum sets Album field to given value.

### HasAlbum

`func (o *ApiV1SongInfoGet200Response) HasAlbum() bool`

HasAlbum returns a boolean if a field has been set.

### GetVideoId

`func (o *ApiV1SongInfoGet200Response) GetVideoId() string`

GetVideoId returns the VideoId field if non-nil, zero value otherwise.

### GetVideoIdOk

`func (o *ApiV1SongInfoGet200Response) GetVideoIdOk() (*string, bool)`

GetVideoIdOk returns a tuple with the VideoId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVideoId

`func (o *ApiV1SongInfoGet200Response) SetVideoId(v string)`

SetVideoId sets VideoId field to given value.


### GetPlaylistId

`func (o *ApiV1SongInfoGet200Response) GetPlaylistId() string`

GetPlaylistId returns the PlaylistId field if non-nil, zero value otherwise.

### GetPlaylistIdOk

`func (o *ApiV1SongInfoGet200Response) GetPlaylistIdOk() (*string, bool)`

GetPlaylistIdOk returns a tuple with the PlaylistId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlaylistId

`func (o *ApiV1SongInfoGet200Response) SetPlaylistId(v string)`

SetPlaylistId sets PlaylistId field to given value.

### HasPlaylistId

`func (o *ApiV1SongInfoGet200Response) HasPlaylistId() bool`

HasPlaylistId returns a boolean if a field has been set.

### GetMediaType

`func (o *ApiV1SongInfoGet200Response) GetMediaType() string`

GetMediaType returns the MediaType field if non-nil, zero value otherwise.

### GetMediaTypeOk

`func (o *ApiV1SongInfoGet200Response) GetMediaTypeOk() (*string, bool)`

GetMediaTypeOk returns a tuple with the MediaType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMediaType

`func (o *ApiV1SongInfoGet200Response) SetMediaType(v string)`

SetMediaType sets MediaType field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


