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
    fetch("/api/recent_change/10").then(function(res) {
        return res.json();
    }).then(function(text) {
        if(temp_save[0] === '') {
            let data = '';
            for(let for_a = 0; for_a < text.length; for_a++) {
                if(text[for_a][6] === '') {
                    data += '<a href="/w/' + ringo_do_url_encode(text[for_a][1]) + '">' + ringo_do_xss_encode(text[for_a][1]) + '</a><br>';
                    data += text[for_a][2] + ' | ' + ringo_do_xss_encode(text[for_a][3]) + '<br>';
                } else {
                    data += '---<br>';
                    data += '--- | ---<br>';
                }
            }

            document.getElementById('side_content').innerHTML = data;
            temp_save[0] = data;
        } else {
            document.getElementById('side_content').innerHTML = temp_save[0];
        }
    }).catch(function(error) {
        document.getElementById('side_content').innerHTML = 'Error';
    });
}

function ringo_do_side_button_2() {
    fetch("/api/recent_discuss/10").then(function(res) {
        return res.json();
    }).then(function(text) {
        if(temp_save[1] === '') {
            let data = '';
            for(let for_a = 0; for_a < text.length; for_a++) {
                data += '<a href="/thread/' + ringo_do_url_encode(text[for_a][3]) + '">' + ringo_do_xss_encode(text[for_a][1]) + '</a><br>';
                data += text[for_a][2] + '<br>';
            }

            document.getElementById('side_content').innerHTML = data;
            temp_save[1] = data;
        } else {
            document.getElementById('side_content').innerHTML = temp_save[1];
        }
    }).catch(function(error) {
        document.getElementById('side_content').innerHTML = 'Error';
    });
}

function ringo_do_side_button_3() {
    if(temp_save[2] === '') {
        if(document.getElementsByClassName('opennamu_TOC').length > 0) {
            temp_save[2] = document.getElementsByClassName('opennamu_TOC')[0].innerHTML;
            document.getElementById('side_content').innerHTML = temp_save[2];
        }
    } else {
        document.getElementById('side_content').innerHTML = temp_save[2];
    }
}

function ringo_do_side_button_4() {
    if(temp_save[3] === '') {
        if(document.getElementsByClassName('opennamu_footnote').length > 0) {
            let data = '';
            for(let for_a = 0; for_a < document.getElementsByClassName('opennamu_footnote').length; for_a++) {
                data += document.getElementsByClassName('opennamu_footnote')[for_a].innerHTML + '<br>';
            }

            document.getElementById('side_content').innerHTML = data;
            temp_save[3] = data;
        }
    } else {
        document.getElementById('side_content').innerHTML = temp_save[3];
    }
}

// init
let temp_save = ['', '', '', ''];

window.addEventListener('DOMContentLoaded', function() {
    document.getElementById("side_button_1").addEventListener("click", ringo_do_side_button_1);
    document.getElementById("side_button_2").addEventListener("click", ringo_do_side_button_2);
    document.getElementById("side_button_3").addEventListener("click", ringo_do_side_button_3);
    document.getElementById("side_button_4").addEventListener("click", ringo_do_side_button_4);

    ringo_do_side_button_1();
});