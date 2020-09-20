function main_css_regex_data(data) {
    return new RegExp('(?:^|; )' + data + '=([^;]*)');
}

function main_css_get_post() {    
    var check = document.getElementById('main_css_strike');
    if(check.value === 'normal') {
        document.cookie = 'main_css_del_strike=0;';
    } else if(check.value === 'change') {
        document.cookie = 'main_css_del_strike=1;';
    } else {
        document.cookie = 'main_css_del_strike=2;';
    }

    check = document.getElementById('main_css_bold');
    if(check.value === 'normal') {
        document.cookie = 'main_css_del_bold=0;';
    } else if(check.value === 'change') {
        document.cookie = 'main_css_del_bold=1;';
    } else {
        document.cookie = 'main_css_del_bold=2;';
    }

    check = document.getElementById('main_css_include');
    if(check.checked) {
        document.cookie = 'main_css_include_link=1;';
    } else {
        document.cookie = 'main_css_include_link=0;';
    }

    check = document.getElementById('main_css_category');
    if(check.value === 'bottom') {
        document.cookie = 'main_css_category_set=0;';
    } else {
        document.cookie = 'main_css_category_set=1;';
    }

    check = document.getElementById('main_css_footnote');
    if(check.value === 'spread') {
        document.cookie = 'main_css_footnote_set=1;';
    } else {
        document.cookie = 'main_css_footnote_set=0;';
    }

    check = document.getElementById('main_css_image');
    if(check.value === 'new_click') {
        document.cookie = 'main_css_image_set=2;';
    } else if(check.value === 'click') {
        document.cookie = 'main_css_image_set=1;';
    } else {
        document.cookie = 'main_css_image_set=0;';
    }

    check = document.getElementById('main_css_image_paste');
    if(check.checked) {
        document.cookie = 'main_css_image_paste=1;';
    } else {
        document.cookie = 'main_css_image_paste=0;';
    }

    check = document.getElementById('main_css_toc');
    if(check.value === 'on') {
        document.cookie = 'main_css_toc_set=2;';
    } else if(check.value === 'off') {
        document.cookie = 'main_css_toc_set=1;';
    } else {
        document.cookie = 'main_css_toc_set=0;';
    }

    check = document.getElementById('main_css_font_size');
    if(check.value.match(/^[0-9]+$/)) {
        document.cookie = 'main_css_font_size=' + check.value + ';';
    } else {
        document.cookie = 'main_css_font_size=;';
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

    if(
        document.cookie.match(main_css_regex_data('main_css_category_set')) &&
        document.cookie.match(main_css_regex_data('main_css_category_set'))[1] === '1'
    ) {
        var get_category = document.getElementById('cate_all');
        if(get_category) {
            var backup_category = get_category.innerHTML;
            var in_data = document.getElementById('in_data_0').innerHTML;
            get_category.innerHTML = '';
            document.getElementById('in_data_0').innerHTML = backup_category + in_data;
            head_data.innerHTML += '<style>#cate { margin-top: 0px; margin-bottom: 20px; }</style>';
        }
    }

    if(
        document.cookie.match(main_css_regex_data('main_css_font_size')) &&
        document.cookie.match(main_css_regex_data('main_css_font_size'))[1] !== ''
    ) {
        head_data.innerHTML += '<style>.all_in_data { font-size: ' + document.cookie.match(main_css_regex_data('main_css_font_size'))[1] + 'px; }</style>';
    }

    if(document.cookie.match(main_css_regex_data('main_css_toc_set'))) {
        if(document.cookie.match(main_css_regex_data('main_css_toc_set'))[1] === '2') {
            head_data.innerHTML += '<style>#auto_toc { display: none; }</style>';
        } else if(document.cookie.match(main_css_regex_data('main_css_toc_set'))[1] === '1') {
            head_data.innerHTML += '<style>#toc { display: none; }</style>';
        }
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
            "font_size" : "font size"
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
            "font_size" : "글자 크기"
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

    if(document.cookie.match(main_css_regex_data('main_css_font_size'))) {
        set_data["font_size"] = document.cookie.match(main_css_regex_data('main_css_font_size'))[1];
    } else {
        set_data["font_size"] = '';
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
        <input ' + set_data["image_paste"] + ' type="checkbox" id="main_css_image_paste" value="image_paste"> 클립보드 이미지 업로드 (ko-KR) \
        <h3>1.7. ' + main_css_load_lang('set_toc') + '</h3> \
        <select id="main_css_toc"> \
            ' + set_data["toc"] + ' \
        </select> \
        <h3>1.8. ' + main_css_load_lang('set_font_size') + '</h3> \
        <input id="main_css_font_size" placeholder="' + main_css_load_lang('font_size') + '" value="' + set_data["font_size"] + '"> \
        <hr class="main_hr"> \
        <button onclick="main_css_get_post();">' + main_css_load_lang('save') + '</button> \
    ';
}