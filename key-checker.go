package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"strings"

	"github.com/Kong/go-pdk"
	xj "github.com/basgys/goxml2json"
)

type Config struct {
	Apikey string
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Response(kong *pdk.PDK) {
	contentType, err := kong.Response.GetHeader("Content-Type")
	if err != nil {
		kong.Log.Err(err.Error())
		return
	}

	if !(strings.Contains(strings.ToLower(contentType), "application/xml") || strings.Contains(strings.ToLower(contentType), "text/xml")) {
		kong.Log.Err(contentType)
		return
	}

	err = kong.Response.ClearHeader("Content-Length")
	if err != nil {
		kong.Log.Err(err.Error())
		return
	}

	err = kong.Response.SetHeader("Content-Type", "text/json")
	if err != nil {
		kong.Log.Err(err.Error())
		return
	}

	headers, err := kong.Response.GetHeaders(1000)
	if err != nil {
		kong.Log.Err(err.Error())
		return
	}

	/*
		// Marshal the map into a JSON string.
		mJson, err := json.Marshal(headers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		jsonStr := string(mJson)
		kong.Log.Alert(jsonStr)
	*/

	rawBody, err := kong.ServiceResponse.GetRawBody()
	if err != nil {
		kong.Log.Err(err.Error())
		return
	}

	var isGzipEncoded = false
	contentEncoding, err := kong.Response.GetHeader("Content-Encoding")
	if err == nil {
		if strings.Contains(strings.ToLower(contentEncoding), "gzip") {
			isGzipEncoded = true
			gzipBody, err := gzip.NewReader(bytes.NewBuffer([]byte(rawBody)))

			if err != nil {
				kong.Log.Err(err.Error())
				return
			}

			//kong.Log.Alert(gzipBody)
			defer gzipBody.Close()

			output, err := ioutil.ReadAll(gzipBody)
			if err != nil {
				kong.Log.Err(err.Error())
				return
			}

			rawBody = string(output)
		}
	}

	xml := strings.NewReader(rawBody)
	json, err := xj.Convert(xml)
	if err != nil {
		kong.Log.Err(err.Error())
		return
	}

	if isGzipEncoded {
		var jsonBytes bytes.Buffer
		gz := gzip.NewWriter(&jsonBytes)

		if _, err := gz.Write([]byte(json.String())); err != nil {
			kong.Log.Err(err.Error())
			return
		}

		if err := gz.Close(); err != nil {
			kong.Log.Err(err.Error())
			return
		}

		//kong.Log.Alert(jsonBytes.String())
		kong.Response.Exit(200, jsonBytes.String(), headers)
	} else {
		kong.Response.Exit(200, json.String(), headers)
	}
}
