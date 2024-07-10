"use strict";

function do_insert_data(data) {
    const name = 'opennamu_edit_textarea';

    if(get_select_editor() === 'textarea') {
        // https://stackoverflow.com/questions/11076975/insert-text-into-textarea-at-cursor-position-javascript
        if(document.selection) {
            document.getElementById(name).focus();

            let sel = document.selection.createRange();
            sel.text = data;
        } else if(
            document.getElementById(name).selectionStart || 
            document.getElementById(name).selectionStart === '0'
        ) {
            let startPos = document.getElementById(name).selectionStart;
            let endPos = document.getElementById(name).selectionEnd;
            let myPos = document.getElementById(name).value;

            document.getElementById(name).value = myPos.substring(0, startPos) + data + myPos.substring(endPos, myPos.length);
        } else {
            document.getElementById(name).value += data;
        }
    } else {
        let selection = editor.getSelection();
        let id = { major: 1, minor: 1 };             
        let text = data;
        let op = {
            identifier: id, 
            range: selection, 
            text: text, 
            forceMoveMarkers: true
        };
        
        editor.executeEdits("my-source", [op]);
    }
}

// 아직 개편이 더 필요함
function do_paste_image() {
    document.getElementById('opennamu_edit_textarea').addEventListener("paste", pasteListener);
}

function pasteListener(e) {
    let lang_data = new FormData();
    lang_data.append('data', 'file_name empty save authority_error same_file_error error');

    if(e.clipboardData && e.clipboardData.items) {
        const items = e.clipboardData.items;
        const formData = new FormData();

        let haveImageInClipboard = false;
        let file_name = '';
        let file;
        
        for(let i = 0; i < items.length; i++) {
            if(items[i].type.indexOf("image") !== -1) {
                file = items[i].getAsFile();
                
                haveImageInClipboard = true;
                e.preventDefault();
                
                break;
            }
        }

        if(!haveImageInClipboard) {
            return;
        }

        fetch('/api/lang', {
            method : 'POST',
            body : lang_data,
        }).then(function(res) {
            return res.json();
        }).then(function(ajax_data) {
            const customName = prompt(ajax_data['data'][0]);
                
            if(!customName) {
                return alert(ajax_data['data'][1]);
            }
            
            file_name = customName + ".png";
            
            const customFile = new File([file], file_name, { type: file.type });
            formData.append("f_data[]", customFile);

            fetch("/upload", {
                method : "POST",
                body : formData,
            }).then((res) => {
                if (res.status === 200 || res.status === 201) {
                    const url = res.url;
                    alert(ajax_data['data'][2] + ' : [[file:' + file_name + ']]');

                    do_insert_data('[[file:' + file_name + ']]');
                } else {
                    console.error("[ERROR] PasteUpload Fail :", res.statusText);

                    if(res.status === 400) {
                        alert(ajax_data['data'][4]);
                    } else if(res.status === 401) {
                        alert(ajax_data['data'][3]);    
                    } else {
                        alert(ajax_data['data'][5]);        
                    }
                }
            }).catch((err) => {
                console.error("[ERROR] PasteUpload Fail :", JSON.stringify(err), err);

                alert(ajax_data['data'][5]);
            });
        });
    }
}

function do_stop_exit() {
    window.onbeforeunload = function() {
        do_sync_monaco_and_textarea();

        let data = document.getElementById('opennamu_edit_textarea').value;
        let origin = document.getElementById('opennamu_edit_origin').value;
        if(data !== origin) {
            return '';
        }
    }
}

function do_stop_exit_release() {
    do_sync_monaco_and_textarea();
    
    window.onbeforeunload = function () {}
}

function opennamu_edit_turn_off_monaco() {
    let now_selected = get_select_editor();
    let editor_list = [
        ['opennamu_edit_textarea', 'none'], 
        ['opennamu_monaco_editor', 'none']
    ];

    if(now_selected === 'textarea') {
        editor_list[0][1] = 'block';
    } else if(now_selected === 'monaco') {
        editor_list[1][1] = 'block';
    } else {
    }

    for(let for_a = 0; for_a < editor_list.length; for_a++) {
        document.getElementById(editor_list[for_a][0]).style.display = editor_list[for_a][1];
    }
}

function do_monaco_to_textarea(set_value) {
    document.getElementById('opennamu_edit_textarea').value = set_value;
}

function do_textarea_to_monaco(set_value) {
    window.editor.setValue(set_value);
}

function get_select_editor() {
    let now_selected = document.getElementById("opennamu_select_editor").value;
    if(now_selected === 'default') {
        return 'textarea';
    } else if(now_selected === 'monaco') {
        return 'monaco';
    } else {
        return '';
    }
}

function get_select_editor_markup() {
    let now_selected = document.getElementById("opennamu_editor_markup").value;
    if(now_selected === 'namumark' || now_selected === 'namumark_beta') {
        return 'namumark';
    } else if(now_selected === 'markdown') {
        return 'markdown';
    } else {
        return 'plaintext';
    }
}

function do_sync_monaco_and_textarea(select = '') {
    let now_selected = get_select_editor();
    if(select === 'textarea_to' || now_selected === 'textarea') {
        let set_value = document.getElementById('opennamu_edit_textarea').value;
        do_textarea_to_monaco(set_value);
    } else if(now_selected === 'monaco') {
        let set_value = window.editor.getValue();
        do_monaco_to_textarea(set_value);
    } else {

    }
}

// https://github.com/microsoft/monaco-editor/issues/568
class PlaceholderContentWidget {
    static ID = 'editor.widget.placeholderHint';

