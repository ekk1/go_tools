var now = new Date().toISOString().split('T')[0];
document.getElementById("dateinput").value = now;

for (var i = 0; i < 5; i++) {
    var _tr = document.createElement("tr");
    for (var j = 0; j < 3; j++) {
        var _data = document.createElement("td");
        _data.innerHTML = "test";
        _tr.appendChild(_data);
    }
    document.getElementById("record-list").appendChild(_tr);
}
