package handlers

type postData struct {
	key   string
	value string
}

var theTests []struct {
	name   string
	url    string
	method string
	params []postData
	expectedStatusCode int 
}{
	{"home", "/", "GET", []postData{},200},
}

func TestHandlers(t *testing.T){
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)

	defer ts.Close()
}
