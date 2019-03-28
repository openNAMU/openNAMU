function get_post() {
    check = document.getElementById('strike');
    if(check.value === 'normal') {
        document.cookie = 'del_strike=0;';
    } else if(check.value === 'change') {
        document.cookie = 'del_strike=1;';
    } else {
        document.cookie = 'del_strike=2;';
    }

    check = document.getElementById('bold');
    if(check.value === 'normal') {
        document.cookie = 'del_bold=0;';
    } else if(check.value === 'change') {
        document.cookie = 'del_bold=1;';
    } else {
        document.cookie = 'del_bold=2;';
    }

    check = document.getElementById('include');
    if(check.checked === true) {
        document.cookie = 'include_link=true;';
    } else {
        document.cookie = 'include_link=false;';
    }

    history.go(0);
}

function regex_data(data) {
    r_data = new RegExp('(?:^|; )' + data + '=([^;]*)')

    return r_data;
}

cookies = document.cookie;

function main_load() {
    head_data = document.querySelector('head');
    if(cookies.match(regex_data('del_strike'))) {
        if(cookies.match(regex_data('del_strike'))[1] === '1') {
            head_data.innerHTML += '<style>s { text-decoration: none; } s:hover { background-color: transparent; }</style>';
        } else if(cookies.match(regex_data('del_strike'))[1] === '2') {
            head_data.innerHTML += '<style>s { display: none; }</style>';
        }
    }

    if(cookies.match(regex_data('del_bold'))) {
        if(cookies.match(regex_data('del_bold'))[1] === '1') {
            head_data.innerHTML += '<style>b { font-weight: normal; }</style>';
        } else if(cookies.match(regex_data('del_bold'))[1] === '2') {
            head_data.innerHTML += '<style>b { display: none; }</style>';
        }
    }

    if(
        cookies.match(regex_data('include_link')) &&
        cookies.match(regex_data('include_link'))[1] === 'true'
    ) {
        head_data.innerHTML += '<style>#include_link { display: inline; }</style>';
    }
}

main_load();

window.onload = function () {
    if(window.location.pathname === '/skin_set') {
        document.getElementById("main_top").innerHTML = '<h1>Skin settings</h1>';
        document.title = document.title.replace(/.*(\- .*)$/, "Skin settings $1");
        
        data = document.getElementById("main_data");
        set_data = {};

        if(cookies.match(regex_data('del_strike'))) {
            if(cookies.match(regex_data('del_strike'))[1] === '0') {
                set_data["strike"] = ' \
                    <option value="normal">Default</option> \
                    <option value="change">Change to normal text</option> \
                    <option value="delete">Delete</option> \
                ';
            } else if(cookies.match(regex_data('del_strike'))[1] === '1') {
                set_data["strike"] = ' \
                    <option value="change">Change to normal text</option> \
                    <option value="normal">Default</option> \
                    <option value="delete">Delete</option> \
                ';
            } else {
                set_data["strike"] = ' \
                    <option value="delete">Delete</option> \
                    <option value="normal">Default</option> \
                    <option value="change">Change to normal text</option> \
                ';
            }
        } else {
            set_data["strike"] = ' \
                <option value="normal">Default</option> \
                <option value="change">Change to normal text</option> \
                <option value="delete">Delete</option> \
            ';
        }

        if(cookies.match(regex_data('del_bold'))) {
            if(cookies.match(regex_data('del_bold'))[1] === '0') {
                set_data["bold"] = ' \
                    <option value="normal">Default</option> \
                    <option value="change">Change to normal text</option> \
                    <option value="delete">Delete</option> \
                ';
            } else if(cookies.match(regex_data('del_bold'))[1] === '1') {
                set_data["bold"] = ' \
                    <option value="change">Change to normal text</option> \
                    <option value="normal">Default</option> \
                    <option value="delete">Delete</option> \
                ';
            } else {
                set_data["bold"] = ' \
                    <option value="delete">Delete</option> \
                    <option value="normal">Default</option> \
                    <option value="change">Change to normal text</option> \
                ';
            }
        } else {
            set_data["bold"] = ' \
                <option value="normal">Default</option> \
                <option value="change">Change to normal text</option> \
                <option value="delete">Delete</option> \
            ';
        }
        
        if(
            cookies.match(regex_data('include_link')) &&
            cookies.match(regex_data('include_link'))[1] === 'true'
        ) {
            set_data["include"] = "checked";
        }

        data.innerHTML = ' \
            <h2>Strike</h2> \
            <hr class="main_hr"> \
            <select id="strike" name="strike"> \
                ' + set_data["strike"] + ' \
            </select> \
            <h2>Bold</h2> \
            <select id="bold" name="bold"> \
                ' + set_data["bold"] + ' \
            </select> \
            <h2>Other</h2> \
            <input ' + set_data["include"] + ' type="checkbox" id="include" name="include" value="include"> Using include link \
            <hr class="main_hr"> \
            <button onclick="get_post();">Save</button> \
        ';
    }
}