# API 对齐索引

此文件由 `go run ./cmd/bpi-probe api-doc` 自动生成，请勿手工编辑。

基准：`bpi-rs` `0.3.0`，提交 `94bcf43e46d6e11b55ec4d36d4848692cecd4213`。

对齐进度：**206/206 条契约**，覆盖 **27 个领域**。

| 风险级别 | 契约数 |
| --- | ---: |
| `public-read` | 115 |
| `authenticated-read` | 32 |
| `private-read` | 52 |
| `login-session` | 7 |

## activity

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `activity.info` | `public-read` | `GET api.bilibili.com/x/activity/subject/info` | `ActivityClient.Info` | `activity.InfoParams` | `activity.Info` |
| `activity.list` | `public-read` | `GET api.bilibili.com/x/activity/page/list` | `ActivityClient.List` | `activity.ListParams` | `activity.ListPage` |

## article

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `article.articles_info` | `public-read` | `GET api.bilibili.com/x/article/list/web/articles` | `ArticleClient.Articles` | `article.ArticlesInfoParams` | `article.Articles` |
| `article.cards` | `authenticated-read` | `GET api.bilibili.com/x/article/cards` | `ArticleClient.Cards` | `article.CardsParams` | `article.CardData` |
| `article.info` | `public-read` | `GET api.bilibili.com/x/article/viewinfo` | `ArticleClient.Info` | `article.InfoParams` | `article.Info` |
| `article.view` | `authenticated-read` | `GET api.bilibili.com/x/article/view` | `ArticleClient.View` | `article.ViewParams` | `article.View` |

## audio

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `audio.coin_count` | `authenticated-read` | `GET www.bilibili.com/audio/music-service-c/web/coin/audio` | `AudioClient.CoinCount` | `audio.SongParams` | `int32` |
| `audio.collection_info` | `authenticated-read` | `GET www.bilibili.com/audio/music-service-c/web/collections/info` | `AudioClient.CollectionInfo` | `audio.CollectionInfoParams` | `*audio.Collection` |
| `audio.collection_status` | `authenticated-read` | `GET www.bilibili.com/audio/music-service-c/web/collections/songs-coll` | `AudioClient.CollectionStatus` | `audio.SongParams` | `bool` |
| `audio.collections_list` | `authenticated-read` | `GET www.bilibili.com/audio/music-service-c/web/collections/list` | `AudioClient.CollectionsList` | `audio.PageParams` | `audio.Page[audio.Collection]` |
| `audio.hot_menu` | `public-read` | `GET www.bilibili.com/audio/music-service-c/web/menu/hit` | `AudioClient.HotMenu` | `audio.PageParams` | `audio.Page[audio.HotMenu]` |
| `audio.info` | `public-read` | `GET www.bilibili.com/audio/music-service-c/web/song/info` | `AudioClient.Info` | `audio.SongParams` | `audio.Info` |
| `audio.lyric` | `public-read` | `GET www.bilibili.com/audio/music-service-c/web/song/lyric` | `AudioClient.Lyric` | `audio.SongParams` | `string` |
| `audio.members` | `public-read` | `GET www.bilibili.com/audio/music-service-c/web/member/song` | `AudioClient.Members` | `audio.SongParams` | `[]audio.MemberGroup` |
| `audio.rank_detail` | `public-read` | `GET api.bilibili.com/x/copyright-music-publicity/toplist/detail` | `AudioClient.RankDetail` | `audio.RankListParams` | `audio.RankDetail` |
| `audio.rank_menu` | `public-read` | `GET www.bilibili.com/audio/music-service-c/web/menu/rank` | `AudioClient.RankMenu` | `audio.PageParams` | `audio.Page[audio.RankMenu]` |
| `audio.rank_music_list` | `public-read` | `GET api.bilibili.com/x/copyright-music-publicity/toplist/music_list` | `AudioClient.RankMusicList` | `audio.RankListParams` | `audio.RankMusicList` |
| `audio.rank_period` | `public-read` | `GET api.bilibili.com/x/copyright-music-publicity/toplist/all_period` | `AudioClient.RankPeriod` | `audio.RankPeriodParams` | `audio.RankPeriods` |
| `audio.status_number` | `public-read` | `GET www.bilibili.com/audio/music-service-c/web/stat/song` | `AudioClient.StatusNumber` | `audio.SongParams` | `audio.StatusNumber` |
| `audio.stream_url` | `public-read` | `GET api.bilibili.com/audio/music-service-c/url` | `AudioClient.StreamURL` | `audio.StreamURLParams` | `audio.StreamURL` |
| `audio.stream_url_web` | `public-read` | `GET www.bilibili.com/audio/music-service-c/web/url` | `AudioClient.StreamURLWeb` | `audio.StreamURLWebParams` | `audio.StreamURLWeb` |
| `audio.tags` | `public-read` | `GET www.bilibili.com/audio/music-service-c/web/tag/song` | `AudioClient.Tags` | `audio.SongParams` | `[]audio.Tag` |

