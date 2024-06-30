const outputElement = document.getElementById('output');

// 使用 Fetch 来获取服务端的流式响应
async function fetchStream() {
    var params = new URLSearchParams();
    params.append("model", document.getElementById('model_input').value)
    params.append("chat", document.getElementById('input').value)
    const response = await fetch('http://127.0.0.1:7777/chat', {
        method: "POST",
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: params,
    }).catch(error => console.error('Error:', error));;
    if (!response.ok) {
        throw new Error(`Network response was not ok: ${response}`);
    }

    const reader = response.body.pipeThrough(new TextDecoderStream()).getReader();
    outputElement.textContent = "";

    while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        outputElement.textContent += value;
    }
}

function SaveInput() {
    var params = new URLSearchParams();
    params.append("chat", document.getElementById('input').value);
    params.append("name", document.getElementById('save_input').value);

    var request = new XMLHttpRequest();
    request.open('POST', 'http://127.0.0.1:7777/save', false);
    request.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');

    try {
        request.send(params);
        if (request.status === 200) {
            alert(request.responseText);
            updateFileList();
        } else {
            throw new Error(`Network response was not ok: ${request.status} ${request.statusText}`);
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Error: ' + error.message);
    }
}

function LoadInput() {
    var params = new URLSearchParams();
    params.append("name", document.getElementById('save_input').value);

    var request = new XMLHttpRequest();
    request.open('POST', 'http://127.0.0.1:7777/load', false);
    request.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');

    try {
        request.send(params);
        if (request.status === 200) {
            document.getElementById('input').value = request.responseText;
        } else {
            throw new Error(`Network response was not ok: ${request.status} ${request.statusText}`);
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Error: ' + error.message);
    }
}

function updateFileList() {
    var request = new XMLHttpRequest();
    request.open('GET', 'http://127.0.0.1:7777/list', true);
    request.onreadystatechange = function () {
        if (this.readyState === 4 && this.status === 200) {
            try {
                try {
                    files = JSON.parse(this.responseText);
                } catch (e) {
                    console.error("Failed to parse responseText:", e);
                }
                var fileListElem = document.getElementById('file_list');
                fileListElem.innerHTML = '';  // 清空文件列表

                if (!files || files.length === 0) {
                    var noFilesMessage = document.createElement('li');
                    noFilesMessage.textContent = 'No files found.';
                    fileListElem.appendChild(noFilesMessage);
                } else {
                    files.forEach(file => {
                        var listItem = document.createElement('li');
                        listItem.textContent = file;
                        listItem.onclick = function() {
                            document.getElementById('save_input').value = file;
                            LoadInput();
                        };
                        fileListElem.appendChild(listItem);
                    });
                }
            } catch (error) {
                console.error('Error parsing JSON:', error);
                alert('Error: ' + error.message);
            }
        }
    };
    request.send();
}

document.addEventListener('DOMContentLoaded', function() {
    updateFileList();
});
