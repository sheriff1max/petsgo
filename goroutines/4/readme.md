# Описание проекта

Три функции: producer(out chan<- int, n int), doubler(in <-chan int, out chan<- int), collector(in <-chan int) []int. Собери цепочку producer → doubler → collector для 1..5. Следи, чтобы каждый этап закрывал свой выходной канал после завершения входа.