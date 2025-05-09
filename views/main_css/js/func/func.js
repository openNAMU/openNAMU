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

function opennamu_xss_filter_decode(str) {
    return str.replace(/&amp;|&lt;|&gt;|&#x27;|&quot;/g, function(match) {
        switch(match) {
            case '&amp;':
                return '&';
            case '&lt;':
                return '<';
            case '&gt;':
                return '>';
            case '&#x27;':
                return "'";
            case '&quot;':
                return '"';
        }
    });
}

function renderSimpleSet(data) {
    let tocData = '';
    const tocRegexAll = /<h([1-6])>([^<>]+)<\/h[1-6]>/g;
    const tocRegex = /<h([1-6])>([^<>]+)<\/h[1-6]>/;
    const tocSearchData = [...data.matchAll(tocRegexAll)];
    let headingStack = [0, 0, 0, 0, 0, 0];

    if (tocSearchData.length > 0) {
        tocData += `
            <div class="opennamu_TOC" id="toc">
                <span class="opennamu_TOC_title">TOC</span>
                <br>
        `;
    }

    tocSearchData.forEach((tocSearchIn) => {
        const headingLevel = parseInt(tocSearchIn[1]);
        const headingLevelStr = headingLevel.toString();

        headingStack[headingLevel - 1] += 1;
        for (let i = headingLevel; i < 6; i++) {
            headingStack[i] = 0;
        }

        const headingStackStr = headingStack
            .map((val) => (val !== 0 ? val + '.' : ''))
            .join('')
            .replace(/\.$/, '');

        tocData += `
            <br>
            <span class="opennamu_TOC_list">
                ${'<span style="margin-left: 10px;"></span>'.repeat(headingStackStr.split('.').length - 1)}
                <a href="#s-${headingStackStr}">${headingStackStr}.</a>
                ${tocSearchIn[2]}
            </span>
        `;

        data = data.replace(
            tocRegex,
            `<h${tocSearchIn[1]} id="s-${headingStackStr}"><a href="#toc">${headingStackStr}.</a> ${tocSearchIn[2]}</h${tocSearchIn[1]}>`
        );
    });

    if (tocData !== '') {
        tocData += '</div><hr class="main_hr">';
    }

    let footnoteData = '';
    const footnoteRegex = /<sup>((?:(?!<sup>|<\/sup>).)+)<\/sup>/g;
    const footnoteSearchData = [...data.matchAll(footnoteRegex)];
    let footnoteCount = 1;

    if (footnoteSearchData.length > 0) {
        footnoteData += '<div class="opennamu_footnote">';
    }

    footnoteSearchData.forEach((footnoteSearch) => {
        const footnoteCountStr = footnoteCount.toString();

        if (footnoteCount !== 1) {
            footnoteData += '<br>';
        }

        footnoteData += `<a id="fn-${footnoteCountStr}" href="#rfn-${footnoteCountStr}">(${footnoteCountStr})</a> ${footnoteSearch[1]}`;
        data = data.replace(
            footnoteRegex,
            `<sup id="rfn-${footnoteCountStr}"><a href="#fn-${footnoteCountStr}">(${footnoteCountStr})</a></sup>`
        );

        footnoteCount += 1;
    });

    if (footnoteData !== '') {
        footnoteData += '</div>';
    }

    data = tocData + data + footnoteData;

    return data;
}

function opennamu_do_id_check(data) {
    if(data.match(/\.|\:/)) {
        return 0;
    } else {
        return 1;
    }
}

