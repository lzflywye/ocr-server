export CPATH="/opt/homebrew/include"
export LIBRARY_PATH="/opt/homebrew/lib"

curl -X POST -F "files=@./hello1.png" http://127.0.0.1:8080/upload | jq
