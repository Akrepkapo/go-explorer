	r.Message = str
}
func (r *Response) Return(dat interface{}, ct CodeType) {
	r.Code = ct.Code
	r.Message = ct.Message
	r.Data = dat
}
