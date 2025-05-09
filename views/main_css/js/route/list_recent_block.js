"use strict";

function opennamu_list_recent_block() {
    const url = window.location.pathname;
    const url_split = url.split('/');

    let tool = 'all';
    let page = '1';
    let user_name = '';
    let why = '';
    let add_path = '';
    if(url_split.length > 2) {
        tool = url_split[2];

        if(tool === 'user' || tool === 'admin') {
            add_path = '_user';

            if(url_split.length > 3) {
                user_name = '/' + url_split[3];

                if(url_split.length > 4) {
                    page = url_split[4];
                }
            }
        } else {
            if(url_split.length > 3) {
                page = url_split[3];

                if(url_split.length > 4) {
                    why = '/' + url_split.slice(4).join('/');
                }
            }
        }
    }

    fetch('/api/v2/recent_block' + add_path + '/' + tool + '/' + page + user_name + why).then(function(res) {
        return res.json();
    }).then(function(data) {
        let lang = data["language"];
        let auth = data["auth"];
        data = data["data"];

        let data_html = '';

        let option_list = [['all', 'all'], ['regex', 'regex'], ['cidr', 'cidr'], ['private', 'private'], ['ongoing', 'in_progress']];
        for(let for_a = 0; for_a < option_list.length; for_a++) {
            data_html += '<a href="/recent_block/' + option_list[for_a][0] + '">(' + lang[option_list[for_a][1]] + ')</a> ';
        }

        option_list = [['/manager/11', 'blocked'], ['/manager/12', 'admin'], ['/manager/19', 'why']];
        for(let for_a = 0; for_a < option_list.length; for_a++) {
            data_html += '<a href="' + option_list[for_a][0] + '">(' + lang[option_list[for_a][1]] + ')</a> ';
        }

        data_html += '<hr class="main_hr">';

        /*
            data_list = append(data_list, []string{
                why,
                ip_pre_block,
                ip_render_block,
                ip_pre_blocker,
                ip_render_blocker,
                end,
                today,
                band,
                ongoing,
            })
        */
        for(let for_a = 0; for_a < data.length; for_a++) {
            let left = '';

            let ban_auth = (auth["owner"] === true || auth["ban"] === true);
            let ip = data[for_a][1];
            if(data[for_a][7] === '') {
                if(ban_auth) {
                    ip = '<a href="/auth/ban/' + opennamu_do_url_encode(data[for_a][1]) + '">' + ip + '</a>';
                }
                
                if(data[for_a][8] === '1') {
                    ip = '<s>' + ip + '</s>';
                }
            } else if(data[for_a][7] === 'private') {
                if(ban_auth) {
                    ip = '<a href="/auth/ban/' + opennamu_do_url_encode(data[for_a][1]) + '">' + ip + '</a>';
                }

                if(data[for_a][8] === '1') {
                    ip = '<s>' + ip + '</s>';
                }

                ip += ' (' + lang['private'] + ')';
            } else if(data[for_a][7] === 'cidr') {
                if(ban_auth) {
                    ip = '<a href="/auth/ban_cidr/' + opennamu_do_url_encode(data[for_a][1]) + '">' + ip + '</a>';
                }

                if(data[for_a][8] === '1') {
                    ip = '<s>' + ip + '</s>';
                }

                ip += ' (' + lang['cidr'] + ')';
            } else {
                if(ban_auth) {
                    ip = '<a href="/auth/ban_regex/' + opennamu_do_url_encode(data[for_a][1]) + '">' + ip + '</a>';
                }

                if(data[for_a][8] === '1') {
                    ip = '<s>' + ip + '</s>';
                }

                ip += ' (' + lang['regex'] + ')';
            }
            
            left += ip + ' ← ' + data[for_a][4];

            let end = "";
            if(data[for_a][5] === "release") {
                end = lang["release"];
            } else if(data[for_a][5] === "") {
                end = lang["limitless"];
            } else {
                end = data[for_a][5];
            }

            let right = end + '<br>' + data[for_a][6];

            let bottom = '';
            if(data[for_a][0] !== "") {
                if(data[for_a][0] === "edit filter") {
                    bottom = '<a href="/edit_filter/' + opennamu_do_url_encode(data[for_a][1]) + '">edit filter</a>'
                } else {
                    bottom = opennamu_send_render(opennamu_xss_filter(data[for_a][0]));
                }
            }

            data_html += opennamu_make_list(left, right, bottom);
        }

        data_html += opennamu_page_control('/recent_block/' + tool + user_name + '/{}' + why, Number(page), data.length);

        document.getElementById('opennamu_list_recent_block').innerHTML = data_html;
    });
}