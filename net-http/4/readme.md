# Описание реализации

startHTTPServer() (addr string, stop func()): http.NewServeMux, маршруты /hello (пишет "hello") и /sum?a=2&b=3 (разбирает query через r.URL.Query(), пишет сумму), сервер запускай как go http.Serve(ln, mux) на том же 127.0.0.1:0. Несуществующий маршрут должен дать 404 автоматически.