## bangumi

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `bangumi.info.review_user` | `public-read` | `GET api.bilibili.com/pgc/review/user` | `BangumiClient.Review` | `bangumi.InfoParams` | `bangumi.Info` |
| `bangumi.info.season_detail_by_ep_id` | `public-read` | `GET api.bilibili.com/pgc/view/web/season` | `BangumiClient.EpisodeDetail` | `ids.EpisodeID` | `bangumi.Detail` |
| `bangumi.info.season_detail_by_season_id` | `public-read` | `GET api.bilibili.com/pgc/view/web/season` | `BangumiClient.SeasonDetail` | `ids.SeasonID` | `bangumi.Detail` |
| `bangumi.info.season_section` | `public-read` | `GET api.bilibili.com/pgc/web/season/section` | `BangumiClient.Sections` | `bangumi.SectionsParams` | `bangumi.Sections` |
| `bangumi.playurl` | `public-read` | `GET api.bilibili.com/pgc/player/web/playurl` | `BangumiClient.PlayURL` | `bangumi.PlayURLParams` | `bangumi.PlayURL` |
| `bangumi.timeline` | `public-read` | `GET api.bilibili.com/pgc/web/timeline` | `BangumiClient.Timeline` | `bangumi.TimelineParams` | `[]bangumi.TimelineDay` |

## cheese

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `cheese.info.ep_list` | `public-read` | `GET api.bilibili.com/pugv/view/web/ep/list` | `CheeseClient.EpisodeList` | `cheese.EpisodeListParams` | `cheese.EpisodeList` |
| `cheese.info.season_detail_by_ep_id` | `public-read` | `GET api.bilibili.com/pugv/view/web/season` | `CheeseClient.EpisodeDetail` | `ids.EpisodeID` | `cheese.Course` |
| `cheese.info.season_detail_by_season_id` | `public-read` | `GET api.bilibili.com/pugv/view/web/season` | `CheeseClient.SeasonDetail` | `ids.SeasonID` | `cheese.Course` |
| `cheese.playurl` | `public-read` | `GET api.bilibili.com/pugv/player/web/playurl` | `CheeseClient.PlayURL` | `cheese.PlayURLParams` | `cheese.PlayURL` |

## clientinfo

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `clientinfo.ip` | `public-read` | `GET api.live.bilibili.com/ip_service/v1/ip_service/get_ip_addr` | `ClientInfoClient.IP` | `clientinfo.IPParams` | `clientinfo.IPInfo` |

## comment

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `comment.read.count` | `public-read` | `GET api.bilibili.com/x/v2/reply/count` | `CommentClient.Count` | `comment.CountParams` | `comment.Count` |
| `comment.read.hot` | `public-read` | `GET api.bilibili.com/x/v2/reply/hot` | `CommentClient.Hot` | `comment.HotParams` | `*comment.Hot` |
| `comment.read.list` | `public-read` | `GET api.bilibili.com/x/v2/reply` | `CommentClient.List` | `comment.ListParams` | `comment.List` |
| `comment.read.replies` | `public-read` | `GET api.bilibili.com/x/v2/reply/reply` | `CommentClient.Replies` | `comment.RepliesParams` | `comment.List` |

## creativecenter

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `creativecenter.railgun.electromagnetic_info` | `private-read` | `GET api.bilibili.com/studio/up-rating/v3/rating/info` | `CreativeCenterClient.ElectromagneticInfo` | `none` | `creativecenter.ElectromagneticInfo` |
| `creativecenter.season.aid` | `private-read` | `GET member.bilibili.com/x2/creative/web/season/aid` | `CreativeCenterClient.SeasonByAID` | `creativecenter.SeasonByAIDParams` | `creativecenter.Season` |
| `creativecenter.season.info` | `private-read` | `GET member.bilibili.com/x2/creative/web/season` | `CreativeCenterClient.SeasonInfo` | `creativecenter.SeasonInfoParams` | `creativecenter.SeasonInfo` |
| `creativecenter.season.list` | `private-read` | `GET member.bilibili.com/x2/creative/web/seasons` | `CreativeCenterClient.SeasonList` | `creativecenter.SeasonListParams` | `creativecenter.SeasonList` |
| `creativecenter.season.section` | `private-read` | `GET member.bilibili.com/x2/creative/web/season/section` | `CreativeCenterClient.SeasonSection` | `creativecenter.SectionParams` | `creativecenter.SectionEpisodes` |
| `creativecenter.statistics.archive_compare` | `private-read` | `GET member.bilibili.com/x/web/data/archive_diagnose/compare` | `CreativeCenterClient.ArchiveCompare` | `creativecenter.ArchiveCompareParams` | `creativecenter.ArchiveCompare` |
| `creativecenter.statistics.article_stat` | `private-read` | `GET member.bilibili.com/x/web/data/article` | `CreativeCenterClient.ArticleStat` | `none` | `creativecenter.ArticleStat` |
| `creativecenter.statistics.article_trend` | `private-read` | `GET member.bilibili.com/x/web/data/article/thirty` | `CreativeCenterClient.ArticleTrend` | `creativecenter.ArticleTrendParams` | `[]creativecenter.Trend` |
| `creativecenter.statistics.play_source` | `private-read` | `GET member.bilibili.com/x/web/data/playsource` | `CreativeCenterClient.PlaySource` | `none` | `*creativecenter.PlaySource` |
| `creativecenter.statistics.up_stat` | `private-read` | `GET member.bilibili.com/x/web/index/stat` | `CreativeCenterClient.UpStat` | `none` | `creativecenter.UpStat` |
| `creativecenter.statistics.video_trend` | `private-read` | `GET member.bilibili.com/x/web/data/pandect` | `CreativeCenterClient.VideoTrend` | `creativecenter.VideoTrendParams` | `[]creativecenter.Trend` |
| `creativecenter.statistics.viewer_data` | `private-read` | `GET member.bilibili.com/x/web/data/base` | `CreativeCenterClient.ViewerData` | `none` | `creativecenter.ViewerData` |
| `creativecenter.videos.archive_videos` | `private-read` | `GET member.bilibili.com/x/web/archive/videos` | `CreativeCenterClient.ArchiveVideos` | `creativecenter.ArchiveVideosParams` | `creativecenter.ArchiveVideos` |
| `creativecenter.videos.archives_list` | `private-read` | `GET member.bilibili.com/x2/creative/web/archives/sp` | `CreativeCenterClient.ArchivesList` | `creativecenter.ArchivesListParams` | `creativecenter.ArchivesList` |

