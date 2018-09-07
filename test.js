const http = require('http');

function testAPI() {
    testTestUpload()
        .then(() => testRun())
        .then(() => console.log("All test passed"))
        .catch(e => console.error(e))
}

function testTestUpload(){
    return new Promise((resolve, reject) => {
        const req = http.request({
            hostname: 'localhost',
            port: 3010,
            path: '/test/upload',
            method: 'POST',
        }, res => {
            res.setEncoding('utf8');
            let body = '';

            res.on('data', chunk => {
                console.log('Response: ' + chunk);
                body += chunk
            });

            res.on('end', () => {
                const data = JSON.parse(body);
                if (data.isSuccessfulAdded) {
                    console.log("Test upload: ok")
                    resolve("Test passed")
                } else {
                    console.log("Test upload: failed")
                    reject("Test failed")
                }
            })
        });

        req.write(JSON.stringify({
            text: "3 2 1\n6"
        }))

        req.end()
    })
}

function testRun(){
    return new Promise((resolve, reject) => {
        const req = http.request({
            hostname: 'localhost',
            port: 3010,
            path: '/run',
            method: 'POST',
        }, res => {
            res.setEncoding('utf8');
            let body = '';

            res.on('data', chunk => {
                console.log('Response: ' + chunk);
                body += chunk
            });

            res.on('end', () => {
                const data = JSON.parse(body);
                if (data.isAllTestsSuccessful) {
                    console.log("Test run: ok")
                    resolve("Test passed")
                } else {
                    console.log("Test run: failed")
                    reject("Test failed")
                }
            })
        });

        req.write(JSON.stringify({
            problemId: "3",
            code: `#include <iostream>
                    using namespace std;
                    
                    int main()
                    {
                        int a, b, c;
                        cin >> a >> b >> c;
                        cout << a + b + c << endl;
                        return 0;
                    }`
        }));

        req.end()
    })
}

testAPI();