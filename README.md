# stdin_analyzer
Осуществляет поиск слова "Go" по тексту, полученному с указанных сайтов или из файлов.
Переключение между поиском по сайтам или файлам осуществляется с помощью флага type, который имеет возожные значения url, file.

```echo -e "https://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org" | go run main.go --type url```

В результате выводит количество нахождений по каждому указанному сайту или файлу и общее количество нахождений:
```
root@Artem:/mnt/e/DEVELOPER/go/test# echo -e "https://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org" | go run main.go --type url
Count for https://golang.org: 9
Count for https://golang.org: 9
Count for https://golang.org: 9
Count for https://golang.org: 9
Total: 36
```
```
echo -e "/etc/passwd\n/etc/hosta" | go run main.go --type file
Error in /etc/hosta: open /etc/hosta: no such file or directory
Count for /etc/passwd: 0
Total: 0
```
