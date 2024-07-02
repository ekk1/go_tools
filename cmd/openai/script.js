const outputElement = document.getElementById('output');

// Constants for original and unique delimiters
const ORIGINAL_SPLITTER = '___SPLITTER_ORIGINAL___';
const ORIGINAL_TEMP_SPLITTER = '___SPLITTER_TEMP___';
const ORIGINAL_FINAL_SPLITTER = '___SPLITTER_FINAL___';

// Helper function to convert a Unicode string to a Base64 string
function unicodeToBase64(str) {
    // Convert the Unicode string to Uint8Array
    const utf8Bytes = new TextEncoder().encode(str);

    // Convert the Uint8Array to a binary string
    const binaryStr = Array.from(utf8Bytes)
        .map(byte => String.fromCharCode(byte))
        .join('');

    // Encode the binary string to Base64
    return btoa(binaryStr);
}

// Helper function to decode a Base64 string to a Unicode string
function base64ToUnicode(base64Str) {
    // Decode the Base64 string to a binary string
    const binaryStr = atob(base64Str);

    // Convert the binary string to Uint8Array
    const utf8Bytes = new Uint8Array(
        [...binaryStr].map(char => char.charCodeAt(0))
    );

    // Decode the Uint8Array to a Unicode string
    return new TextDecoder().decode(utf8Bytes);
}


// 使用 Fetch 来获取服务端的流式响应
async function fetchStream() {

    const textareas = document.querySelectorAll('.textarea-container textarea');
    let combinedInput = '';
    textareas.forEach(textarea => combinedInput += textarea.value);

    var params = new URLSearchParams();
    params.append("model", document.getElementById('model_input').value)
    params.append("chat", combinedInput);

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
    const textareas = document.querySelectorAll('.textarea-container textarea');
    let combinedInput = '';
    textareas.forEach(textarea => {
        let value = unicodeToBase64(textarea.value);
        combinedInput += `${ORIGINAL_SPLITTER}${value}`;
    });

    const tempInputValue = unicodeToBase64(document.getElementById('temp_input').value);
    combinedInput += `${ORIGINAL_TEMP_SPLITTER}${tempInputValue}`;

    const outputContent = unicodeToBase64(document.getElementById('output').textContent);
    combinedInput += `${ORIGINAL_FINAL_SPLITTER}${outputContent}`;

    var params = new URLSearchParams();
    params.append("chat", combinedInput);
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
            const textareas = document.querySelectorAll('.textarea-container textarea');
            let combinedInput = request.responseText;

            let parts = combinedInput.split(ORIGINAL_TEMP_SPLITTER);
            let splitInputs = parts[0].split(ORIGINAL_SPLITTER);
            splitInputs.shift(); // The first element will be an empty string due to the initial splitter

            for (var i = 0; i < textareas.length; i++) {
                if (splitInputs[i] !== undefined) {
                    try{
                        textareas[i].value = base64ToUnicode(splitInputs[i]);
                    } catch (error) {
                        textareas[i].value = splitInputs[i];
                    }
                }
            }
            if (parts.length > 1) {
                let split2 = parts[1].split(ORIGINAL_FINAL_SPLITTER);
                try{
                    document.getElementById('temp_input').value = base64ToUnicode(split2[0]);
                } catch (error) {
                    document.getElementById('temp_input').value = split2[0];
                }
                if (split2.length > 1) {
                    let outputContent = split2[1];
                    try{
                        document.getElementById('output').textContent = base64ToUnicode(outputContent);
                    } catch (error) {
                        document.getElementById('output').textContent = outputContent;
                    }
                }
            }
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

function MergeInput() {
    const textareas = document.querySelectorAll('.textarea-container textarea');
    let combinedInput = '';
    textareas.forEach(textarea => {
        combinedInput += textarea.value;
    });

    document.getElementById('output').value = combinedInput
}

function ClearInput() {
    const textareas = document.querySelectorAll('.textarea-container textarea');
    textareas.forEach(textarea => {
        textarea.value = ""
    });

    document.getElementById('output').value = ""
    document.getElementById('temp_input').value = ""
}

document.addEventListener('DOMContentLoaded', function() {
    updateFileList();
});
