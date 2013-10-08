package controllers

func (ar *ActionResponse) AsJson() *ActionResponse {
	ar.ContentType = "application/javascript; charset=utf-8"
	return ar
}

func (ar *ActionResponse) AsText() *ActionResponse {
	ar.ContentType = "plain/text; charset=utf-8"
	return ar
}

func (ar *ActionResponse) AsHtml() *ActionResponse {
	ar.ContentType = "text/html; charset=utf-8"
	return ar
}

