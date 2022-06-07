function main_css_regex_data(data) {
    return new RegExp('(?:^|; )' + data + '=([^;]*)');
}

function main_css_get_post() {    
    var check = document.getElementById('main_css_strike');
    if(check.value === 'normal') {
        document.cookie = 'main_css_del_strike=0; path=/;';
    } else if(check.value === 'change') {
        document.cookie = 'main_css_del_strike=1; path=/;';
    } else {
        document.cookie = 'main_css_del_strike=2; path=/;';
    }

    check = document.getElementById('main_css_bold');
    if(check.value === 'normal') {
        document.cookie = 'main_css_del_bold=0; path=/;';
    } else if(check.value === 'change') {
        document.cookie = 'main_css_del_bold=1; path=/;';
    } else {
        document.cookie = 'main_css_del_bold=2; path=/;';
    }

    check = document.getElementById('main_css_include');
    if(check.checked) {
        document.cookie = 'main_css_include_link=1; path=/;';
    } else {
        document.cookie = 'main_css_include_link=0; path=/;';
    }

    check = document.getElementById('main_css_category');
    if(check.value === 'bottom') {
        document.cookie = 'main_css_category_set=0; path=/;';
    } else {
        document.cookie = 'main_css_category_set=1; path=/;';
    }

    check = document.getElementById('main_css_footnote');
    if(check.value === 'spread') {
        document.cookie = 'main_css_footnote_set=1; path=/;';
    } else {
        document.cookie = 'main_css_footnote_set=0; path=/;';
    }

    check = document.getElementById('main_css_image');
    if(check.value === 'new_click') {
        document.cookie = 'main_css_image_set=2; path=/;';
    } else if(check.value === 'click') {
        document.cookie = 'main_css_image_set=1; path=/;';
    } else {
        document.cookie = 'main_css_image_set=0; path=/;';
    }

    check = document.getElementById('main_css_image_paste');
    if(check.checked) {
        document.cookie = 'main_css_image_paste=1; path=/;';
    } else {
        document.cookie = 'main_css_image_paste=0; path=/;';
    }

    check = document.getElementById('main_css_toc');
    if(check.value === 'on') {
        document.cookie = 'main_css_toc_set=2; path=/;';
    } else if(check.value === 'off') {
        document.cookie = 'main_css_toc_set=1; path=/;';
    } else {
        document.cookie = 'main_css_toc_set=0; path=/;';
    }
    
    check = document.getElementById('main_css_font_size');
    if(check.value.match(/^[0-9]+$/)) {
        document.cookie = 'main_css_font_size=' + check.value + '; path=/;';
    } else {
        document.cookie = 'main_css_font_size=; path=/;';
    }

    check = document.getElementById('main_css_monaco');
    if(check.checked) {
        document.cookie = 'main_css_monaco=1; path=/;';
    } else {
        document.cookie = 'main_css_monaco=0; path=/;';
    }
    
    check = document.getElementById('main_css_exter_link');
    if(check.value === 'self') {
        document.cookie = 'main_css_exter_link=1; path=/;';
    } else {
        document.cookie = 'main_css_exter_link=0; path=/;';
    }
    
    check = document.getElementById('main_css_link_delimiter');
    if(check.checked) {
        document.cookie = 'main_css_link_delimiter=1; path=/;';
    } else {
        document.cookie = 'main_css_link_delimiter=0; path=/;';
    }
    
    history.go(0);
}

