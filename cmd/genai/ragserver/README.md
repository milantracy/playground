##### steps

```
$ docker compose -f docker-compose-ragserver.yml up -d
$ go run cmd/genai/ragserver/main.go 
$ curl \                                          
    -X POST \
    -H 'Content-Type: application/json' \
    -d @- \
    http://localhost:8080/add/ << EOF  
{                                                               
        "documents": [
        {"text": "google stock price on 2025 Nov 28 is 320"},
        {"text": "google stock price on 2025 Nov 27 is 310"}
]}
EOF

$ curl \
    -X POST \
    -H 'Content-Type: application/json' \
    -d @- \
    http://localhost:8080/query/ << EOF
{ "query": "what is google stock price on 2025 November 28th?" }
EOF     
"According to the provided context, the Google stock price on 2025 November 28th is 320."
```
