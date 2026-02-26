package storage

func InitStorage(baseURL string) Storage {
	//暂时初始化本地存储，后续可以根据配置切换到OSS等云存储
	return NewLocalStorage("./uploads", baseURL)
}