## danmaku

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `danmaku.adv.state` | `authenticated-read` | `GET api.bilibili.com/x/dm/adv/state` | `DanmakuClient.AdvState` | `danmaku.AdvStateParams` | `danmaku.AdvState` |
| `danmaku.history.dates` | `authenticated-read` | `GET api.bilibili.com/x/v2/dm/history/index` | `DanmakuClient.HistoryDates` | `danmaku.HistoryDatesParams` | `[]string` |
| `danmaku.history.xml` | `authenticated-read` | `GET api.bilibili.com/x/v2/dm/history` | `DanmakuClient.HistoryXMLBytes` | `danmaku.HistoryBytesParams` | `[]byte` |
| `danmaku.mobile.seg` | `public-read` | `GET api.bilibili.com/x/v2/dm/list/seg.so` | `DanmakuClient.MobileSegment` | `danmaku.SegmentParams` | `[]byte` |
| `danmaku.snapshot` | `public-read` | `GET api.bilibili.com/x/v2/dm/ajax` | `DanmakuClient.Snapshot` | `danmaku.SnapshotParams` | `[]string` |
| `danmaku.thumbup.stats` | `public-read` | `GET api.bilibili.com/x/v2/dm/thumbup/stats` | `DanmakuClient.ThumbupStats` | `danmaku.ThumbupStatsParams` | `danmaku.ThumbupStats` |
| `danmaku.web.history_seg` | `authenticated-read` | `GET api.bilibili.com/x/v2/dm/web/history/seg.so` | `DanmakuClient.WebHistorySegment` | `danmaku.HistoryBytesParams` | `[]byte` |
| `danmaku.web.seg` | `public-read` | `GET api.bilibili.com/x/v2/dm/web/seg.so` | `DanmakuClient.WebSegment` | `danmaku.SegmentParams` | `[]byte` |
| `danmaku.web.seg_wbi` | `public-read` | `GET api.bilibili.com/x/v2/dm/wbi/web/seg.so` | `DanmakuClient.WebSegmentWBI` | `danmaku.SegmentParams` | `[]byte` |
| `danmaku.web.view` | `public-read` | `GET api.bilibili.com/x/v2/dm/web/view` | `DanmakuClient.WebView` | `danmaku.WebViewParams` | `[]byte` |
| `danmaku.xml.comment_xml` | `public-read` | `GET comment.bilibili.com/16546.xml` | `DanmakuClient.XMLList` | `danmaku.XMLListParams` | `danmaku.XML` |
| `danmaku.xml.list_so` | `public-read` | `GET api.bilibili.com/x/v1/dm/list.so` | `DanmakuClient.XMLListSO` | `danmaku.XMLListParams` | `danmaku.XML` |

## dynamic

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `dynamic.detail` | `public-read` | `GET api.bilibili.com/x/polymer/web-dynamic/v1/detail` | `DynamicClient.Detail` | `dynamic.DetailParams` | `dynamic.Detail` |
| `dynamic.detail_forward` | `public-read` | `GET api.bilibili.com/x/polymer/web-dynamic/v1/detail/forward` | `DynamicClient.Forwards` | `dynamic.OffsetParams` | `dynamic.Forwards` |
| `dynamic.detail_forward_item` | `authenticated-read` | `GET api.bilibili.com/x/polymer/web-dynamic/v1/detail/forward/item` | `DynamicClient.ForwardItem` | `dynamic.ItemParams` | `dynamic.ForwardInfo` |
| `dynamic.detail_pic` | `public-read` | `GET api.bilibili.com/x/polymer/web-dynamic/v1/detail/pic` | `DynamicClient.Pictures` | `dynamic.ItemParams` | `[]dynamic.Picture` |
| `dynamic.detail_reaction` | `public-read` | `GET api.bilibili.com/x/polymer/web-dynamic/v1/detail/reaction` | `DynamicClient.Reactions` | `dynamic.OffsetParams` | `dynamic.Reactions` |
| `dynamic.feed_all` | `authenticated-read` | `GET api.bilibili.com/x/polymer/web-dynamic/v1/feed/all` | `DynamicClient.All` | `dynamic.AllParams` | `dynamic.Feed` |
| `dynamic.feed_all_update` | `authenticated-read` | `GET api.bilibili.com/x/polymer/web-dynamic/v1/feed/all/update` | `DynamicClient.CheckNew` | `dynamic.CheckNewParams` | `dynamic.Update` |
| `dynamic.feed_banner` | `public-read` | `GET api.bilibili.com/x/dynamic/feed/dyn/banner` | `DynamicClient.FeedBanner` | `none` | `dynamic.BannerFeed` |
| `dynamic.feed_nav` | `authenticated-read` | `GET api.bilibili.com/x/polymer/web-dynamic/v1/feed/nav` | `DynamicClient.NavFeed` | `dynamic.NavFeedParams` | `dynamic.NavFeed` |
| `dynamic.live_users` | `authenticated-read` | `GET api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/w_live_users` | `DynamicClient.LiveUsers` | `dynamic.LiveUsersParams` | `dynamic.LiveUsers` |
| `dynamic.lottery_notice` | `public-read` | `GET api.vc.bilibili.com/lottery_svr/v1/lottery_svr/lottery_notice` | `DynamicClient.LotteryNotice` | `dynamic.LotteryNoticeParams` | `dynamic.LotteryNotice` |
| `dynamic.recent_up` | `authenticated-read` | `GET api.bilibili.com/x/polymer/web-dynamic/v1/portal` | `DynamicClient.RecentUp` | `none` | `dynamic.RecentUp` |
| `dynamic.up_users` | `authenticated-read` | `GET api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/w_dyn_uplist` | `DynamicClient.UpUsers` | `dynamic.UpUsersParams` | `dynamic.UpUsers` |