function main_css_skin_load() {    
    var head_data = document.querySelector('head');
    if(document.cookie.match(main_css_regex_data('main_css_del_strike'))) {
        if(document.cookie.match(main_css_regex_data('main_css_del_strike'))[1] === '1') {
            head_data.innerHTML += '<style>s { text-decoration: none; } s:hover { background-color: transparent; }</style>';
        } else if(document.cookie.match(main_css_regex_data('main_css_del_strike'))[1] === '2') {
            head_data.innerHTML += '<style>s { display: none; }</style>';
        }
    }

    if(document.cookie.match(main_css_regex_data('main_css_del_bold'))) {
        if(document.cookie.match(main_css_regex_data('main_css_del_bold'))[1] === '1') {
            head_data.innerHTML += '<style>b { font-weight: normal; }</style>';
        } else if(document.cookie.match(main_css_regex_data('main_css_del_bold'))[1] === '2') {
            head_data.innerHTML += '<style>b { display: none; }</style>';
        }
    }

    if(
        document.cookie.match(main_css_regex_data('main_css_include_link')) &&
        document.cookie.match(main_css_regex_data('main_css_include_link'))[1] === '1'
    ) {
        head_data.innerHTML += '<style>#include_link { display: inline; }</style>';
    }

    if(document.cookie.match(main_css_regex_data('main_css_toc_set'))) {
        if(document.cookie.match(main_css_regex_data('main_css_toc_set'))[1] === '2') {
            head_data.innerHTML += '<style>#auto_toc { display: none; }</style>';
        } else if(document.cookie.match(main_css_regex_data('main_css_toc_set'))[1] === '1') {
            head_data.innerHTML += '<style>#toc { display: none; }</style>';
        }
    }
    
    if(
        document.cookie.match(main_css_regex_data('main_css_font_size')) &&
        document.cookie.match(main_css_regex_data('main_css_font_size'))[1] !== ''
    ) {
        head_data.innerHTML += '<style>body, input, textarea { font-size: ' + document.cookie.match(main_css_regex_data('main_css_font_size'))[1] + 'px; }</style>';
    }
    
    if(
        document.cookie.match(main_css_regex_data('main_css_darkmode')) &&
        document.cookie.match(main_css_regex_data('main_css_darkmode'))[1] === '1'
    ) {
        head_data.innerHTML += '' +
            '<link rel="stylesheet" href="/views/main_css/css/sub/dark.css?ver=5">' +
        '';
    }
    
    if(
        document.cookie.match(main_css_regex_data('main_css_link_delimiter')) &&
        document.cookie.match(main_css_regex_data('main_css_link_delimiter'))[1] === '1'
    ) {
        head_data.innerHTML += '<style>#real_normal_link::before, #not_thing::before, #inside::before { content: \'🅸\'; font-weight: lighter; background: transparent; }</style>';
    }
}

function main_css_load_lang(name) {
    var set_language = {
        "en-US" : {
            "default" : "Default",
            "change_to_normal" : "Change to normal text",
            "delete" : "Delete",
            "include_link" : "Using include link",
            "save" : "Save",
            "strike" : "Strike",
            "bold" : "Bold",
            "other" : "Other",
            "where_category" : "Set category location",
            "bottom" : "Bottom",
            "top" : "Top",
            "set_footnote" : "Set footnote",
            "renderer" : "Renderer",
            "spread" : "Spread",
            "set_image" : "Set image",
            "set_toc" : "Set TOC",
            "click_load" : "Load on click",
            "in_content" : "Only when TOC is in the document",
            "all_off" : "Always off",
            "set_font_size" : "Set font size",
            "change_to_link" : "Change to link",
            "font_size" : "font size",
            "editor" : "Editor",
            "main" : "Main",
            "clipboard_upload" : "Clipboard upload",
            "only_korean" : "Supported in korean only",
            "except_ie" : "Not supported for Internet Explorer",
            "use_monaco" : "Use monaco editor",
            "self_tab" : "Current tab",
            "exter_link_open_method" : "External link",
            "link_delimiter" : "Add link delimiter"
        }, "ko-KR" : {
            "default" : "기본값",
            "change_to_normal" : "일반 텍스트로 변경",
            "delete" : "삭제",
            "include_link" : "틀 링크 사용",
            "save" : "저장",
            "strike" : "취소선",
            "bold" : "볼드체",
            "other" : "기타",
            "where_category" : "분류 위치 설정",
            "bottom" : "아래",
            "top" : "위",
            "set_footnote" : "각주 설정",
            "renderer" : "렌더러",
            "spread" : "펼치기",
            "set_image" : "이미지 설정",
            "set_toc" : "목차 설정",
            "click_load" : "클릭시 불러오기",
            "in_content" : "문서 안에 있을 때만",
            "all_off" : "항상 끔",
            "set_font_size" : "글자 크기 설정",
            "change_to_link" : "링크로 변경",
            "font_size" : "글자 크기",
            "editor" : "편집기",
            "main" : "메인",
            "clipboard_upload" : "클립보드 파일 올리기",
            "only_korean" : "한국어로만 지원됨",
            "except_ie" : "인터넷 익스플로러에선 지원되지 않음",
            "use_monaco" : "모나코 에디터 사용",
            "self_tab" : "현재 탭",
            "exter_link_open_method" : "외부 링크",
            "link_delimiter" : "링크 구분자 추가"
        }
    }

    var server_language = document.cookie.match(main_css_regex_data('language'))[1];
    var user_language = document.cookie.match(main_css_regex_data('user_language'))[1];
    if(user_language in set_language) {
        language = user_language;
    } else {
        if(server_language in set_language) {
            language = server_language;
        } else {
            language = 'en-US';
        }
    }

    if(name in set_language[language]) {
        return set_language[language][name];
    } else {
        return name + ' (' + language + ')';
    }
}

