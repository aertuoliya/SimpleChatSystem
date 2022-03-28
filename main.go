package main

func main() {
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}

// // import "net/http"

// // type helloHandler struct{}
// // type aboutHandler struct{}

// // func (m *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// // 	w.Write([]byte("hhhh"))
// // }
// // func (m *aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// // 	w.Write([]byte("nothing"))

// // }
// func main() {
// 	// mh := helloHandler{}
// 	// ma := aboutHandler{}
// 	// server := http.Server{
// 	// 	Addr:    "localhost:8081",
// 	// 	Handler: nil, //为nil则默认多路复用器
// 	// }
// 	// http.Handle("/hello", &mh)
// 	// http.Handle("/about", &ma)
// 	// http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
// 	// 	w.Write([]byte("home"))
// 	// })
// 	// http.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {
// 	// 	w.Write([]byte("welcome"))
// 	// })
// 	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 	// 	w.Write([]byte("really home"))
// 	// })
// 	// server.ListenAndServe()
// 	//http.ListenAndServe("localhost:808 1", nil) //defaultServeMux
// 	// http.ListenAndServe(":8080", http.FileServer(http.Dir("html_study")))
// }
