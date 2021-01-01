package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Get() (proxy string) {
	requestURL := "http://icanhazip.com"
	proxys := "http://27.220.49.108:9000"

	p, _ := url.Parse(proxys)
	fmt.Println(p)
	netTransport := &http.Transport{
		Proxy: http.ProxyURL(p),
	}

	httpClient := &http.Client{
		Transport: netTransport,
	}

	resp, err := httpClient.Get(requestURL)
	
	defer resp.Body.Close()
	
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 && string(body) != "" {
		fmt.Println(200)
	}
	return proxys
}

	for _, proxy = range proxys {
	p, _ := url.Parse("http://" + proxy)
	fmt.Println(p)
	netTransport := &http.Transport{
		Proxy:                 http.ProxyURL(p),
		ResponseHeaderTimeout: time.Second * time.Duration(5),
		MaxIdleConns:          10,
	}

	httpClient := &http.Client{
		timeout: time.Second*60
		Transport: netTransport,
	}

	resp, err := httpClient.Get(requestURL)
	if err != nil {
		fmt.Println(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 && string(body) != "" {
		fmt.Printf("%s request %s return %s  可用\n", proxy, requestURL, string(body))
		resp.Body.Close()
		return "http://" + proxy
	} else {
		fmt.Printf("%s request %s return %s  不可用\n", proxy, requestURL, string(body))
		resp.Body.Close()
		proxy = ""
	}
}

return proxy
}

func get() []string {
	reg := regexp.MustCompile("[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+:[0-9]+")
	//config.Conf.ProxySource
	resp, _ := http.Get("http://www.66ip.cn/mo.php?tqsl=100")
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	proxys := reg.FindAllString(string(b), 100)

	return proxys
}
