"use strict";

function opennamu_xss_filter(str) {
    return str.replace(/[&<>"']/g, function(match) {
        switch(match) {
            case '&':
                return '&amp;';
            case '<':
                return '&lt;';
            case '>':
                return '&gt;';
            case "'":
                return '&#x27;';
            case '"':
                return '&quot;';
        }
    });
}

function opennamu_do_ip_click(obj) {
    if (obj.id === "") {
        let user_name = obj.name;

        fetch('/api/v2/ip_menu/' + user_name).then(response => {
            if(!response.ok) {
                throw new Error(`API 호출 실패: ${response.status}`);
            }

            return response.json();
        }).then(data => {
            data = data["data"];

            let data_html = '';
            for(let key in data) {
                for(let for_a = 0; for_a < data[key].length; for_a++) {
                    data_html += '<a href="' + data[key][for_a][0] + '">' + data[key][for_a][1] + '</a> | ';
                }
            }

            data_html = data_html.replace(/ \| $/g, '');

            let for_a;
            for(for_a = 0; document.getElementById("opennamu_ip_render_" + String(for_a) + "_load"); for_a++) {}

            let popup_html = '<span class="opennamu_popup_footnote" id="opennamu_ip_render_' + String(for_a) + '_load" style="display: none;"></span>';
            popup_html += '<span style="display: none;" id="opennamu_ip_tool_' + String(for_a) + '">';
            popup_html += data_html;
            popup_html += '</span>';

            obj.innerHTML += popup_html;
            obj.id = 'opennamu_ip_render_' + String(for_a);
            obj.onclick = '';

            document.getElementById('opennamu_ip_render_' + String(for_a)).addEventListener("click", function () {
                opennamu_do_footnote_popover('opennamu_ip_render_' + String(for_a), '', 'opennamu_ip_tool_' + String(for_a), 'open');
            });
            document.addEventListener("click", function () {
                opennamu_do_footnote_popover('opennamu_ip_render_' + String(for_a), '', 'opennamu_ip_tool_' + String(for_a), 'close');
            });

            obj.click();
        }).catch(err => {
            console.error('IP 메뉴 호출 중 오류 발생:', err);
            obj.innerHTML = 'IP 정보를 불러오는 데 실패했습니다.';
        });
    }
}

function opennamu_list_hidden_remove() {
    const style = document.querySelector('#opennamu_list_hidden_style');
    if(style !== null) {
        if(style.innerHTML !== "") {
            style.innerHTML = '';
        } else {
            style.innerHTML = '.opennamu_list_hidden { display: none; }';
        }
    }
}

function opennamu_do_footnote_popover(set_name, load_name, sub_obj = undefined, do_type = 'open') {
    if(document.getElementById(set_name + '_load')) {
        document.getElementById(set_name + '_load').onclick = function (event) {
            event.stopPropagation();
        };
        if(do_type === 'open') {
            if(sub_obj !== undefined) {
                document.getElementById(set_name + '_load').innerHTML = document.getElementById(sub_obj).innerHTML;
            } else {
                document.getElementById(set_name).title = '';
                document.getElementById(set_name + '_load').innerHTML = '<a href="#' + load_name + '">(Go)</a> ' + document.getElementById(load_name + '_title').innerHTML;
            }
            document.getElementById(set_name + '_load').style.display = "inline-block";
            document.getElementById(set_name + '_load').count = 0;

            let width = document.getElementById(set_name + '_load').clientWidth;
            let screen_width = window.innerWidth;
            let left = document.getElementById(set_name).getBoundingClientRect().left;
            let top = window.pageYOffset + document.getElementById(set_name).getBoundingClientRect().top;

            document.getElementById(set_name + '_load').style.top = String(top) + "px";
            if(screen_width - (left + width) < 50) {
                if(left > 350) {
                    document.getElementById(set_name + '_load').style.left = String(left - 300) + "px";
                } else {
                    document.getElementById(set_name + '_load').style.left = "0px";
                }

                left = document.getElementById(set_name + '_load').getBoundingClientRect().left;
                width = document.getElementById(set_name + '_load').clientWidth;
                if(300 > width) {
                    document.getElementById(set_name + '_load').style.left = String(left + (300 - width)) + "px";
                } else {
                    document.getElementById(set_name + '_load').style.marginTop = "20px";
                }
            }
        } else {
            if(document.getElementById(set_name + '_load').count === 1) {
                document.getElementById(set_name + '_load').style.display = "none";
            } else {
                document.getElementById(set_name + '_load').count = 1;
            }
        }
    }
}

document.addEventListener("click", function () {
    document.querySelectorAll('span[id$="_over"] > .opennamu_popup_footnote').forEach(function (obj) {
        opennamu_do_footnote_popover(obj.id.slice(0, -5), '', undefined, 'close');
    });
});

