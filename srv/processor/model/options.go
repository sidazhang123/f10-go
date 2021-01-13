package model

type Opts struct {
	Flags         string `json:"flags"`
	PluginExcl    string `json:"plugin_excl"`
	DbName        string `json:"db_name"`
	RawDbName     string `json:"raw_db_name"`
	PluginSrcPath string `json:"plugin_src_path"`
	PluginSoPath  string `json:"plugin_so_path"`
	PluginGoPath  string `json:"plugin_debug_path"`
	PluginLevel   string `json:"plugin_level"`
}
