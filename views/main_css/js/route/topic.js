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

function opennamu_get_thread_ui(user_id, date, data, code, color = '', blind = '', add_style = '', topic_num = '') {
    let color_b, class_b;
    if(blind == 'O') {
        color_b = data == '' ? 'opennamu_comment_blind' : 'opennamu_comment_blind_admin';
        class_b = 'opennamu_comment_blind_js';
    } else {
        color_b = 'opennamu_comment_blind_not';
        class_b = '';
    }

    let admin_check_box = ''
    if(topic_num != '') {
        admin_check_box = '<input type="checkbox" class="opennamu_blind_button" id="opennamu_blind_' + topic_num + '_' + code + '">';
    }
        
    return `
        <span class="` + class_b + `">
            <table class="opennamu_comment" style="` + add_style + `">
                <tr>
                    <td class="opennamu_comment_color_` + color + `">
                        ` + admin_check_box + `
                        <a href="#thread_shortcut" id="` + code + `">#` + code + `</a>
                        ` + user_id + `
                        <span style="float: right;">` + date + `</span>
                    </td>
                </tr>
                <tr>
                    <td class="` + color_b + ` opennamu_comment_data_main" id="thread_` + code + `">
                        ` + data + `
                    </td>
                </tr>
            </table>
            <hr class="main_hr">
        </span>
    `;
}

function opennamu_get_thread(topic_num = "", do_type = "") {
    let url, to_obj, color;    
    if(do_type === "top") {
        url = "/api/thread/" + topic_num + "/top";
        to_obj = 'opennamu_top_thread';
        color = 'red';
    } else {
        url = "/api/thread/" + topic_num;
        to_obj = 'opennamu_main_thread';
        color = 'default';
    }

    fetch(url).then(function(res) {
        return res.json();
    }).then(function(data) {
        let end_data = '';
        let end_render = [];

        let lang = data["language"];
        data = data["data"];

        let first = '';
        for(let for_a = 0; for_a < data.length; for_a++) {
            if(first === '') {
                first = data[for_a]["ip"];
            }

            let real_color = color;
            if(color !== 'red') {
                if(data[for_a]["blind"] === '1') {
                    real_color = 'blue';
                } else if(first === data[for_a]["ip"]) {
                    real_color = 'green';
                } else {
                    real_color = 'default';
                }
            }

            let date = '<a href="/thread/' + topic_num + '/comment/' + data[for_a]["id"] + '/tool">(' + lang["tool"] + ')</a> ' + data[for_a]["date"];
            let render_button = ' <a href="javascript:void(0);" id="opennamu_' + color + '_thread_render_' + data[for_a]["id"] + '_button">(' + lang["render"] + ")</a>"
            let render_data = data[for_a]["data"] !== "" ? data[for_a]["data"] : "[br]";

            end_data += opennamu_get_thread_ui(
                data[for_a]["ip_render"] + render_button, 
                date, 
                '<div class="opennamu_comment_scroll" id="opennamu_' + color + '_thread_render_' + data[for_a]["id"] + '">' + opennamu_xss_filter(render_data) + '</div>',
                data[for_a]["id"],
                real_color,
                data[for_a]["blind"],
                '',
                topic_num
            )

            end_render.push([
                render_data,
                data[for_a]["id"]
            ]);
        }

        document.getElementById(to_obj).innerHTML = end_data;

        for(let for_a = 0; for_a < end_render.length; for_a++) {
            let observer = new IntersectionObserver(entries => {
                entries.forEach(entry => {
                    if(entry.isIntersecting) {
                        opennamu_do_render(
                            'opennamu_' + color + '_thread_render_' + end_render[for_a][1],
                            end_render[for_a][0], 
                            '',
                            'thread'
                        );

                        observer.unobserve(entry.target);
                    }
                });
            });

            observer.observe(document.getElementById('opennamu_' + color + '_thread_render_' + end_render[for_a][1] + '_button'));
        }
    });
}