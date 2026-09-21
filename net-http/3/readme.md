# Описание реализации

startSlowServer(d time.Duration): сервер читает строку, спит d, отвечает "ok". requestWithTimeout(addr string, t time.Duration) (string, error): dial, conn.SetDeadline(time.Now().Add(t)), запись запроса, чтение ответа. Медленный сервер (300 мс) при лимите 100 мс обязан отдать ошибку дедлайна, быстрый (10 мс) — успех.