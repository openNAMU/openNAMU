from .tool.func import *

def view_w_raw(name = '', rev = '', doc_acl = ''):
    with get_db_connect() as conn:
        rev_str = str(rev)

        sub = '(' + get_lang(conn, 'raw') + ')'
        sub += ' (' + rev_str + ')' if rev != '' else ''

        if rev != '':
            menu = [['history_tool/' + rev_str + '/' + url_pas(name), get_lang(conn, 'return')]]
        else:
            menu = [['w/' + url_pas(name), get_lang(conn, 'return')]]

        p_data = ''
        p_data += '''
            <div id="opennamu_preview_area">
                <script defer src="/views/main_css/js/route/w_raw.js''' + cache_v() + '''"></script>
                <textarea id="opennamu_editor_doc_name" style="display: none;">''' + html.escape(name) + '''</textarea>
                <textarea id="opennamu_editor_rev" style="display: none;">''' + rev_str + '''</textarea>
                <button id="opennamu_preview_button" type="button" onclick="opennamu_w_raw_preview();">''' + get_lang(conn, 'preview') + '''</button>
                <hr class="main_hr">
                <textarea readonly id="opennamu_edit_textarea" class="opennamu_textarea_500"></textarea>
                <script>window.addEventListener("DOMContentLoaded", function() { opennamu_w_raw(); });</script>
            </div>
        '''
        
        if doc_acl == 'on':
            p_data = '' + \
                get_lang(conn, 'authority_error') + \
                '<hr class="main_hr">' + \
                p_data
            ''
            
            sub = ' (' + get_lang(conn, 'edit') + ')'

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [name, wiki_set(conn), wiki_custom(conn), wiki_css([sub, 0])],
            data = p_data,
            menu = menu
        ))