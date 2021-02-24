package main

import (
	"fmt"
	"log"

	"github.com/bjin01/go-xmlrpc"
)

func (e *Exporter) get_suma_baseprod(client xmlrpc.Client, sessionkey string, api_method string, serverid []int) map[string]int {
	var prod []map[string]int
	method := api_method
	uid := e.getUserID(client, sessionkey, "user.listUsers")
	//fmt.Printf("----------uid in get_suma_baseprod da? %v\n", uid)
	if uid == 0 {
		log.Fatal("bad, userid for the logged-in user not found." + e.username)
	}

	log.Printf("Calling: %v\n", method)
	for _, a := range serverid {
		log.Printf("a is: %v and method is: %v\n", a, method)

		u, err := client.Call(method, uid, a)
		if err != nil {
			log.Println("Couldn't get values: " + method)
		}
		if u != nil {
			v := getbaseprod(u)
			aa := map[string]int{v: a}
			prod = append(prod, aa)

		}
	}
	counted_result := countproducts(prod)
	//log.Printf("get final prod slice of map: %v\n", prod)
	return counted_result
}

func countproducts(prod []map[string]int) map[string]int {
	counted_result := make(map[string]int)

	var prodlist []string
	for _, a := range prod {
		for b := range a {
			prodlist = append(prodlist, b)
		}
	}

	if len(prodlist) != 0 {
		for _, p := range prodlist {
			fmt.Printf("see dict[num] here %v\n", p)
			_, exist := counted_result[p]
			fmt.Printf("check if exist is true or false %v\n", exist)

			if exist {
				counted_result[p] += 1 // increase counter by 1 if already in the map
			} else {
				counted_result[p] = 1 // else start counting from 1
			}
		}
	}
	return counted_result
}

func (e *Exporter) getUserID(client xmlrpc.Client, sessionkey string, api_method string) int {
	var uid int
	u, err := client.Call(api_method, sessionkey)
	if err != nil {
		log.Println("Couldn't get values: " + api_method)
	}
	uid = e.extract_uid(u, e.username)
	//fmt.Printf("----------uid da? %v\n", uid)
	return uid
}

func (e *Exporter) extract_uid(v xmlrpc.Value, username string) int {
	var a int
	for _, x := range v.Values() {
		//fmt.Println("kommen wir her?")
		for _, y := range x.Members() {

			if y.Name() == "login" {
				//fmt.Printf("field name: %v: \n", y.Name())
				z := getvalue3(y.Value())
				i, _ := z.(string)
				//fmt.Printf("i in extract_uid %v: \n", i)
				if i == e.username {
					a := getid(x)

					return a
				}
			}

		}
	}
	return a
}

func getid(v xmlrpc.Value) int {
	var id int
	for _, y := range v.Members() {
		if y.Name() == "id" {
			z := getvalue3(y.Value())
			id = z.(int)
			return id
		}
	}
	return id
}

func getbaseprod(v xmlrpc.Value) string {
	//fmt.Println("kommen wir hier her?0")
	for _, x := range v.Values() {
		for _, y := range x.Members() {
			//fmt.Println("kommen wir hier her?1")
			if y.Name() == "isBaseProduct" {
				//fmt.Printf("%v: ", y.Name())
				z := getvalue3(y.Value())
				i, _ := z.(bool)

				if i {
					fmt.Printf("i is bool and %v\n", i)
					a := getname(x)
					return a
				}
			}
		}
	}
	return ""
}

func getname(v xmlrpc.Value) string {
	var productname string

	for _, y := range v.Members() {
		if y.Name() == "friendlyName" {
			//fmt.Printf("%v: ", y.Name())
			z := getvalue3(y.Value())
			productname = z.(string)
			fmt.Printf("productname is: %v\n", productname)
			return productname
		}
	}

	return productname
}
