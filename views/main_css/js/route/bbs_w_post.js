"use strict";

function opennamu_change_comment(get_id) {
    const input = document.querySelector('#opennamu_comment_select');
    if(input !== null) {
        input.value = get_id;
        document.getElementById('opennamu_comment_select')?.focus();
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
            let lang = data["language"];
            data = data["data"];

            let end_render = [];
            let select = '<select id="opennamu_comment_select" name="comment_select">';
            select += '<option value="0">' + lang["normal"] + '</option>';

            for(let for_a in data) {
                let data_in = data[for_a];

                let code_id = data_in["id"] + '-' + data_in["code"];
                code_id = code_id.replace(/^[0-9]+-[0-9]+-/, '');

                let count = 0;
                for(let for_a = 0; for_a < code_id.length; for_a++) {
                    if(code_id[for_a] === '-') {
                        count += 1;
                    }
                }
                
                select += '<option value="' + code_id + '">' + code_id + '</option>';

                let color = 'default';
                let date = '';

                date += '<a href="javascript:opennamu_change_comment(\'' + code_id + '\');">(' + lang["comment"] + ')</a> ';
                date += '<a href="/bbs/tool/' + bbs_id + '/' + bbs_code + '/' + code_id + '">(' + lang["tool"] + ')</a> ';
                date += data_in["comment_date"];

                data_html += '<span style="padding-left: ' + String(20 * count) + 'px;"></span>';
                data_html += opennamu_get_thread_ui(
                    data_in["comment_user_id_render"], 
                    date, 
                    '<div class="opennamu_comment_scroll" id="opennamu_thread_render_' + code_id + '">' + opennamu_xss_filter(data_in["comment"]) + '</div>',
                    code_id,
                    color,
                    '',
                    'width: calc(100% - ' + String(20 * count) + 'px);',
                    ''
                );

                end_render.push([
                    data_in["comment"],
                    code_id
                ]);
            }

            select += '</select> <a href="javascript:opennamu_return_comment();">(' + lang["return"] + ')</a>';
            select += '<hr class="main_hr">';

            if(document.getElementById('opennamu_bbs_w_post_select')) {
                document.getElementById('opennamu_bbs_w_post_select').innerHTML = select;
            }

            document.getElementById('opennamu_bbs_w_post').innerHTML = data_html;

            for(let for_a = 0; for_a < end_render.length; for_a++) {
                let observer = new IntersectionObserver(entries => {
                    entries.forEach(entry => {
                        if(entry.isIntersecting) {
                            opennamu_do_render(
                                'opennamu_thread_render_' + end_render[for_a][1],
                                end_render[for_a][0], 
                                '',
                                'thread'
                            );

                            observer.unobserve(entry.target);
                        }
                    });
                });

                observer.observe(document.getElementById('opennamu_thread_render_' + end_render[for_a][1]));
            }
        }
    });
}