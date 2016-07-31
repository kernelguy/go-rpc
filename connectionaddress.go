package gorpc

import (

)


type ConnectionAddress struct {
	src, dest string
	options interface{}
}

func (this *ConnectionAddress) SetAddress(src, dest string, options interface{}) {
	this.src = src
	this.dest = dest
	this.options = options
}

func (this *ConnectionAddress) Source() string {
	return this.src
}

func (this *ConnectionAddress) Destination() string {
	return this.dest
}

func (this *ConnectionAddress) Options() interface{} {
	return this.options
}
