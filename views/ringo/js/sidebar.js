"use strict";

// func
function ringo_do_xss_encode(data) {
    data = data.replace(/'/g, '&#x27;');
    data = data.replace(/"/g, '&quot;');
    data = data.replace(/</g, '&lt;');
    data = data.replace(/</g, '&gt;');

    return data;
}

function ringo_do_url_encode(data) {
    return encodeURIComponent(data);
}

// event
function ringo_do_side_button_1() {
    if(temp_save[0] === '') {
        fetch("/api/recent_change/10").then(function(res) {
            return res.json();
        }).then(function(text) {
            let data = '';
            for(let for_a = 0; for_a < text.length; for_a++) {
                if(text[for_a][6] === '') {
                    data += '<a href="/w/' + ringo_do_url_encode(text[for_a][1]) + '">' + ringo_do_xss_encode(text[for_a][1]) + '</a><br>';
                    data += text[for_a][2] + ' | ' + ringo_do_xss_encode(text[for_a][3]) + '<br>';
                }
            }

            document.getElementById('side_content').innerHTML = data;
            temp_save[0] = data;
        }).catch(function(error) {
            document.getElementById('side_content').innerHTML = 'Error';
        });
    } else {
        document.getElementById('side_content').innerHTML = temp_save[0];
    }
}

function ringo_do_side_button_2() {
    if(temp_save[1] === '') {
        fetch("/api/recent_discuss/10").then(function(res) {
            return res.json();
        }).then(function(text) {
            let data = '';
            for(let for_a = 0; for_a < text.length; for_a++) {
                data += '<a href="/thread/' + ringo_do_url_encode(text[for_a][3]) + '">' + ringo_do_xss_encode(text[for_a][1]) + '</a><br>';
                data += text[for_a][2] + ' | ' + text[for_a][5] +'<br>';
            }

            document.getElementById('side_content').innerHTML = data;
            temp_save[1] = data;
        }).catch(function(error) {
            document.getElementById('side_content').innerHTML = 'Error';
        });
    } else {
        document.getElementById('side_content').innerHTML = temp_save[1];
    }
}

function ringo_do_side_button_3() {
    if(temp_save[2] === '') {
        fetch("/api/v2/bbs/main").then(function(res) {
            return res.json();
        }).then(function(data) {
            let end_data = '';

            let text = data['data'];
            for(let for_a = 0; for_a < text.length; for_a++) {
                end_data += '<a href="/bbs/w/' + text[for_a].set_id + '/' + text[for_a].set_code + '">' + ringo_do_xss_encode(text[for_a].title) + '</a><br>';
                end_data += text[for_a].date + ' | ' + text[for_a].user_id +'<br>';
            }

            document.getElementById('side_content').innerHTML = end_data;
            temp_save[2] = end_data;
        });
    } else {
        document.getElementById('side_content').innerHTML = temp_save[2];
    }
}

// init
let temp_save = ['', '', ''];

window.addEventListener('DOMContentLoaded', function() {
    if(document.getElementById("side_button_1")) {
        document.getElementById("side_button_1").addEventListener("click", ringo_do_side_button_1);
        document.getElementById("side_button_2").addEventListener("click", ringo_do_side_button_2);
        document.getElementById("side_button_3").addEventListener("click", ringo_do_side_button_3);

        if(window.localStorage.getItem('main_css_off_sidebar') && window.localStorage.getItem('main_css_off_sidebar') === '0') {
            ringo_do_side_button_1();
        }
    }
});