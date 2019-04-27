function load_ver() {
    var n_ver = document.getElementById('ver_send');

    var url = "/api/version";
    
    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            n_ver.innerHTML += JSON.parse(this.responseText)['lastest_version'];
            n_ver.style.display = "list-item";
        }
    }
}