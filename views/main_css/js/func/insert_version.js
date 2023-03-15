"use strict";

function opennamu_do_insert_version(
    dom_name_version_now, 
    dom_name_version_new
) {
    let url = "/api/version";
    let xhr = new XMLHttpRequest();
    xhr.open("GET", url);
    xhr.send();

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            let get_data = JSON.parse(this.responseText);
            document.getElementById(dom_name_version_now).innerHTML += get_data['version'];
            
            let url_2 = 'https://raw.githubusercontent.com/openNAMU/openNAMU/' + get_data['build'] + '/version.json';
            let xhr_2 = new XMLHttpRequest();
            xhr_2.open("GET", url_2);
            xhr_2.send();
            
            xhr_2.onreadystatechange = function() {
                if(this.readyState === 4 && this.status === 200) {
                    document.getElementById(dom_name_version_new).innerHTML += JSON.parse(this.responseText)['beta']['r_ver'];
                }
            }
        }
    }
}

let opennamu_do_insert_version_url = [
    '/manager/1',
    '/manager',
    '/update'
];
if(opennamu_do_insert_version_url.includes(window.location.pathname)) {
    opennamu_do_insert_version('ver_send_2', 'ver_send');
}