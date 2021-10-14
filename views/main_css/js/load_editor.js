function do_insert_data(name, data, monaco = 0) {
    if(monaco === 0) {
        // https://stackoverflow.com/questions/11076975/insert-text-into-textarea-at-cursor-position-javascript
        if(document.selection) {
            document.getElementById(name).focus();

            var sel = document.selection.createRange();
            sel.text = data;
        } else if(
            document.getElementById(name).selectionStart || 
            document.getElementById(name).selectionStart == '0'
        ) {
            var startPos = document.getElementById(name).selectionStart;
            var endPos = document.getElementById(name).selectionEnd;
            var myPos = document.getElementById(name).value;

            document.getElementById(name).value = myPos.substring(0, startPos) + data + myPos.substring(endPos, myPos.length);
        } else {
            document.getElementById(name).value += data;
        }
    } else {
        var selection = editor.getSelection();
        var id = { major: 1, minor: 1 };             
        var text = data;
        var op = {
            identifier: id, 
            range: selection, 
            text: text, 
            forceMoveMarkers: true
        };
        
        editor.executeEdits("my-source", [op]);
    }
}

function monaco_to_content() {
    try {
        document.getElementById('textarea_edit_view').value = window.editor.getValue();
    } catch(e) {}
}

function do_not_out() {
    window.addEventListener('DOMContentLoaded', function() {
        window.onbeforeunload = function() {
            monaco_to_content();
            section_edit_do();
            
            data = document.getElementById('content').value;
            origin = document.getElementById('origin').value;
            if(data !== origin) {
                return '';
            }
        }
    });
}

function save_stop_exit() {
    window.onbeforeunload = function () {}
}

function do_paste_image() {
    window.addEventListener('DOMContentLoaded', function() {
        if(
            document.cookie.match(main_css_regex_data('main_css_image_paste')) &&
            document.cookie.match(main_css_regex_data('main_css_image_paste'))[1] === '1'
        ) {
            const textarea = document.querySelector("textarea");
            if (textarea) {
                textarea.addEventListener("paste", pasteListener);
            }
        }
    });
}

function pasteListener(e) {
    // find file
    if(e.clipboardData && e.clipboardData.items) {
        const items = e.clipboardData.items;
        let haveImageInClipboard = false;
        const formData = new FormData();
        for(let i = 0; i < items.length; i++) {
            if(items[i].type.indexOf("image") !== -1) {
                const file = items[i].getAsFile();
                const customName = prompt("파일 이름 (확장자 제외)");
                
                if (!customName) {
                    return alert("파일 이름 없음");
                }
                
                var file_name = customName + ".png";
                const customFile = new File([file], file_name, { type: file.type });
                formData.append("f_data[]", customFile);
                haveImageInClipboard = true;
                e.preventDefault();
                
                break;
            }
        }

        if(!haveImageInClipboard) {
            return;
        }

        // send to server
        fetch("/upload", {
            method: "POST",
            body: formData,
        }).then((res) => {
            if (res.status === 200 || res.status === 201) {
                const url = res.url;
                alert(
                    '업로드 완료 : ' +
                    '[[파일:' + file_name + ']]'
                );
            } else {
                console.error("[ERROR] PasteUpload Fail :", res.statusText);
                if(res.status === 400) {
                    alert("파일 이름 중복");
                } else if(res.status === 401) {
                    alert("권한 부족");    
                } else {
                    alert("업로드 실패");        
                }
            }
        }).catch((err) => {
            console.error("오류 내역 :", JSON.stringify(err), err);
            alert("업로드 실패");
        });
    }
}

function load_preview(name) {
    var s_data = new FormData();
    s_data.append('data', document.getElementById('textarea_edit_view').value);

    var url = "/api/w/" + name;
    var url_2 = "/api/markup";

    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.send(s_data);

    var xhr_2 = new XMLHttpRequest();
    xhr_2.open("GET", url_2, true);
    xhr_2.send();

    xhr.onreadystatechange = function() {
        if(xhr.readyState === 4 && xhr.status === 200) {
            var o_p_data = JSON.parse(xhr.responseText);
            document.getElementById('see_preview').innerHTML = o_p_data['data'];
            eval(o_p_data['js_data'])
        }
    }
}

function load_raw_preview(name_1, name_2) {
    document.getElementById(name_2).innerHTML = document.getElementById(name_1).value;
}

function section_edit_init() {
    var data_server = JSON.parse(
        document.getElementById('server_set').innerHTML
    );
    
    if(data_server['markup'] === 'namumark') {
        var data = document.getElementById('textarea_edit_view').value;
        var data_org = data;
        var data_section = Number(data_server['section']);
        var re_heading = /(^|\n)(={1,6})(#)? ?([^=]+) ?#?={1,6}(\n|$)/;
        for(i = 1; data.match(re_heading); i++) {
            if(i === data_section) {
                var start_point = data.search(re_heading);
                if(data[start_point] === '\n') {
                    start_point += 1;
                }
                
                data = data.replace(re_heading, function(x) {
                    return '.'.repeat(x.length - 1) + '\n';
                });
                
                var end_point = data.search(re_heading);
                if(end_point === -1) {
                    end_point = data.length;
                }
                
                data = data_org.slice(start_point, end_point);
                data = data.replace(/\n$/, '');
                
                document.getElementById('textarea_edit_view').value = data;
                
                data_server['start_point'] = start_point;
                data_server['end_point'] = end_point;
                
                document.getElementById('server_set').innerHTML = JSON.stringify(data_server);
                
                break;
            } else {
                data = data.replace(re_heading, function(x) {
                    return '.'.repeat(x.length - 1) + '\n';
                });
            }
        }
    }
}

function section_edit_do() {
    var data_server = JSON.parse(
        document.getElementById('server_set').innerHTML
    );
    
    if(data_server['start_point'] !== undefined) {
        var data = document.getElementById('origin').value;
        var data_section = document.getElementById('textarea_edit_view').value;
        
        var start_point = data_server['start_point'];
        var end_point = data_server['end_point'];
        
        if(data.length >= end_point) {
            var data_new = '';
            data_new += data.slice(0, start_point);
            data_new += data_section;
            data_new += data.slice(end_point, data.length);
            
            document.getElementById('content').value = data_new;
        }
    } else {
        document.getElementById('content').value = document.getElementById('textarea_edit_view').value;
    }
}