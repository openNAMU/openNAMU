"use strict";

function opennamu_change_comment(get_id) {
    const input = document.querySelector('#opennamu_comment_select');
    if(input !== null) {
        input.value = get_id;
        document.getElementById('opennamu_edit_textarea')?.focus();
    }
}

function opennamu_return_comment() {
    const input = document.querySelector('#opennamu_comment_select');
    if(input !== null) {
        document.getElementById(input.value)?.focus();
    }
}

function opennamu_load_comment() {
    const url = window.location.pathname;
    const url_split = url.split('/');

    let bbs_id = url_split[3];
    let bbs_code = url_split[4];

    fetch('/api/v2/bbs/w/comment/' + bbs_id + '-' + bbs_code + '/normal').then(function(res) {
        return res.json();
    }).then(function(data) {
        let data_html = '';

        if(data) {
            for(let for_a in data.data) {
                let data_in = data.data[for_a];
                console.log(data_in);

                data_html += opennamu_get_thread_ui(
                    data[for_a]["ip_render"], 
                    date, 
                    '<div class="opennamu_comment_scroll" id="opennamu_' + color + '_thread_render_' + data[for_a]["id"] + '">' + opennamu_xss_filter(render_data) + '</div>',
                    data[for_a]["id"],
                    real_color,
                    data[for_a]["blind"],
                    '',
                    topic_num
                )
            }
        }

        document.getElementById('opennamu_bbs_w_post').innerHTML = data_html;
    });
}