"use strict";

function do_insert_user_info() {
    if(document.getElementById('opennamu_get_user_info')) {
        let name = document.getElementById('opennamu_get_user_info').innerHTML;
        
        let url = "/api/user_info/" + encodeURI(name) + "?render=1";
        let xhr = new XMLHttpRequest();
        xhr.open("GET", url, true);
        xhr.send(null);

        xhr.onreadystatechange = function() {
            if(this.readyState === 4 && this.status === 200) {
                document.getElementById('opennamu_get_user_info').innerHTML = JSON.parse(this.responseText)['data'];
                
                opennamu_do_ip_parser();
            }
        }
    }
}

do_insert_user_info();