## electric

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `electric.charge_follow_info` | `private-read` | `GET api.bilibili.com/x/upower/charge/follow/info` | `ElectricClient.ChargeFollowInfo` | `electric.UpMIDParams` | `electric.FollowInfo` |
| `electric.charge_record` | `private-read` | `GET api.live.bilibili.com/xlive/revenue/v1/guard/getChargeRecord` | `ElectricClient.ChargeRecord` | `electric.ChargeRecordParams` | `electric.ChargeRecord` |
| `electric.month_up_list` | `public-read` | `GET api.bilibili.com/x/ugcpay-rank/elec/month/up` | `ElectricClient.MonthUpList` | `electric.MonthUpListParams` | `electric.MonthUpList` |
| `electric.rank_recent` | `private-read` | `GET member.bilibili.com/x/h5/elec/rank/recent` | `ElectricClient.RankRecent` | `electric.PaginationParams` | `electric.RecentRank` |
| `electric.recharge_list` | `private-read` | `GET pay.bilibili.com/bk/brokerage/listForCustomerRechargeRecord` | `ElectricClient.RechargeList` | `electric.RechargeListParams` | `electric.RechargeList` |
| `electric.remark_detail` | `private-read` | `GET member.bilibili.com/x/web/elec/remark/detail` | `ElectricClient.RemarkDetail` | `electric.RemarkDetailParams` | `electric.RemarkDetail` |
| `electric.remark_list` | `private-read` | `GET member.bilibili.com/x/web/elec/remark/list` | `ElectricClient.RemarkList` | `electric.RemarkListParams` | `electric.RemarkList` |
| `electric.upower_item_detail` | `public-read` | `GET api.bilibili.com/x/upower/item/detail` | `ElectricClient.ItemDetail` | `electric.UpMIDParams` | `electric.ItemDetail` |
| `electric.upower_member_rank` | `public-read` | `GET api.bilibili.com/x/upower/up/member/rank/v2` | `ElectricClient.MemberRank` | `electric.MemberRankParams` | `electric.MemberRank` |
| `electric.video_show` | `public-read` | `GET api.bilibili.com/x/web-interface/elec/show` | `ElectricClient.VideoShow` | `electric.VideoShowParams` | `electric.VideoShow` |

## fav

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `fav.collected_list` | `private-read` | `GET api.bilibili.com/x/v3/fav/folder/collected/list` | `FavClient.CollectedList` | `fav.CollectedListParams` | `fav.CollectedList` |
| `fav.created_list` | `private-read` | `GET api.bilibili.com/x/v3/fav/folder/created/list-all` | `FavClient.CreatedList` | `fav.CreatedListParams` | `fav.CreatedList` |
| `fav.folder_info` | `private-read` | `GET api.bilibili.com/x/v3/fav/folder/info` | `FavClient.FolderInfo` | `fav.FolderInfoParams` | `fav.FolderInfo` |
| `fav.list_detail` | `private-read` | `GET api.bilibili.com/x/v3/fav/resource/list` | `FavClient.ListDetail` | `fav.ListDetailParams` | `fav.ListDetail` |
| `fav.resource_ids` | `private-read` | `GET api.bilibili.com/x/v3/fav/resource/ids` | `FavClient.ResourceIDs` | `fav.ResourceIDsParams` | `[]fav.ResourceID` |
| `fav.resource_infos` | `private-read` | `GET api.bilibili.com/x/v3/fav/resource/infos` | `FavClient.ResourceInfos` | `fav.ResourceInfosParams` | `[]fav.ResourceInfo` |

## historytoview

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `historytoview.history_list` | `private-read` | `GET api.bilibili.com/x/web-interface/history/cursor` | `HistoryToViewClient.HistoryList` | `historytoview.ListParams` | `historytoview.HistoryList` |
| `historytoview.history_shadow` | `private-read` | `GET api.bilibili.com/x/v2/history/shadow` | `HistoryToViewClient.HistoryShadow` | `none` | `bool` |
| `historytoview.toview_list` | `private-read` | `GET api.bilibili.com/x/v2/history/toview` | `HistoryToViewClient.ToViewList` | `none` | `historytoview.ToViewList` |

