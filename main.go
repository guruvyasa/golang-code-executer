package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	// "io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type ExecutionRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	Input    string `json:"input"`
}

type OutputResponse struct {
	Err    bool
	Output string
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func RunCode(w http.ResponseWriter, r *http.Request) {
	// done := make(chan bool)
	//cors handling
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if (*r).Method == "OPTIONS" {
		return
	}
	output := make(chan OutputResponse)
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	reqObj := ExecutionRequest{}
	err := decoder.Decode(&reqObj)
	if err != nil {
		http.Error(w, "Bad input given", http.StatusBadRequest)
		return
	}
	//create temp file and store code in it
	dir, _ := os.Getwd()
	// dir := os.TempDir()
	tmpFile, err := ioutil.TempFile(dir, "*.c")
	if err != nil {
		log.Fatal("Error creating tmp file")
		http.Error(w, "Server error!", http.StatusInternalServerError)
	}
	defer os.Remove(tmpFile.Name())

	fmt.Println("created tmp file")

	text := []byte(reqObj.Code)
	fmt.Println(text)
	if _, err = tmpFile.Write(text); err != nil {
		log.Fatal("Failed to write to tmp file")
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	// cmd := exec.Command("gcc", "-v")
	fmt.Println(tmpFile.Name())
	fname := strings.Split(tmpFile.Name(), ".")[0]
	defer os.Remove(fname)

	timer := time.NewTimer(10 * time.Second)
	go func() {
		cmd := exec.Command("gcc", tmpFile.Name(), "-o", fname)
		s, e := cmd.CombinedOutput()
		if e != nil {
			fmt.Println("compiler problem", e, s)
			output <- OutputResponse{Err: true, Output: "Compilation Error!"}
			// http.Error(w, e.Error(), http.StatusInternalServerError)
			return

		}
		// cmd.Wait()
		cmd = exec.Command(fname)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, reqObj.Input)
		}()

		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
		}
		output <- OutputResponse{Err: false, Output: string(out)}
	}()
	// out := s
	var resp OutputResponse
	select {
	case <-timer.C:
		resp = OutputResponse{Err: true, Output: "Timeout error!!"}
		// return
	case resp = <-output:
		fmt.Println(resp)

		// return
	}
	fmt.Println(resp.Output)
	tmpFile.Close()
	encoder.Encode(resp)
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to code runner!!")
	})

	http.HandleFunc("/execute", RunCode)
	http.ListenAndServe(":5000", nil)
}
