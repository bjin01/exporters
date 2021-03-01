package main

import (
	"log"

	"github.com/bjin01/go-xmlrpc"
)

func (e *Exporter) get_suma_systemid(client xmlrpc.Client, sessionkey string, api_method string) interface{} {
	method := api_method
	log.Printf("Calling: %v\n", method)
	u, err := client.Call(method, sessionkey)
	if err != nil {
		log.Fatal("Couldn't get values: " + method)
	}
	a := getID(u, "id")
	return a
}

func getID(v xmlrpc.Value, fname string) interface{} {
	var id_slice []int
	for _, x := range v.Values() {
		for _, y := range x.Members() {
			if y.Name() == fname {
				z := getvalue3(y.Value())
				i, _ := z.(int)
				id_slice = append(id_slice, i)
			}
		}
	}
	return id_slice
}

func (e *Exporter) get_suma_values(client xmlrpc.Client, sessionkey string, api_method string) interface{} {

	method := api_method
	log.Printf("Calling: %v\n", method)
	u, err := client.Call(method, sessionkey)
	if err != nil {
		log.Fatal("Couldn't get values: " + method)
	}
	a := getVal(u)
	return a
}

func getVal(v xmlrpc.Value) interface{} {
	for _, x := range v.Values() {
		for _, y := range x.Members() {
			getvalue3(y.Value())
		}
	}
	return len(v.Values())
}

func getvalue3(v xmlrpc.Value) interface{} {
	z := v.Kind()
	y := v
	var return_val interface{}

	switch f := z; f {
	case 1:
	case 2:
	case 3:
		return y.Bool()
	case 4:
		return y.Time
	case 5:
	case 6: //this is a int type
		return y.Int()

	case 7: //this is a string type
		return y.String()

	case 8: //this is a member type
	default:
	}
	return return_val

}
