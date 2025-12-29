from .tool.func import *

from .go_api_func_language import api_func_language
from .go_api_setting import api_setting

async def main_setting_404_page():
    with get_db_connect() as conn:
        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 0)
            
        if flask.request.method == 'POST':
            select_data = flask.request.form.get('select', '404_page')
            form_data = flask.request.form.get('data', 'Test')

            await api_setting('manage_404_page', 'PUT', select_data)
            await api_setting('manage_404_page_content', 'PUT', form_data)

            await acl_check(tool = 'owner_auth', memo = 'edit_set (404_page)')

            return redirect(conn, '/setting/404_page')
        else:
            lang = await api_func_language('', 'enter_html save 404_file 404_page preview')

            set_type = await api_setting('manage_404_page')
            set_data = await api_setting('manage_404_page_content')

            data_html = ''
            select_list = [
                ['404_page', lang['data']['404_page']],
                ['404_file', lang['data']['404_file']]
            ]

            data_html += '<select class="__ON_SELECT__" name="select">'
            for for_a in select_list:
                selected = ''
                if set_type['data'] == for_a[0]:
                    selected = 'selected'
                    
                data_html += '<option value="' + for_a[0] + '" ' + selected + '>' + for_a[1] + '</option>'

            data_html += '</select>'
            data_html += '<hr class="main_hr">'

            form_data = ''
            if len(set_data['data']) != 0:
                form_data = set_data['data'][0][0]

            data_html += (
                '<form method="post">' +
                    '<textarea class="opennamu_textarea_500 __ON_TEXTAREA__" name="data" placeholder="' + lang['data']['enter_html'] + '">' + html.escape(form_data) + '</textarea>' +
                    '<hr class="main_hr">' +
                    '<button class="__ON_BUTTON__" id="opennamu_save_button" type="submit">' + lang['data']['save'] + '</button>' +
                '</form>'
            )

            return await render_template(
                await get_lang('404_page_setting'),
                data_html,
                0,
                [['setting', await get_lang('return')]]
            )
