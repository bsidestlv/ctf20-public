package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang/gddo/httputil/header"
)

func Listen() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}
	fmt.Println("Started serving...")
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	http.HandleFunc("/encrypt", Encrypt)
	http.HandleFunc("/decrypt", Decrypt)
	_, https := os.LookupEnv("HTTPS")
	if https {
		fmt.Println("Starting HTTPS connection on port:", port)
		err := srv.ListenAndServeTLS("server.crt", "server.key")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Starting HTTP connection on port:", port)
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}
}

type Payload struct {
	Message   string
	Key       string
	Signature string
}

type RespE struct {
	Result     string `json:"result"`
	Signature  string
	StatusCode int
}

type RespD struct {
	Result     string `json:"result"`
	StatusCode int
}

var globKey []byte

func Decrypt(w http.ResponseWriter, r *http.Request) {
	var p Payload
	err := decodeJSONBody(w, r, &p)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := &RespD{}
	resp.StatusCode = http.StatusOK

	if len(p.Message) == 0 {
		resp.Result = "Empty message"
		resp.StatusCode = http.StatusBadRequest
	}

	if (len(p.Message) % 20) != 0 {
		resp.Result = "Unpadded message"
		resp.StatusCode = http.StatusBadRequest
	}

	if len(p.Signature) == 0 {
		resp.Result = "signature missing"
		resp.StatusCode = http.StatusBadRequest
	}

	if len(p.Key) != 0 {
		k, err := hex.DecodeString(p.Key)
		if err != nil {
			resp.Result = "bad key format"
			resp.StatusCode = http.StatusBadRequest
		}
		globKey = k
	} else {
		globKey = []byte("I have the same combination on my luggage!")
	}

	if resp.StatusCode == http.StatusOK {
		log.Printf("decrypting message")
		m, errm := hex.DecodeString(p.Message)
		h := hash(p.Message)
		if h != p.Signature {
			resp.Result = "bad signature"
			resp.StatusCode = http.StatusBadRequest
		}
		if errm != nil {
			resp.Result = "corrupted/non-hex message or key"
			resp.StatusCode = http.StatusBadRequest
		}
		if resp.StatusCode == http.StatusOK {
			resp.Result = decryptMsg(m)
		}
	}

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
	}
	json.NewEncoder(w).Encode(resp)
}

func Encrypt(w http.ResponseWriter, r *http.Request) {
	var p Payload
	err := decodeJSONBody(w, r, &p)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	log.Printf("msg: \"%s\"", p.Message)
	w.Header().Set("Content-Type", "application/json")
	resp := &RespE{}
	resp.StatusCode = http.StatusOK

	if len(p.Message) == 0 {
		resp.Result = "Empty message"
		resp.StatusCode = http.StatusBadRequest
	}

	if len(p.Key) != 0 {
		k, err := hex.DecodeString(p.Key)
		if err != nil {
			resp.Result = "bad key format"
			resp.StatusCode = http.StatusBadRequest
		}
		globKey = k
	} else {
		globKey = []byte("I have the same combination on my luggage!")
	}

	if resp.StatusCode == http.StatusOK {
		padMsg := pad(p.Message)
		e := encryptMsg(padMsg)
		resp.Signature = hash(e)
		resp.Result = e
		log.Printf("sig: %s", resp.Signature)
	}

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
	}
	json.NewEncoder(w).Encode(resp)
}

func hash(m string) string {
	h := sha1.New()
	h.Write([]byte(m))
	return hex.EncodeToString(h.Sum(nil)[:2])
}

func decryptMsg(m []byte) string {
	log.Printf("decrypting \"%s\"", hex.EncodeToString(m))
	d := make([]byte, len(m))
	for i := 0; i < len(m); i++ {
		pos := i % len(globKey)
		d[i] = m[i] ^ globKey[pos]
	}
	unpadd := unpad(d)
	log.Printf("decrypted: \"%s\"", unpadd)
	return unpadd
}

func encryptMsg(m []byte) string {
	log.Printf("message: \"%s\"", m)
	e := make([]byte, len(m))
	for i := 0; i < len(m); i++ {
		pos := i % len(globKey)
		e[i] = m[i] ^ globKey[pos]
	}
	res := hex.EncodeToString(e)
	log.Printf("encrypted: \"%s\", %d bytes", res, len(res)/2)
	return res
}

func unpad(msg []byte) string {
	lastChar := int(msg[len(msg)-1])
	return string(msg[:len(msg)-lastChar])
}

func pad(msg string) []byte {
	var block int = 20
	l := len(msg)
	padlen := byte(block - (l % block))
	if padlen == 0 {
		padlen = byte(block)
	}
	res := make([]byte, len(msg)+int(padlen))
	for i := 0; i < (len(msg) + int(padlen)); i++ {
		if i < len(msg) {
			res[i] = msg[i]
		} else {
			res[i] = padlen
		}
	}
	return res
}

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

func test() {
	msg := "this is a test i would like to test"
	pad := pad(msg)
	e := encryptMsg([]byte(pad))
	h := hash(e)
	d, _ := hex.DecodeString(e)
	msg2 := decryptMsg(d)
	if msg != msg2 {
		log.Println("test fail", h)
	}
}

func main() {
	log.SetOutput(os.Stdout)
	Listen()
	//test()
}