## live

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `live.area_list` | `public-read` | `GET api.live.bilibili.com/room/v1/Area/getList` | `LiveClient.AreaList` | `none` | `[]live.ParentArea` |
| `live.banned_users` | `private-read` | `GET api.live.bilibili.com/xlive/app-ucenter/v2/xbanned/banned/GetBlackList` | `LiveClient.BannedUsers` | `live.BannedUsersParams` | `live.BannedUsers` |
| `live.blind_gift_info` | `authenticated-read` | `GET api.live.bilibili.com/xlive/general-interface/v1/blindFirstWin/getInfo` | `LiveClient.BlindGiftInfo` | `live.BlindGiftInfoParams` | `live.BlindGiftInfo` |
| `live.danmu_info` | `authenticated-read` | `GET api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo` | `LiveClient.DanmuInfo` | `live.DanmuInfoParams` | `live.DanmuInfo` |
| `live.emoticons` | `authenticated-read` | `GET api.live.bilibili.com/xlive/web-ucenter/v2/emoticon/GetEmoticons` | `LiveClient.Emoticons` | `live.EmoticonsParams` | `live.EmoticonData` |
| `live.follow_up_list` | `private-read` | `GET api.live.bilibili.com/xlive/web-ucenter/user/following` | `LiveClient.FollowUpList` | `live.FollowUpListParams` | `live.FollowUpList` |
| `live.follow_up_web_list` | `private-read` | `GET api.live.bilibili.com/xlive/web-ucenter/v1/xfetter/GetWebList` | `LiveClient.FollowUpWebList` | `live.FollowUpWebListParams` | `live.FollowUpWebList` |
| `live.gift_types` | `authenticated-read` | `GET api.live.bilibili.com/gift/v1/master/getGiftTypes` | `LiveClient.GiftTypes` | `none` | `[]live.GiftType` |
| `live.guard_list` | `public-read` | `GET api.live.bilibili.com/xlive/app-room/v2/guardTab/topListNew` | `LiveClient.GuardList` | `live.GuardListParams` | `live.GuardList` |
| `live.lottery_info` | `authenticated-read` | `GET api.live.bilibili.com/xlive/lottery-interface/v1/lottery/getLotteryInfoWeb` | `LiveClient.LotteryInfo` | `live.LotteryInfoParams` | `live.LotteryInfo` |
| `live.my_medals` | `private-read` | `GET api.live.bilibili.com/xlive/app-ucenter/v1/user/GetMyMedals` | `LiveClient.MyMedals` | `live.MyMedalsParams` | `live.MyMedals` |
| `live.recommend` | `public-read` | `GET api.live.bilibili.com/xlive/web-interface/v1/webMain/getMoreRecList` | `LiveClient.Recommend` | `none` | `live.Recommend` |
| `live.replay_list` | `private-read` | `GET api.live.bilibili.com/xlive/app-blink/v1/anchorVideo/AnchorGetReplayList` | `LiveClient.ReplayList` | `live.ReplayListParams` | `live.ReplayList` |
| `live.room_gift_list` | `public-read` | `GET api.live.bilibili.com/xlive/web-room/v1/giftPanel/roomGiftList` | `LiveClient.RoomGiftList` | `live.RoomGiftListParams` | `live.RoomGiftList` |
| `live.room_info` | `public-read` | `GET api.live.bilibili.com/room/v1/Room/get_info` | `LiveClient.RoomInfo` | `live.RoomInfoParams` | `live.RoomInfo` |
| `live.shield_keywords` | `private-read` | `POST api.live.bilibili.com/xlive/app-ucenter/v1/banned/GetShieldKeywordList` | `LiveClient.ShieldKeywords` | `live.ShieldKeywordsParams` | `live.ShieldKeywords` |
| `live.silent_users` | `private-read` | `POST api.live.bilibili.com/xlive/web-ucenter/v1/banned/GetSilentUserList` | `LiveClient.SilentUsers` | `live.SilentUsersParams` | `live.SilentUsers` |
| `live.stream` | `public-read` | `GET api.live.bilibili.com/room/v1/Room/playUrl` | `LiveClient.Stream` | `live.StreamParams` | `live.Stream` |
| `live.version` | `public-read` | `GET api.live.bilibili.com/xlive/app-blink/v1/liveVersionInfo/getHomePageLiveVersion` | `LiveClient.Version` | `none` | `live.Version` |
| `live.web_heart_beat` | `public-read` | `GET live-trace.bilibili.com/xlive/rdata-interface/v1/heartbeat/webHeartBeat` | `LiveClient.WebHeartBeat` | `live.WebHeartBeatParams` | `live.HeartBeat` |

## login

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `login.account_info` | `private-read` | `GET api.bilibili.com/x/member/web/account` | `LoginClient.AccountInfo` | `none` | `login.AccountInfo` |
| `login.captcha_generate` | `login-session` | `GET passport.bilibili.com/x/passport-login/captcha` | `LoginClient.GenerateCaptcha` | `none` | `login.Captcha` |
| `login.coin` | `private-read` | `GET account.bilibili.com/site/getCoin` | `LoginClient.Coin` | `none` | `login.CoinBalance` |
| `login.daily_reward` | `authenticated-read` | `GET api.bilibili.com/x/member/web/exp/reward` | `LoginClient.DailyReward` | `none` | `login.DailyReward` |
| `login.log` | `private-read` | `GET api.bilibili.com/x/member/web/login/log` | `LoginClient.Log` | `login.LogParams` | `login.Log` |
| `login.nav` | `private-read` | `GET api.bilibili.com/x/web-interface/nav` | `LoginClient.Nav` | `none` | `login.Nav` |
| `login.notice` | `private-read` | `GET api.bilibili.com/x/safecenter/login_notice` | `LoginClient.Notice` | `login.NoticeParams` | `login.Notice` |
| `login.qr.flow` | `login-session` | `FLOW` | `LoginClient.QRFlow` | `none` | `login.QRFlow` |
| `login.qr_generate` | `login-session` | `GET passport.bilibili.com/x/passport-login/web/qrcode/generate` | `LoginClient.GenerateQR` | `none` | `login.QRGenerate` |
| `login.qr_poll` | `login-session` | `GET passport.bilibili.com/x/passport-login/web/qrcode/poll` | `LoginClient.PollQR` | `login.QRPollParams` | `login.QRStatus` |
| `login.stat` | `private-read` | `GET api.bilibili.com/x/web-interface/nav/stat` | `LoginClient.Stat` | `none` | `login.Stats` |
| `login.today_coin_exp` | `private-read` | `GET api.bilibili.com/x/web-interface/coin/today/exp` | `LoginClient.TodayCoinExp` | `none` | `login.TodayCoinExp` |
| `login.vip_info` | `authenticated-read` | `GET api.bilibili.com/x/vip/web/user/info` | `LoginClient.VIPInfo` | `none` | `login.VIPInfo` |

