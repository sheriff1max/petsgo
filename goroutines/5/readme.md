# Описание проекта

delayed(v int, d time.Duration) <-chan int — канал, в который значение придёт через d. fetchWithTimeout(ch <-chan int, timeout time.Duration) (int, bool) — select между каналом и time.After: не успело → false.