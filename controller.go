package gorpc

import (
	log "github.com/Sirupsen/logrus"
)

/* Controler objects holds the callable RPC methods.
 * All RPC method names must be prefixed with RPC_.
 * Method parameters can be omitted if not needed. Otherwise there should be only one struct parameter
 * encapsulating the needed parameters. This way named parameter transfers can be supported.
 * e.g, all the following gives a correct parameter declaration for the RPC_Echo method:
 *
 *	1 params := struct{Value string `json:"value"`}{"Hello World"}
 *	2 params := []interface{}{"Hello World"}
 *	3 params := make(map[string]interface{})
 *	  params["value"] = "Hello World"
 *
 * Example 1 uses an anonymous struct to declare a named variable. Notice the json tag renaming Value to value.
 * Example 2 is the fastest and shortest one, but has limited parameter check in the receiver.
 * Example 3 is also possible, but the others are faster.
 */
type Controller struct {
	connection IConnection
}

type EchoParams struct {
	Value string
}
func (this *Controller) RPC_Echo(params EchoParams) string {
	log.Infof("RPC Echo called with: (%T)%s", params.Value, params.Value)
	return params.Value
}


func (this *Controller) Echo(value string) string {
	p := struct{Value string}{Value: value}
	r, err := this.connection.Call("Echo", p)
	if err != nil {
		return err.Error()
	}
	return r.(string)
}