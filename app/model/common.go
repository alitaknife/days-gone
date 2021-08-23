package model

// Page 分页实体
type Page struct {
	Current int `p:"current" v:"required|length:1,1000#请输入页数|页数长度为:min到:max位" json:"current" form:"current"`
	Size    int `p:"size" v:"required|length:1,1000#请输入每页大小|每页大小为:min到:max位" json:"size" form:"size"`
}

func (p *Page) Paginate() (limit, offset int) {
	limit = p.Size
	offset = p.Size * (p.Current - 1)
	return limit, offset
}
