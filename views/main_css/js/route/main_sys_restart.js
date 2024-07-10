"use strict";

function opennamu_main_sys_restart() {
    let lang_data = new FormData();
    lang_data.append('data', 'restart');
    
    fetch('/api/lang', {
        method : 'POST',
        body : lang_data,
    }).then(function(res) {
        return res.json();
    }).then(function(lang) {
        lang = lang["data"];

        document.getElementById('opennamu_main_sys_restart').innerHTML = '<button id="opennamu_main_sys_restart_button">' + lang[0] + '</button>';

        document.getElementById('opennamu_main_sys_restart_button').addEventListener('click', function() {
            fetch('/restart', {
                method : 'POST'
            });

            window.location.href = '/';
        });
    });
}