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

function do_skin_ver_check() {
    var url = "/api/skin_info?all=true";
    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.send(null);

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            var json_data = JSON.parse(this.responseText);
            var all_need_update = [];
            for(var key in json_data) {
                if(json_data[key]['lastest_version']) {
                    var new_skin_ver = json_data[key]['lastest_version']['skin_ver'];
                    var old_skin_ver = json_data[key]['skin_ver'];
                    var skin_name = json_data[key]['name'];
                    if(new_skin_ver !== old_skin_ver) {
                        all_need_update.push(skin_name);
                    }
                }
            }
            
            if(all_need_update.length !== 0) {
                document.getElementById('need_skin_update').innerHTML += ' (' + (all_need_update.join(', ')) + ')';
            }
        }
    }
}

function do_twofa_check(init = 0) {
    var data_check = document.getElementById('twofa_check_input').checked;
    document.getElementById('fa_plus_content').style.display = data_check === true ? "block" : "none";
}

function send_render(i = 0) {
    var get_class = document.getElementsByClassName('send_content')[i];
    if(get_class) {
        send_render(i + 1);
        
        var data = get_class.innerHTML;
        
        if(data === '&lt;br&gt;') {
            document.getElementsByClassName('send_content')[i].innerHTML = '<br>';
        } else {
            data = data.replace(/javascript:/i, '');
            
            data = data.replace(/&lt;a&gt;((?:(?!&lt;\/a&gt;).)+)&lt;\/a&gt;/g, function(x, x_1) {
                return '<a href="/w/' + encodeURIComponent(x_1) + '">' + x_1 + '</a>';
            });
            
            document.getElementsByClassName('send_content')[i].innerHTML = data;
        }
    }
}

function ie_end_support() {
    if(document.currentScript === undefined) {
        window.location = 'microsoft-edge:' + window.location;
        setTimeout(function() {
            window.location = 'https://go.microsoft.com/fwlink/?linkid=2135547';
        }, 1);
    }
}

ie_end_support();