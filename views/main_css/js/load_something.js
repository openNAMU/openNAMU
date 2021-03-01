function load_user_info(name) {
    var url = "/api/user_info/" + encodeURI(name) + "?render=1";
    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            document.getElementById('get_user_info').innerHTML += JSON.parse(this.responseText)['data'];
        }
    }
}

function load_ver() {
    var url = "/api/version";
    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            document.getElementById('ver_send').innerHTML += JSON.parse(this.responseText)['lastest_version'];
            document.getElementById('ver_send').style.display = "list-item";
        }
    }
}

function do_twofa_check(init = 0) {
    var data_check = document.getElementById('twofa_check_input').checked;
    document.getElementById('fa_plus_content').style.display = data_check === true ? "block" : "none";
}