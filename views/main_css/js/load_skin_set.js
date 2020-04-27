function main_css_regex_data(data) {
    return new RegExp('(?:^|; )' + data + '=([^;]*)');
}

function main_css_get_post() {    
    var check = document.getElementById('strike');
    if(check.value === 'normal') {
        document.cookie = 'main_css_del_strike=0;';
    } else if(check.value === 'change') {
        document.cookie = 'main_css_del_strike=1;';
    } else {
        document.cookie = 'main_css_del_strike=2;';
    }

    check = document.getElementById('bold');
    if(check.value === 'normal') {
        document.cookie = 'main_css_del_bold=0;';
    } else if(check.value === 'change') {
        document.cookie = 'main_css_del_bold=1;';
    } else {
        document.cookie = 'main_css_del_bold=2;';
    }

    check = document.getElementById('include');
    if(check.checked) {
        document.cookie = 'main_css_include_link=1;';
    } else {
        document.cookie = 'main_css_include_link=0;';
    }

    check = document.getElementById('category');
    if(check.value === 'bottom') {
        document.cookie = 'main_css_category_set=0;';
    } else {
        document.cookie = 'main_css_category_set=1;';
    }

    check = document.getElementById('footnote');
    if(check.value === 'spread') {
        document.cookie = 'main_css_footnote_set=1;';
    } else {
        document.cookie = 'main_css_footnote_set=0;';
    }

    check = document.getElementById('image_paste');
    if(check.checked) {
        document.cookie = 'main_css_image_paste=1;';
    } else {
        document.cookie = 'main_css_image_paste=0;';
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

            backup_category = backup_category.replace('<hr>', '') + '<hr>';

            document.getElementById('in_data_0').innerHTML = backup_category + in_data;
        }
    }
}

function main_css_skin_set() {    
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
            "set_footnote" : "Set footnote"
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
            "set_footnote" : "각주 설정"
        }
    }

    var language = document.cookie.match(main_css_regex_data('language'))[1];
    var user_language = document.cookie.match(main_css_regex_data('user_language'))[1];
    if(user_language in set_language) {
        language = user_language;
    }

    if(!language in set_language) {
        language = "en-US";
    }

    var set_data = {};
    var strike_list = [
        ['0', 'normal', set_language[language]['default']],
        ['1', 'change', set_language[language]['change_to_normal']],
        ['2', 'delete', set_language[language]['delete']]
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
        ['0', 'normal', set_language[language]['default']],
        ['1', 'change', set_language[language]['change_to_normal']],
        ['2', 'delete', set_language[language]['delete']]
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
        ['0', 'bottom', set_language[language]['bottom']],
        ['1', 'top', set_language[language]['top']],
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
        ['0', 'normal', set_language[language]['default']],
        ['1', 'spread', 'spread (beta)']
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

    document.getElementById("main_skin_set").innerHTML = ' \
        <h2>1. renderer</h2> \
        <h3>1.1. ' + set_language[language]['strike'] + '</h3> \
        <select id="strike" name="strike"> \
            ' + set_data["strike"] + ' \
        </select> \
        <h3>1.2. ' + set_language[language]['bold'] + '</h3> \
        <select id="bold" name="bold"> \
            ' + set_data["bold"] + ' \
        </select> \
        <h3>1.3. ' + set_language[language]['where_category'] + '</h3> \
        <select id="category" name="category"> \
            ' + set_data["category"] + ' \
        </select> \
        <h3>1.4. ' + set_language[language]['set_footnote'] + '</h3> \
        <select id="footnote" name="footnote"> \
            ' + set_data["footnote"] + ' \
        </select> \
        <h3>1.4. ' + set_language[language]['other'] + '</h3> \
        <input ' + set_data["include"] + ' type="checkbox" id="include" name="include" value="include"> ' + set_language[language]['include_link'] + ' \
        <hr class="main_hr"> \
        <input ' + set_data["image_paste"] + ' type="checkbox" id="image_paste" name="image_paste" value="image_paste"> image_paste (beta) \
        <hr class="main_hr"> \
        <button onclick="main_css_get_post();">' + set_language[language]['save'] + '</button> \
    ';
}