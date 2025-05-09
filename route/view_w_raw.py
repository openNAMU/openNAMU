from .tool.func import *

from .go_api_w_raw import api_w_raw

async def view_w_raw(name = '', rev = '', doc_acl = ''):
    with get_db_connect() as conn:
        rev_str = str(rev)

        sub = '(' + get_lang(conn, 'raw') + ')'
        sub += ' (' + rev_str + ')' if rev != '' else ''

        if rev != '':
            menu = [['history_tool/' + rev_str + '/' + url_pas(name), get_lang(conn, 'return')]]
        else:
            menu = [['w/' + url_pas(name), get_lang(conn, 'return')]]

        data = await api_w_raw(name, rev)
        if data["response"] == "ok":
            data_in = data["data"]
        else:
            data_in = ''

        p_data = ''
        p_data += '''
            <div id="opennamu_preview_area">
                <textarea readonly id="opennamu_edit_textarea" class="opennamu_textarea_500">''' + html.escape(data_in) + '''</textarea>
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
            imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css([sub, 0])],
            data = p_data,
            menu = menu
        ))