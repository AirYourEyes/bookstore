package store

type Book struct {
	Id      string   `json:"id"`      // 图书 ID
	Name    string   `json:"name"`    // 图书名
	Authors []string `json:"authors"` // 作者
	Press   string   `json:"press"`   // 出版社
}

type Store interface {
	Create(*Book) error       // 添加图书
	Update(*Book) error       // 更新图书
	Get(string) (Book, error) // 获取图书详情
	GetAll() ([]Book, error)  // 获取所有图书
	Delete(string) error      // 删除图书
}
