const outputElement = document.getElementById('output');

// 使用 Fetch 来获取服务端的流式响应
async function fetchStream() {
    const response = await fetch('http://127.0.0.1:7777/chat', {
        method: "POST",
    }).catch(error => console.error('Error:', error));;
    if (!response.ok) {
        throw new Error(`Network response was not ok: ${response}`);
    }
  
    const reader = response.body.pipeThrough(new TextDecoderStream()).getReader();
    
    while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        outputElement.textContent += value.replace(/(\r\n|\n|\r)/gm, "");
    }
}

// 调用 fetchStream 函数以开始获取和显示流式响应
fetchStream();
