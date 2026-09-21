# Описание реализации

startEchoServer() (addr string, stop func()): net.Listen("tcp", "127.0.0.1:0"), в горутине цикл Accept, на каждое соединение — своя горутина с io.Copy(conn, conn) (эхо) и defer Close. Адрес верни из ln.Addr().String(). Клиентом выступит тест: net.Dial, запись строки, чтение ответа через bufio.