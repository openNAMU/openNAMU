from .tool.func import *

def view_w_raw(name = '', rev = '', doc_acl = ''):
    with get_db_connect() as conn:
        rev_str = str(rev)

        sub = '(' + load_lang('raw') + ')'
        sub += ' (' + rev_str + ')' if rev != '' else ''

        if rev != '':
            menu = [['history_tool/' + rev_str + '/' + url_pas(name), load_lang('return')]]
        else:
            menu = [['w/' + url_pas(name), load_lang('return')]]

        p_data = ''
        p_data += '''
            <div id="opennamu_preview_area">
                <textarea id="opennamu_editor_doc_name" style="display: none;">''' + html.escape(name) + '''</textarea>
                <textarea id="opennamu_editor_rev" style="display: none;">''' + rev_str + '''</textarea>
                <button id="opennamu_preview_button" type="button" onclick="opennamu_view_w_raw_preview();">''' + load_lang('preview') + '''</button>
                <hr class="main_hr">
                <textarea readonly id="opennamu_edit_textarea" class="opennamu_textarea_500"></textarea>
                <script>opennamu_view_w_raw();</script>
            </div>
        '''
        
        if doc_acl == 'on':
            p_data = '' + \
                load_lang('authority_error') + \
                '<hr class="main_hr">' + \
                p_data
            ''
            
            sub = ' (' + load_lang('edit') + ')'

        return easy_minify(flask.render_template(skin_check(),
            imp = [name, wiki_set(), wiki_custom(), wiki_css([sub, 0])],
            data = p_data,
            menu = menu
        ))