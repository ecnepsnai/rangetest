package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
)

const fail = "\033[31m[FAIL]\033[0m"
const pass = "\033[32m[PASS]\033[0m"

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func performTestSuite(url string) {
	type tf func(string) error

	runTest := func(name string, t tf) {
		if err := t(url); err != nil {
			fmt.Printf("%s %s\n       %s\n", fail, name, err.Error())
		} else {
			fmt.Printf("%s %s\n", pass, name)
		}
	}

	runTest("HEAD request", testHEADRequest)
	runTest("Get all data without range", testGetAllDataWithoutRange)
	runTest("Get all data with range", testGetAllDataWithRange)
	runTest("Get single absolute range", testGetSingleAbsoluteRange)
	runTest("Get single relative range with start index", testGetSingleRelativeRangeStart)
	runTest("Get single relative range with end index", testGetSingleRelativeRangeEnd)
	runTest("Get multiple absolute ranges", testGetMultipleAbsoluteRanges)
	runTest("Get multiple absolute and relative ranges", testGetMultipleAbsoluteAndRelativeRanges)
	runTest("Unsupported unit type", testUnsupportedUnitType)
	runTest("Index out of range", testIndexOutOfRange)
}

func testHEADRequest(url string) error {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "rangetest/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}

	if resp.ContentLength != int64(len(sampleData)) {
		return fmt.Errorf("incorrect value of content length header. Expected %d got %d", len(sampleData), resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if ct != "text/plain" {
		return fmt.Errorf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	if value := resp.Header.Get("Accept-Ranges"); value != "bytes" {
		return fmt.Errorf("incorrect or missing Accept-Ranges header. Expected bytes got %s", value)
	}

	return nil
}

func testGetAllDataWithoutRange(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "rangetest/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}

	if resp.ContentLength != int64(len(sampleData)) {
		return fmt.Errorf("incorrect value of content length header. Expected %d got %d", len(sampleData), resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if ct != "text/plain" {
		return fmt.Errorf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(respData) != len(sampleData) {
		return fmt.Errorf("incorrect data length. Expected %d got %d", len(sampleData), len(respData))
	}

	if !bytes.Equal(respData, sampleData) {
		return fmt.Errorf("invalid data returned")
	}

	return nil
}

func testGetAllDataWithRange(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=0-")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 206 {
		return fmt.Errorf("unexpected HTTP status code. Expected %d got %d", 206, resp.StatusCode)
	}

	if resp.ContentLength != int64(len(sampleData)) {
		return fmt.Errorf("incorrect value of content length header. Expected %d got %d", len(sampleData), resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if ct != "text/plain" {
		return fmt.Errorf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	if value := resp.Header.Get("Content-Range"); value != "bytes 0-499/500" {
		return fmt.Errorf("incorrect or missing Accept-Ranges header. Expected 'bytes 0-499/500' got '%s'", value)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(respData) != len(sampleData) {
		return fmt.Errorf("incorrect data length. Expected %d got %d", len(sampleData), len(respData))
	}

	if !bytes.Equal(respData, sampleData) {
		return fmt.Errorf("invalid data returned")
	}

	return nil
}

func testGetSingleAbsoluteRange(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=0-99")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 206 {
		return fmt.Errorf("unexpected HTTP status code. Expected %d got %d", 206, resp.StatusCode)
	}

	if resp.ContentLength != int64(100) {
		return fmt.Errorf("incorrect value of content length header. Expected %d got %d", 100, resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if ct != "text/plain" {
		return fmt.Errorf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	if value := resp.Header.Get("Content-Range"); value != "bytes 0-99/500" {
		return fmt.Errorf("incorrect or missing Accept-Ranges header. Expected 'bytes 0-99/500' got '%s'", value)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(respData) != 100 {
		return fmt.Errorf("incorrect data length. Expected %d got %d", 100, len(respData))
	}

	if !bytes.Equal(respData, sampleData[0:100]) {
		return fmt.Errorf("invalid data returned")
	}

	return nil
}

func testGetSingleRelativeRangeStart(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=400-")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 206 {
		return fmt.Errorf("unexpected HTTP status code. Expected %d got %d", 206, resp.StatusCode)
	}

	if resp.ContentLength != int64(100) {
		return fmt.Errorf("incorrect value of content length header. Expected %d got %d", 100, resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if ct != "text/plain" {
		return fmt.Errorf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	if value := resp.Header.Get("Content-Range"); value != "bytes 400-499/500" {
		return fmt.Errorf("incorrect or missing Accept-Ranges header. Expected 'bytes 400-499/500' got '%s'", value)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(respData) != 100 {
		return fmt.Errorf("incorrect data length. Expected %d got %d", 100, len(respData))
	}

	if !bytes.Equal(respData, sampleData[400:500]) {
		return fmt.Errorf("invalid data returned")
	}

	return nil
}

func testGetSingleRelativeRangeEnd(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=-100")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 206 {
		return fmt.Errorf("unexpected HTTP status code. Expected %d got %d", 206, resp.StatusCode)
	}

	if resp.ContentLength != int64(100) {
		return fmt.Errorf("incorrect value of content length header. Expected %d got %d", 100, resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if ct != "text/plain" {
		return fmt.Errorf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	if value := resp.Header.Get("Content-Range"); value != "bytes 400-499/500" {
		return fmt.Errorf("incorrect or missing Accept-Ranges header. Expected 'bytes 400-499/500' got '%s'", value)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(respData) != 100 {
		return fmt.Errorf("incorrect data length. Expected %d got %d", 100, len(respData))
	}

	if !bytes.Equal(respData, sampleData[400:500]) {
		return fmt.Errorf("invalid data returned")
	}

	return nil
}

func testGetMultipleAbsoluteRanges(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=0-99,200-299,400-499")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 206 {
		return fmt.Errorf("unexpected HTTP status code. Expected %d got %d", 206, resp.StatusCode)
	}

	ct, args, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if ct != "multipart/byteranges" {
		return fmt.Errorf("incorrect value of content typ header. Expected 'multipart/byteranges' got '%s'", ct)
	}
	boundary := args["boundary"]
	if boundary == "" {
		return fmt.Errorf("missing multipart boundary")
	}
	mpReader := multipart.NewReader(resp.Body, boundary)

	expectedContentRangeHeaders := []string{
		"bytes 0-99/500",
		"bytes 200-299/500",
		"bytes 400-499/500",
	}
	expectedData := [][]byte{
		sampleData[0:100],
		sampleData[200:300],
		sampleData[400:500],
	}
	expectedNumberOfParts := 3

	partIdx := 0
	for {
		part, err := mpReader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if partIdx > expectedNumberOfParts-1 {
			return fmt.Errorf("unexpected number of data parts returned. Expected %d but got at least %d", expectedNumberOfParts, partIdx)
		}

		if partType := part.Header.Get("Content-Type"); partType != "text/plain" {
			return fmt.Errorf("invalid content type header value in part %d. Expected 'text/plain' got '%s'", partIdx+1, partType)
		}

		if contentRange := part.Header.Get("Content-Range"); contentRange != expectedContentRangeHeaders[partIdx] {
			return fmt.Errorf("invalid content range header value in part %d. Expected '%s' got '%s'", partIdx+1, expectedContentRangeHeaders[partIdx], contentRange)
		}

		partData, err := io.ReadAll(part)
		if err != nil {
			return fmt.Errorf("error reading data from part %d: %s", partIdx+1, err.Error())
		}

		if !bytes.Equal(partData, expectedData[partIdx]) {
			return fmt.Errorf("invalid data returned in part %d", partIdx+1)
		}

		partIdx++
	}

	return nil
}

func testGetMultipleAbsoluteAndRelativeRanges(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=0-99,-100")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 206 {
		return fmt.Errorf("unexpected HTTP status code. Expected %d got %d", 206, resp.StatusCode)
	}

	ct, args, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if ct != "multipart/byteranges" {
		return fmt.Errorf("incorrect value of content typ header. Expected 'multipart/byteranges' got '%s'", ct)
	}
	boundary := args["boundary"]
	if boundary == "" {
		return fmt.Errorf("missing multipart boundary")
	}
	mpReader := multipart.NewReader(resp.Body, boundary)

	expectedContentRangeHeaders := []string{
		"bytes 0-99/500",
		"bytes 400-499/500",
	}
	expectedData := [][]byte{
		sampleData[0:100],
		sampleData[400:500],
	}
	expectedNumberOfParts := 2

	partIdx := 0
	for {
		part, err := mpReader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if partIdx > expectedNumberOfParts-1 {
			return fmt.Errorf("unexpected number of data parts returned. Expected %d but got at least %d", expectedNumberOfParts, partIdx)
		}

		if partType := part.Header.Get("Content-Type"); partType != "text/plain" {
			return fmt.Errorf("invalid content type header value in part %d. Expected 'text/plain' got '%s'", partIdx+1, partType)
		}

		if contentRange := part.Header.Get("Content-Range"); contentRange != expectedContentRangeHeaders[partIdx] {
			return fmt.Errorf("invalid content range header value in part %d. Expected '%s' got '%s'", partIdx+1, expectedContentRangeHeaders[partIdx], contentRange)
		}

		partData, err := io.ReadAll(part)
		if err != nil {
			return fmt.Errorf("error reading data from part %d: %s", partIdx+1, err.Error())
		}

		if !bytes.Equal(partData, expectedData[partIdx]) {
			return fmt.Errorf("invalid data returned in part %d", partIdx+1)
		}

		partIdx++
	}

	return nil
}

func testUnsupportedUnitType(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "centimeters=0-")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}

	if resp.ContentLength != int64(len(sampleData)) {
		return fmt.Errorf("incorrect value of content length header. Expected %d got %d", len(sampleData), resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if ct != "text/plain" {
		return fmt.Errorf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(respData) != len(sampleData) {
		return fmt.Errorf("incorrect data length. Expected %d got %d", len(sampleData), len(respData))
	}

	if !bytes.Equal(respData, sampleData) {
		return fmt.Errorf("invalid data returned")
	}

	return nil
}

func testIndexOutOfRange(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=700-800")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 416 {
		return fmt.Errorf("unexpected HTTP status code. Expected %d got %d", 416, resp.StatusCode)
	}

	return nil
}
