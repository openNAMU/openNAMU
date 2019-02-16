// 쿠키 생성
function setCookie(name, value, expiredays) {
    var cookie = name + "=" + escape(value) + "; path=/;"
    if (typeof expiredays != 'undefined') {
        var todayDate = new Date();
        todayDate.setDate(todayDate.getDate() + expiredays);
        cookie += "expires=" + todayDate.toGMTString() + ";"
    }
    document.cookie = cookie;
}
 
// 쿠키 획득
function getCookie(name) {
    name += "=";
    var cookie = document.cookie;
    var startIdx = cookie.indexOf(name);
    if (startIdx != -1) {
        startIdx += name.length;
        var endIdx = cookie.indexOf(";", startIdx);
        if (endIdx == -1) {
            endIdx = cookie.length;
            return unescape(cookie.substring(startIdx, endIdx));
        }
    }
    return null;
}
 
// 쿠키 삭제
function deleteCookie(name) {
    setCookie(name, "", -1);
}

// http://vip00112.tistory.com/33

function get_post() {
    check = document.getElementById('strike');
    if(check.checked == true) {
        setCookie("set_strike", "1");
    } else {
        deleteCookie("set_strike");
    }
    
    window.location.reload(true);
}

head_data = document.querySelector('head');
if(getCookie("set_strike") == "1") {
    head_data.innerHTML += '<style>s { display: none; }';
}

window.onload = function () {
    if(window.location.pathname == '/skin_set') {
        document.getElementById("main_top").innerHTML = '<h1>skin setting</h1>';
        data = document.getElementById("main_data")
        
        set_data = {};
        if(getCookie("set_strike") == "1") {
            set_data["strike"] = "checked";
        } 
        
        data.innerHTML =    `
                            <input ` + set_data["strike"] + ` type="checkbox" id="strike" name="strike" value="strike"> remove strikethrough
                            <hr>
                            <button onclick="get_post();">Save</button>
                            `;
                            
        document.title = document.title.replace(/.*(\- .*)$/, "skin setting $1");
    }
}