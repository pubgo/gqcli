package jsonq

type IJsonQ interface {
	// 存储路径
	// 缓存路径
	//	 基于目录结构的方式存储数据
	// 通过缓存的方式加快文件的查询和搜索
	// 每一个目录代表了一张表的信息
	// 缓存信息可以删除，后期通过源文件可以构建出来
	// 包含搜索，关键词搜索，联合搜索，id查询等
	// 内容存储和meta存储分开，内容就是实体的文件，图片，视频等信息
	//cache meta data

	// 保存数据
	Insert(string, map[string]interface{})
	Delete(string, map[string]interface{})
	Update(string, map[string]interface{})
	Query(string, map[string]interface{})
}