function opennamu_do_ip_click(obj) {
    if (obj.id === "") {
        let user_name = obj.name;

        fetch('/api/v2/ip_menu/' + user_name)
            .then(response => {
                if (!response.ok) {
                    throw new Error(`API 호출 실패: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                data = data["data"];

                let data_html = '';
                for (let key in data) {
                    for (let for_a = 0; for_a < data[key].length; for_a++) {
                        data_html += '<a href="' + data[key][for_a][0] + '">' + data[key][for_a][1] + '</a> | ';
                    }
                }

                data_html = data_html.replace(/ \| $/g, '');

                let for_a;
                for (for_a = 0; document.getElementById("opennamu_ip_render_" + String(for_a) + "_load"); for_a++) {}

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
            })
            .catch(err => {
                console.error('IP 메뉴 호출 중 오류 발생:', err);
                obj.innerHTML = 'IP 정보를 불러오는 데 실패했습니다.';
            });
    }
}

function opennamu_do_ip_render() {
    for (let for_a = 0; for_a < document.getElementsByClassName('opennamu_render_ip').length; for_a++) {
        let ip = document.getElementsByClassName('opennamu_render_ip')[for_a].innerHTML.replace(/&amp;/g, '&');

        fetch('/api/v2/ip/' + opennamu_do_url_encode(ip))
            .then(response => {
                if (!response.ok) {
                    throw new Error(`API 호출 실패: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                if (document.getElementsByClassName('opennamu_render_ip')[for_a].id !== "opennamu_render_end") {
                    document.getElementsByClassName('opennamu_render_ip')[for_a].innerHTML = data["data"];
                    document.getElementsByClassName('opennamu_render_ip')[for_a].id = "opennamu_render_end";
                }
            })
            .catch(err => {
                console.error('IP 렌더링 호출 중 오류 발생:', err);
                document.getElementsByClassName('opennamu_render_ip')[for_a].innerHTML = 'IP 정보를 불러오는 데 실패했습니다.';
            });
    }
}


function opennamu_do_url_encode(data) {
    return encodeURIComponent(data);
}

function opennamu_cookie_split_regex(data) {
    return new RegExp('(?:^|; )' + data + '=([^;]*)');
}

function opennamu_send_render(data) {
    if(data === '&lt;br&gt;' || data === '' || data.match(/^ +$/)) {
        data = '<br>';
    } else {
        data = data.replace(/( |^)(https?:\/\/(?:[^ ]+))/g, function(m0, m1, m2) {
            let link_main = m2;
            link_main = link_main.replace('"', '&quot;');

            return m1 + '<a href="' + link_main + '">' + link_main + '</a>';
        });
        data = data.replace(/&lt;a(?:(?:(?!&gt;).)*)&gt;((?:(?!&lt;\/a&gt;).)+)&lt;\/a&gt;/g, function(m0, m1) {
            let data_unescape = opennamu_xss_filter_decode(m1)

            return '<a href="/w/' + opennamu_do_url_encode(data_unescape) + '">' + m1 + '</a>'
        })
    }

    return data;
}

function opennamu_insert_v(name, data) {
    document.getElementById(name).value = data;
}

function opennamu_do_trace_spread() {
    if(document.getElementsByClassName('opennamu_trace')) {
        document.getElementsByClassName('opennamu_trace')[0].innerHTML = '' +
            '<style>.opennamu_trace_button { display: none; } .opennamu_trace { white-space: pre-wrap; overflow-x: unset; text-overflow: unset; }</style>' +
        '' + document.getElementsByClassName('opennamu_trace')[0].innerHTML
    }
}

function opennamu_do_render(to_obj, data, name = '', do_type = '', option = '', callback = undefined) {
    let url;
    if (do_type === '') {
        url = "/api/render";
    } else {
        url = "/api/render/" + do_type;
    }

    fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({
            'name': name,
            'data': data,
            'option': option
        })
    })
        .then(response => {
            if (!response.ok) {
                throw new Error(`API 호출 실패: ${response.status}`);
            }
            return response.json();
        })
        .then(text => {
            if (document.getElementById(to_obj)) {
                if (text["data"]) {
                    document.getElementById(to_obj).innerHTML = text["data"];
                    eval(text["js_data"]);
                } else {
                    document.getElementById(to_obj).innerHTML = '';
                }

                if (callback) {
                    callback();
                }
            }
        })
        .catch(err => {
            console.error('렌더링 호출 중 오류 발생:', err);
            if (document.getElementById(to_obj)) {
                document.getElementById(to_obj).innerHTML = '렌더링에 실패했습니다. 잠시 후 다시 시도해 주세요.';
            }
        });
}


function opennamu_page_control(url, page, data_length, data_length_max = 50) {
    let next = function() {
        if(data_length_max === data_length) {
            return '<a href="' + url.replace('{}', String(page + 1)) + '">(+)</a>';
        } else {
            return '';
        }
    };

    let back = function() {
        if(page !== 1) {
            return '<a href="' + url.replace('{}', String(page - 1)) + '">(-)</a>';
        } else {
            return '';
        }
    };

    return (back() + ' ' + next()).replace(/^ /, '');
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

function opennamu_make_list(left = '', right = '', bottom = '', class_name = '') {
    let data_html = '<span class="' + class_name + '">';
    data_html += '<div class="opennamu_recent_change">';
    data_html += left;
    
    data_html += '<div style="float: right;">';
    data_html += right;
    data_html += '</div>'

    data_html += '<div style="clear: both;"></div>';

    if(bottom !== "") {
        data_html += '<hr>'
        data_html += bottom;
    }

    data_html += '</div>';
    data_html += '<hr class="main_hr">';
    data_html += '</span>';

    return data_html;
}



