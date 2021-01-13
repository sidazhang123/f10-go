package model

type Opts struct {
	SkipDays               int    `json:"skip_days"`
	RefinedDbName          string `json:"refined_db_name"`
	FeedDbName             string `json:"feed_db_name"`
	RulesCollectionName    string `json:"rules_collection_name"`
	FocusCollectionName    string `json:"focus_collection_name"`
	OutdatedCollectionName string `json:"outdated_collection_name"`
	JPush0                 string `json:"j_push_0"`
	JPush1                 string `json:"j_push_1"`
	LocalFilePathPrefix    string `json:"local_filepath_prefix"`
}