    constructor(placeholder, editor) {
        this.placeholder = placeholder;
        this.editor = editor;
        editor.onDidChangeModelContent(() => this.onDidChangeModelContent());
        this.onDidChangeModelContent();
    }

    onDidChangeModelContent() {
        if(this.editor.getValue() === '') {
            this.editor.addContentWidget(this);
        } else {
            this.editor.removeContentWidget(this);
        }
    }

    getId() {
        return PlaceholderContentWidget.ID;
    }

    getDomNode() {
        if(!this.domNode) {
            this.domNode = document.createElement('div');
            this.domNode.style.width = 'max-content';
            this.domNode.textContent = this.placeholder;
            this.domNode.style.fontStyle = 'italic';
            this.editor.applyFontInfo(this.domNode);
        }

        return this.domNode;
    }

    getPosition() {
        return {
            position: { lineNumber: 1, column: 1 },
            preference: [monaco.editor.ContentWidgetPositionPreference.EXACT],
        };
    }

    dispose() {
        this.editor.removeContentWidget(this);
    }
}

function do_monaco_init(monaco_thema) {
    require.config({ paths: { 'vs' : 'https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.48.0/min/vs' }});
    require.config({ 'vs/nls' : { availableLanguages: { '*' : 'ko' } }});
    require(["vs/editor/editor.main"], function () {
        monaco.languages.register({ id : "namumark" });
        monaco.languages.setMonarchTokensProvider("namumark", {
            tokenizer : {
                root : [
                    [/\[/, "namumark-color"],
                    [/\]/, "namumark-color"],
                    
                    [/\{/, "namumark-color"],
                    [/\}/, "namumark-color"],

                    [/'/, "namumark-color"],
                    [/-/, "namumark-color"],
                    [/~/, "namumark-color"],
                    [/=/, "namumark-color"],
                    [/_/, "namumark-color"],
                    [/\^/, "namumark-color"],
                    [/,/, "namumark-color"],

                    [/\\/, "namumark-color"],
                    [/\*/, "namumark-color"],
                ],
            },
        });

        let thema_set = [["namumark", "vs"], ["namumark-vs-dark", "vs-dark"]]
        for(let for_a = 0; for_a < thema_set.length; for_a++) {
            monaco.editor.defineTheme(thema_set[for_a][0], {
                base : thema_set[for_a][1],
                inherit : true,
                rules : [
                    { token : "namumark-color", foreground : "d94844" },
                ],
                colors : {},
            });
        }

        window.editor = monaco.editor.create(document.getElementById('opennamu_monaco_editor'), {
            value : document.getElementById('opennamu_edit_textarea').value,
            language : 'namumark',
            automaticLayout : true,
            wordWrap : true,
            theme : "namumark" + (monaco_thema === "" ? "" : "-" + monaco_thema)
        });

        if(typeof opennamu_monaco_custom === 'function') {
            opennamu_monaco_custom();
        }

        new PlaceholderContentWidget(document.getElementById('opennamu_edit_textarea').placeholder, window.editor);

        opennamu_do_sync_monaco_markup();
    });
}

function opennamu_do_editor_preview() {
    do_sync_monaco_and_textarea();

    const input = document.querySelector('#opennamu_edit_textarea');
    if(input !== null) {
        let name = "test";
        if(document.getElementById('opennamu_editor_doc_name')) {
            name = opennamu_xss_filter_decode(document.getElementById('opennamu_editor_doc_name').value);
        }

        opennamu_do_render('opennamu_preview_area', input.value, name);
    }
}

function opennamu_do_sync_monaco_markup() {
    let now_selected = get_select_editor_markup();
    monaco.editor.setModelLanguage(window.editor.getModel(), now_selected);
}

function opennamu_do_editor_temp_save() {
    do_sync_monaco_and_textarea();

    const input = document.querySelector('#opennamu_edit_textarea');
    if(input !== null) {
        localStorage.setItem("key", input.value);
    }
}

function opennamu_do_editor_temp_save_load() {
    const data = localStorage.getItem("key");
    if(data !== null) {
        const input = document.querySelector('#opennamu_edit_textarea');
        if(input !== null) {
            input.value = data;
        }
        
        do_sync_monaco_and_textarea('textarea_to');
    }
}

function opennamu_do_user_editor_insert() {
    let data = prompt();
    if(data !== null && data !== "") {
        let form_data = new FormData();
        form_data.append('data', data);

        fetch('/api/v2/user/setting/editor', {
            method : 'POST',
            body : form_data,
        }).then(function() {
            opennnamu_do_user_editor();
        });
    }
}

function opennamu_do_user_editor_delete() {
    let data = prompt();
    if(data !== null && data !== "") {
        let form_data = new FormData();
        form_data.append('data', data);

        fetch('/api/v2/user/setting/editor', {
            method : 'DELETE',
            body : form_data,
        }).then(function() {
            opennnamu_do_user_editor();
        });
    }
}

function opennnamu_do_user_editor() {
    fetch('/api/v2/user/setting/editor').then(function(res) {
        return res.json();
    }).then(function(data) {
        if(data["response"] === "ok") {
            let data_html = '';

            for(let for_a = 0; for_a < data["data"].length; for_a++) {
                data_html += '<a href="javascript:do_insert_data(\'' + opennamu_xss_filter(data["data"][for_a]) + '\');">(' + opennamu_xss_filter(data["data"][for_a]) + ')</a> ';
            }

            data_html += '<a href="javascript:opennamu_do_user_editor_insert();">(+)</a> ';
            data_html += '<a href="javascript:opennamu_do_user_editor_delete();">(-)</a>';
            data_html += '<hr class="main_hr">';

            document.getElementById("opennamu_editor_user_button").innerHTML = data_html;
        }
    });
}