package files

type UploadItem struct {
	Source      string
	Destination string
	Note        string
	NoteColor   string
}

type UploadedFile struct {
	File       string
	RemotePath string
}