## manga

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `manga.clock_in_info` | `public-read` | `POST manga.bilibili.com/twirp/activity.v1.Activity/GetClockInInfo` | `MangaClient.ClockInInfo` | `none` | `manga.ClockInInfo` |
| `manga.coupons` | `public-read` | `POST manga.bilibili.com/twirp/user.v1.User/GetCoupons` | `MangaClient.Coupons` | `manga.CouponsParams` | `manga.Coupons` |
| `manga.point_products` | `public-read` | `POST manga.bilibili.com/twirp/pointshop.v1.Pointshop/ListProduct` | `MangaClient.PointProducts` | `none` | `[]manga.Product` |
| `manga.season_info` | `public-read` | `POST manga.bilibili.com/twirp/user.v1.Season/GetSeasonInfo` | `MangaClient.SeasonInfo` | `none` | `manga.SeasonInfo` |
| `manga.user_point` | `public-read` | `POST manga.bilibili.com/twirp/pointshop.v1.Pointshop/GetUserPoint` | `MangaClient.UserPoint` | `none` | `manga.UserPoint` |

## message

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `message.reply_feed` | `private-read` | `GET api.bilibili.com/x/msgfeed/reply` | `MessageClient.ReplyFeed` | `message.ReplyFeedParams` | `message.ReplyFeed` |
| `message.single_unread` | `private-read` | `GET api.vc.bilibili.com/session_svr/v1/session_svr/single_unread` | `MessageClient.SingleUnread` | `message.SingleUnreadParams` | `message.SingleUnread` |
| `message.unread_count` | `private-read` | `GET api.vc.bilibili.com/x/im/web/msgfeed/unread` | `MessageClient.UnreadCount` | `message.UnreadCountParams` | `message.UnreadCount` |

## misc

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `misc.b23tv.short_link` | `public-read` | `POST api.biliapi.net/x/share/click` | `MiscClient.ShortLink` | `misc.ShortLinkParams` | `misc.ShortLink` |
| `misc.bili_ticket` | `login-session` | `POST api.bilibili.com/bapis/bilibili.api.ticket.v1.Ticket/GenWebTicket` | `MiscClient.BiliTicket` | `none` | `misc.Ticket` |
| `misc.buvid` | `login-session` | `GET api.bilibili.com/x/frontend/finger/spi` | `MiscClient.Buvid` | `none` | `misc.Buvid` |
| `misc.buvid3` | `login-session` | `GET api.bilibili.com/x/web-frontend/getbuvid` | `MiscClient.Buvid3` | `none` | `misc.Buvid3` |

## note

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `note.archive_list` | `authenticated-read` | `GET api.bilibili.com/x/note/list/archive` | `NoteClient.ArchiveList` | `note.ArchiveListParams` | `note.ArchiveList` |
| `note.is_forbid` | `public-read` | `GET api.bilibili.com/x/note/is_forbid` | `NoteClient.IsForbid` | `note.IsForbidParams` | `note.IsForbid` |
| `note.private_info` | `private-read` | `GET api.bilibili.com/x/note/info` | `NoteClient.PrivateInfo` | `note.PrivateInfoParams` | `note.PrivateInfo` |
| `note.public_archive_list` | `public-read` | `GET api.bilibili.com/x/note/publish/list/archive` | `NoteClient.PublicArchiveList` | `note.PublicArchiveListParams` | `note.PublicArchiveList` |
| `note.public_info` | `public-read` | `GET api.bilibili.com/x/note/publish/info` | `NoteClient.PublicInfo` | `note.PublicInfoParams` | `note.PublicInfo` |
| `note.user_private_list` | `private-read` | `GET api.bilibili.com/x/note/list` | `NoteClient.UserPrivateList` | `note.Pagination` | `note.PrivateList` |
| `note.user_public_list` | `public-read` | `GET api.bilibili.com/x/note/publish/list/user` | `NoteClient.UserPublicList` | `note.Pagination` | `note.PublicUserList` |

## opus

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `opus.space_feed` | `public-read` | `GET api.bilibili.com/x/polymer/web-dynamic/v1/opus/feed/space` | `OpusClient.SpaceFeed` | `opus.SpaceFeedParams` | `opus.SpaceFeed` |

