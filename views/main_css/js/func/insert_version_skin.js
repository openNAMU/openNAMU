"use strict";

function opennamu_do_insert_version_skin(
    dom_name_version
) {
    let url = "/api/skin_info?all=true";
    let xhr = new XMLHttpRequest();
    xhr.open("GET", url);
    xhr.send();

    xhr.onreadystatechange = function() {
        if(this.readyState === 4 && this.status === 200) {
            var json_data = JSON.parse(this.responseText);
            for(let key in json_data) {
                document.getElementById(dom_name_version).innerHTML += '<li>' +
                    json_data[key]['name'] + ' : ' + json_data[key]['skin_ver'] +
                    (json_data[key]['lastest_version'] ? ' (' + json_data[key]['lastest_version']['skin_ver'] + ')' : '') +
                '</li>'
            }
        }
    }
}

let opennamu_do_insert_version_skin_url = [
    '/manager/1',
    '/manager'
];
if(opennamu_do_insert_version_skin_url.includes(window.location.pathname)) {
    opennamu_do_insert_version_skin('ver_send_3');
}