function main_css_skin_set() {    
    var set_data = {};
    var strike_list = [
        ['0', 'normal', main_css_load_lang('default')],
        ['1', 'change', main_css_load_lang('change_to_normal')],
        ['2', 'delete', main_css_load_lang('delete')]
    ];
    set_data["strike"] = '';
    var i = 0;
    while(strike_list[i]) {
        if(
            document.cookie.match(main_css_regex_data('main_css_del_strike')) && 
            document.cookie.match(main_css_regex_data('main_css_del_strike'))[1] === strike_list[i][0]
        ) {
            set_data["strike"] = '<option value="' + strike_list[i][1] + '">' + strike_list[i][2] + '</option>' + set_data["strike"];
        } else {
            set_data["strike"] += '<option value="' + strike_list[i][1] + '">' + strike_list[i][2] + '</option>';
        }

        i += 1;
    }

    var bold_list = [
        ['0', 'normal', main_css_load_lang('default')],
        ['1', 'change', main_css_load_lang('change_to_normal')],
        ['2', 'delete', main_css_load_lang('delete')]
    ];
    set_data["bold"] = '';
    i = 0;
    while(bold_list[i]) {
        if(
            document.cookie.match(main_css_regex_data('main_css_del_bold')) && 
            document.cookie.match(main_css_regex_data('main_css_del_bold'))[1] === bold_list[i][0]
        ) {
            set_data["bold"] = '<option value="' + bold_list[i][1] + '">' + bold_list[i][2] + '</option>' + set_data["bold"];
        } else {
            set_data["bold"] += '<option value="' + bold_list[i][1] + '">' + bold_list[i][2] + '</option>';
        }

        i += 1;
    }

    if(
        document.cookie.match(main_css_regex_data('main_css_include_link')) &&
        document.cookie.match(main_css_regex_data('main_css_include_link'))[1] === '1'
    ) {
        set_data["include"] = "checked";
    } else {
        set_data["include"] = "";
    }

    if(
        document.cookie.match(main_css_regex_data('main_css_image_paste')) &&
        document.cookie.match(main_css_regex_data('main_css_image_paste'))[1] === '1'
    ) {
        set_data["image_paste"] = "checked";
    } else {
        set_data["image_paste"] = "";
    }

    var category_list = [
        ['0', 'bottom', main_css_load_lang('bottom')],
        ['1', 'top', main_css_load_lang('top')],
    ];
    set_data["category"] = '';
    i = 0;
    while(category_list[i]) {
        if(
            document.cookie.match(main_css_regex_data('main_css_category_set')) && 
            document.cookie.match(main_css_regex_data('main_css_category_set'))[1] === category_list[i][0]
        ) {
            set_data["category"] = '<option value="' + category_list[i][1] + '">' + category_list[i][2] + '</option>' + set_data["category"];
        } else {
            set_data["category"] += '<option value="' + category_list[i][1] + '">' + category_list[i][2] + '</option>';
        }

        i += 1;
    }

    var footnote_list = [
        ['0', 'normal', main_css_load_lang('default')],
        ['1', 'spread', main_css_load_lang('spread')]
    ];
    set_data["footnote"] = '';
    i = 0;
    while(footnote_list[i]) {
        if(
            document.cookie.match(main_css_regex_data('main_css_footnote_set')) && 
            document.cookie.match(main_css_regex_data('main_css_footnote_set'))[1] === footnote_list[i][0]
        ) {
            set_data["footnote"] = '<option value="' + footnote_list[i][1] + '">' + footnote_list[i][2] + '</option>' + set_data["footnote"];
        } else {
            set_data["footnote"] += '<option value="' + footnote_list[i][1] + '">' + footnote_list[i][2] + '</option>';
        }

        i += 1;
    }

    var image_list = [
        ['0', 'normal', main_css_load_lang('default')],
        ['1', 'click', main_css_load_lang('change_to_link')],
        ['2', 'new_click', main_css_load_lang('click_load')]
    ];
    set_data["image"] = '';
    i = 0;
    while(image_list[i]) {
        if(
            document.cookie.match(main_css_regex_data('main_css_image_set')) && 
            document.cookie.match(main_css_regex_data('main_css_image_set'))[1] === image_list[i][0]
        ) {
            set_data["image"] = '<option value="' + image_list[i][1] + '">' + image_list[i][2] + '</option>' + set_data["image"];
        } else {
            set_data["image"] += '<option value="' + image_list[i][1] + '">' + image_list[i][2] + '</option>';
        }

        i += 1;
    }

    var toc_list = [
        ['0', 'normal', main_css_load_lang('default')],
        ['1', 'off', main_css_load_lang('all_off')],
        ['2', 'on', main_css_load_lang('in_content')]
    ];
    set_data["toc"] = '';
    i = 0;
    while(toc_list[i]) {
        if(
            document.cookie.match(main_css_regex_data('main_css_toc_set')) && 
            document.cookie.match(main_css_regex_data('main_css_toc_set'))[1] === toc_list[i][0]
        ) {
            set_data["toc"] = '<option value="' + toc_list[i][1] + '">' + toc_list[i][2] + '</option>' + set_data["toc"];
        } else {
            set_data["toc"] += '<option value="' + toc_list[i][1] + '">' + toc_list[i][2] + '</option>';
        }

        i += 1;
    }

    if(
        document.cookie.match(main_css_regex_data('main_css_monaco')) &&
        document.cookie.match(main_css_regex_data('main_css_monaco'))[1] === '1'
    ) {
        set_data["monaco"] = "checked";
    } else {
        set_data["monaco"] = "";
    }
    
    if(document.cookie.match(main_css_regex_data('main_css_font_size'))) {
        set_data["font_size"] = document.cookie.match(main_css_regex_data('main_css_font_size'))[1];
    } else {
        set_data["font_size"] = '';
    }
    
    let exter_link_list = [
        ['0', 'blank', main_css_load_lang('default')],
        ['1', 'self', main_css_load_lang('self_tab')]
    ];
    set_data["exter_link"] = '';
    for(let i = 0; exter_link_list[i]; i++) {
        if(
            document.cookie.match(main_css_regex_data('main_css_exter_link')) && 
            document.cookie.match(main_css_regex_data('main_css_exter_link'))[1] === exter_link_list[i][0]
        ) {
            set_data["exter_link"] = '<option value="' + exter_link_list[i][1] + '">' + exter_link_list[i][2] + '</option>' + set_data["exter_link"];
        } else {
            set_data["exter_link"] += '<option value="' + exter_link_list[i][1] + '">' + exter_link_list[i][2] + '</option>';
        }
    }
    
    if(
        document.cookie.match(main_css_regex_data('main_css_link_delimiter')) &&
        document.cookie.match(main_css_regex_data('main_css_link_delimiter'))[1] === '1'
    ) {
        set_data["link_delimiter"] = "checked";
    } else {
        set_data["link_delimiter"] = "";
    }

    document.getElementById("main_skin_set").innerHTML = ' \
        <h2>1. ' + main_css_load_lang('renderer') + '</h2> \
        <h3>1.1. ' + main_css_load_lang('strike') + '</h3> \
        <select id="main_css_strike"> \
            ' + set_data["strike"] + ' \
        </select> \
        <h3>1.2. ' + main_css_load_lang('bold') + '</h3> \
        <select id="main_css_bold"> \
            ' + set_data["bold"] + ' \
        </select> \
        <h3>1.3. ' + main_css_load_lang('where_category') + '</h3> \
        <select id="main_css_category"> \
            ' + set_data["category"] + ' \
        </select> \
        <h3>1.4. ' + main_css_load_lang('set_footnote') + '</h3> \
        <select id="main_css_footnote"> \
            ' + set_data["footnote"] + ' \
        </select> \
        <h3>1.5. ' + main_css_load_lang('set_image') + '</h3> \
        <select id="main_css_image"> \
            ' + set_data["image"] + ' \
        </select> \
        <h3>1.6. ' + main_css_load_lang('other') + '</h3> \
        <input ' + set_data["include"] + ' type="checkbox" id="main_css_include" value="include"> ' + main_css_load_lang('include_link') + ' \
        <hr class="main_hr"> \
        <input ' + set_data["link_delimiter"] + ' type="checkbox" id="main_css_link_delimiter" value="link_delimiter"> ' + main_css_load_lang('link_delimiter') + '<sup>(1)</sup> \
        <h3>1.7. ' + main_css_load_lang('set_toc') + '</h3> \
        <select id="main_css_toc"> \
            ' + set_data["toc"] + ' \
        </select> \
        <h3>1.8. ' + main_css_load_lang('set_font_size') + '</h3> \
        <input id="main_css_font_size" placeholder="' + main_css_load_lang('font_size') + ' (EX : 11)" value="' + set_data["font_size"] + '"> \
        <h3>1.9. ' + main_css_load_lang('exter_link_open_method') + '</h3> \
        <select id="main_css_exter_link"> \
            ' + set_data["exter_link"] + ' \
        </select> \
        <h2>2. ' + main_css_load_lang('editor') + '</h2> \
        <h3>2.1. ' + main_css_load_lang('main') + '</h3> \
        <input ' + set_data["monaco"] + ' type="checkbox" id="main_css_monaco" value="monaco"> ' + main_css_load_lang('use_monaco') + '<sup>(1)</sup> \
        <hr class="main_hr"> \
        <input ' + set_data["image_paste"] + ' type="checkbox" id="main_css_image_paste" value="image_paste"> ' + 
            main_css_load_lang('clipboard_upload') + '<sup>(ko-KR)</sup><sup>(1)</sup> \
        <hr class="main_hr"> \
        <button onclick="main_css_get_post();">' + main_css_load_lang('save') + '</button> \
        <hr class="main_hr"> \
        <ul id="footnote_data"> \
            <li><a id="note_1_end" href="#note_1">(1)</a> ' + main_css_load_lang('except_ie') + '</li> \
            <li><a href="#note_1_1">(1.1)</a></li> \
            <li><a id="note_2_end" href="#note_2">(ko-KR)</a> ' + main_css_load_lang('only_korean') + '</li> \
        </ul> \
    ';
 
    simple_render('main_skin_set');
}

document.addEventListener("DOMContentLoaded", main_css_skin_load);