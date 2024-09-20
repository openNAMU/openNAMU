"use strict";

function opennamu_main_redirect() {
    const url = window.location.pathname;
    const url_split = url.split('/');

    const post_n = [];
    const n = url_split[2];
    
    if(post_n.includes(n)) {
        opennamu_main_redirect_post(n);
    } else {
        opennamu_main_redirect_get(n);
    }
}

function opennamu_main_redirect_send() {
    const url = window.location.pathname;
    const url_split = url.split('/');

    const n = url_split[2];
    const data = document.getElementById('opennamu_main_redirect_input').value;

    if(n === '1') {
        window.location.pathname = '/recent_block/all/1/' + opennamu_do_url_encode(data);
    }
}

function opennamu_main_redirect_get(n) {
    let lang_data = new FormData();
    lang_data.append('data', 'send search');
    
    fetch('/api/lang', {
        method : 'POST',
        body : lang_data,
    }).then(function(res) {
        return res.json();
    }).then(function(lang) {
        lang = lang["data"];

        let data_html = '';
        
        if(n === '1') {
            data_html += '<input id="opennamu_main_redirect_input" placeholder="' + lang[1] + '">';
        } else {
            data_html += '<input id="opennamu_main_redirect_input">';
        }

        data_html += '<hr class="main_hr">';
        data_html += '<button onclick="opennamu_main_redirect_send();">' + lang[0] + '</button>';

        document.getElementById('opennamu_main_redirect').innerHTML = data_html;
    });
}

function opennamu_main_redirect_post(n) {
    let url = '';
    let input = '';

    let data_html = '<form action="' + url + '">';
    data_html += input;
    data_html += '</form>';

    document.getElementById('opennamu_main_redirect').innerHTML = '';
}