# Code Sandbox Service
A code sandbox service that exposes HTTP API for running code in a sandboxed environment.
Sandboxes using [isolate](https://github.com/ioi/isolate).

**Request sample:**

***request.json***
```json
{
    "code": "#include <iostream> \n using namespace std; int main() {int a,b;cin>>a>>b;cout<<a+b<<endl;}",
    "language" : "cpp",
    "max_time": 1,
    "max_mem": 10000,
    "stdin": ["3 4", "23 6"]
}
```
**Response sample:**
***response.json***
Success:
```json
{
    "code": 200,
    "message": "",
    "meta": [
        {
            "time": 0.004,
            "time_wall": 0.017,
            "max_rss": 3636,
            "killed": false,
            "message": "",
            "status": "",
            "exitsig": 0,
            "stderr": "",
            "stdout": "7"
        },
        {
            "time": 0.004,
            "time_wall": 0.014,
            "max_rss": 3568,
            "killed": false,
            "message": "",
            "status": "",
            "exitsig": 0,
            "stderr": "",
            "stdout": "29"
        }
    ]
}
```
Error:
```json
{
    "code": 400,
    "message": "/root/project/sandbox/c1bdd4b6-17c2-4614-88c9-f613e5ce0089/build.cpp:2:21: error: expected ';' before 'int'\n    2 |  using namespace std int main() {int a,b;cin>>a>>b;cout<<a+b<<endl;}\n      |                     ^~~~\n      |                     ;\n",
    "meta": null
}
```