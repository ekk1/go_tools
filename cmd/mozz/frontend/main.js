var now = new Date().toISOString().split('T')[0];
document.getElementById("dateinput").value = now;

const base_url = "http://127.0.0.1:5011";

for (var i = 0; i < 5; i++) {
    var _tr = document.createElement("tr");
    for (var j = 0; j < 3; j++) {
        var _data = document.createElement("td");
        _data.innerHTML = "test";
        _tr.appendChild(_data);
    }
    document.getElementById("record-list").appendChild(_tr);
}

function refreshAccounts() {
    var xhr = new XMLHttpRequest();
    xhr.open("GET", base_url + "/accounts", true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            var accounts = JSON.parse(xhr.responseText);
            document.getElementById("account-list").innerHTML = "";
            for (var i = 0; i < accounts.length; i++) {
                var _tr = document.createElement("tr");
                var _td = document.createElement("td");
                _td.innerHTML = accounts[i].name;
                _tr.appendChild(_td);
                var _td = document.createElement("td");
                _td.innerHTML = accounts[i].balance;
                _tr.appendChild(_td);
                var _td = document.createElement("td");
                _td.innerHTML = accounts[i].is_credit;
                _tr.appendChild(_td);
                var _td = document.createElement("td");
                _td.innerHTML = accounts[i].billing_date;
                _tr.appendChild(_td);
                document.getElementById("account-list").appendChild(_tr);
            }
        } else if (xhr.readyState == 4 && xhr.status != 200) {
            alert(xhr.status, xhr.responseText);
        }
    }
    xhr.send();
}

function AddAccounts() {
    var name = document.getElementById("account_name").value;
    var initial_balance = document.getElementById("initial_balance").value;
    var xhr = new XMLHttpRequest();
    xhr.open("POST", base_url + "/accounts", true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            refreshAccounts();
        } else if (xhr.readyState == 4 && xhr.status != 200) {
            alert(xhr.status, xhr.responseText);
        }
    }
    xhr.send(JSON.stringify({ name: name, initial_balance: Number(initial_balance) }));
}

// Get accounts info from server
refreshAccounts();