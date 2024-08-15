# Code Sandbox Service
A code sandbox service that exposes HTTP API for running code in a sandboxed environment.

Sandboxes using [isolate](https://github.com/ioi/isolate).

Now supports c/c++/go/python/java.

## TODO
Docker image building.

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
            "time": 0.003,
            "time_wall": 0.013,
            "max_rss": 3620,
            "killed": false,
            "message": "",
            "status": "",
            "exitsig": 0,
            "stdin": "3 4",
            "stderr": "",
            "stdout": "7"
        },
        {
            "time": 0.003,
            "time_wall": 0.013,
            "max_rss": 3484,
            "killed": false,
            "message": "",
            "status": "",
            "exitsig": 0,
            "stdin": "6 7",
            "stderr": "",
            "stdout": "13"
        },
        {
            "time": 0.003,
            "time_wall": 0.013,
            "max_rss": 3556,
            "killed": false,
            "message": "",
            "status": "",
            "exitsig": 0,
            "stdin": "4 5",
            "stderr": "",
            "stdout": "9"
        },
        {
            "time": 0.003,
            "time_wall": 0.018,
            "max_rss": 3524,
            "killed": false,
            "message": "",
            "status": "",
            "exitsig": 0,
            "stdin": "33 44",
            "stderr": "",
            "stdout": "77"
        },
        {
            "time": 0.003,
            "time_wall": 0.017,
            "max_rss": 3528,
            "killed": false,
            "message": "",
            "status": "",
            "exitsig": 0,
            "stdin": "2 3",
            "stderr": "",
            "stdout": "5"
        }
    ]
}
```
Error:

```json
{
    "code": 400,
    "message": "/root/project/sandbox/running/eadd8240-3f0d-45d7-88dc-5249df5cac92/code.c:2:21: error: expected ';' before 'int'\n    2 |  using namespace std int main() {int a,b;cin>>a>>b;cout<<a+b<<endl;}\n      |                     ^~~~\n      |                     ;\n",
    "meta": null
}
```