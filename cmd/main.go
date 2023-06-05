package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FormData struct {
	URL         string
	Duration    int
	Concurrency int
}

func main() {
	log.Print("Starting server")
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/submit", handleSubmit)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	form := `<html>
	<head>
		<title>HTML Form Example</title>
	</head>
	<body>
		<form method="POST" action="/submit">
			<label for="url">URL:</label>
			<input type="text" id="url" name="url" required value="http://haproxy:6001/api"><br><br>

			<label for="duration">Duration:</label>
			<input type="number" id="duration" name="duration" required value=1><br><br>

			<label for="concurrency">Concurrency:</label>
			<input type="number" id="concurrency" name="concurrency" required value=1><br><br>

			<label for="hosts">Choose a host:</label>
			<select name="host" id="host">
  				<optgroup label="Single Host">
    				<option value="worker1">worker1</option>
    				<option value="worker2">worker2</option>
  				</optgroup>
				<optgroup label="Multiple">
    				<option value="all">all</option>
  				</optgroup>
			</select>
			<input type="submit" value="Submit">
		</form>
	</body>
	</html>`

	w.Write([]byte(form))
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	log.Print("Start Submit")
	if r.Method == "POST" {
		url := r.FormValue("url")
		sDuration := r.FormValue("duration")
		sConcurrency := r.FormValue("concurrency")
		sHosts := r.FormValue("host")

		duration, _ := strconv.Atoi(sDuration)
		concurrency, _ := strconv.Atoi(sConcurrency)

		formData := FormData{
			URL:         url,
			Duration:    duration,
			Concurrency: concurrency,
		}

		marshalled, err := json.Marshal(formData)
		if err != nil {
			log.Fatalf("impossible to marshall teacher: %s", err)
		}

		hosts := []string{}
		if sHosts == "all" {
			hosts = append(hosts, "worker1")
			hosts = append(hosts, "worker2")
		} else {
			hosts = append(hosts, sHosts)
		}
		for _, host := range hosts {

			req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:8081/run", host), bytes.NewReader(marshalled))
			if err != nil {
				log.Fatalf("impossible to build request: %s", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// create http client
			// do not forget to set timeout; otherwise, no timeout!
			client := http.Client{Timeout: 10 * time.Second}
			// send the request
			res, err := client.Do(req)
			if err != nil {
				log.Fatalf("impossible to send request: %s", err)
			}
			log.Printf("status Code: %d", res.StatusCode)

			// we do not forget to close the body to free resources
			// defer will execute that at the end of the current function
			defer res.Body.Close()
			log.Print("End Submit")
		}
	} else {
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
	}
}
