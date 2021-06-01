# golang-code-executer
## This project is an attempt to spin up a minimal golang web server, served via docker to execute user submitted c/c++ code

I have tested this on C programs.
To run the server in docker:
```
docker build . -t go-code-executer
docker run -it -p 5000:5000 go-code-executer
```
After this you can send json input to the rest server at port 5000
The following is an example request using python requests library.

```
r = requests.post("http://localhost:5000/execute", json={"code":"sample code", "language":"C", "input":"10 20"})
r.json()
```


