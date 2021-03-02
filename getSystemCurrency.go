package main

import (
	"log"
	"sort"

	"github.com/bjin01/go-xmlrpc"
)

type Currency_result struct {
	systemname           string
	serverid             int
	critical_patches     int
	important_patches    int
	moderate_patches     int
	low_security_patches int
	bug_fix_patches      int
	enhancement_patches  int
	total_scores         int
}

func (e *Exporter) makeithapopen(metric_desc string) []Currency_result {
	client := xmlrpc.NewClient(e.suma_server_url)

	f, err := client.Call("auth.login", e.username, e.password)

	if err != nil {
		log.Fatal("Couldn't login to suse manager host.")
	}
	result := get_currency(client, f.String(), "system.getSystemCurrencyScores")
	result_final := getTop10(result)
	/* for _, b := range result_final {
		fmt.Printf("%v: \t%v, criticals: %v, important: %v\n", b.systemname, b.total_scores, b.critical_patches, b.important_patches)
	} */
	client.Call("auth.logout", f.String())
	return result_final
}

func getTop10(s []Currency_result) []Currency_result {
	var new_list []Currency_result

	var tempscores = make(map[int]int)
	for _, k := range s {
		tempscores[k.serverid] = k.total_scores
	}

	sid_list := make([]int, 0, len(tempscores))
	for sid := range tempscores {
		sid_list = append(sid_list, sid)
	}

	sort.Slice(sid_list, func(i, j int) bool {
		return tempscores[sid_list[i]] > tempscores[sid_list[j]]
	})
	var counter int

	for _, sid := range sid_list {
		if counter < 10 {
			for c, d := range s {
				if d.serverid == sid {
					new_list = append(new_list, s[c])
				}
			}
			counter++
		}
	}
	return new_list
}

func get_currency(client xmlrpc.Client, sessionkey string, api_method string) []Currency_result {
	var currency Currency_result
	var finallist []Currency_result
	log.Printf("Calling: %v\n", api_method)

	u, err := client.Call(api_method, sessionkey)
	if err != nil {
		log.Println("Couldn't get values: " + api_method)
	}

	x, err := client.Call("system.listSystems", sessionkey)
	if err != nil {
		log.Println("Couldn't get values: " + api_method)
	}

	abc := getSystemName(x, "name")

	if u != nil && abc != nil {
		for _, b := range u.Values() {
			a := b.Members()
			for _, z := range a {
				switch fielname := z.Name(); fielname {
				case "sid":
					intval := getvalue3(z.Value())
					for h, i := range abc {
						if h == intval.(int) {
							currency.systemname = i
						}
					}
					currency.serverid = intval.(int)
				case "score":
					intval := getvalue3(z.Value())
					currency.total_scores = intval.(int)
				case "mod":
					intval := getvalue3(z.Value())
					currency.moderate_patches = intval.(int)
				case "enh":
					intval := getvalue3(z.Value())
					currency.enhancement_patches = intval.(int)
				case "imp":
					intval := getvalue3(z.Value())
					currency.important_patches = intval.(int)
				case "crit":
					intval := getvalue3(z.Value())
					currency.critical_patches = intval.(int)
				case "low":
					intval := getvalue3(z.Value())
					currency.low_security_patches = intval.(int)
				}
			}
			finallist = append(finallist, currency)
		}
	}
	return finallist
}

func getSystemName(v xmlrpc.Value, fname string) map[int]string {
	var system = make(map[int]string)
	var temp_id int
	var temp_sname string
	for _, x := range v.Values() {
		for _, y := range x.Members() {
			if y.Name() == "id" {

				z := getvalue3(y.Value())
				temp_id, _ = z.(int)
			}
			if y.Name() == fname {
				z := getvalue3(y.Value())
				temp_sname, _ = z.(string)
			}
		}
		if temp_id != 0 && temp_sname != "" {
			system[temp_id] = temp_sname
		}
	}
	return system
}
