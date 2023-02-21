package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
)

const QUOTES_BYTE = 34

type IsPrime bool
type HttpRequestParser struct {
}

func (p HttpRequestParser) ParseJson(r *http.Request) ([]IsPrime, error) {
	var retSlice []IsPrime
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &retSlice)
	return retSlice, err
}

func (slice *IsPrime) UnmarshalJSON(data []byte) error {
	var i int64  = -1
	errMsg := "Element is not a number"
	if data[0] == QUOTES_BYTE {
		n, err := strconv.Atoi(string(data[1 : len(data)-1]))
		i = int64(n)
		if err != nil {
			return errors.New(errMsg)
		}
	} else {
		n, err := strconv.Atoi(string(data))
		i = int64(n)
		if err != nil {
			return errors.New(errMsg)
		}
	}
	*slice = IsPrime(big.NewInt(i).ProbablyPrime(0))
	return nil
}

func readRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":
		parser := HttpRequestParser{}
		retIsPrime, err := parser.ParseJson(r)

		if err != nil {
			count := len(retIsPrime)
			retVal := "the given input is invalid. Element on index " + strconv.Itoa( count ) + " is not a number"
			http.Error(w, retVal, http.StatusInternalServerError)
		} else {
			retVal := fmt.Sprintf("%v", retIsPrime)
			io.WriteString(w, retVal)
		}
	default:
		http.Error(w, "Sorry, only POST method is supported.", http.StatusNotFound)
	}
}

