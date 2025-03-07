Test results:

1. 
samirsingh@Samirs-MacBook-Pro temp % go run main.go
server running on port 8080...


2.

samirsingh@Samirs-MacBook-Pro temp % curl -v -X POST http://localhost:8080/scan -F "file=@sample.pdf"
Note: Unnecessary use of -X or --request, POST is already inferred.
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> POST /scan HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: */*
> Content-Length: 19015
> Content-Type: multipart/form-data; boundary=------------------------Mpjf3I3ee4kPCaLGGTlxyA
> 
* upload completely sent off: 19015 bytes
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Fri, 07 Mar 2025 16:00:41 GMT
< Content-Length: 76
< 
{"hash":"229defbb0cee6f02673a5cde290d0673e75a0dc31cec43989c8ab2a4eca7e1bb"}
* Connection #0 to host localhost left intact

3.

samirsingh@Samirs-MacBook-Pro temp % curl -v "http://localhost:8080/lookup?hash=229defbb0cee6f02673a5cde290d0673e75a0dc31cec43989c8ab2a4eca7e1bb"
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /lookup?hash=229defbb0cee6f02673a5cde290d0673e75a0dc31cec43989c8ab2a4eca7e1bb HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: */*
> 
* Request completely sent off
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Fri, 07 Mar 2025 16:00:52 GMT
< Content-Length: 290
< 
{"hash":"229defbb0cee6f02673a5cde290d0673e75a0dc31cec43989c8ab2a4eca7e1bb","pdf_version":"1.3","producer":"Mac OS X 10.5.4 Quartz PDFContext","author":"Philip Hutchison","created_date":"D:20080701052447+00'00'","updated_date":"D:20080701052447+00'00'","uploaded_at":"2025-03-07T16:00:41Z"}
* Connection #0 to host localhost left intact
