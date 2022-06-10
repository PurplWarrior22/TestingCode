package utilities

//a simple image path builder interface

//PathBuilder is an interface which given an image identifier builds the full path to the image
//for simple file system appending this probably seems a bit like overkill, but will be a very
//useful interface if there is ever a request for image storage somewhere other than a local file system
type PathBuilder interface {
	//BuildFullPath takes an image and builds the full path associated with that image
	BuildFullPath(image string) string
}

//implementation of PathBuilder, simply adds a predefined path to the image
type simpleAppender struct {
	//the basePath to append to the image
	basePath string
}

func NewSimpleAppender(basePath string) PathBuilder {
	return &simpleAppender{basePath: basePath}
}

func (appender *simpleAppender) BuildFullPath(image string) string {
	return appender.basePath + image
}