## search

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `search.article` | `public-read` | `GET api.bilibili.com/x/web-interface/wbi/search/type` | `SearchClient.Article` | `search.ArticleParams` | `search.Data[[]search.Article]` |
| `search.bangumi` | `public-read` | `GET api.bilibili.com/x/web-interface/wbi/search/type` | `SearchClient.Bangumi` | `search.BangumiParams` | `search.Data[[]search.Bangumi]` |
| `search.bili_user` | `public-read` | `GET api.bilibili.com/x/web-interface/wbi/search/type` | `SearchClient.Users` | `search.UserParams` | `search.Data[[]search.User]` |
| `search.default` | `public-read` | `GET api.bilibili.com/x/web-interface/wbi/search/default` | `SearchClient.Default` | `none` | `search.Default` |
| `search.hotwords` | `public-read` | `GET s.search.bilibili.com/main/hotword` | `SearchClient.HotWords` | `none` | `search.HotWords` |
| `search.live` | `public-read` | `GET api.bilibili.com/x/web-interface/wbi/search/type` | `SearchClient.Live` | `search.LiveParams` | `search.Data[search.LiveData]` |
| `search.live_room` | `public-read` | `GET api.bilibili.com/x/web-interface/wbi/search/type` | `SearchClient.LiveRooms` | `search.LiveRoomParams` | `search.Data[[]search.LiveRoom]` |
| `search.live_user` | `public-read` | `GET api.bilibili.com/x/web-interface/wbi/search/type` | `SearchClient.LiveUsers` | `search.LiveUserParams` | `search.Data[[]search.LiveUser]` |
| `search.movie` | `public-read` | `GET api.bilibili.com/x/web-interface/wbi/search/type` | `SearchClient.Movies` | `search.MovieParams` | `search.Data[[]search.Movie]` |
| `search.suggest` | `public-read` | `GET s.search.bilibili.com/main/suggest` | `SearchClient.Suggest` | `search.SuggestParams` | `search.Suggest` |
| `search.video` | `public-read` | `GET api.bilibili.com/x/web-interface/wbi/search/type` | `SearchClient.Videos` | `search.VideoParams` | `search.Data[[]search.Video]` |

## user

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `user.album_count` | `public-read` | `GET api.vc.bilibili.com/link_draw/v1/doc/upload_count` | `UserClient.AlbumCount` | `user.AlbumCountParams` | `user.AlbumCount` |
| `user.bangumi_follow_list` | `public-read` | `GET api.bilibili.com/x/space/bangumi/follow/list` | `UserClient.BangumiFollowList` | `user.BangumiFollowListParams` | `user.BangumiFollowList` |
| `user.card` | `public-read` | `GET api.bilibili.com/x/web-interface/card` | `UserClient.Card` | `user.CardParams` | `user.CardProfile` |
| `user.cards` | `authenticated-read` | `GET api.vc.bilibili.com/account/v1/user/cards` | `UserClient.Cards` | `user.CardsParams` | `[]user.BatchCard` |
| `user.follow_tags` | `private-read` | `GET api.bilibili.com/x/relation/tags` | `UserClient.FollowTags` | `none` | `[]user.FollowTag` |
| `user.followers` | `private-read` | `GET api.bilibili.com/x/relation/fans` | `UserClient.Followers` | `user.FollowersParams` | `user.Followers` |
| `user.followings` | `private-read` | `GET api.bilibili.com/x/relation/followings` | `UserClient.Followings` | `user.FollowingsParams` | `user.Followings` |
| `user.infos` | `authenticated-read` | `GET api.vc.bilibili.com/x/im/user_infos` | `UserClient.Infos` | `user.InfosParams` | `[]user.BatchInfo` |
| `user.medal_wall` | `authenticated-read` | `GET api.live.bilibili.com/xlive/web-ucenter/user/MedalWall` | `UserClient.MedalWall` | `user.MedalWallParams` | `user.MedalWall` |
| `user.name_to_uid` | `authenticated-read` | `GET api.bilibili.com/x/polymer/web-dynamic/v1/name-to-uid` | `UserClient.NameToUID` | `user.NameToUIDParams` | `user.NameToUID` |
| `user.nav_stat` | `public-read` | `GET api.bilibili.com/x/space/navnum` | `UserClient.NavStat` | `user.NavStatParams` | `user.NavStat` |
| `user.relation_stat` | `public-read` | `GET api.bilibili.com/x/relation/stat` | `UserClient.RelationStat` | `user.RelationStatParams` | `user.RelationStat` |
| `user.space_info` | `authenticated-read` | `GET api.bilibili.com/x/space/wbi/acc/info` | `UserClient.SpaceInfo` | `user.SpaceParams` | `user.SpaceProfile` |
| `user.space_notice` | `public-read` | `GET api.bilibili.com/x/space/notice` | `UserClient.SpaceNotice` | `user.SpaceNoticeParams` | `user.SpaceNotice` |
| `user.up_stat` | `public-read` | `GET api.bilibili.com/x/space/upstat` | `UserClient.UpStat` | `user.UpStatParams` | `user.UpStat` |
| `user.uploaded_videos` | `public-read` | `GET api.bilibili.com/x/space/wbi/arc/search` | `UserClient.UploadedVideos` | `user.UploadedVideosParams` | `user.UploadedVideos` |

