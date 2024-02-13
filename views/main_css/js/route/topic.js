"use strict";

function opennamu_do_remove_blind_thread() {
    const style = document.querySelector('#opennamu_remove_blind');
    if(style !== null) {
        if(style.innerHTML !== "") {
            style.innerHTML = '';
        } else {
            style.innerHTML = `
                .opennamu_comment_blind_js {
                    display: none;
                }
            `;
        }
    }
}

function opennamu_thread_blind() {
    let do_true = 0;
    for(let for_a = 0; for_a < document.getElementsByClassName("opennamu_blind_button").length; for_a++) {
        let id = document.getElementsByClassName("opennamu_blind_button")[for_a].id;
        id = id.replace(/^opennamu_blind_/, '');
        id = id.split('_');

        let checked = document.getElementsByClassName("opennamu_blind_button")[for_a].checked;
        if(checked) {
            fetch("/thread/" + id[0] + '/comment/' + id[1] + '/blind', { method : 'GET' });
            do_true = 1;
        }
    }

    if(do_true === 1) {
        history.go(0);
    }
}

function opennamu_thread_where() {
    for(let for_a = 0; for_a < document.getElementsByClassName('opennamu_comment_data_main').length; for_a++) {
        let data = document.getElementsByClassName('opennamu_comment_data_main')[for_a].innerHTML;
        let id = document.getElementsByClassName('opennamu_comment_data_main')[for_a].id.replace(/^thread_/, '');

        let match = Array.from(data.matchAll(/<a href="#(?:[0-9]+)">#([0-9]+)<\/a>/g)).map(match => match[1]);
        console.log(data);
        if(match) {
            console.log(match);
            for(let for_b = 0; for_b < match.length; for_b++) {
                if(document.getElementById('opennamu_topic_req_' + match[for_b])) {
                    if(document.getElementById('opennamu_topic_req_' + match[for_b]).innerHTML === '') {
                        document.getElementById('opennamu_topic_req_' + match[for_b]).innerHTML += '<hr>'
                    }

                    document.getElementById('opennamu_topic_req_' + match[for_b]).innerHTML += '<a href="#' + id + '">#' + id + '</a> ';
                }
            }
        }
    }
}