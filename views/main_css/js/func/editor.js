"use strict";

const opennamu_monaco_worker_base = 'https://cdn.jsdelivr.net/npm/monaco-editor@0.56.0/esm/vs';
const opennamu_monaco_codicon_font = 'https://cdn.jsdelivr.net/npm/monaco-editor@0.56.0/esm/vs/base/browser/ui/codicons/codicon/codicon.ttf';

if(document.getElementById('opennamu_monaco_editor') !== null) {
    window.MonacoEnvironment = window.MonacoEnvironment || {};
    window.MonacoEnvironment.getWorker = function(module_id, label) {
        let worker_path = '/editor/editor.worker.js';

        if(label === 'json') {
            worker_path = '/language/json/json.worker.js';
        } else if(label === 'css' || label === 'scss' || label === 'less') {
            worker_path = '/language/css/css.worker.js';
        } else if(label === 'html' || label === 'handlebars' || label === 'razor') {
            worker_path = '/language/html/html.worker.js';
        } else if(label === 'typescript' || label === 'javascript') {
            worker_path = '/language/typescript/ts.worker.js';
        }

        const worker_name = label || 'editorWorkerService';
        const worker_url = opennamu_monaco_worker_base + worker_path;
        const worker_blob = URL.createObjectURL(new Blob([
            'import ' + JSON.stringify(worker_url) + ';'
        ], {
            type: 'text/javascript',
        }));
        const worker = new Worker(worker_blob, {
            type: 'module',
            name: worker_name,
        });

        setTimeout(function() {
            URL.revokeObjectURL(worker_blob);
        }, 60000);

        return worker;
    };

    window.opennamu_monaco_ready = import('https://cdn.jsdelivr.net/npm/monaco-editor@0.56.0/+esm').then(function(monaco) {
        document.querySelectorAll('style').forEach(function(style) {
            style.textContent = style.textContent.split('url(./codicon.ttf)').join('url(' + opennamu_monaco_codicon_font + ')');
        });
        window.monaco = monaco;
        return monaco;
    }).catch(function(error) {
        console.error('Monaco ESM load failed:', error);
        throw error;
    });
}

function do_insert_data(data) {
    const name = 'opennamu_edit_textarea';

    if(get_select_editor() === 'textarea' || window.editor === undefined || window.editor === null) {
        let textarea = document.getElementById(name);
        textarea.focus();

        let startPos = textarea.selectionStart;
        let endPos = textarea.selectionEnd;
        let myPos = textarea.value;

        textarea.value = myPos.substring(0, startPos) + data + myPos.substring(endPos, myPos.length);
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
    const textarea = document.getElementById('opennamu_edit_textarea');
    textarea.addEventListener("paste", pasteListener);
    textarea.addEventListener("dragover", function(e) {
        e.preventDefault();
    });
    textarea.addEventListener("drop", pasteListener);
}

function pasteListener(e) {
    const clipboard_data = e.clipboardData || e.dataTransfer;
    if(clipboard_data && clipboard_data.items) {
        const items = clipboard_data.items;
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

        let lang_data = new FormData();
        lang_data.append('data', 'file_name empty save error');

        fetch('/api/v2/lang', {
            method : 'POST',
            body : lang_data,
        }).then(function(res) {
            return res.json();
        }).then(function(lang) {
            lang = lang["data"];

            const customName = prompt(lang['file_name']);
                
            if(!customName) {
                return alert(lang['empty']);
            }
            
            file_name = customName + ".png";
            
            const customFile = new File([file], file_name, { type: file.type });
            formData.append("f_data[]", customFile);

            fetch("/upload", {
                method : "POST",
                body : formData,
            }).then((res) => {
                if (res.ok) {
                    const url = res.url;
                    alert(lang['save'] + ' : [[file:' + file_name + ']]');

                    do_insert_data('[[file:' + file_name + ']]');
                } else {
                    return res.text().then(function(data) {
                        const body = new DOMParser().parseFromString(data, 'text/html');
                        const error = body.querySelector('.opennamu_main li');
                        alert(error ? error.textContent.trim() : data.trim() || lang['error']);
                    });
                }
            }).catch((err) => {
                console.error("[ERROR] PasteUpload Fail :", JSON.stringify(err), err);

                alert(lang['error']);
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

    if(now_selected === 'monaco' && window.editor !== undefined && window.editor !== null) {
        window.requestAnimationFrame(function() {
            window.editor.layout();
        });
    }
}

function do_monaco_to_textarea(set_value) {
    document.getElementById('opennamu_edit_textarea').value = set_value;
}

function do_textarea_to_monaco(set_value) {
    if(window.editor === undefined || window.editor === null) {
        return;
    }
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
    if(window.editor === undefined || window.editor === null) {
        return;
    }
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
    static ID = 'opennamu.editor.widget.placeholderHint';

    constructor(placeholder, editor) {
        this.placeholder = placeholder;
        this.editor = editor;
        this.content_widget_visible = false;
        editor.onDidChangeModelContent(() => this.onDidChangeModelContent());
        this.onDidChangeModelContent();
    }

    onDidChangeModelContent() {
        if(this.editor.getValue() === '') {
            if(!this.content_widget_visible) {
                this.editor.addContentWidget(this);
                this.content_widget_visible = true;
            }
        } else if(this.content_widget_visible) {
            this.editor.removeContentWidget(this);
            this.content_widget_visible = false;
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
        if(this.content_widget_visible) {
            this.editor.removeContentWidget(this);
            this.content_widget_visible = false;
        }
    }
}

function do_monaco_init(monaco_thema) {
    if(window.opennamu_monaco_ready === undefined || window.opennamu_monaco_ready === null) {
        console.error('Monaco ESM loader is not available.');
        return;
    }

    window.opennamu_monaco_ready.then(function(monaco) {
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
    }).catch(function(error) {
        console.error('Monaco ESM initialization failed:', error);
    });
}


function opennamu_do_editor_preview() {
    do_sync_monaco_and_textarea();

    const input = document.getElementById('opennamu_edit_textarea');
    const preview = document.getElementById('opennamu_preview_area');
    const doc_name = document.getElementById('opennamu_editor_doc_name');
    if(input === null || preview === null) {
        return;
    }

    const body = new URLSearchParams({
        name: doc_name === null ? 'test' : doc_name.value,
        data: input.value,
        option: '',
    });

    fetch('/api/render', {
        method: 'POST',
        headers: {'Content-Type': 'application/x-www-form-urlencoded'},
        body: body,
    }).then(function(response) {
        if(!response.ok) {
            throw new Error('render failed: ' + response.status);
        }
        return response.json();
    }).then(function(result) {
        preview.innerHTML = result.data || '';
    }).catch(function(error) {
        console.error('Preview failed:', error);
        preview.textContent = 'Preview failed.';
    });
}

function opennamu_do_sync_monaco_markup() {
    if(window.editor === undefined || window.editor === null || window.monaco === undefined) {
        return;
    }
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