## video

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `video.ai_summary` | `authenticated-read` | `GET api.bilibili.com/x/web-interface/view/conclusion/get` | `VideoClient.AISummary` | `video.AISummaryParams` | `video.AISummary` |
| `video.collection.home_seasons_series` | `public-read` | `GET api.bilibili.com/x/polymer/web-space/home/seasons_series` | `VideoClient.HomeSeasonsSeries` | `video.HomeSeasonsSeriesParams` | `video.SeasonsSeries` |
| `video.collection.seasons_archives_list` | `public-read` | `GET api.bilibili.com/x/polymer/web-space/seasons_archives_list` | `VideoClient.SeasonsArchives` | `video.SeasonsArchivesParams` | `video.SeasonsArchives` |
| `video.collection.seasons_series_list` | `public-read` | `GET api.bilibili.com/x/polymer/web-space/seasons_series_list` | `VideoClient.SeasonsSeries` | `video.SeasonsSeriesParams` | `video.SeasonsSeries` |
| `video.collection.series_archives` | `public-read` | `GET api.bilibili.com/x/series/archives` | `VideoClient.SeriesArchives` | `video.SeriesArchivesParams` | `video.SeriesArchives` |
| `video.collection.series_info` | `public-read` | `GET api.bilibili.com/x/series/series` | `VideoClient.SeriesInfo` | `video.SeriesInfoParams` | `video.SeriesInfo` |
| `video.desc` | `public-read` | `GET api.bilibili.com/x/web-interface/archive/desc` | `VideoClient.Desc` | `video.DescParams` | `string` |
| `video.detail` | `public-read` | `GET api.bilibili.com/x/web-interface/view/detail` | `VideoClient.Detail` | `video.DetailParams` | `video.Detail` |
| `video.homepage_recommendations` | `public-read` | `GET api.bilibili.com/x/web-interface/wbi/index/top/feed/rcmd` | `VideoClient.HomepageRecommendations` | `video.HomepageRecommendationsParams` | `video.HomepageRecommendations` |
| `video.interactive_video_info` | `public-read` | `GET api.bilibili.com/x/stein/edgeinfo_v2` | `VideoClient.InteractiveVideoInfo` | `video.InteractiveInfoParams` | `video.InteractiveInfo` |
| `video.online_total` | `public-read` | `GET api.bilibili.com/x/player/online/total` | `VideoClient.OnlineTotal` | `video.OnlineTotalParams` | `video.OnlineTotal` |
| `video.pagelist` | `public-read` | `GET api.bilibili.com/x/player/pagelist` | `VideoClient.PageList` | `video.PageListParams` | `[]video.Page` |
| `video.play_url` | `public-read` | `GET api.bilibili.com/x/player/wbi/playurl` | `VideoClient.PlayURL` | `video.PlayURLParams` | `video.PlayURL` |
| `video.player_info_v2` | `public-read` | `GET api.bilibili.com/x/player/wbi/v2` | `VideoClient.PlayerInfoV2` | `video.PlayerInfoParams` | `video.PlayerInfo` |
| `video.related_videos` | `public-read` | `GET api.bilibili.com/x/web-interface/archive/related` | `VideoClient.RelatedVideos` | `video.RelatedParams` | `[]video.RelatedVideo` |
| `video.tags` | `public-read` | `GET api.bilibili.com/x/web-interface/view/detail/tag` | `VideoClient.Tags` | `video.TagsParams` | `[]video.VideoTag` |
| `video.view` | `public-read` | `GET api.bilibili.com/x/web-interface/view` | `VideoClient.View` | `video.ViewParams` | `video.View` |

## video_ranking

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `video_ranking.popular_list` | `public-read` | `GET api.bilibili.com/x/web-interface/popular` | `VideoRankingClient.PopularList` | `videoranking.PopularListParams` | `videoranking.PopularList` |
| `video_ranking.popular_precious` | `public-read` | `GET api.bilibili.com/x/web-interface/popular/precious` | `VideoRankingClient.Precious` | `none` | `videoranking.PreciousVideos` |
| `video_ranking.popular_series_list` | `public-read` | `GET api.bilibili.com/x/web-interface/popular/series/list` | `VideoRankingClient.PopularSeriesList` | `none` | `videoranking.PopularSeriesList` |
| `video_ranking.popular_series_one` | `public-read` | `GET api.bilibili.com/x/web-interface/popular/series/one` | `VideoRankingClient.PopularSeries` | `videoranking.PopularSeriesParams` | `videoranking.PopularSeries` |
| `video_ranking.ranking_list` | `public-read` | `GET api.bilibili.com/x/web-interface/ranking/v2` | `VideoRankingClient.RankingList` | `videoranking.RankingListParams` | `videoranking.RankingList` |
| `video_ranking.region_dynamic` | `authenticated-read` | `GET api.bilibili.com/x/web-interface/dynamic/region` | `VideoRankingClient.RegionDynamic` | `videoranking.RegionDynamicParams` | `videoranking.RegionArchives` |
| `video_ranking.region_newlist` | `public-read` | `GET api.bilibili.com/x/web-interface/newlist` | `VideoRankingClient.RegionNewList` | `videoranking.RegionNewListParams` | `videoranking.RegionArchives` |
| `video_ranking.region_newlist_rank` | `public-read` | `GET api.bilibili.com/x/web-interface/newlist_rank` | `VideoRankingClient.RegionNewListRank` | `videoranking.RegionNewListRankParams` | `videoranking.NewListRank` |
| `video_ranking.region_tag_dynamic` | `public-read` | `GET api.bilibili.com/x/web-interface/dynamic/tag` | `VideoRankingClient.RegionTagDynamic` | `videoranking.RegionTagDynamicParams` | `videoranking.RegionArchives` |

## vip

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `vip.center_info` | `public-read` | `GET api.bilibili.com/x/vip/web/vip_center/combine` | `VIPClient.Center` | `vip.CenterParams` | `vip.Center` |

## wallet

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `wallet.info` | `private-read` | `POST pay.bilibili.com/paywallet/wallet/getUserWallet` | `WalletClient.Info` | `wallet.InfoParams` | `wallet.Info` |

## web_widget

| 契约 | 风险级别 | 请求 | Go 方法 | 参数 | 模型 |
| --- | --- | --- | --- | --- | --- |
| `web_widget.header_page` | `public-read` | `GET api.bilibili.com/x/web-show/page/header` | `WebWidgetClient.HeaderPage` | `webwidget.HeaderPageParams` | `webwidget.HeaderData` |
| `web_widget.online` | `public-read` | `GET api.bilibili.com/x/web-interface/online` | `WebWidgetClient.Online` | `none` | `webwidget.OnlineData` |
| `web_widget.region_banner` | `public-read` | `GET api.bilibili.com/x/web-show/region/banner` | `WebWidgetClient.RegionBanner` | `webwidget.RegionBannerParams` | `webwidget.RegionBannerData` |
