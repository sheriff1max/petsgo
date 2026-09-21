# Описание реализации

На том же HTTP-сервере добавь маршруты /echo (читает тело и заголовок X-User, отвечает user=%s;body=%s) и /slow (спит 300 мс). В тесте создай client := &http.Client{Timeout: 100 * time.Millisecond}: GET на /slow должен упасть по таймауту; POST на /echo оформи через http.NewRequest + req.Header.Set + client.Do, тело прочитай до конца и закрой.