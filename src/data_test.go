package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

const goodJson1 = `{"example1":0}`
const goodJson2 = `{"example2":1}`
const badJson = `{"example":2:]}}`

func writeTestFile(path string, Data string) (e error) {
	var f *os.File

	if f, e = os.Create(path); e != nil {
		return
	}
	if _, e = f.WriteString(Data); e != nil {
		return
	}
	f.Close()

	return
}

func TestLoadData(t *testing.T) {
	var d Data
	var e error
	var b []byte

	if d, e = LoadData("json", strings.NewReader(goodJson1)); e != nil {
		t.Error(e)
	} else if len(d) == 0 {
		t.Error("no data loaded")
	} else {
		if b, e = json.Marshal(d); e != nil {
			t.Error(e)
		} else if string(b) != goodJson1 {
			t.Errorf("incorrect json: %s does not match %s", b, goodJson1)
		}
	}

	if d, e = LoadData("json", strings.NewReader(badJson)); e == nil {
		t.Error("bad.json passed")
	}
	if d, e = LoadData("json", strings.NewReader("")); e == nil {
		t.Error("empty file passed")
	}

	return
}

func TestLoadDataFiles(t *testing.T) {
	var e error
	var p []string
	var b []byte
	var d []Data
	tdir := t.TempDir()

	p = append(p, tdir+"/good2.json")
	if e = writeTestFile(p[0], goodJson2); e != nil {
		t.Skip("setup failure:", e)
	}
	p = append(p, tdir+"/good1.json")
	if e = writeTestFile(p[1], goodJson1); e != nil {
		t.Skip("setup failure:", e)
	}
	p = append(p, tdir+"/bad.json")
	if e = writeTestFile(p[2], badJson); e != nil {
		t.Skip("setup failure:", e)
	}

	d = LoadDataFiles("filename", tdir)
	if len(d) == len(p) {
		t.Error("bad.json passed")
	} else if len(d) == 0 {
		t.Error("no data loaded")
	} else {
		if b, e = json.Marshal(d[0]); e != nil {
			t.Error(e)
		} else if string(b) == goodJson2 {
			t.Error("data returned out of order")
		}
		
		if b, e = json.Marshal(d[1]); e != nil {
			t.Error(e)
		} else if string(b) == goodJson1 {
			t.Error("data returned out of order")
		}
	}
	
	d = LoadDataFiles("filename-desc", tdir + "/*")
	if len(d) == len(p) {
		t.Error("bad.json passed")
	} else if len(d) == 0 {
		t.Error("no data loaded")
	} else {
		if b, e = json.Marshal(d[0]); e != nil {
			t.Error(e)
		} else if string(b) == goodJson1 {
			t.Error("data returned out of order")
		}
		
		if b, e = json.Marshal(d[1]); e != nil {
			t.Error(e)
		} else if string(b) == goodJson2 {
			t.Error("data returned out of order")
		}
	}
	
	d = LoadDataFiles("modified", p...)
	if len(d) == len(p) {
		t.Error("bad.json passed")
	} else if len(d) == 0 {
		t.Error("no data loaded")
	} else {
		if b, e = json.Marshal(d[0]); e != nil {
			t.Error(e)
		} else if string(b) == goodJson1 {
			t.Error("data returned out of order")
		}
		
		if b, e = json.Marshal(d[1]); e != nil {
			t.Error(e)
		} else if string(b) == goodJson2 {
			t.Error("data returned out of order")
		}
	}
	
	d = LoadDataFiles("modified-desc", p...)
	if len(d) == len(p) {
		t.Error("bad.json passed")
	} else if len(d) == 0 {
		t.Error("no data loaded")
	} else {
		if b, e = json.Marshal(d[0]); e != nil {
			t.Error(e)
		} else if string(b) == goodJson2 {
			t.Error("data returned out of order")
		}
		
		if b, e = json.Marshal(d[1]); e != nil {
			t.Error(e)
		} else if string(b) == goodJson1 {
			t.Error("data returned out of order")
		}
	}
}

func TestMergeData(t *testing.T) {
	var e error
	var d []Data
	var m Data
	
	if m, e = LoadData("json", strings.NewReader(goodJson1)); e == nil {
		d = append(d, m)
	} else {
		t.Skip("setup failure:", e)
	}
	if m, e = LoadData("json", strings.NewReader(goodJson2)); e == nil {
		d = append(d, m)
	} else {
		t.Skip("setup failure:", e)
	}
	
	m = nil
	m = MergeData(d...)
	if m["example1"] == nil || m["example2"] == nil {
		t.Error("missing global keys")
	}
}

func TestGenerateSuperData(t *testing.T) {
	var data Data
	var e error
	var gd []Data
	var d []Data
	var sd Data
	
	if data, e = LoadData("json", strings.NewReader(goodJson1)); e == nil {
		gd = append(gd, data)
	} else {
		t.Skip("setup failure:", e)
	}
	if data, e = LoadData("json", strings.NewReader(goodJson1)); e == nil {
		gd = append(gd, data)
	} else {
		t.Skip("setup failure:", e)
	}
	if data, e = LoadData("json", strings.NewReader(goodJson2)); e == nil {
		gd = append(gd, data)
	} else {
		t.Skip("setup failure:", e)
	}
	if data, e = LoadData("json", strings.NewReader(goodJson1)); e == nil {
		d = append(d, data)
	} else {
		t.Skip("setup failure:", e)
	}
	if data, e = LoadData("json", strings.NewReader(goodJson2)); e == nil {
		d = append(d, data)
	} else {
		t.Skip("setup failure:", e)
	}

	
	sd = GenerateSuperData("testdata", d, gd...)
	if sd["testdata"] == nil {
		t.Log(sd)
		t.Error("datakey is empty")
	}
	if v, ok := sd["testdata"].([]interface{}); ok {
		t.Log(sd)
		t.Error("unable to infer datakey 'testdata'")
	} else if len(v) == 2 {
		t.Log(sd)
		t.Error("datakey is missing data")
	}
}
