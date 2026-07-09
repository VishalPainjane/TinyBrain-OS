const http = require('http');

async function runTest() {
  const reqData = JSON.stringify({
    prompt: [1, 2, 3],
    max_tokens: 10,
    eos_token_id: 9999
  });

  const options = {
    hostname: 'localhost',
    port: 8080,
    path: '/v2/generate',
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Content-Length': Buffer.byteLength(reqData)
    }
  };

  const req = http.request(options, (res) => {
    let receivedTokens = 0;
    
    res.on('data', (chunk) => {
      const dataStr = chunk.toString();
      // Split by newlines and count how many "data: {" tokens we got
      const lines = dataStr.split('\n');
      for (const line of lines) {
        if (line.startsWith('data: {"token"')) {
          receivedTokens++;
          console.log(`Received token ${receivedTokens}`);
        }
      }
    });

    res.on('end', () => {
      console.log(`\nConnection closed. Total tokens received: ${receivedTokens}`);
      if (receivedTokens === 10) {
        console.log("SUCCESS: End-to-end execution path validated successfully.");
        process.exit(0);
      } else {
        console.error(`ERROR: Expected 10 tokens, but received ${receivedTokens}`);
        process.exit(1);
      }
    });
  });

  req.on('error', (e) => {
    console.error(`Problem with request: ${e.message}`);
    process.exit(1);
  });

  req.write(reqData);
  req.end();
}

runTest();
