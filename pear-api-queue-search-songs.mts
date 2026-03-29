import something from "./search-podcast-episode_l-txZMzBmbA.json" with { type: "json" };

const EMUSIC_VIDEO_TYPE = {
  ATV: "MUSIC_VIDEO_TYPE_ATV", // Album Music
  OMV: "MUSIC_VIDEO_TYPE_OMV", // Original Music Video
  UGC: "MUSIC_VIDEO_TYPE_UGC", // User Generated Content
  PODCAST_EPISODE: "MUSIC_VIDEO_TYPE_PODCAST_EPISODE", // Podcast
  OTHER_VIDEO: "MUSIC_VIDEO_TYPE_OTHER_VIDEO", // Other
} as const;

type TMusicVideoType =
  (typeof EMUSIC_VIDEO_TYPE)[keyof typeof EMUSIC_VIDEO_TYPE];

const EMUSIC_PAGE_TYPE = {
  USER_CHANNEL: "MUSIC_PAGE_TYPE_USER_CHANNEL",
  ARTIST: "MUSIC_PAGE_TYPE_ARTIST",
} as const;

const contents = something.contents.tabbedSearchResultsRenderer.tabs;

// Check for thumbnail element for videoId

for (const tab of contents) {
  const tabContent = tab.tabRenderer.content.sectionListRenderer.contents;
  if (!tabContent) continue;
  let videoId: string = "";
  for (const content of tabContent) {
    if (content.musicCardShelfRenderer) {
      // This is the main "pushed" result from yt, usually the more popular click
      // Not always a video or music
      let title: string | null = null;
      let artistOrUploader: string | null = null;
      // Get videoId from
      const validRun =
        content.musicCardShelfRenderer.thumbnailOverlay
          .musicItemThumbnailOverlayRenderer; // im assuming musicItemThumbnailOverlayRenderer might be undefined because this is optional when artists show up in music card shelf renderer
      if (!validRun) continue;
      videoId =
        validRun.content.musicPlayButtonRenderer.playNavigationEndpoint
          .watchEndpoint.videoId;
      const artistData = content.musicCardShelfRenderer.subtitle.runs.find(
        (v) => {
          const pageType =
            v.navigationEndpoint?.browseEndpoint
              ?.browseEndpointContextSupportedConfigs
              ?.browseEndpointContextMusicConfig?.pageType;
          if (
            pageType === EMUSIC_PAGE_TYPE.ARTIST ||
            pageType === EMUSIC_PAGE_TYPE.USER_CHANNEL
          ) {
            return true;
          }
        },
      );
      artistOrUploader = artistData ? artistData.text : null;
      if (title) console.log(`${title} - ${artistOrUploader} = ${videoId}`);
    }

    if (content.musicShelfRenderer) {
      // This is the list of other results
      const { contents } = content.musicShelfRenderer;
      for (const content of contents) {
        let mediaTitle = "";
        let videoId = "";
        let artistOrUploader = "";
        let mediaType: TMusicVideoType | null = null;

        if (
          content.musicResponsiveListItemRenderer.overlay
            ?.musicItemThumbnailOverlayRenderer.content.musicPlayButtonRenderer
            .playNavigationEndpoint.watchEndpoint
            ?.watchEndpointMusicSupportedConfigs.watchEndpointMusicConfig
            .musicVideoType === EMUSIC_VIDEO_TYPE.ATV
        )
          mediaType = EMUSIC_VIDEO_TYPE.ATV;
        if (
          content.musicResponsiveListItemRenderer.overlay
            ?.musicItemThumbnailOverlayRenderer.content.musicPlayButtonRenderer
            .playNavigationEndpoint.watchEndpoint
            ?.watchEndpointMusicSupportedConfigs.watchEndpointMusicConfig
            .musicVideoType === EMUSIC_VIDEO_TYPE.UGC
        )
          mediaType = EMUSIC_VIDEO_TYPE.UGC;
        if (
          content.musicResponsiveListItemRenderer.overlay
            ?.musicItemThumbnailOverlayRenderer.content.musicPlayButtonRenderer
            .playNavigationEndpoint.watchEndpoint
            ?.watchEndpointMusicSupportedConfigs.watchEndpointMusicConfig
            .musicVideoType === EMUSIC_VIDEO_TYPE.OMV
        )
          mediaType = EMUSIC_VIDEO_TYPE.OMV;
        if (
          content.musicResponsiveListItemRenderer.overlay
            ?.musicItemThumbnailOverlayRenderer.content.musicPlayButtonRenderer
            .playNavigationEndpoint.watchEndpoint
            ?.watchEndpointMusicSupportedConfigs.watchEndpointMusicConfig
            .musicVideoType === EMUSIC_VIDEO_TYPE.PODCAST_EPISODE
        )
          mediaType = EMUSIC_VIDEO_TYPE.PODCAST_EPISODE;
        if (!mediaType) {
          continue;
        }

        // get media title and artist / uploader
        for (const flexColumn of content.musicResponsiveListItemRenderer
          .flexColumns) {
          for (const run of flexColumn.musicResponsiveListItemFlexColumnRenderer
            .text.runs) {
            // get title
            if (
              (run as any).navigationEndpoint?.watchEndpoint
                ?.watchEndpointMusicSupportedConfigs?.watchEndpointMusicConfig
                ?.musicVideoType === EMUSIC_VIDEO_TYPE.ATV ||
              (run as any).navigationEndpoint?.watchEndpoint
                ?.watchEndpointMusicSupportedConfigs?.watchEndpointMusicConfig
                ?.musicVideoType === EMUSIC_VIDEO_TYPE.UGC ||
              (run as any).navigationEndpoint?.watchEndpoint
                ?.watchEndpointMusicSupportedConfigs?.watchEndpointMusicConfig
                ?.musicVideoType === EMUSIC_VIDEO_TYPE.PODCAST_EPISODE ||
              (run as any).navigationEndpoint?.watchEndpoint
                ?.watchEndpointMusicSupportedConfigs?.watchEndpointMusicConfig
                ?.musicVideoType === EMUSIC_VIDEO_TYPE.OMV
            ) {
              // This is the title text
              mediaTitle = run.text;
              videoId = (run as any).navigationEndpoint?.watchEndpoint?.videoId;
            }
            if (
              (run as any).navigationEndpoint?.browseEndpoint
                ?.browseEndpointContextSupportedConfigs
                ?.browseEndpointContextMusicConfig?.pageType ===
                EMUSIC_PAGE_TYPE.USER_CHANNEL ||
              (run as any).navigationEndpoint?.browseEndpoint
                ?.browseEndpointContextSupportedConfigs
                ?.browseEndpointContextMusicConfig?.pageType ===
                EMUSIC_PAGE_TYPE.ARTIST
            ) {
              // channel name
              artistOrUploader = run.text;
            }
          }
        }
        console.log(`${mediaTitle} - ${artistOrUploader} = ${videoId}`);
      }
    }
  }
}
