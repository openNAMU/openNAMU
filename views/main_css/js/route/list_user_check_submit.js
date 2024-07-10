"use strict";

function opennamu_list_user_check_submit_post() {
    window.location.pathname = window.location.pathname.replace('/check_submit/', '/check/');
}

function opennamu_list_user_check_submit() {
    let lang_data = new FormData();
    lang_data.append('data', 'check');
    
    fetch('/api/lang', {
        method : 'POST',
        body : lang_data,
    }).then(function(res) {
        return res.json();
    }).then(function(lang) {
        lang = lang["data"];

        document.getElementById('opennamu_list_user_check_submit').innerHTML = '<button onclick="opennamu_list_user_check_submit_post();">' + lang[0] + '</button>';
    });
}