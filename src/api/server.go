package api

import (
	"encoding/json"
	"net/http"

	"rce/src/docker"
)

func Serve(PORT string, RUNS_DIR string) {
	// handle the /run endpoint
	http.HandleFunc("/run", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			code := req.FormValue("code")
			language := req.FormValue("language")

			filePath := docker.CreateFile(code, language, RUNS_DIR)

			// convert the output to json
			output, err := json.Marshal(map[string]string{
				"output": string(docker.Run(filePath, language)),
			})
			if err != nil {
				panic(err)
			}

			// set the content type to json, enable CORS and write the output
			res.Header().Set("Content-Type", "application/json")
			res.Header().Set("Access-Control-Allow-Origin", "*")
			res.Write(output)
		case http.MethodGet:
			res.Write([]byte("Hello World"))
		}
	})

	// start the server
	if err := http.ListenAndServe("0.0.0.0:"+PORT, nil); err != nil {
		panic(err)
	}
}
