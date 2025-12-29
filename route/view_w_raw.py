from .tool.func import *

from .go_api_w_raw import api_w_raw

async def view_w_raw(name = '', rev = '', doc_acl = ''):
    with get_db_connect() as conn:
        rev_str = str(rev)

        sub = '(' + await get_lang('raw') + ')'
        sub += ' (' + rev_str + ')' if rev != '' else ''

        if rev != '':
            menu = [['history_tool/' + rev_str + '/' + url_pas(name), await get_lang('return')]]
        else:
            menu = [['w/' + url_pas(name), await get_lang('return')]]

        data = await api_w_raw(name, rev)
        if data["response"] == "ok":
            data_in = data["data"]
        else:
            data_in = ''

        p_data = ''
        p_data += '''
            <div id="opennamu_preview_area">
                <textarea readonly id="opennamu_edit_textarea" class="opennamu_textarea_500 __ON_TEXTAREA__">''' + html.escape(data_in) + '''</textarea>
            </div>
        '''
        
        if doc_acl == 'on':
            p_data = '' + \
                await get_lang('authority_error') + \
                '<hr class="main_hr">' + \
                p_data
            ''
            
            sub = ' (' + await get_lang('edit') + ')'

        return await render_template(
            name,
            p_data,
            sub,
            menu
        )
