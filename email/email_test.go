package email

type multipartDesc struct {
	ContentType string
	Content     []byte
	Parts       []*multipartDesc
}
