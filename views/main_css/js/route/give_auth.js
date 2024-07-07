"use strict";

function opennamu_give_auth() {
    const url = window.location.pathname;
    const url_split = url.split('/');

    let mode = 0;
    if(url_split.length === 3) {
        if(url_split[2] === 'auth_total') {
            mode = 1;
        }
    }

    let html_data = '';
    if(mode === 0) {
        if(url_split.length === 3) {
            html_data += '<textarea></textarea>';
        }


    } else {
        
    }

    fetch('/api/v2/list/auth').then(function(res) {
        return res.json();
    }).then(function(data) {
        let data_list = data["data"];
        console.log(data_list);

        document.getElementById('opennamu_give_auth').innerHTML = html_data;
    });
}