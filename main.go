package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"text/tabwriter"
	"text/template"
)

var (
	certFile = flag.String("cert", "", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "", "A PEM encoded private key file.")
	testFile = flag.String("tests", "./tests.json", "A file listing the endpoints to benchmark")
	tmpl     = flag.String("template", "", "A template file that will be used to render results")
	client   *http.Client
)

func main() {
	flag.Parse()

	client = http.DefaultClient
	if *certFile != "" && *keyFile != "" {
		client = tlsClient(*certFile, *keyFile)
	}

	tests, err := loadTests(*testFile)
	if err != nil {
		log.Fatal(err)
	}

	results := map[string]testing.BenchmarkResult{}

	w := tabwriter.NewWriter(os.Stderr, 0, 0, 1, ' ', 0)
	for t, req := range tests {
		results[t] = testing.Benchmark(benchmarkUrlTest(req))
		fmt.Fprintf(w, "%s\t%d\t%s\n", t, results[t].N, results[t].T)
	}
	w.Flush()

	if *tmpl != "" {
		t := template.Must(template.ParseFiles(*tmpl))
		t.Execute(os.Stdout, results)
	}
}

func loadTests(testFile string) (tests map[string]*http.Request, err error) {
	f, err := os.Open(testFile)
	if err != nil {
		return tests, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&tests); err != nil {
		return tests, err
	}
	return tests, err
}

func benchmarkUrlTest(req *http.Request) func(*testing.B) {
	return func(b *testing.B) {
		benchmarkUrl(req, b)
	}
}

func benchmarkUrl(req *http.Request, b *testing.B) {
	for i := 0; i < b.N; i++ {
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Println(resp.StatusCode)
			continue
		}
	}
}

func tlsClient(certFile, keyFile string) *http.Client {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	//caCert, err := ioutil.ReadFile(*caFile)
	//if err != nil {
	//log.Fatal(err)
	//}
	//caCertPool := x509.NewCertPool()
	//caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		//RootCAs:      caCertPool,
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{Transport: transport}
}
